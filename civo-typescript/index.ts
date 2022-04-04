import * as pulumi from "@pulumi/pulumi";
import * as civo from "@pulumi/civo";

const firewall = new civo.Firewall("civo-firewall", {
    name: "myFirstFirewall",
    region: "LON1",
    createDefaultRules: true,
});

const cluster = new civo.KubernetesCluster("civo-k3s-cluster", {
    name: "myFirstCivoCluster",
    numTargetNodes: 3,
    targetNodesSize: "g3.k3s.medium",
    region: "LON1",
    firewallId: firewall.id,
})

export const clusterName = cluster.name
