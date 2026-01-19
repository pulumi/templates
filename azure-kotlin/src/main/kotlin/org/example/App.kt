package org.example

import com.pulumi.Pulumi
import com.pulumi.azurenative.resources.ResourceGroup
import com.pulumi.azurenative.storage.StorageAccount
import com.pulumi.azurenative.storage.StorageAccountArgs
import com.pulumi.azurenative.storage.inputs.SkuArgs

fun main() {
    Pulumi.run { ctx ->
        val rg = ResourceGroup("my-rg")

        val storage = StorageAccount(
            "mystorageacct",
            StorageAccountArgs.builder()
                .resourceGroupName(rg.name())
                .sku(SkuArgs.builder().name("Standard_LRS").build())
                .kind("StorageV2")
                .build()
        )

        ctx.export("storageAccountName", storage.name())
    }
}
