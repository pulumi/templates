import * as github from "@pulumi/github";

const repo = new github.Repository("demo-repo", {
    description: "Demo Repository for GitHub",
});

export const clusterName = repo.name