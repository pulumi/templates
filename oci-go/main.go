package main

import (
	"github.com/pulumi/pulumi-oci/sdk/go/oci/identity"
	"github.com/pulumi/pulumi-oci/sdk/go/oci/objectstorage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		myCompartment, err := identity.NewCompartment(ctx, "my-compartment", &identity.CompartmentArgs{
			Name:         pulumi.String("my-compartment"),
			Description:  pulumi.String("My description text"),
			EnableDelete: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		myNamespace := pulumi.All(myCompartment.CompartmentId).ApplyT(
			func(args []interface{}) (string, error) {
				namespace, err := objectstorage.GetNamespace(ctx, &objectstorage.GetNamespaceArgs{
					CompartmentId: pulumi.StringRef(args[0].(string)),
				})
				if err != nil {
					return "", err
				}
				return namespace.Namespace, nil
			},
		).(pulumi.StringOutput)

		myBucket, err := objectstorage.NewBucket(ctx, "my-bucket", &objectstorage.BucketArgs{
			Name:          pulumi.String("my-bucket"),
			Namespace:     myNamespace,
			CompartmentId: myCompartment.ID(),
		})
		if err != nil {
			return err
		}

		ctx.Export("name", myBucket.Name)

		return nil
	})
}
