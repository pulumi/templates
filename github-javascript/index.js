"use strict";
const pulumi = require("@pulumi/pulumi");
const github = require("@pulumi/github")

const repo = new github.Repository("demo-repo", {
    description: "Demo Repository for GitHub",
});

exports.repositoryName = repo.name