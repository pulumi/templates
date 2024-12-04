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
    var ovhRegion = config.Get("ovhRegion") ?? "GRA";
    var planName = config.Get("planName") ?? "SMALL";
    var registryName = config.Get("registryName") ?? "my-registry";

    var registryUserName = config.Get("registryUserName") ?? "user";
    var registryUserEmail = config.Get("registryUserEmail") ?? "myuser@ovh.com";
    var registryUserLogin = config.Get("registryUserLogin") ?? "myuser";

    // Initiate the configuration of the registry
    var regcap = Ovh.CloudProject.GetCapabilitiesContainerFilter.Invoke(new()
    {
        ServiceName = ovhServiceName,
        PlanName = planName,
        Region = ovhRegion,
    });

    // Deploy a new Managed private registry
    var myRegistry = new Ovh.CloudProject.ContainerRegistry(registryName, new()
    {
        ServiceName = regcap.Apply(getCapabilitiesContainerFilterResult => getCapabilitiesContainerFilterResult.ServiceName),
        PlanId = regcap.Apply(getCapabilitiesContainerFilterResult => getCapabilitiesContainerFilterResult.Id),
        Region = regcap.Apply(getCapabilitiesContainerFilterResult => getCapabilitiesContainerFilterResult.Region),
    });

    // Create a Private Registry User
    var myRegistryUser = new Ovh.CloudProject.ContainerRegistryUser(registryUserName, new()
    {
        ServiceName = ovhServiceName,
        RegistryId = myRegistry.Id,
        Email = registryUserEmail,
        Login = registryUserLogin,
    });

    //  Export some values for use elsewhere
    return new Dictionary<string, object?>
    {
        ["registryURL"] = myRegistry.Url,
        ["registryUser"] = myRegistryUser.User,
        ["registryPassword"] = myRegistryUser.Password,
    };
});
