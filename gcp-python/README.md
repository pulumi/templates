# Pulumi GCP Python Storage Bucket Template

 A minimal Pulumi template for provisioning a Google Cloud Storage bucket using Python.

 This template helps you get started with Pulumi and the `pulumi-gcp` provider to create a simple storage bucket and export its URL.

 ## When to Use

 - You need a quick example of using Pulumi with Google Cloud in Python.
 - You want to manage a Google Cloud Storage bucket as code.
 - You’re looking for a minimal scaffold to build more complex GCP infrastructure.

 ## Prerequisites

 - A Google Cloud account and a target GCP project.
 - Authentication set up via `gcloud auth login` or the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.
 - Python 3.7 or later installed on your machine.
 - Pulumi CLI installed and logged in to your Pulumi account.

 ## Usage

 Run the following command to scaffold a new project from this template:

 ```bash
 pulumi new gcp-python
 ```

 Follow the interactive prompts:
 - **Project Name**: Your project name.
 - **Project Description**: A short description of your project.
 - **gcp:project**: The ID of the Google Cloud project where resources will be created.

 After the project is generated, navigate into your project directory and deploy:

 ```bash
 cd <project-name>
 pulumi up
 ```

 Confirm the changes to provision your storage bucket.

 ## Project Layout

 ```
 .
 ├── __main__.py        # Pulumi program defining resources
 ├── Pulumi.yaml        # Project configuration and template metadata
 └── requirements.txt   # Python dependencies for Pulumi and GCP
 ```

 ## Configuration

 - **gcp:project**: (Required) The Google Cloud project ID where resources will be created.

 ## Resources Created

 - **Storage Bucket** (`pulumi_gcp.storage.Bucket`): A bucket named `my-bucket` in the `US` location.

 ## Outputs

 - **bucket_name**: The URL of the created storage bucket.

 ## Next Steps

 - Modify `__main__.py` to customize the bucket:
   - Change the bucket name.
   - Adjust the `location` or add bucket labels and IAM policies.
 - Add more GCP resources such as Pub/Sub topics, Compute instances, or BigQuery datasets.
 - Integrate with CI/CD pipelines using `pulumi preview` and `pulumi up --yes`.
 - Explore the [Pulumi GCP Provider Documentation](https://www.pulumi.com/registry/packages/gcp/) for more examples.

 ## Need Help?

 - Pulumi Docs: https://www.pulumi.com/docs/
 - GCP Provider Docs: https://www.pulumi.com/registry/packages/gcp/
 - Community Slack: https://slack.pulumi.com/
 - GitHub Issues: https://github.com/pulumi/pulumi/issues