import * as pulumi from "@pulumi/pulumi";
import * as metal from "@pulumi/equinix-metal";

// Create an Equinix Metal resource (Project)
const project = new metal.Project("my-test-project", {
  name: "My Test Project",
});

// Export the name of the project
export const projectName = project.name;
