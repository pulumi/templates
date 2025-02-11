package myproject;

import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.oci.Identity.Compartment;
import com.pulumi.oci.Identity.CompartmentArgs;
import com.pulumi.oci.ObjectStorage.Bucket;
import com.pulumi.oci.ObjectStorage.BucketArgs;
import com.pulumi.oci.ObjectStorage.ObjectStorageFunctions;
import com.pulumi.oci.ObjectStorage.inputs.GetNamespacePlainArgs;
import com.pulumi.oci.ObjectStorage.outputs.GetNamespaceResult;

import java.util.concurrent.CompletableFuture;


public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {

            Compartment myCompartment = new Compartment("my-compartment",
                    CompartmentArgs.builder()
                            .name("my-compartment")
                            .enableDelete(true)
                            .description("My description text").build()
            );

            Output<String> namespace = Output.all(myCompartment.id()).apply(values -> {
                try {
                    CompletableFuture<GetNamespaceResult> result = ObjectStorageFunctions.getNamespacePlain(GetNamespacePlainArgs.builder()
                            .compartmentId(values.get(0))
                            .build());
                    return Output.of(result.get().namespace());
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
            });

            Bucket myBucket = new Bucket("my-bucket",
                    BucketArgs.builder()
                            .name("my-bucket")
                            .compartmentId(myCompartment.id())
                            .namespace(namespace)
                            .build());

            ctx.export("name", myBucket.name());
        });
    }
}
