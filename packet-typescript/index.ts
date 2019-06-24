import * as pulumi from "@pulumi/pulumi";
import * as packet from "@pulumi/packet";

// Create a Packet resource (Project)
const domain = new packet.Project("my-test-project", {
    name: "My Test Project"
});

// Export the name of the project
export const projectName = packet.name;
