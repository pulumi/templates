package main

import (
	"strconv"

	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		appPath := "./app"
		if param := cfg.Get("appPath"); param != "" {
			appPath = param
		}
		imageName := "my-app"
		if param := cfg.Get("imageName"); param != "" {
			imageName = param
		}
		containerPort := 8080
		if param := cfg.GetInt("containerPort"); param != 0 {
			containerPort = param
		}
		cpu := 1
		if param := cfg.GetInt("cpu"); param != 0 {
			cpu = param
		}
		memory := "1Gi"
		if param := cfg.Get("memory"); param != "" {
			memory = param
		}
		concurrency := 50
		if param := cfg.GetInt("concurrency"); param != 0 {
			concurrency = param
		}

		// Import the provider's configuration settings.
		providerConfig := config.New(ctx, "gcp")
		location := providerConfig.Require("region")
		project := providerConfig.Require("project")

		// Generate a unique Artifact Registry repository ID
		uniqueString, err := random.NewRandomString(ctx, "unique-string", &random.RandomStringArgs{
			Length:  pulumi.Int(4),
			Lower:   pulumi.Bool(true),
			Upper:   pulumi.Bool(false),
			Numeric: pulumi.Bool(true),
			Special: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		repoId := pulumi.Sprintf("repo-%s", uniqueString.Result)

		// Create an Artifact Registry repository
		repository, err := artifactregistry.NewRepository(ctx, "repository", &artifactregistry.RepositoryArgs{
			Description:  pulumi.String("Repository for container image"),
			Format:       pulumi.String("DOCKER"),
			Location:     pulumi.String(location),
			RepositoryId: repoId,
		})
		if err != nil {
			return err
		}

		// Form the repository URL
		repoUrl := pulumi.Sprintf("%s-docker.pkg.dev/%s/%s", repository.Location, project, repository.RepositoryId)

		// Create a container image for the service.
		// Before running `pulumi up`, configure Docker for authentication to Artifact Registry as
		// described here: https://cloud.google.com/artifact-registry/docs/docker/authentication
		image, err := dockerbuild.NewImage(ctx, "image", &dockerbuild.ImageArgs{
			Tags: pulumi.StringArray{pulumi.Sprintf("%s/%s", repoUrl, imageName)},
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String(appPath),
			},
			// Cloud Run currently requires x86_64 images
			// https://cloud.google.com/run/docs/container-contract#languages
			Platforms: dockerbuild.PlatformArray{"linux/amd64"},
		})
		if err != nil {
			return err
		}

		// Create a Cloud Run service definition.
		service, err := cloudrun.NewService(ctx, "service", &cloudrun.ServiceArgs{
			Location: pulumi.String(location),
			Template: cloudrun.ServiceTemplateArgs{
				Spec: cloudrun.ServiceTemplateSpecArgs{
					Containers: cloudrun.ServiceTemplateSpecContainerArray{
						cloudrun.ServiceTemplateSpecContainerArgs{
							Image: image.Ref,
							Resources: cloudrun.ServiceTemplateSpecContainerResourcesArgs{
								Limits: pulumi.ToStringMap(map[string]string{
									"memory": memory,
									"cpu":    strconv.Itoa(cpu),
								}),
							},
							Ports: cloudrun.ServiceTemplateSpecContainerPortArray{
								cloudrun.ServiceTemplateSpecContainerPortArgs{
									ContainerPort: pulumi.Int(containerPort),
								},
							},
						},
					},
					ContainerConcurrency: pulumi.Int(concurrency),
				},
			},
		})
		if err != nil {
			return err
		}

		// Create an IAM member to make the service publicly accessible.
		_, err = cloudrun.NewIamMember(ctx, "invoker", &cloudrun.IamMemberArgs{
			Service:  service.Name,
			Location: pulumi.String(location),
			Role:     pulumi.String("roles/run.invoker"),
			Member:   pulumi.String("allUsers"),
		})
		if err != nil {
			return err
		}

		// Export the URL of the service.
		ctx.Export("url", service.Statuses.Index(pulumi.Int(0)).Url())

		return nil
	})
}
