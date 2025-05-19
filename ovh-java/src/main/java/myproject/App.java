package myproject;

import com.ovhcloud.pulumi.ovh.CloudProject.CloudProjectFunctions;
import com.ovhcloud.pulumi.ovh.CloudProject.inputs.GetCapabilitiesContainerFilterArgs;
import com.ovhcloud.pulumi.ovh.CloudProject.inputs.GetContainerRegistryArgs;
import com.ovhcloud.pulumi.ovh.CloudProject.ContainerRegistry;
import com.ovhcloud.pulumi.ovh.CloudProject.ContainerRegistryArgs;
import com.ovhcloud.pulumi.ovh.CloudProject.ContainerRegistryUser;
import com.ovhcloud.pulumi.ovh.CloudProject.ContainerRegistryUserArgs;
import com.pulumi.Context;
import com.pulumi.Pulumi;

public class App {

    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {

        var config = ctx.config();
        var ovhServiceName = config.require("ovhServiceName");
        var ovhRegion = config.get("ovhRegion").orElse("GRA");
        var planName = config.get("planName").orElse("SMALL");

        var registryName = config.get("registryName").orElse("my-registry");
        var registryUserName = config.get("registryUserName").orElse("user");
        var registryUserEmail = config.get("registryUserEmail").orElse("myuser@ovh.com");
        var registryUserLogin = config.get("registryUserLogin").orElse("myuser");

        // Initiate the configuration of the registry
        final var regcap = CloudProjectFunctions.getCapabilitiesContainerFilter(GetCapabilitiesContainerFilterArgs.builder()
            .serviceName(ovhServiceName)
            .planName(planName)
            .region(ovhRegion)
            .build());

        // Deploy a new Managed private registry
        var myRegistry = new ContainerRegistry("myRegistry", ContainerRegistryArgs.builder()
            .serviceName(regcap.applyValue(getCapabilitiesContainerFilterResult -> getCapabilitiesContainerFilterResult.serviceName()))
            .planId(regcap.applyValue(getCapabilitiesContainerFilterResult -> getCapabilitiesContainerFilterResult.id()))
            .region(regcap.applyValue(getCapabilitiesContainerFilterResult -> getCapabilitiesContainerFilterResult.region()))
            .build());

        // Create a Private Registry User
        var myRegistryUser = new ContainerRegistryUser(registryUserName, ContainerRegistryUserArgs.builder()
            .serviceName(ovhServiceName)
            .registryId(myRegistry.id().asPlaintext())
            .email(registryUserEmail)
            .login(registryUserLogin)
            .build());

		// Add as an output registry information
		ctx.export("registryURL", myRegistry.url().asPlaintext());
        ctx.export("registryUser", myRegistryUser.user().asPlaintext());
		ctx.export("registryPassword", myRegistryUser.password().asPlaintext());
    }
}