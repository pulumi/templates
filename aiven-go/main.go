package main

import (
	"github.com/pulumi/pulumi-aiven/sdk/v5/go/aiven"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		kafka, err := aiven.NewKafka(ctx, "aiven-kafka", &aiven.KafkaArgs{
			Project:     pulumi.String("<YOUR_AIVEN_PROJECT_NAME>"),
			CloudName:   pulumi.String("azure-westeurope"),
			Plan:        pulumi.String("startup-2"),
			ServiceName: pulumi.String("kafka-azure-eu"),
			KafkaUserConfig: &aiven.KafkaKafkaUserConfigArgs{
				KafkaVersion: pulumi.String("3.2"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("name", kafka.ServiceName)
		return nil
	})
}
