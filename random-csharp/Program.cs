using Pulumi;
using Pulumi.Random;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    var username = new RandomPet("username", new RandomPetArgs{});

    return new Dictionary<string, object?>
    {
        ["name"] = username.Id
    };
});
