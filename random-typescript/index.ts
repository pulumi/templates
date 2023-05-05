import * as random from "@pulumi/random";

const username = new random.RandomPet("demo-repo", {});

export const name = repo.id