import * as pulumi from "@pulumi/pulumi";
import * as linode from "@pulumi/linode";

// Create a Linode resource (Linode Instance)
const instance = new linode.Instance("my-instance", {
    type: "g6-nanode-1",
    region: "us-east",
    image: "linode/ubuntu18.04",
});

// Export the Instance label of the instance
export const instanceLabel = instance.label;
