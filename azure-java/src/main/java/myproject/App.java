package myproject;

import com.pulumi.Pulumi;
import com.pulumi.azurenative.resources.ResourceGroup;
import com.pulumi.azurenative.storage.StorageAccount;
import com.pulumi.azurenative.storage.StorageAccountArgs;
import com.pulumi.azurenative.storage.StorageFunctions;
import com.pulumi.azurenative.storage.enums.Kind;
import com.pulumi.azurenative.storage.enums.SkuName;
import com.pulumi.azurenative.storage.inputs.ListStorageAccountKeysArgs;
import com.pulumi.azurenative.storage.inputs.SkuArgs;
import com.pulumi.core.Either;
import com.pulumi.core.Output;
import com.pulumi.deployment.InvokeOptions;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var resourceGroup = new ResourceGroup("resourceGroup");
            var storageAccount = new StorageAccount("sa", StorageAccountArgs.builder()
                    .resourceGroupName(resourceGroup.name())
                    .sku(SkuArgs.builder()
                            .name(SkuName.Standard_LRS)
                            .build())
                    .kind(Kind.StorageV2)
                    .build());

            var primaryStorageKey = getStorageAccountPrimaryKey(
                    resourceGroup.name(),
                    storageAccount.name());

            ctx.export("primaryStorageKey", primaryStorageKey);
        });
    }

    private static Output<String> getStorageAccountPrimaryKey(Output<String> resourceGroupName,
                                                              Output<String> accountName) {
        return Output.tuple(resourceGroupName, accountName).apply(tuple -> {
            var actualResourceGroupName = tuple.t1;
            var actualAccountName = tuple.t2;
            var invokeResult = StorageFunctions.listStorageAccountKeys(ListStorageAccountKeysArgs.builder()
                    .resourceGroupName(actualResourceGroupName)
                    .accountName(actualAccountName)
                    .build(), InvokeOptions.Empty);
            return invokeResult
                    .applyValue(r -> r.keys().get(0).value())
                    .asSecret();
        });
    }
}
