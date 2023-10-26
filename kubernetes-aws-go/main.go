package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/eks"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	awsx "github.com/pulumi/pulumi-awsx/sdk/v2/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get some configuration values or set default values
		cfg := config.New(ctx, "")
		minClusterSize, err := cfg.TryInt("minClusterSize")
		if err != nil {
			minClusterSize = 3
		}
		maxClusterSize, err := cfg.TryInt("maxClusterSize")
		if err != nil {
			maxClusterSize = 6
		}
		desiredClusterSize, err := cfg.TryInt("desiredClusterSize")
		if err != nil {
			desiredClusterSize = 3
		}
		eksNodeInstanceType, err := cfg.Try("eksNodeInstanceType")
		if err != nil {
			eksNodeInstanceType = "t3.medium"
		}
		vpcNetworkCidr, err := cfg.Try("vpcNetworkCidr")
		if err != nil {
			vpcNetworkCidr = "10.0.0.0/16"
		}

		// Create a new VPC, subnets, and associated infrastructure
		eksVpc, err := awsx.NewVpc(ctx, "eks-vpc", &awsx.VpcArgs{
			CidrBlock:          &vpcNetworkCidr,
			EnableDnsHostnames: pulumi.Bool(true),
			EnableDnsSupport:   pulumi.Bool(true),
			SubnetSpecs: []awsx.SubnetSpecArgs{
				{
					Type: awsx.SubnetTypePrivate,
					Tags: pulumi.StringMap{
						"kubernetes.io/role/internal-elb": pulumi.String("1"),
					},
				},
				{
					Type: awsx.SubnetTypePublic,
					Tags: pulumi.StringMap{
						"kubernetes.io/role/elb": pulumi.String("1"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Get policy documents for cluster and node assume role statements
		// First for the cluster
		clusterAssumeRole, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
			Statements: []iam.GetPolicyDocumentStatement{
				{
					Effect: pulumi.StringRef("Allow"),
					Principals: []iam.GetPolicyDocumentStatementPrincipal{
						{
							Type: "Service",
							Identifiers: []string{
								"eks.amazonaws.com",
							},
						},
					},
					Actions: []string{
						"sts:AssumeRole",
					},
				},
			},
		}, nil)
		if err != nil {
			return err
		}
		// Second for the nodes
		nodeAssumeRole, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
			Statements: []iam.GetPolicyDocumentStatement{
				{
					Effect: pulumi.StringRef("Allow"),
					Principals: []iam.GetPolicyDocumentStatementPrincipal{
						{
							Type: "Service",
							Identifiers: []string{
								"ec2.amazonaws.com",
							},
						},
					},
					Actions: []string{
						"sts:AssumeRole",
					},
				},
			},
		}, nil)
		if err != nil {
			return err
		}

		// Define the cluster IAM role
		clusterIamRole, err := iam.NewRole(ctx, "cluster-iam-role", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(clusterAssumeRole.Json),
		})
		if err != nil {
			return err
		}

		// Define the node IAM role
		nodeIamRole, err := iam.NewRole(ctx, "node-iam-role", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(nodeAssumeRole.Json),
		})
		if err != nil {
			return err
		}

		// Attach the cluster IAM role to necessary policies
		clusterPolicies := []string{
			"arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
			"arn:aws:iam::aws:policy/AmazonEKSVPCResourceController",
			"arn:aws:iam::aws:policy/AmazonEKSServicePolicy",
		}
		for i, policy := range clusterPolicies {
			_, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("cluster-pa-%d", i), &iam.RolePolicyAttachmentArgs{
				PolicyArn: pulumi.String(policy),
				Role:      clusterIamRole.Name,
			})
			if err != nil {
				return err
			}
		}

		// Attach the node IAM role to necessary policies
		nodePolicies := []string{
			"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
			"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
			"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
			"arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy",
		}
		for i, policy := range nodePolicies {
			_, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("node-pa-%d", i), &iam.RolePolicyAttachmentArgs{
				PolicyArn: pulumi.String(policy),
				Role:      nodeIamRole.Name,
			})
			if err != nil {
				return err
			}
		}

		// Create a Security Group that we can use to actually connect to our cluster
		clusterSg, err := ec2.NewSecurityGroup(ctx, "cluster-sg", &ec2.SecurityGroupArgs{
			VpcId: eksVpc.VpcId,
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(80),
					ToPort:     pulumi.Int(80),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(443),
					ToPort:     pulumi.Int(443),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		})
		if err != nil {
			return err
		}

		// Create an EKS cluster
		eksCluster, err := eks.NewCluster(ctx, "eks-cluster", &eks.ClusterArgs{
			RoleArn: clusterIamRole.Arn,
			VpcConfig: &eks.ClusterVpcConfigArgs{
				// Change these values for a private cluster (VPN access required)
				EndpointPrivateAccess: pulumi.Bool(false),
				EndpointPublicAccess:  pulumi.Bool(true),
				SecurityGroupIds:      pulumi.StringArray{clusterSg.ID()},
				SubnetIds:             eksVpc.PrivateSubnetIds,
			},
		})
		if err != nil {
			return err
		}

		// Create a node group for the EKS cluster
		_, err = eks.NewNodeGroup(ctx, "node-group", &eks.NodeGroupArgs{
			ClusterName:   eksCluster.Name,
			InstanceTypes: pulumi.StringArray{pulumi.String(eksNodeInstanceType)},
			NodeRoleArn:   nodeIamRole.Arn,
			SubnetIds:     eksVpc.PrivateSubnetIds,
			ScalingConfig: &eks.NodeGroupScalingConfigArgs{
				DesiredSize: pulumi.Int(desiredClusterSize),
				MaxSize:     pulumi.Int(maxClusterSize),
				MinSize:     pulumi.Int(minClusterSize),
			},
			UpdateConfig: &eks.NodeGroupUpdateConfigArgs{
				MaxUnavailable: pulumi.Int(1),
			},
		})
		if err != nil {
			return err
		}

		// Install the AWS EBS CSI addon
		_, err = eks.NewAddon(ctx, "aws-ebs-csi", &eks.AddonArgs{
			ClusterName:              eksCluster.Name,
			AddonName:                pulumi.String("aws-ebs-csi-driver"),
			AddonVersion:             pulumi.String("v1.24.0-eksbuild.1"),
			ResolveConflictsOnUpdate: pulumi.String("PRESERVE"),
		})
		if err != nil {
			return err
		}

		// Generate a Kubeconfig to access the cluster and make it accessible
		clusterKubeconfig := generateKubeconfig(eksCluster.Endpoint, eksCluster.CertificateAuthority.Data().Elem(), eksCluster.Name)
		ctx.Export("kubeconfig", clusterKubeconfig)
		ctx.Export("vpcId", eksVpc.VpcId)

		return nil
	})
}

// Create the KubeConfig structure as per https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html
func generateKubeconfig(clusterEndpoint pulumi.StringOutput, certData pulumi.StringOutput, clusterName pulumi.StringOutput) pulumi.StringOutput {
	return pulumi.Sprintf(`{
        "apiVersion": "v1",
        "clusters": [{
            "cluster": {
                "server": "%s",
                "certificate-authority-data": "%s"
            },
            "name": "kubernetes",
        }],
        "contexts": [{
            "context": {
                "cluster": "kubernetes",
                "user": "aws",
            },
            "name": "aws",
        }],
        "current-context": "aws",
        "kind": "Config",
        "users": [{
            "name": "aws",
            "user": {
                "exec": {
                    "apiVersion": "client.authentication.k8s.io/v1beta1",
                    "command": "aws-iam-authenticator",
                    "args": [
                        "token",
                        "-i",
                        "%s",
                    ],
                },
            },
        }],
    }`, clusterEndpoint, certData, clusterName)
}
