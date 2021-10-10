"use strict";
const pulumi = require("@pulumi/pulumi");
const civo = require("@pulumi/civo");

const cluster = new civo.KubernetesCluster("civo-k3s-cluster", {
    name: "myFirstCivoCluster",
    numTargetNodes: 3,
    targetNodesSize: "g3.k3s.medium",
    region: "LON1"
})

exports.clusterName = cluster.name

