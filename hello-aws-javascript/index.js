// Import the [cloud-aws](https://pulumi.io/packages/pulumi-cloud/) package
const cloud = require("@pulumi/cloud-aws");

// Create a public HTTP endpoint (using AWS APIGateway)
const endpoint = new cloud.API("hello");

// Serve static files from the `www` folder (using AWS S3)
endpoint.static("/", "www");

// Serve a simple REST API on `GET /name` (using AWS Lambda)
endpoint.get("/source", (req, res) => res.json({name: "AWS"}));

// Export the public URL for the HTTP service
exports.url = endpoint.publish().url;
