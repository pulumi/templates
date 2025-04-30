import * as pulumi from "@pulumi/pulumi";
import * as pulumiservice from "@pulumi/pulumiservice";
import * as aws from "@pulumi/aws";

const ARCHIVE_BUCKET_PREFIX = "public-esc-rotator-lambdas-production";
const ARCHIVE_KEY = "aws-lambda/latest.zip";
const ARCHIVE_SIGNING_PROFILE_VERSION_ARN = "arn:aws:signer:us-west-2:388588623842:/signing-profiles/pulumi_esc_production_20250325212043887700000001/jva5X9nqMa";
const organization = pulumi.getOrganization()

// Load configs
const templateConfig = new pulumi.Config();
const awsConfig = new pulumi.Config("aws");
const awsRegion = awsConfig.require("region");
const rdsDbIdentifier = templateConfig.require("rdsDbIdentifier");
const rotatorEnvironmentName = templateConfig.require("rotatedSecretsEnvironmentName");
const credsEnvironmentName = templateConfig.get("managingCredsEnvironmentName") ?? rotatorEnvironmentName + "ManagingCreds";
const backendUrl = templateConfig.get("backendUrl") ?? "https://api.pulumi.com";
const oidcUrl = new URL(`oidc`, backendUrl).toString();

// Parse environment names
const envNameSplit = rotatorEnvironmentName.split("/");
if (envNameSplit.length != 2) {
    throw Error(`Invalid environmentName supplied "${rotatorEnvironmentName}" - needs to be in format "myProject/myEnvironment"`)
}
const environment = {
    organization: organization,
    project: envNameSplit[0],
    name: envNameSplit[1],
};
const credsNameSplit = credsEnvironmentName.split("/");
if (credsNameSplit.length != 2) {
    throw Error(`Invalid managingCredsEnvironmentName supplied "${credsEnvironmentName}" - needs to be in format "myProject/myEnvironment"`)
}
const credsEnvironment = {
    organization: organization,
    project: credsNameSplit[0],
    name: credsNameSplit[1],
};

// Retrieve reference to current code artifact from trusted pulumi bucket
const lambdaArchiveBucket = `${ARCHIVE_BUCKET_PREFIX}-${awsRegion}`
const codeArtifact = aws.s3.getObjectOutput({bucket: lambdaArchiveBucket, key: ARCHIVE_KEY});

// Introspect RDS to discover network settings
const database = aws.rds.getClusterOutput({
    clusterIdentifier: rdsDbIdentifier,
});
const subnetGroup = aws.rds.getSubnetGroupOutput({
    name: database.dbSubnetGroupName,
});
const databaseSecurityGroupId = database.vpcSecurityGroupIds[0];
const databasePort = database.port;
const vpcId = subnetGroup.vpcId;
let validatedSubnetIds = subnetGroup.subnetIds.apply(async ids => {
    let subnetIds: string[] = [];
    for (const id of ids) {
        await aws.ec2.getSubnet({id: id}, {async: false}).then(
            _ => subnetIds.push(id),
            _ => console.log("bad subnet found: "+id),
        );
    }
    return subnetIds;
});

