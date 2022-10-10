"use strict";
const auth0 = require("@pulumi/auth0");
const pulumi = require("@pulumi/pulumi");

const client = new auth0.Client("client", {
    allowedLogoutUrls: ["https://www.example.com/logout"],
    allowedOrigins: ["https://www.example.com"],
    callbacks: ["https://example.com/auth/callback"],
    appType: "regular_web",
    jwtConfiguration: {
        alg: "RS256"
    },

})

exports.clientId = client.clientId
exports.clientSecret = client.clientSecret
