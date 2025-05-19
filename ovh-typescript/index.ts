import * as pulumi from "@pulumi/pulumi";
import * as ovh from "@ovhcloud/pulumi-ovh";

// Get configuration values or use defaults
const config = new pulumi.Config();
const ovhServiceName = config.require("ovhServiceName");
const ovhRegion = config.get("ovhRegion") || "GRA";
const planName = config.get("planName") || "SMALL";

const registryName = config.get("registryName") || "my-registry";
const registryUserName = config.get("registryUserName") || "user";
const registryUserEmail = config.get("registryUserEmail") || "myuser@ovh.com";
const registryUserLogin = config.get("registryUserLogin") || "myuser";

// Initiate the configuration of the registry
const regcap = ovh.cloudproject.getCapabilitiesContainerFilter({
    serviceName: ovhServiceName,
    planName: planName,
    region: ovhRegion,
});

// Deploy a new Managed private registry
const myRegistry = new ovh.cloudproject.ContainerRegistry(registryName, {
    serviceName: regcap.then((regcap: { serviceName: any; }) => regcap.serviceName),
    planId: regcap.then((regcap: { id: any; }) => regcap.id),
    region: regcap.then((regcap: { region: any; }) => regcap.region),
});

// Create a Private Registry User
const myRegistryUser = new ovh.cloudproject.ContainerRegistryUser(registryUserName, {
    serviceName: ovhServiceName,
    registryId: myRegistry.id,
    email: registryUserEmail,
    login: registryUserLogin,
});

// Add as an output registry information
export const registryURL = myRegistry.url;
export const registryUser = myRegistryUser.user;
export const registryPassword = myRegistryUser.password;
