package main

import (
	"github.com/pulumi/pulumi-packet/sdk/v2/go/packet"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a Packet resource (Project)
		project, err := packet.NewProject(ctx, "my-project", &packet.ProjectArgs{
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
