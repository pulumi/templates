const fs = require("fs");
const yaml = require("js-yaml");
const glob = require("glob");

// Load architecture metadata.
const architectures = yaml.load(fs.readFileSync("metadata/architectures.yaml", "utf-8"));

// Load template-group metadata.
const templateGroups = glob.sync("metadata/groups/*.yaml")
    .map(groupYamlFile => {
        return yaml.load(fs.readFileSync(groupYamlFile, "utf-8"));
    });

// Load all templates.
const allTemplates = glob.sync("*/Pulumi.yaml")

    // For each one, parse it into JSON and change the name and description.
    .map(pulumiYamlFile => {
        const template = yaml.load(fs.readFileSync(pulumiYamlFile, "utf-8"));
        const templateName = pulumiYamlFile.replace("/Pulumi.yaml", "");
        return {
            name: templateName,
            runtime: template.runtime,
            template: template.template,
        };
    })
    .reduce((templates, template) => {
        templates[template.name] = template;
        return templates;
    }, {});

// Build the list of architectures and groups.
const architectureTemplates = Object.keys(architectures)
    .map(arch => architectures[arch])
    .map(arch => {
        return {
            ...arch,
            ...{
                groups: templateGroups.filter(group => {
                    return group.kind === "architecture" && group.parent === arch.slug;
                }),
            },
        };
    });

// Assemble the list.
const result = {
    architectures: architectureTemplates,
    templates: allTemplates,
};

// Write the list to an output file.
fs.writeFileSync("metadata/dist/metadata.json", JSON.stringify(result, null, 4), "utf-8");
