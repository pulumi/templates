package main

import (
	"github.com/pulumi/pulumi-packet/sdk/go/packet"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a Packet resource (Project)
		project, err := packet.NewProject(ctx, "my-project", &packet.ProjectArgs{
			Name: "TestProject1",
		})
		if err != nil {
			return err
		}

		// Export the name of the project
		ctx.Export("projectName", project.Name())
		return nil
	})
}
