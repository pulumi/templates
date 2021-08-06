package main

import (
	"github.com/pulumi/pulumi-aiven/sdk/v4/go/aiven"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		kafka, err := aiven.NewKafka(ctx, "aiven-kafka", &aiven.KafkaArgs{
			Project:     pulumi.String("<project-name>"),
			CloudName:   pulumi.String("azure-westeurope"),
			Plan:        pulumi.String("startup-2"),
			ServiceName: pulumi.String("myAivenKafkaService"),
			KafkaUserConfig: &aiven.KafkaKafkaUserConfigArgs{
				KafkaVersion: pulumi.String("2.7"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("name", kafka.ServiceName)
		return nil
	})
}
