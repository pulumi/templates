import * as pulumi from "@pulumi/pulumi";
import * as oci from "@pulumi/oci";

const myCompartment = new oci.identity.Compartment("myCompartment", {
  name: "my-compartment",
  description: "My description text",
  enableDelete: true
});

const myNamespace = pulumi.all([myCompartment.id]).apply(([id]) => {
  return oci.objectstorage.getNamespace({
    compartmentId: id,
  });
})

const myBucket = new oci.objectstorage.Bucket("myBucket", {
  compartmentId: myCompartment.id,
  namespace: myNamespace.namespace,
  name: "my-bucket"
});

export const name = myBucket.name
