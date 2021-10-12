import * as pulumi from "@pulumi/pulumi";
import * as civo from "@pulumi/civo";

const cluster = new civo.KubernetesCluster("civo-k3s-cluster", {
    name: "myFirstCivoCluster",
    numTargetNodes: 3,
    targetNodesSize: "g3.k3s.medium",
    region: "LON1"
})

export const clusterName = cluster.name
