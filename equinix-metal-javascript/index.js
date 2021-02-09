"use strict";
const pulumi = require("@pulumi/pulumi");
const metal = require("@pulumi/equinix-metal");

// Create an Equinix Metal resource (Project)
const project = new metal.Project("my-test-project", {
    name: "my-project",
});

// Export the name of the project
exports.projectName = project.name;
