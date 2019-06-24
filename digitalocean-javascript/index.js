"use strict";
const pulumi = require("@pulumi/pulumi");
const digitalocean = require("@pulumi/digitalocean");

// Create a DigitalOcean resource (Domain)
const domain = new digitalocean.Domain("my-test-domain");

// Export the name of the domain
exports.domainName = domain.name;
