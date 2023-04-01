import * as pulumi from "@pulumi/pulumi";
import * as scaleway from "@lbrlabs/pulumi-scaleway";

const kapsule = new scaleway.KubernetesCluster("pulumi-kapsule", {
      name: "pulumi-kapsule",
      version: "1.23",
      region: "fr-par",
      cni: "cilium",
      tags: [
        "pulumi",
        "scaleway",
      ],
      autoUpgrade: {
        enable: true,
        maintenanceWindowStartHour: 3,
        maintenanceWindowDay: "monday"
      },
    }
)

new scaleway.KubernetesNodePool("pulumi-kapsule-pool", {
  zone: "fr-par-1",
  name: "pulumi-kapsule-pool",
  nodeType: "DEV1-L",
  size: 1,
  autoscaling: true,
  minSize: 1,
  maxSize: 3,
  autohealing: true,
  clusterId: kapsule.id,
})

export const kapsuleId = kapsule.id;
export const kubeconfig = pulumi.secret(kapsule.kubeconfigs[0].configFile);
