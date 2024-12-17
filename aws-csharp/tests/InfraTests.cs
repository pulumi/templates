using System.Collections.Generic;
using System.Collections.Immutable;
using System.Linq;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Aws.S3;
using Pulumi.Testing;
using Xunit;

public static class TestingExtensions
{
    // Test helper to convert a Pulumi Output to a Task.
    // This should only be used in tests.
    public static Task<T> GetValueAsync<T>(this Output<T> output)
        => Pulumi.Utilities.OutputUtilities.GetValueAsync(output);
}

class Mocks : IMocks
{
    // Mock calls to create new resources and return a canned response.
    public Task<(string? id, object state)> NewResourceAsync(MockResourceArgs args)
    {
        // Here, we're returning a same-shaped object for all resource types.
        // We could, however, use the arguments passed into this function to
        // customize the mocked-out properties of a particular resource.
        // See the unit-testing docs for details:
        // https://www.pulumi.com/docs/iac/concepts/testing/unit/
        return Task.FromResult<(string?, object)>(($"{args.Name}_id", (object)args.Inputs));
    }

    // Mock function calls and return an empty response.
    public Task<object> CallAsync(MockCallArgs args)
    {
        return Task.FromResult((object)ImmutableDictionary<string, object>.Empty);
    }
}

public class InfraTests
{
    class TestStack : Stack
    {
        public TestStack()
        {
            Outputs = Deploy.Infra();
        }

        public Dictionary<string, object?> Outputs { get; set; }
    }

    private static Task<ImmutableArray<Resource>> TestAsync()
        => Deployment.TestAsync<TestStack>(new Mocks(), new TestOptions { IsPreview = false });


    // Example test. To run, uncomment and run `dotnet test` from the tests directory.
    // [Fact]
    // public async Task TestBucketTags()
    // {
    //     var resources = await TestAsync();
    //     var bucket = resources.OfType<BucketV2>().Single();

    //     var tags = await bucket.Tags.GetValueAsync();
    //     Assert.NotNull(tags);
    //     Assert.Contains("Name", tags);
    // }
}
