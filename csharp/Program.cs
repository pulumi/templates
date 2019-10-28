using System.Threading.Tasks;

using Pulumi;

class Program
{
    static Task<int> Main()
    {
        return Deployment.RunAsync(() => {

            // Add you resources here

        });
    }
}
