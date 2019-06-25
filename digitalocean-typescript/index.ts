import * as pulumi from "@pulumi/pulumi";
import * as digitalocean from "@pulumi/digitalocean";

// Create a DigitalOcean resource (Domain)
const domain = new digitalocean.Domain("my-domain", {
  name: "my-domain.io"
});

// Export the name of the domain
export const domainName = domain.name;
