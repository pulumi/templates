![Build Status](https://github.com/pulumi/templates/actions/workflows/ci.yml/badge.svg?branch=master)

# Pulumi Templates

This repo contains the templates for `pulumi new`, which make it easy to quickly get started building new Pulumi projects.

## Adding a new template

1. Create a new directory for the template, e.g. `my-template-javascript`. By convention, hyphens are used to separate words and the language is included as a suffix.

1. Add template files in the new directory.

2. If the template is an architecture template, include the requisite supplemental metadata at `./metadata`:

    * If the template adds to an existing template group (for example, if it's a new a `static-website-aws` template), add a new line for the template in the `templates` section of that group:

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

    * If the template introduces a new architecture, make a new entry in `./metadata/architectures.yaml` using the others as a guide (the keys and `slug` values correspond with these templates' eventual paths at <https://pulumi.com/templates>), then add a new file that lists the new template at `./metadata/groups/{architecture}-{cloud}-{language}.yaml`. Set the new group's `parent` property to match the key/name of the item you added to `architectures.yaml`.

3. Request a review from the @pulumi/content-engineering team.

## Text replacement

The following special strings can be included in any template file; these will be replaced by the CLI when laying down the template files.

 - `${PROJECT}` - The name of the project.
 - `${DESCRIPTION}` - The description of the project.