// Create resources
const namePrefix = "PulumiEscSecretConnectorLambda-"
const codeSigningConfig = new aws.lambda.CodeSigningConfig(namePrefix + "CodeSigningConfig", {
    description: "Pulumi ESC rotation connector lambda signature - https://github.com/pulumi/esc-rotator-lambdas",
    allowedPublishers: {
        signingProfileVersionArns: [ARCHIVE_SIGNING_PROFILE_VERSION_ARN],
    },
    policies: {
        untrustedArtifactOnDeployment: "Enforce",
    },
});
const lambdaExecRole = new aws.iam.Role(namePrefix + "ExecutionRole", {
    assumeRolePolicy: JSON.stringify({
        Version: "2012-10-17",
        Statement: [{
            Action: "sts:AssumeRole",
            Effect: "Allow",
            Principal: {
                Service: "lambda.amazonaws.com",
            },
        }],
    }),
    managedPolicyArns: ["arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"],
});
const lambdaSecurityGroup = new aws.ec2.SecurityGroup(namePrefix + "SecurityGroup", {
    vpcId: vpcId,
    description: "Security group for Pulumi ESC rotation lambda",
});
const lambdaEgressRule = new aws.ec2.SecurityGroupRule(namePrefix + "ToDatabaseEgressRule", {
    description: "Allow connections to database",
    type: "egress",
    protocol: "tcp",
    fromPort: databasePort,
    toPort: databasePort,
    securityGroupId: lambdaSecurityGroup.id,
    sourceSecurityGroupId: databaseSecurityGroupId,
});
const databaseIngressRule = new aws.ec2.SecurityGroupRule(namePrefix + "FromDatabaseIngressRule", {
    description: "Allow connections from rotation lambda",
    type: "ingress",
    protocol: "tcp",
    fromPort: databasePort,
    toPort: databasePort,
    sourceSecurityGroupId: lambdaSecurityGroup.id,
    securityGroupId: databaseSecurityGroupId,
});
const lambda = new aws.lambda.Function(namePrefix + "Function", {
    description: "The connector lambda proxies a secret rotation request from Pulumi ESC to a service within your VPC.",
    s3Bucket: codeArtifact.bucket,
    s3Key: codeArtifact.key,
    s3ObjectVersion: codeArtifact.versionId,
    codeSigningConfigArn: codeSigningConfig.arn,
    runtime: "provided.al2023",
    handler: "bootstrap",
    role: lambdaExecRole.arn,
    vpcConfig: {
        subnetIds: validatedSubnetIds,
        securityGroupIds: [lambdaSecurityGroup.id],
    },
});
let oidcProviderArn: pulumi.Output<String>
const oidcAudience = "aws:"+organization;
const oidcUrlNoProtocol = oidcUrl.replace("https://", "");
oidcProviderArn = pulumi.output(
    aws.iam.getOpenIdConnectProvider({ url: oidcUrl }, { async: false })
    .then(
        res => { 
            if (!res.clientIdLists.includes(oidcAudience)) {
                throw Error(`Unable to create OIDC identity provider, because OIDC provider for ${oidcUrlNoProtocol} already exists for the AWS Account.
                    Please manually add "${oidcAudience}" to the list of audiences within the ${oidcUrlNoProtocol} identity provider`)
            }
            return res.arn 
        },
        _ => {
            return new aws.iam.OpenIdConnectProvider(namePrefix + "OidcProvider", {
                url: oidcUrl,
                clientIdLists: [oidcAudience],
            }, {
                retainOnDelete: true,
            }).arn;
        }
    )
);
const assumedRole = new aws.iam.Role(namePrefix + "InvocationRole", {
    description: "Allow Pulumi ESC to invoke/manage the connector lambda",
    assumeRolePolicy: pulumi.jsonStringify({
        Version: "2012-10-17",
        Statement: [{
            Action: "sts:AssumeRoleWithWebIdentity",
            Effect: "Allow",
            Principal: {
                Federated: oidcProviderArn,
            },
            Condition: {
                StringEquals: {
                    [`${oidcUrlNoProtocol}:aud`]: oidcAudience,
                },
            }
        }],
    }),
    inlinePolicies: [{
        policy: pulumi.jsonStringify({
            Version: "2012-10-17",
            Statement: [
                {
                    Sid: "AllowPulumiToInvokeLambda",
                    Effect: "Allow",
                    Action: [
                        "lambda:GetFunction",
                        "lambda:InvokeFunction",
                    ],
                    Resource: lambda.arn,
                },
                {
                    Sid: "AllowPulumiToUpdateLambda",
                    Effect: "Allow",
                    Action: "lambda:UpdateFunctionCode",
                    Resource: lambda.arn,
                },
                {
                    Sid: "AllowPulumiToFetchUpdatedLambdaArchives",
                    Effect: "Allow",
                    Action: "s3:GetObject",
                    Resource: `arn:aws:s3:::${lambdaArchiveBucket}/*`,
                },
            ],
        }),
    }],
});
const psp = new pulumiservice.Provider(namePrefix + "PSP", {
    apiUrl: backendUrl
})
const credsYaml = pulumi.interpolate
    `values:
       managingUser:
         username: managing_user # Replace with your user value
         # Replace ciphertext below with your password, keeping fn::secret to encrypt it, like so "fn::secret: <password>"
         password:
           fn::secret: manager_password
       awsLogin:
         fn::open::aws-login:
           oidc:
             duration: 1h
             roleArn: ${assumedRole.arn}
             sessionName: pulumi-esc-secret-rotator`
const creds = new pulumiservice.Environment(namePrefix + "RotatorEnvironmentManagingCreds", {
    organization: credsEnvironment.organization,
    project: credsEnvironment.project,
    name: credsEnvironment.name,
    yaml: credsYaml,
}, {
    deleteBeforeReplace: true,
    provider: psp,
})
const rotatorType = databasePort.apply(port => port === 5432 ? "postgres" : "mysql");
const managingUserImport = "${environments." + `${credsEnvironment.project}.${credsEnvironment.name}.managingUser}`
const awsLoginImport = "${environments." + `${credsEnvironment.project}.${credsEnvironment.name}.awsLogin}`
const yaml = pulumi.interpolate
    `values:
       dbRotator:
         fn::rotate::${rotatorType}:
           inputs:
             database:
               connector:
                 awsLambda:
                   login: ${awsLoginImport}
                   lambdaArn: ${lambda.arn}
               database: rotator_db # Replace with your DB name
               host: ${database.endpoint}
               port: ${databasePort}
               managingUser: ${managingUserImport}
             rotateUsers:
               username1: user1 # Replace with your user value
               username2: user2 # Replace with your user value`
const _ = new pulumiservice.Environment(namePrefix + "RotatorEnvironment", {
    organization: environment.organization,
    project: environment.project,
    name: environment.name,
    yaml: yaml,
}, {
    deleteBeforeReplace: true,
    dependsOn: creds,
    provider: psp,
})

export const lambdaArn = lambda.arn;
export const assumedRoleArn = assumedRole.arn;