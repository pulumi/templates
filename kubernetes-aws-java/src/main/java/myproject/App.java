package myproject;

import com.pulumi.Pulumi;
import com.pulumi.awsx.ec2.Vpc;
import com.pulumi.awsx.ec2.VpcArgs;
import com.pulumi.eks.Cluster;
import com.pulumi.eks.ClusterArgs;
import com.pulumi.eks.enums.AuthenticationMode;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            // Grab some values from the Pulumi configuration (or use default values)
            var config = ctx.config();
            var minClusterSize = config.getInteger("minClusterSize").orElse(3);
            var maxClusterSize = config.getInteger("maxClusterSize").orElse(6);
            var desiredClusterSize = config.getInteger("desiredClusterSize").orElse(3);
            var eksNodeInstanceType = config.get("eksNodeInstanceType").orElse("t3.medium");
            var vpcNetworkCidr = config.get("vpcNetworkCidr").orElse("10.0.0.0/16");

            // Create a VPC for the EKS cluster
            var eksVpc = new Vpc("eks-vpc", VpcArgs.builder()
                    .enableDnsHostnames(true)
                    .cidrBlock(vpcNetworkCidr)
                    .build());

            // Create the EKS cluster
            var eksCluster = new Cluster("eks-cluster", ClusterArgs.builder()
                    // Put the cluster in the new VPC created earlier
                    .vpcId(eksVpc.vpcId())
                    // Use the "API" authentication mode to support access entries
                    .authenticationMode(AuthenticationMode.Api)
                    // Public subnets will be used for load balancers
                    .publicSubnetIds(eksVpc.publicSubnetIds())
                    // Private subnets will be used for cluster nodes
                    .privateSubnetIds(eksVpc.privateSubnetIds())
                    // Change configuration values to change any of the following settings
                    .instanceType(eksNodeInstanceType)
                    .desiredCapacity(desiredClusterSize)
                    .minSize(minClusterSize)
                    .maxSize(maxClusterSize)
                    // Do not give the worker nodes public IP addresses
                    .nodeAssociatePublicIpAddress(false)
                    // Change these values for a private cluster (VPN access required)
                    .endpointPrivateAccess(false)
                    .endpointPublicAccess(true)
                    .build());

            // Export some values for use elsewhere
            ctx.export("kubeconfig", eksCluster.kubeconfig());
            ctx.export("vpcId", eksVpc.vpcId());
        });
    }
}
