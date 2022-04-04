import pulumi
import pulumi_civo as civo

firewall = civo.Firewall("civo-firewall", create_default_rules=True, region='LON1')

cluster = civo.KubernetesCluster('civo-k3s-cluster', name='myFirstCivoCluster', region='LON1',
                                 num_target_nodes=3, target_nodes_size='g3.k3s.medium', firewall_id=firewall.id)

pulumi.export('cluster_name', cluster.name)
