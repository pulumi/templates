 # Pulumi GCP Go: Minimal Storage Bucket Template

 This template provisions a Google Cloud Storage bucket using Pulumi and Go. It demonstrates how to:
   - Use the Pulumi GCP provider in a Go program
   - Create a simple GCP resource (a Storage Bucket)
   - Export resource outputs for use in your stacks

 It’s a great starting point for learning Pulumi with Go on GCP or bootstrapping a project that needs object storage.

 ## Providers

 - Google Cloud Platform via the Pulumi GCP SDK for Go (`github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp`)

 ## Resources

 - **Storage Bucket** (`gcp.storage.Bucket`)
   - Logical name: `my-bucket`
   - Location: `US`

 ## Outputs

 - **bucketName**: The URL of the newly created bucket (e.g., `https://storage.googleapis.com/my-bucket`)

 ## When to Use This Template

 - You want a minimal Pulumi program in Go targeting GCP
 - You need a simple object storage bucket for assets or data
 - You’re exploring Pulumi’s Go SDK for cloud provisioning

 ## Prerequisites

 - Go 1.20 or later installed
 - A Google Cloud account with billing enabled
 - GCP credentials configured for Pulumi (for example, via `gcloud auth application-default login`)

 ## Usage

 1. Scaffold a new project from this template:
    ```bash
    pulumi new gcp-go
    ```
 2. When prompted, fill in:
    - **Project name**: your desired project identifier
    - **Description**: a short description of your stack
    - **gcp:project**: your target GCP project ID
 3. Change into your project directory:
    ```bash
    cd <your-project-name>
    ```
 4. Preview and deploy your stack:
    ```bash
    pulumi preview
    pulumi up
    ```

 ## Project Layout

 ```
 ├── Pulumi.yaml   Pulumi project definition and template settings
 ├── go.mod        Go module declaration and dependencies
 └── main.go       Pulumi program defining the Storage Bucket
 ```

 ## Configuration

 The following Pulumi configuration values are available:

 | Name           | Description                             | Default    |
 | -------------- | --------------------------------------- | ---------- |
 | `gcp:project`  | The Google Cloud project to deploy into | _required_ |

 Set configuration with:
 ```bash
 pulumi config set gcp:project YOUR_PROJECT_ID
 ```

 ## Next Steps

 - Add more GCP resources (e.g., Compute Engine, Pub/Sub, Cloud Functions)
 - Parameterize bucket settings such as versioning, access control, and lifecycle rules
 - Integrate IAM bindings for fine-grained permission management
 - Connect this bucket to other services or CI/CD pipelines

 ## Getting Help

 - Pulumi Documentation: https://www.pulumi.com/docs/
 - GCP Provider Reference: https://www.pulumi.com/registry/packages/gcp/
 - Community Slack: https://slack.pulumi.com/
 - GitHub Issues: https://github.com/pulumi/pulumi/issues