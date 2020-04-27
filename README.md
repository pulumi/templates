![Build Status](https://github.com/pulumi/templates/workflows/Run%20Template%20Tests/badge.svg)

# Pulumi Templates       

This repo contains the templates for `pulumi new`, which make it easy to quickly get started building new Pulumi projects.

## Adding a new template

 1. Create a new directory for the template, e.g. `my-template-javascript`. By convention, hyphens are used to separate words and the language is included as a suffix.

 2. Add template files in the new directory.

## Text replacement

The following special strings can be included in any template file; these will be replaced by the CLI when laying down the template files.

 - `${PROJECT}` - The name of the project.
 - `${DESCRIPTION}` - The description of the project.
