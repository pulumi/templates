using System;
using System.Net;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Azure.WebJobs;
using Microsoft.Azure.WebJobs.Extensions.Http;
using Microsoft.AspNetCore.Http;

namespace App
{
    public static class App
    {
        [FunctionName("data")]
        public static async Task<IActionResult> Run([HttpTrigger(AuthorizationLevel.Anonymous, "get", Route = null)] HttpRequest req)
        {
            var res = req.HttpContext.Response;
            res.Headers.Append("Access-Control-Allow-Origin", "*");
            res.Headers.Append("Access-Control-Allow-Methods", "GET");

            if (HttpMethods.IsOptions(req.Method)) {
                return new NoContentResult();
            }

            res.Headers.Append("Content-Type", "application/json");
            var now = new { now = DateTimeOffset.Now.ToUnixTimeMilliseconds() };
            return new OkObjectResult(now);
        }
    }
}
