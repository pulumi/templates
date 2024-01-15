package main

import (
	"github.com/pinecone-io/pulumi-pinecone/sdk/go/pinecone"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		myExampleIndex, err := pinecone.NewPineconeIndex(ctx, "my-example-index", &pinecone.PineconeIndexArgs{
			Name:   pulumi.String("example-index-go"),
			Metric: pinecone.IndexMetricCosine,
			Spec: &pinecone.PineconeSpecArgs{
				Serverless: &pinecone.PineconeServerlessSpecArgs{
					Cloud:  pinecone.ServerlessSpecCloudAws,
					Region: pulumi.String("us-west-2"),
				},
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("myPineconeIndexHost", myExampleIndex.Host)

		return nil
	})
}
