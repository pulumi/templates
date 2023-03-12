const metadata = require("../../metadata/dist/metadata.json");
const fs = require("fs");

describe("Template metadata", () => {

    describe("architectures list", () => {
        const architectures = metadata.architectures;

        it("is a non-empty array", () => {
            expect(Array.isArray(architectures)).toBe(true);
            expect(architectures.length).toBeGreaterThan(0);
        });

        it("contains valid items", () => {
            const item = architectures[0];
            expect(item.name).toBeDefined();
            expect(item.slug).toBeDefined();
            expect(Array.isArray(item.groups)).toBe(true);
            expect(item.groups.length).toBeGreaterThan(1);
            expect(item.groups[0].name).toBeDefined();
            expect(item.groups[0].kind).toBeDefined();
            expect(item.groups[0].parent).toBeDefined();
            expect(item.groups[0].slug).toBeDefined();
            expect(item.groups[0].clouds.length).toBeGreaterThan(0);
            expect(item.groups[0].templates.length).toBeGreaterThan(0);
        });

        it("contains the architectures we expect", () => {
            expect(architectures.find(arch => arch.name === "Static Website")).toBeDefined();
            expect(architectures.find(arch => arch.name === "Serverless")).toBeDefined();
            expect(architectures.find(arch => arch.name === "Container Service")).toBeDefined();
            expect(architectures.find(arch => arch.name === "Kubernetes Cluster")).toBeDefined();
            expect(architectures.find(arch => arch.name === "Kubernetes Application")).toBeDefined();
            expect(architectures.find(arch => arch.name === "Virtual Machine")).toBeDefined();
        });

        it("only contains templates that exist in the repository", () => {
            architectures.forEach(arch => {
                arch.groups.forEach(group => {
                    group.templates.forEach(template => {
                        expect(() => fs.readdirSync(`${template}`)).not.toThrow();
                    });
                });
            });
        });
    });

    describe("templates list", () => {
        const templates = metadata.templates;

        it("is a map keyed by template name", () => {
            const keys = Object.keys(templates);
            expect(keys.length).toBeGreaterThan(0);
            expect(templates["aws-typescript"]).toBeDefined();
        });

        it("contains valid items", () => {
            const [ name, template ] = Object.entries(templates)[1];
            expect(template.name).toBeDefined();
            expect(template.runtime).toBeDefined();
            expect(template.template.description).toBeDefined();
            expect(template.template.config).toBeDefined();
        });

        it("contains the templates we expect", () => {
            expect(templates["aws-typescript"]).toBeDefined();
            expect(templates["aws-python"]).toBeDefined();
            expect(templates["aws-go"]).toBeDefined();
            expect(templates["aws-csharp"]).toBeDefined();
            expect(templates["aws-yaml"]).toBeDefined();
            expect(templates["aws-java"]).toBeDefined();

            expect(templates["azure-typescript"]).toBeDefined();
            expect(templates["azure-python"]).toBeDefined();
            expect(templates["azure-go"]).toBeDefined();
            expect(templates["azure-csharp"]).toBeDefined();
            expect(templates["azure-yaml"]).toBeDefined();
            expect(templates["azure-java"]).toBeDefined();

            expect(templates["gcp-typescript"]).toBeDefined();
            expect(templates["gcp-python"]).toBeDefined();
            expect(templates["gcp-go"]).toBeDefined();
            expect(templates["gcp-csharp"]).toBeDefined();
            expect(templates["gcp-yaml"]).toBeDefined();
            expect(templates["azure-java"]).toBeDefined();
        });

        it("only contains templates that exist in the repository", () => {
            Object.keys(templates).forEach(template => {
                expect(() => fs.readdirSync(`${template}`)).not.toThrow();
            });
        });
    });
});
