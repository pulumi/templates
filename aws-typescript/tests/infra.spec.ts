import * as pulumi from "@pulumi/pulumi";
import "jest";

// Test helper to convert a Pulumi Output to a Promise.
// This should only be used in tests.
function promiseOf<T>(output: pulumi.Output<T>): Promise<T> {
    return new Promise(resolve => output.apply(resolve));
}

describe("infrastructure", () => {
    // Define the infra variable as a type whose shape matches the that of the
    // to-be-defined infra module.
    let infra: typeof import("../infra");

    beforeAll(() => {
        // Put Pulumi in unit-test mode, mocking all calls to cloud-provider APIs.
        pulumi.runtime.setMocks({
            // Mock calls to create new resources and return a canned response.
            newResource: (args: pulumi.runtime.MockResourceArgs) => {
                // Here, we're returning a same-shaped object for all resource types.
                // We could, however, use the arguments passed into this function to
                // customize the mocked-out properties of a particular resource.
                // See the unit-testing docs for details:
                // https://www.pulumi.com/docs/iac/concepts/testing/unit/
                return {
                    id: `${args.name}-id`,
                    state: args.inputs,
                };
            },

            // Mock function calls and return an empty response.
            call: (args: pulumi.runtime.MockCallArgs) => {
                return {};
            },
        });
    });

    beforeEach(async () => {
        // Dynamically import the infra module.
        infra = await import("../infra");
    });

    // Example test. To run, uncomment and run `npm test`.
    describe("bucket", () => {
        it("must have a name tag", async () => {
            const tags = await promiseOf(infra.bucket.tags);
            expect(tags).toBeDefined();
            expect(tags).toHaveProperty("Name");
        });
        it("must have the right tags", async () => {
            await (expect(infra.bucket.tags) as any).toEqualOutputOf({ "Name": "My bucket oops" })
        })
        it("some apply thing", async () => {
            await (expect(infra.bucket.tags) as any).apply((tags: any) => {
                expect(tags).toBeDefined();
                expect(tags).toHaveProperty("Name2");
            })
        })
    });
});
