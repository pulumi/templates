using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Ovh = Pulumi.Ovh;
using System;

return await Deployment.RunAsync(() => 
{
    // Get some configuration values (or use defaults)
    var config = new Pulumi.Config();
    var ovhServiceName = config.Require("ovhServiceName");
    var ovhRegion = config.Get("ovhRegion") ?? "GRA9";
    var clusterName = config.Get("clusterName") ?? "my_desired_cluster";
    var nodePoolName = config.Get("nodePoolName") ?? "my-desired-pool";
    var nodePoolDesiredNodes = config.GetInt32("nodePoolDesiredNodes") ?? 1;
    var nodePoolMaxNodes = config.GetInt32("nodePoolMaxNodes") ?? 3;
    var nodePoolMinNodes = config.GetInt32("nodePoolMinNodes") ?? 1;
    var flavorName = config.Get("flavorName") ?? "b3-8";

	// Deploy a new Kubernetes cluster
    var myKubeCluster = new Ovh.CloudProject.Kube(clusterName, new()
    {
        Region = ovhRegion,
        ServiceName = ovhServiceName,
        Name = clusterName,
    });

	//Create a Node Pool
    var nodePool = new Ovh.CloudProject.KubeNodePool(nodePoolName, new()
    {
        DesiredNodes = nodePoolDesiredNodes,
        FlavorName = flavorName,
        KubeId = myKubeCluster.Id,
        MaxNodes = nodePoolMaxNodes,
        MinNodes = nodePoolMinNodes,
        ServiceName = ovhServiceName,
    });

    // Export some values for use elsewhere
    return new Dictionary<string, object?>
    {
        ["kubeconfig"] = myKubeCluster.Kubeconfig,
        ["nodePoolID"] = nodePool.Id,
    };
});
