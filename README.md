# Pulumi Templates

This repo contains the templates for `pulumi new`, which make it easy to quickly get started building new Pulumi projects.

## Adding a new template

> [!IMPORTANT]
> We are currently re-evaluating what content belongs in github.com/pulumi/templates, how
> it should be organized, and how it should be maintained. During this evaluation process,
> we will not be accepting new third-party templates.

1. Create a new directory for the template, e.g. `my-template-typescript`. By convention, hyphens are used to separate words and the language is included as a suffix.

1. Add template files in the new directory. Note that when new projects are created from templates, all of the files contained in the template directory are copied into the resulting Pulumi project. Be sure to exclude any unnecessary files.

   Node templates also need a committed lockfile. Run `make lockfiles` and commit the result. See [Lockfiles](#lockfiles) below.

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

## Lockfiles

Every Node template commits a lockfile alongside its `package.json`: `package-lock.json` for a `nodejs` runtime, `bun.lock` for a `bun` runtime. A template's dependency ranges are open (`^7.0.0`), so without a lockfile every `pulumi new` resolves fresh from the registry and installs whatever was published most recently, including anything published by an attacker who compromised a package in the tree. The lockfile pins the whole transitive tree, so everyone who creates a project from the template gets the same versions, and changing those versions becomes a reviewable diff.

- **After editing a template's `package.json`**, run `make lockfiles` and commit the regenerated lockfile. It resolves each lockfile from its `package.json` and nothing else, so it gives the same answer on your machine as it does in CI, and is a no-op when nothing changed.
- **Keeping dependencies fresh** is Renovate's job, not a manual one. `renovate.json5` enables `lockFileMaintenance`, which refreshes the pinned versions weekly within their existing ranges. The shared Pulumi config makes Renovate wait three days before adopting a release, and skips versions that OSV reports as malicious.
- **A pull request cannot introduce a vulnerability.** `make audit_lockfiles` compares the advisories in your lockfiles against the base branch, so a vulnerability somebody else introduced never blocks you.
- **npm and bun are both first-class.** Each is audited by the tool that wrote its lockfile, since `npm audit` cannot read `bun.lock` and `bun audit` cannot read `package-lock.json`. Renovate keeps a separate `bun` manager, so `renovate.json5` enables both.
- **The lockfile has to match the runtime the template declares.** The Pulumi CLI picks a project's package manager by looking for these files, so a stray `yarn.lock` or `pnpm-lock.yaml` would hand the template to a package manager it never declared. A `nodejs` template gets npm's lockfile; a `bun` template always installs with bun, so it gets bun's.
- **Lockfiles must stay out of `.gitignore`.** Publishing a template archives its directory with git-ignore rules applied, so an ignored lockfile stays in the repo but never reaches the users who need it. `make test_lockfiles` guards against this.

Lockfiles for the other languages, including `go.sum`, are git-ignored, so those templates resolve their dependencies when the project is created.

## Text replacement

The following special strings can be included in any template file; these will be replaced by the CLI when laying down the template files.

- `${PROJECT}` - The name of the project.
- `${DESCRIPTION}` - The description of the project.
