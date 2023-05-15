package main

import (
	"strconv"

	"github.com/RedisLabs/pulumi-rediscloud/sdk/go/rediscloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		card, err := rediscloud.GetPaymentMethod(ctx, &rediscloud.GetPaymentMethodArgs{
			CardType:        pulumi.StringRef("Visa"),
			LastFourNumbers: pulumi.StringRef("1234"),
		}, nil)

		if err != nil {
			return err
		}

		cardId := card.Id

		subscription, err := rediscloud.NewSubscription(ctx, "subscription", &rediscloud.SubscriptionArgs{
			PaymentMethod:   pulumi.String("credit-card"),
			PaymentMethodId: pulumi.String(cardId),
			CloudProvider: &rediscloud.SubscriptionCloudProviderArgs{
				Regions: rediscloud.SubscriptionCloudProviderRegionArray{
					&rediscloud.SubscriptionCloudProviderRegionArgs{
						Region:                    pulumi.String("us-east-1"),
						MultipleAvailabilityZones: pulumi.Bool(false),
						NetworkingDeploymentCidr:  pulumi.String("10.0.0.0/24"),
						PreferredAvailabilityZones: pulumi.StringArray{
							pulumi.String("use1-az1"),
							pulumi.String("use1-az2"),
							pulumi.String("use1-az3"),
						},
					},
				},
			},
			CreationPlan: &rediscloud.SubscriptionCreationPlanArgs{
				MemoryLimitInGb:            pulumi.Float64(10),
				Quantity:                   pulumi.Int(1),
				Replication:                pulumi.Bool(true),
				SupportOssClusterApi:       pulumi.Bool(false),
				ThroughputMeasurementBy:    pulumi.String("operations-per-second"),
				ThroughputMeasurementValue: pulumi.Int(20000),
				Modules:                    pulumi.StringArray{pulumi.String("RedisJSON")},
			},
		})
		if err != nil {
			return err
		}
		_, err = rediscloud.NewSubscriptionDatabase(ctx, "database", &rediscloud.SubscriptionDatabaseArgs{
			SubscriptionId:             subscription.ID().ApplyT(strconv.Atoi).(pulumi.IntOutput),
			Protocol:                   pulumi.String("redis"),
			MemoryLimitInGb:            pulumi.Float64(10),
			DataPersistence:            pulumi.String("aof-every-1-second"),
			ThroughputMeasurementBy:    pulumi.String("operations-per-second"),
			ThroughputMeasurementValue: pulumi.Int(20000),
			SupportOssClusterApi:       pulumi.Bool(false),
			Replication:                pulumi.Bool(true),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
