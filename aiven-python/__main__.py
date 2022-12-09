import pulumi
import pulumi_aiven as aiven

kafka1 = aiven.Kafka("kafka1",
    project="<YOUR_AIVEN_PROJECT_NAME>",
    cloud_name="google-europe-west1",
    plan="business-4",
    service_name="kafka-gcp-eu",
    maintenance_window_dow="saturday",
    maintenance_window_time="10:00:00",
    kafka_user_config=aiven.KafkaKafkaUserConfigArgs(
        kafka_rest="true",
        kafka_connect="true",
        schema_registry="true",
        kafka_version="3.2",
        kafka=aiven.KafkaKafkaUserConfigKafkaArgs(
            group_max_session_timeout_ms="70000",
            log_retention_bytes="1000000000",
        ),
        public_access=aiven.KafkaKafkaUserConfigPublicAccessArgs(
            kafka_rest="true",
            kafka_connect="true",
        ),
    ))