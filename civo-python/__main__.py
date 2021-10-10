import pulumi
import pulumi_civo

cluster = pulumi_civo.KubernetesCluster('civo-k3s-cluster', name='myFirstCivoCluster', region='LON1',
                                        num_target_nodes=3, target_nodes_size='g3.k3s.medium')

pulumi.export('cluster_name', cluster.name)
