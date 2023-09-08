"""
A Pulumi program to deploy a Hugging Face Language Model (LLM) on Amazon SageMaker.

This Python module defines a Pulumi component resource for deploying a Hugging Face
LLM model on Amazon SageMaker. It also sets up CloudWatch monitoring for the deployed model.
"""

import json
import pulumi
from pulumi import Output
from pulumi_aws import config, iam, sagemaker, cloudwatch
from sagemaker import huggingface
from typing import Mapping, Optional


class HuggingFaceLlm(pulumi.ComponentResource):
    """
    A Pulumi component for deploying a Hugging Face LLM model on Amazon SageMaker.

    Attributes:
        endpoint: The deployed SageMaker endpoint for the model.

    Methods:
        setup_cloudwatch_alarms: Sets up CloudWatch alarms for monitoring the deployed model.
    """
    def __init__(
            self,
            name: str,
            instance_type: str,
            environment_variables: Mapping[str, str],
            tgi_version: str = '0.9.3',
            pytorch_version: str = '2.0.1',
            startup_health_check_timeout_in_seconds: int = 600,
            opts: Optional[pulumi.ResourceOptions] = None
    ):
        """
        Initialize a new Hugging Face LLM deployment.

        Parameters:
            name: Name of the deployment.
            instance_type: AWS instance type for the deployment.
            environment_variables: Environment variables for the model.
            tgi_version: Version for the backend technology (default is '0.9.3').
            pytorch_version: PyTorch version for the model (default is '2.0.1').
            startup_health_check_timeout_in_seconds: Health check timeout (default is 600 seconds).
            opts: Additional options for the Pulumi resource.
        """
        super().__init__('huggingface:llm:HuggingFaceLlm', name, None, opts)

        # Merge environment variables with optional version specifications
        extended_env_vars = {
            **environment_variables,
            'TGI_VERSION': tgi_version,
            'PYTORCH_VERSION': pytorch_version
        }

        # Fetch the container image URI
        container_image = huggingface.get_huggingface_llm_image_uri(
            backend='huggingface',
            region=config.region,
            version=tgi_version
        )

        # Create an IAM role for SageMaker
        role = iam.Role(f'{name}-role',
            assume_role_policy=json.dumps({
                'Version': '2012-10-17',
                'Statement': [{
                    'Effect': 'Allow',
                    'Principal': {'Service': 'sagemaker.amazonaws.com'},
                    'Action': 'sts:AssumeRole',
                }],
            }),
            managed_policy_arns=['arn:aws:iam::aws:policy/AmazonSageMakerFullAccess'],
            opts=pulumi.ResourceOptions(parent=self)
        )

        # Create a SageMaker model
        sage_maker_model = sagemaker.Model(f'{name}-model',
            execution_role_arn=role.arn,
            primary_container=sagemaker.ModelContainerArgs(
                image=container_image,
                environment=extended_env_vars
            ),
            opts=pulumi.ResourceOptions(parent=self)
        )

        # Configure the SageMaker endpoint
        cfn_endpoint_config = sagemaker.EndpointConfiguration(f'{name}-config',
            production_variants=[
                sagemaker.EndpointConfigurationProductionVariantArgs(
                    model_name=sage_maker_model.name,
                    variant_name='primary',
                    initial_variant_weight=1.0,
                    initial_instance_count=1,
                    instance_type=instance_type,
                    container_startup_health_check_timeout_in_seconds=startup_health_check_timeout_in_seconds
                ),
            ],
            opts=pulumi.ResourceOptions(parent=self)
        )

        # Create the SageMaker endpoint
        self.endpoint = sagemaker.Endpoint(f'{name}-endpoint',
            endpoint_config_name=cfn_endpoint_config.name,
            opts=pulumi.ResourceOptions(parent=self)
        )

        # Set up CloudWatch alarms for monitoring
        self.setup_cloudwatch_alarms(self.endpoint.name)

    def setup_cloudwatch_alarms(self, endpoint_name: Output[str]):
        """
        Set up CloudWatch alarms for monitoring the SageMaker endpoint.

        Creates two alarms:
        - One for high latency
        - One for high error rates

        Parameters:
            endpoint_name: The name of the SageMaker endpoint to monitor.
        """

        def create_high_latency_alarm(name: str):
            """Create a CloudWatch alarm for high latency."""
            cloudwatch.MetricAlarm(f"{name}-HighLatency",
                metric_name="ModelLatency",
                namespace="AWS/SageMaker",
                statistic="Average",
                period=60,
                evaluation_periods=1,
                threshold=100000,  # in microseconds
                comparison_operator="GreaterThanOrEqualToThreshold",
                alarm_description=f"High latency alarm for {name}",
                opts=pulumi.ResourceOptions(parent=self)
            )

        def create_high_errors_alarm(name: str):
            """Create a CloudWatch alarm for high error rates."""
            cloudwatch.MetricAlarm(f"{name}-HighErrors",
                metric_name="ModelError",
                namespace="AWS/SageMaker",
                statistic="SampleCount",
                period=60,
                evaluation_periods=1,
                threshold=5,
                comparison_operator="GreaterThanOrEqualToThreshold",
                alarm_description=f"High error rate alarm for {name}",
                opts=pulumi.ResourceOptions(parent=self)
            )

        # Apply the alarms to the endpoint
        endpoint_name.apply(create_high_latency_alarm)
        endpoint_name.apply(create_high_errors_alarm)
