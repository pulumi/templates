"""A Python Pulumi program"""
"""
This Pulumi program deploys a Hugging Face Language Model (LLM)
on an Amazon SageMaker instance. The deployment includes various configurations
such as instance type, environment variables, and health check timeouts.

Modules:
    pulumi: Pulumi Infrastructure as Code library
    HuggingFaceLlm: Component for Hugging Face LLM deployment

Attributes:
    llm (HuggingFaceLlm): An instance of the HuggingFaceLlm component
    representing the deployed LLM model.

Exports:
    EndpointName (pulumi.Output): The AWS name of the deployed SageMaker endpoint.
"""

# Import required modules
import pulumi
from huggingface_llm import HuggingFaceLlm

# Initialize the HuggingFaceLlm component with the required configurations
# Note: 'Llama2Llm' is a custom name given to this particular LLM instance.
llm = HuggingFaceLlm(
    'Llama2Llm',  # Custom name for the LLM model
    instance_type='ml.g5.2xlarge',  # AWS instance type for SageMaker deployment
    environment_variables={
        'HF_MODEL_ID': 'NousResearch/Llama-2-7b-chat-hf',  # HuggingFace model ID
        'SM_NUM_GPUS': '1',  # Number of GPUs to use
        'MAX_INPUT_LENGTH': '2048',  # Maximum input length for the model
        'MAX_TOTAL_TOKENS': '4096',  # Maximum number of tokens
        'MAX_BATCH_TOTAL_TOKENS': '8192',  # Maximum tokens in a batch
    },
    tgi_version='0.9.3',  # TGI version for the backend
    pytorch_version='2.0.1',  # PyTorch version to use
    startup_health_check_timeout_in_seconds=600  # Health check timeout in seconds
)

# Export the endpoint name for external access
# This will be available as an output after successful pulumi up
pulumi.export('EndpointName', llm.endpoint.name)
