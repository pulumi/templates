"use strict";
import rediscloud from "@rediscloud/pulumi-rediscloud";

const card = await rediscloud.getPaymentMethod(
	{
		cardType: "Visa",
		lastFourNumbers: "1234",
	},
);

const subscription = new rediscloud.Subscription(
	"my-subscription",
	{
		name: "my-subscription",
		paymentMethod: "credit-card",
		paymentMethodId: card.id,
		cloudProvider: {
			regions: [
				{
					region: "us-east-1",
					multipleAvailabilityZones: false,
					networkingDeploymentCidr: "10.0.0.0/24",
					preferredAvailabilityZones: ["use1-az1", "use1-az2", "use1-az5"],
				},
			],
		},

		creationPlan: {
			memoryLimitInGb: 10,
			quantity: 1,
			replication: true,
			supportOssClusterApi: false,
			throughputMeasurementBy: "operations-per-second",
			throughputMeasurementValue: 20000,
			modules: ["RedisJSON"],
		},
	},
);

const database = new rediscloud.SubscriptionDatabase("my-db", {
	name: "my-db",
	subscriptionId: subscription.id,
	protocol: "redis",
	memoryLimitInGb: 10,
	dataPersistence: "aof-every-1-second",
	throughputMeasurementBy: "operations-per-second",
	throughputMeasurementValue: 20000,
	replication: true,
});
