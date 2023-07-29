package main

import (
	metal "github.com/pulumi/pulumi-equinix-metal/sdk/v3/go/equinix"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an Equinix Metal resource (Project)
		project, err := metal.NewProject(ctx, "my-project", &metal.ProjectArgs{
			Name: pulumi.String("TestProject1"),
		})
		if err != nil {
			return err
		}

		// Export the name of the project
		ctx.Export("projectName", project.Name)
		return nil
	})
}
