import asyncio
from typing import TypeVar

import pulumi


# Test helper to convert a Pulumi Output to a Future.
# This should only be used in tests.
T = TypeVar("T")
def future_of(output: pulumi.Output[T]) -> asyncio.Future[T]:
    loop = asyncio.get_running_loop()
    future = loop.create_future()
    output.apply(lambda x: future.set_result(x))
    return future


class MyMocks(pulumi.runtime.Mocks):
    # Mock calls to create new resources and return a canned response.
    def new_resource(self, args: pulumi.runtime.MockResourceArgs):
        # Here, we're returning a same-shaped object for all resource types.
        # We could, however, use the arguments passed into this function to
        # customize the mocked-out properties of a particular resource.
        # See the unit-testing docs for details:
        # https://www.pulumi.com/docs/iac/concepts/testing/unit/
        return [args.name + "_id", args.inputs]

    # Mock function calls and return an empty response.
    def call(self, args: pulumi.runtime.MockCallArgs):
        return {}

# Put Pulumi in unit-test mode, mocking all calls to cloud-provider APIs.
pulumi.runtime.set_mocks(MyMocks())

# Now import the code that creates resources, and then test it.
import infra

# Example test. To run, uncomment and run `python -m pytest --disable-pytest-warnings`.
# @pulumi.runtime.test
# async def test_bucket_tags():
#     tags = await future_of(infra.bucket.tags)
#     assert tags, "bucket must have tags"
#     assert "Name" in tags, "bucket must have a Name tag"
