"use strict";
const pulumi = require("@pulumi/pulumi");
const packet = require("@pulumi/packet");

// Create a Packet resource (Project)
const project = new packet.Project("my-test-project", {
    name: "my-project",
});

// Export the name of the project
exports.projectName = project.name;
