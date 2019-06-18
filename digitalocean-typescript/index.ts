import * as pulumi from "@pulumi/pulumi";
import * as digitalocean from "@pulumi/digitalocean";

// Create a DigitalOcean resource (Domain)
const domain = new digitalocean.Domain("my-test-domain");

// Export the name of the domain
export const domainName = domain.name;
