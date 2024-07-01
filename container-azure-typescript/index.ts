import * as containerinstance from "@pulumi/azure-native/containerinstance";
import * as containerregistry from "@pulumi/azure-native/containerregistry";
import * as dockerBuild from "@pulumi/docker-build";
import * as pulumi from "@pulumi/pulumi";
import * as random from "@pulumi/random";
import * as resources from "@pulumi/azure-native/resources";

// Import the program's configuration settings.
const config = new pulumi.Config();
const appPath = config.get("appPath") || "./app";
const imageName = config.get("imageName") || "my-app";
const imageTag = config.get("imageTag") || "latest";
const containerPort = config.getNumber("containerPort") || 80;
const cpu = config.getNumber("cpu") || 1;
const memory = config.getNumber("memory") || 2;

// Create a resource group for the container registry.
const resourceGroup = new resources.ResourceGroup("resource-group");

// Create a container registry.
const registry = new containerregistry.Registry("registry", {
    resourceGroupName: resourceGroup.name,
    adminUserEnabled: true,
    sku: {
        name: containerregistry.SkuName.Basic,
    },
});

// Fetch login credentials for the registry.
const credentials = containerregistry.listRegistryCredentialsOutput({
    resourceGroupName: resourceGroup.name,
    registryName: registry.name,
}).apply(creds => {
    return {
        username: creds.username!,
        password: creds.passwords![0].value!,
    };
});

// Create a container image for the service.
const image = new dockerBuild.Image("image", {
    push: true,
    tags: [pulumi.interpolate`${registry.loginServer}/${imageName}:${imageTag}`],
    platforms: [dockerBuild.Platform.Linux_amd64],
    context: {
        location: appPath
    },
    registries: [{
        address: registry.loginServer,
        username: credentials.username,
        password: credentials.password,
    }],
});

// Use a random string to give the service a unique DNS name.
const dnsName = new random.RandomString("dns-name", {
    length: 8,
    special: false,
}).result.apply((result: string) => `${imageName}-${result.toLowerCase()}`);

// Create a container group for the service that makes it publicly accessible.
const containerGroup = new containerinstance.ContainerGroup("container-group", {
    resourceGroupName: resourceGroup.name,
    osType: "linux",
    restartPolicy: "always",
    imageRegistryCredentials: [
        {
            server: registry.loginServer,
            username: credentials.username,
            password: credentials.password,
        },
    ],
    containers: [
        {
            name: imageName,
            image: image.ref,
            ports: [
                {
                    port: containerPort,
                    protocol: "tcp",
                },
            ],
            environmentVariables: [
                {
                    name: "PORT",
                    value: containerPort.toString(),
                },
            ],
            resources: {
                requests: {
                    cpu: cpu,
                    memoryInGB: memory,
                },
            },
        },
    ],
    ipAddress: {
        type: containerinstance.ContainerGroupIpAddressType.Public,
        dnsNameLabel: dnsName,
        ports: [
            {
                port: containerPort,
                protocol: "tcp",
            },
        ],
    },
});

// Export the service's IP address, hostname, and fully-qualified URL.
export const hostname = containerGroup.ipAddress.apply((addr: any) => addr!.fqdn!);
export const ip = containerGroup.ipAddress.apply((addr: any) => addr!.ip!);
export const url = containerGroup.ipAddress.apply((addr: any) => `http://${addr!.fqdn!}:${containerPort}`);
