import pulumi
import lbrlabs_scaleway as scaleway

kapsule = scaleway.KubernetesCluster("pulumi-kapsule",
                                     name="pulumi-kapsule",
                                     version="1.23",
                                     region="fr-par",
                                     cni="cilium",
                                     tags=["pulumi", "kapsule"],
                                     auto_upgrade=scaleway.KubernetesClusterAutoUpgradeArgs(
                                         enable=True,
                                         maintenance_window_start_hour=3,
                                         maintenance_window_day="monday")),

scaleway.KubernetesNodePool("pulumi-kapsule-nodepool",
                            zone="fr-par-1",
                            name="pulumi-kapsule-nodepool",
                            node_type="DEV1-L",
                            size=1,
                            autoscaling=True,
                            min_size=1,
                            max_size=3,
                            autohealing=True,
                            cluster_id=kapsule[0].id)

pulumi.export('kapsule_id', kapsule[0].id)
pulumi.export('kubeconfig', pulumi.Output.secret(kapsule[0].kubeconfigs[0].config_file))
