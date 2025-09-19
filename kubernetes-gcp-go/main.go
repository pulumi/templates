package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get some provider-namespaced configuration values
		providerCfg := config.New(ctx, "gcp")
		gcpProject := providerCfg.Require("project")
		gcpRegion, err := providerCfg.Try("region")
		if err != nil {
			gcpRegion = "us-central1"
		}
		// Get some additional configuration values or use defaults
		cfg := config.New(ctx, "")
		nodesPerZone, err := cfg.TryInt("nodesPerZone")
		if err != nil {
			nodesPerZone = 1
		}

		// Create a new network
		gkeNetwork, err := compute.NewNetwork(ctx, "gke-network", &compute.NetworkArgs{
			AutoCreateSubnetworks: pulumi.Bool(false),
			Description:           pulumi.String("A virtual network for your GKE cluster(s)"),
		})
		if err != nil {
			return err
		}

		// Create a subnet in the network
		gkeSubnet, err := compute.NewSubnetwork(ctx, "gke-subnet", &compute.SubnetworkArgs{
			IpCidrRange:           pulumi.String("10.128.0.0/12"),
			Network:               gkeNetwork.ID(),
			PrivateIpGoogleAccess: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Create a new GKE cluster
		gkeCluster, err := container.NewCluster(ctx, "gke-cluster", &container.ClusterArgs{
			AddonsConfig: &container.ClusterAddonsConfigArgs{
				DnsCacheConfig: &container.ClusterAddonsConfigDnsCacheConfigArgs{
					Enabled: pulumi.Bool(true),
				},
			},
			BinaryAuthorization: &container.ClusterBinaryAuthorizationArgs{
				EvaluationMode: pulumi.String("PROJECT_SINGLETON_POLICY_ENFORCE"),
			},
			DatapathProvider: pulumi.String("ADVANCED_DATAPATH"),
			Description:      pulumi.String("A GKE cluster"),
			InitialNodeCount: pulumi.Int(1),
			IpAllocationPolicy: &container.ClusterIpAllocationPolicyArgs{
				ClusterIpv4CidrBlock:  pulumi.String("/14"),
				ServicesIpv4CidrBlock: pulumi.String("/20"),
			},
			Location: pulumi.String(gcpRegion),
			MasterAuthorizedNetworksConfig: &container.ClusterMasterAuthorizedNetworksConfigArgs{
				CidrBlocks: container.ClusterMasterAuthorizedNetworksConfigCidrBlockArray{
					&container.ClusterMasterAuthorizedNetworksConfigCidrBlockArgs{
						CidrBlock:   pulumi.String("0.0.0.0/0"),
						DisplayName: pulumi.String("All networks"),
					},
				},
			},
			Network:        gkeNetwork.Name,
			NetworkingMode: pulumi.String("VPC_NATIVE"),
			PrivateClusterConfig: &container.ClusterPrivateClusterConfigArgs{
				EnablePrivateNodes:    pulumi.Bool(true),
				EnablePrivateEndpoint: pulumi.Bool(false),
				MasterIpv4CidrBlock:   pulumi.String("10.100.0.0/28"),
			},
			RemoveDefaultNodePool: pulumi.Bool(true),
			ReleaseChannel: &container.ClusterReleaseChannelArgs{
				Channel: pulumi.String("STABLE"),
			},
			Subnetwork: gkeSubnet.Name,
			WorkloadIdentityConfig: &container.ClusterWorkloadIdentityConfigArgs{
				WorkloadPool: pulumi.String(fmt.Sprintf("%v.svc.id.goog", gcpProject)),
			},
		})
		if err != nil {
			return err
		}

		// Create a GCP Service Account for the node pool
		gkeNodepoolSa, err := serviceaccount.NewAccount(ctx, "gke-nodepool-sa", &serviceaccount.AccountArgs{
			AccountId:   pulumi.Sprintf("%v-np-1-sa", gkeCluster.Name),
			DisplayName: pulumi.String("Nodepool 1 Service Account"),
		})
		if err != nil {
			return err
		}

		// Create a new node pool
		_, err = container.NewNodePool(ctx, "gke-nodepool", &container.NodePoolArgs{
			Cluster:   gkeCluster.ID(),
			NodeCount: pulumi.Int(nodesPerZone),
			NodeConfig: &container.NodePoolNodeConfigArgs{
				OauthScopes: pulumi.StringArray{
					pulumi.String("https://www.googleapis.com/auth/cloud-platform"),
				},
				ServiceAccount: gkeNodepoolSa.Email,
			},
		})
		if err != nil {
			return err
		}

		// Build Kubeconfig for accessing the cluster
		clusterKubeconfig := pulumi.Sprintf(`apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: %[3]s
    server: https://%[2]s
  name: %[1]s
contexts:
- context:
    cluster: %[1]s
    user: %[1]s
  name: %[1]s
current-context: %[1]s
kind: Config
preferences: {}
users:
- name: %[1]s
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: gke-gcloud-auth-plugin
      installHint: Install gke-gcloud-auth-plugin for use with kubectl by following
        https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke
      provideClusterInfo: true
        `, gkeCluster.Name, gkeCluster.Endpoint, gkeCluster.MasterAuth.ClusterCaCertificate().Elem())

		// Export some values for use elsewhere
		ctx.Export("networkName", gkeNetwork.Name)
		ctx.Export("networkId", gkeNetwork.ID())
		ctx.Export("clusterName", gkeCluster.Name)
		ctx.Export("clusterId", gkeCluster.ID())
		ctx.Export("kubeconfig", clusterKubeconfig)
		return nil
	})
}
