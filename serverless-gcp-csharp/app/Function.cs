using Google.Cloud.Functions.Framework;
using Microsoft.AspNetCore.Http;
using System.Threading.Tasks;
using System.Text.Json;
using System.Net;
using System;

namespace App
{
    public class Data : IHttpFunction
    {
        public async Task HandleAsync(HttpContext context)
        {
            HttpRequest request = context.Request;
            HttpResponse response = context.Response;

            response.Headers.Append("Access-Control-Allow-Origin", "*");
            response.Headers.Append("Access-Control-Allow-Methods", "GET");

            if (HttpMethods.IsOptions(request.Method)) {
                response.StatusCode = (int) HttpStatusCode.NoContent;
                return;
            }

            response.Headers.Append("Content-Type", "application/json");
            var now = new { now = DateTimeOffset.Now.ToUnixTimeMilliseconds() };
            await response.WriteAsync(JsonSerializer.Serialize(now));
        }
    }
}
