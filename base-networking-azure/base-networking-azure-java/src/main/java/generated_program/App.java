package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.azurenative.resources.ResourceGroup;
import com.pulumi.azurenative.resources.ResourceGroupArgs;
import com.pulumi.azurenative.network.VirtualNetwork;
import com.pulumi.azurenative.network.VirtualNetworkArgs;
import com.pulumi.azurenative.network.inputs.AddressSpaceArgs;
import com.pulumi.azurenative.network.Subnet;
import com.pulumi.azurenative.network.SubnetArgs;
import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.io.File;
import java.nio.file.Files;
import java.nio.file.Paths;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        final var config = ctx.config();
        final var azureNativeLocation = config.get("azureNativeLocation");
        final var mgmtGroupId = config.get("mgmtGroupId");
        var resourceGroup = new ResourceGroup("resourceGroup", ResourceGroupArgs.builder()        
            .location(azureNativeLocation)
            .resourceGroupName("rg")
            .build());

        var virtualNetwork = new VirtualNetwork("virtualNetwork", VirtualNetworkArgs.builder()        
            .addressSpace(AddressSpaceArgs.builder()
                .addressPrefixes("10.0.0.0/16")
                .build())
            .location(azureNativeLocation)
            .resourceGroupName(resourceGroup.name())
            .virtualNetworkName("vnet")
            .build());

        var subnet1 = new Subnet("subnet1", SubnetArgs.builder()        
            .addressPrefix("10.0.0.0/22")
            .name("subnet-1")
            .resourceGroupName(resourceGroup.name())
            .virtualNetworkName(virtualNetwork.name())
            .build());

        var subnet2 = new Subnet("subnet2", SubnetArgs.builder()        
            .addressPrefix("10.0.4.0/22")
            .name("subnet-2")
            .resourceGroupName(resourceGroup.name())
            .virtualNetworkName(virtualNetwork.name())
            .build());

        var subnet3 = new Subnet("subnet3", SubnetArgs.builder()        
            .addressPrefix("10.0.8.0/22")
            .name("subnet-3")
            .resourceGroupName(resourceGroup.name())
            .virtualNetworkName(virtualNetwork.name())
            .build());

    }
}
