import * as pulumi from "@pulumi/pulumi";
import * as resources from "@pulumi/azure-native/resources";
import * as containerregistry from "@pulumi/azure-native/containerregistry";
import * as containerinstance from "@pulumi/azure-native/containerinstance";
import * as random from "@pulumi/random";
import * as docker from "@pulumi/docker";

const config = new pulumi.Config();
const imageName = config.get("imageName") || "my-app";
const appPath = config.get("appPath") || "./app";
const containerPort = config.getNumber("containerPort") || 80;

const resourceGroup = new resources.ResourceGroup("resource-group");

const registry = new containerregistry.Registry("registry", {
    resourceGroupName: resourceGroup.name,
    adminUserEnabled: true,
    sku: {
        name: containerregistry.SkuName.Basic,
    },
});

const credentials = containerregistry.listRegistryCredentialsOutput({
    resourceGroupName: resourceGroup.name,
    registryName: registry.name,
}).apply(creds => {
    return {
        username: creds.username!,
        password: creds.passwords![0].value!,
    };
});

const image = new docker.Image("image", {
    imageName: pulumi.interpolate`${registry.loginServer}/${imageName}`,
    build: {
        context: appPath,
    },
    registry: {
        server: registry.loginServer,
        username: credentials.username,
        password: credentials.password,
    },
});

const dnsName = new random.RandomString("dns-name", {
    length: 8,
    special: false,
}).result.apply(result => `${imageName}-${result.toLowerCase()}`);

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
            image: image.imageName,
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
                    cpu: 1.0,
                    memoryInGB: 1.5,
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

export const ipAddress = containerGroup.ipAddress.apply(addr => addr!.ip!);
export const hostname = containerGroup.ipAddress.apply(addr => addr!.fqdn!);
export const url = containerGroup.ipAddress.apply(addr => `http://${addr!.fqdn!}:${containerPort}`);
