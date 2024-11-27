package myproject;

import com.ovhcloud.pulumi.ovh.CloudProject.Kube;
import com.ovhcloud.pulumi.ovh.CloudProject.KubeArgs;
import com.ovhcloud.pulumi.ovh.CloudProject.KubeNodePool;
import com.ovhcloud.pulumi.ovh.CloudProject.KubeNodePoolArgs;
import com.pulumi.Context;
import com.pulumi.Pulumi;

public class App {

    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {

        var config = ctx.config();
        var ovhServiceName = config.require("ovhServiceName");
        var ovhRegion = config.get("ovhRegion").orElse("GRA9");
        var clusterName = config.get("clusterName").orElse("my_desired_cluster");
        var nodePoolName = config.get("nodePoolName").orElse("my-desired-pool");

        var nodePoolDesiredNodes = config.getInteger("nodePoolDesiredNodes").orElse(1);
        var nodePoolMaxNodes = config.getInteger("nodePoolMaxNodes").orElse(3);
        var nodePoolMinNodes = config.getInteger("nodePoolMinNodes").orElse(1);
        var flavorName = config.get("flavorName").orElse("b3-8");

		// Deploy a new Kubernetes cluster
        var myKube = new Kube(clusterName, KubeArgs.builder()
            .region(ovhRegion)
            .serviceName(ovhServiceName)
            .name(clusterName)
            .build());

		// Export kubeconfig file to a secret
        ctx.export("kubeconfig", myKube.kubeconfig());

		//Create a Node Pool
        var nodePool = new KubeNodePool(nodePoolName, KubeNodePoolArgs.builder()
            .desiredNodes(nodePoolDesiredNodes)
            .flavorName(flavorName)
            .kubeId(myKube.id().asPlaintext())
            .maxNodes(nodePoolMaxNodes)
            .minNodes(nodePoolMinNodes)
            .serviceName(ovhServiceName)
            .build());
        
        ctx.export("nodePoolID", nodePool.id().asPlaintext());
    }
}