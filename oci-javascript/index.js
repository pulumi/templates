"use strict";
const pulumi = require("@pulumi/pulumi");
const oci = require("@pulumi/oci");

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

exports.name = myBucket.name
