import * as pulumi from "@pulumi/pulumi";
import * as gcp from "@pulumi/gcp";

// Get some provider-namespaced configuration values
const providerCfg = new pulumi.Config("gcp");
const gcpProject = providerCfg.require("project");
const gcpRegion = providerCfg.get("region") || "us-central1";
// Get some other configuration values
const cfg = new pulumi.Config();
const nodePoolCount = cfg.getNumber("nodePoolCount") || 1;
// Create a new network
const gkeNetwork = new gcp.compute.Network("gke-network", {
    autoCreateSubnetworks: false,
    description: "A virtual network for your GKE cluster(s)",
    name: "gke-network",
    project: gcpProject,
});
// Create a new subnet in the network created above
const gkeSubnet = new gcp.compute.Subnetwork("gke-subnet", {
    ipCidrRange: "10.128.0.0/12",
    name: "gke-subnet",
    network: gkeNetwork.id,
    privateIpGoogleAccess: true,
    region: gcpRegion,
});
// Create a new GKE cluster
const gkeCluster = new gcp.container.Cluster("gke-cluster", {
    addonsConfig: {
        dnsCacheConfig: {
            enabled: true,
        },
    },
    binaryAuthorization: {
        evaluationMode: "PROJECT_SINGLETON_POLICY_ENFORCE",
    },
    datapathProvider: "ADVANCED_DATAPATH",
    description: "A GKE cluster",
    initialNodeCount: 1,
    ipAllocationPolicy: {
        clusterIpv4CidrBlock: "/14",
        servicesIpv4CidrBlock: "/20",
    },
    location: gcpRegion,
    masterAuthorizedNetworksConfig: {
        cidrBlocks: [{
            cidrBlock: "0.0.0.0/0",
            displayName: "All networks",
        }],
    },
    name: "gke-cluster",
    network: gkeNetwork.name,
    networkingMode: "VPC_NATIVE",
    privateClusterConfig: {
        enablePrivateNodes: true,
        enablePrivateEndpoint: false,
        masterIpv4CidrBlock: "10.100.0.0/28",
    },
    project: gcpProject,
    removeDefaultNodePool: true,
    releaseChannel: {
        channel: "STABLE",
    },
    subnetwork: gkeSubnet.name,
    workloadIdentityConfig: {
        workloadPool: `${gcpProject}.svc.id.goog`,
    },
});
// Build a Kubeconfig for accessing the cluster
const clusterKubeconfig = pulumi.interpolate `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ${gkeCluster.masterAuth.clusterCaCertificate}
    server: https://${gkeCluster.endpoint}
  name: ${gkeCluster.name}
contexts:
- context:
    cluster: ${gkeCluster.name}
    user: ${gkeCluster.name}
  name: ${gkeCluster.name}
current-context: ${gkeCluster.name}
kind: Config
preferences: {}
users:
- name: ${gkeCluster.name}
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: gke-gcloud-auth-plugin
      installHint: Install gke-gcloud-auth-plugin for use with kubectl by following
        https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke
      provideClusterInfo: true
`;
// Create a service account for the node pool
const gkeNodepoolSa = new gcp.serviceaccount.Account("gke-nodepool-sa", {
    accountId: "nodepool-1-sa",
    displayName: "Nodepool 1 Service Account",
});
// Create a nodepool for the GKE cluster
const gkeNodepool = new gcp.container.NodePool("gke-nodepool", {
    cluster: gkeCluster.id,
    nodeCount: nodePoolCount,
    location: gcpRegion,
    name: "gke-nodepool-1",
    nodeConfig: {
        oauthScopes: ["https://www.googleapis.com/auth/cloud-platform"],
        serviceAccount: gkeNodepoolSa.email,
    },
});
// Export some values for use elsewhere
export const networkName = gkeNetwork.name;
export const networkId = gkeNetwork.id;
export const clusterName = gkeCluster.name;
export const clusterId = gkeCluster.id;
export const kubeconfig = clusterKubeconfig;
