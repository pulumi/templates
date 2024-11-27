import * as pulumi from "@pulumi/pulumi";
import * as ovh from "@ovhcloud/pulumi-ovh";

// Get configuration values or use defaults
const config = new pulumi.Config();
const ovhServiceName = config.require("ovhServiceName");
const ovhRegion = config.get("ovhRegion") || "GRA9";
const clusterName = config.get("clusterName") || "my_desired_cluster";
const nodePoolName = config.get("nodePoolName") || "my-desired-pool";

const nodePoolDesiredNodes = config.getNumber("nodePoolDesiredNodes") || 1;
const nodePoolMaxNodes = config.getNumber("nodePoolMaxNodes") || 3;
const nodePoolMinNodes = config.getNumber("nodePoolMinNodes") || 1;
const flavorName = config.get("flavorName") || "b3-8";

// Deploy a new Kubernetes cluster
const myKubeCluster = new ovh.cloudproject.Kube(clusterName, {
    region: ovhRegion,
    serviceName: ovhServiceName,
    name: clusterName
});

// Export kubeconfig file
export const kubeconfig = pulumi.secret(myKubeCluster.kubeconfig)

//Create a Node Pool
const nodePool = new ovh.cloudproject.KubeNodePool(nodePoolName, {
    desiredNodes: nodePoolDesiredNodes,
    flavorName: flavorName,
    kubeId: myKubeCluster.id,
    maxNodes: nodePoolMaxNodes,
    minNodes: nodePoolMinNodes,
    serviceName: ovhServiceName,
});

export const nodePoolID = nodePool.id;
