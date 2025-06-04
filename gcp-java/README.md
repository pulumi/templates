# Minimal Google Cloud Java Pulumi Template

A minimal Pulumi template for provisioning a Google Cloud Storage bucket using Java.

This template demonstrates how to use the Pulumi Java SDK and the GCP provider to create cloud resources in a typed, reusable manner.

---

## Overview

This template performs the following actions:

- Creates a Google Cloud Storage bucket in the `US` region.
- Exports the bucket URL as an output (`bucketName`).

It is a great starting point for building Java applications that manage GCP resources with Pulumi.

## Prerequisites

Before using this template, ensure you have:

- A Google Cloud Platform project with appropriate permissions.
- Credentials configured (for example, via `gcloud auth application-default login`).
- Java Development Kit (JDK) 11 or later.
- Apache Maven 3.6 or later (or use the included Maven Wrapper).
- Pulumi CLI installed and authenticated.

## Usage

1. Create a new Pulumi project from this template:

   ```bash
   pulumi new gcp-java
   ```

2. Navigate into your new project directory:

   ```bash
   cd <your-project-name>
   ```

3. Configure your GCP project:

   ```bash
   pulumi config set gcp:project YOUR_GCP_PROJECT_ID
   ```

4. (Optional) Build your Java program:

   ```bash
   mvn package
   # or, using the Maven Wrapper
   ./mvnw package
   ```

5. Deploy your stack:

   ```bash
   pulumi up
   ```

## Project Layout

```text
./
├── Pulumi.yaml       # Pulumi project and template metadata
├── pom.xml           # Maven build definition
└── src
    └── main/java
        └── myproject
            └── App.java   # Pulumi program entrypoint
```

## Configuration

| Name          | Description                                |
| ------------- | ------------------------------------------ |
| `gcp:project` | The Google Cloud project ID to deploy into |

Use `pulumi config` to view or change these values.

## Resources Created

- **Google Cloud Storage Bucket** (`my-bucket`)
  - Region: US

## Outputs

| Name         | Description                   |
| ------------ | ----------------------------- |
| `bucketName` | The URL of the created bucket |

## When to Use This Template

- You are kickstarting a Java-based Pulumi project targeting GCP.
- You want a minimal example showing how to provision and export a basic resource.
- You prefer working in Java with full IDE support and type safety.

## Next Steps

- Extend `App.java` to create additional GCP resources (Compute Engine, Cloud Functions, Pub/Sub, etc.).
- Organize your code into packages and modules for larger projects.
- Integrate with CI/CD pipelines to automate deployments.
- Explore Pulumi’s Java API documentation: https://www.pulumi.com/docs/reference/pkg/java/

## Getting Help

If you encounter any issues or have questions:

- Check the Pulumi documentation: https://www.pulumi.com/docs/
- Browse or file issues on GitHub: https://github.com/pulumi/pulumi/issues
- Join the Pulumi community Slack: https://slack.pulumi.com/