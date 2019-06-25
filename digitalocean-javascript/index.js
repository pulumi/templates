"use strict";
const pulumi = require("@pulumi/pulumi");
const digitalocean = require("@pulumi/digitalocean");

// Create a DigitalOcean resource (Domain)
const domain = new digitalocean.Domain("my-domain", {
    name: "my-domain.io"
});

// Export the name of the domain
exports.domainName = domain.name;
