import pulumi
import pulumi_rediscloud

config = pulumi.Config()

card = pulumi_rediscloud.get_payment_method(
    card_type=config.require("cardType"),
    last_four_numbers=config.require("lastFourNumbers"),
)

subscription = pulumi_rediscloud.Subscription(
    "my-subscription",
    name="my-subscription",
    payment_method="credit-card",
    payment_method_id=card.id,
    cloud_provider={
        "regions": [
            {
                "region": "us-east-1",
                "multiple_availability_zones": True,
                "networking_deployment_cidr": "10.0.0.0/24",
                "preferred_availability_zones": ["use1-az1", "use1-az2", "use1-az5"],
            }
        ]
    },
    creation_plan={
        "memory_limit_in_gb": 10,
        "quantity": 1,
        "replication": True,
        "support_oss_cluster_api": False,
        "throughput_measurement_by": "operations-per-second",
        "throughput_measurement_value": 20000,
        "modules": ["RedisJSON"],
    },
)

database = pulumi_rediscloud.SubscriptionDatabase(
    "my-db",
    name="my-db",
    subscription_id=subscription.id,
    protocol="redis",
    memory_limit_in_gb=10,
    data_persistence="aof-every-1-second",
    throughput_measurement_by="operations-per-second",
    throughput_measurement_value=20000,
    replication=True,
    modules=[
        {
            "name": "RedisJSON",
        }
    ],
)
