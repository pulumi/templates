"use strict";
const pulumi = require("@pulumi/pulumi");
const civo = require("@pulumi/civo");

const firewall = new civo.Firewall("civo-firewall", {
    name: "myFirstFirewall",
    region: "LON1",
    createDefaultRules: true
})

const cluster = new civo.KubernetesCluster("civo-k3s-cluster", {
    name: "myFirstCivoCluster",
    pools: {
        nodeCount: 3,
        size: "g4s.kube.medium"
    },
    region: "LON1",
    firewallId: firewall.id,
})

exports.clusterName = cluster.name

