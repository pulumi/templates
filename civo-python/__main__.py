import pulumi
import pulumi_civo as civo

firewall = civo.Firewall("civo-firewall", create_default_rules=True, region="LON1")

cluster = civo.KubernetesCluster(
    "civo-k3s-cluster",
    name="myFirstCivoCluster",
    region="LON1",
    firewall_id=firewall.id,
    pools={
        "node_count": 3,
        "size": "g4s.kube.medium",
    },
)

pulumi.export("cluster_name", cluster.name)
