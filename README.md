[![Build Status](https://travis-ci.com/pulumi/templates.svg?token=as6W2KPEwXJYiS5Jt2wi&branch=master)](https://travis-ci.com/pulumi/templates)

# Pulumi Templates

This repo contains the templates for `pulumi new [template]`, which make it easy to quickly get started building new Pulumi projects.

## Using a template

You can use the CLI to create a new project from a template:

```
$ mkdir myproj
$ cd myproj
$ pulumi new [template]
```

If `[template]` isn't specified, the CLI will offer a list of available templates to choose from:

```
$ pulumi new
> javascript
  python
  typescript
```

## Adding a new template

 1. Create a new directory under `templates/`, e.g. `templates/my-template-javascript`. By convention, hyphens are used to separate words and the language is included as a suffix.

 2. Add template files in the new directory.

Travis publishes a tarball of each template directory under `templates/` in `master` to S3.

Once the template has been published, it will be included in the JSON list of templates at `https://api.pulumi.com/releases/templates` and available to download from `https://api.pulumi.com/releases/templates/<template-name>.tar.gz`

## Text replacement

The following special strings can be included in any template file; these will be replaced by the CLI when laying down the template files.

 - `${PROJECT}` - The name of the project.
 - `${DESCRIPTION}` - The description of the project.
