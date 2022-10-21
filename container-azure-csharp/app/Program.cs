var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/", async (context) =>
{
    await context.Response.WriteAsJsonAsync(new
    {
        message = "Hello, world!"
    });
});

app.Run();
