![Build Status](https://github.com/pulumi/templates/actions/workflows/ci.yml/badge.svg?branch=master)

# Pulumi Templates

This repo contains the templates for `pulumi new`, which make it easy to quickly get started building new Pulumi projects.

## Adding a new template

1. Create a new directory for the template, e.g. `my-template-javascript`. By convention, hyphens are used to separate words and the language is included as a suffix.

1. Add template files in the new directory. Note that when new projects are created from templates, all of the files contained in the template directory are copied into the resulting Pulumi project. Be sure to exclude any unnecessary files.

   Also note that dependency lockfiles like `package-lock.json` and `go.sum` are deliberately git-ignored to ensure that new projects always track with the latest Pulumi and provider SDKs.

1. If the template is an architecture template, include the requisite supplemental metadata at `./metadata`:

   - If the template adds to an existing template group (for example, if it's a new a `static-website-aws` template), add a new line for the template in the `templates` section of that group:

     ```diff
       name: AWS Static Website
       ...
       templates:
         ...
         - static-website-go
         - static-website-csharp
     -   - static-website-yaml
     +   - static-website-yaml
     +   - static-website-java
     ```

   - If the template introduces a new architecture, make a new entry in `./metadata/architectures.yaml` using the others as a guide (the keys and `slug` values correspond with these templates' eventual paths at <https://pulumi.com/templates>), then add a new file that lists the new template at `./metadata/groups/{architecture}-{cloud}-{language}.yaml`. Set the new group's `parent` property to match the key/name of the item you added to `architectures.yaml`.

1. Ensure the template applies sensible, conservative defaults for all configuration values. Ideally, users should be able to run `pulumi new --yes` with your template and get an immediately deployable project out of the box.

1. Ensure the template supports the _minimum_ runtime version for its associated language. Consult the [Languages & SDKs documentation](https://www.pulumi.com/docs/iac/languages-sdks/) for reference. (This is why our CI workflows use older runtimes. Every template in this repository should comply with this requirement.)

1. Request a review from the @pulumi/content-engineering team.

## Text replacement

The following special strings can be included in any template file; these will be replaced by the CLI when laying down the template files.

- `${PROJECT}` - The name of the project.
- `${DESCRIPTION}` - The description of the project.
