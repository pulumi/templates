import * as pulumi from "@pulumi/pulumi";
import * as auth0 from "@pulumi/auth0";

const client = new auth0.Client("client", {
    allowedLogoutUrls: ["https://www.example.com/logout"],
    allowedOrigins: ["https://www.example.com"],
    callbacks: ["https://example.com/auth/callback"],
    appType: "regular_web",
    jwtConfiguration: {
        alg: "RS256"
    },

})

export const clientId = client.clientId
export const clientSecret = client.clientSecret
