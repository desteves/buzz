package main

import (
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/cloudrun"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	// Create a new Pulumi project
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a Docker image from a Dockerfile and push it to Docker Hub.
		image, err := docker.NewImage(ctx, "my-image", &docker.ImageArgs{
			Build: &docker.DockerBuildArgs{
				Context: pulumi.String("./app"), // Path to the directory with Dockerfile and source.
			},
			ImageName: pulumi.Sprintf("<DOCKER_HUB_USERNAME>/myapp:v1.0.0"), // Replace with your Docker Hub username.
			SkipPush:  pulumi.Bool(false),                                   // Set to false to push to Docker Hub.
			Registry: &docker.RegistryArgs{
				Server:   pulumi.String("docker.io"),             // Docker Hub server.
				Username: pulumi.String("<DOCKER_HUB_USERNAME>"), // Docker Hub username.
				Password: pulumi.String("<DOCKER_HUB_PASSWORD>"), // Docker Hub password.
			},
		})
		if err != nil {
			return err
		}

		// Create a new Cloud Run Service using the image
		service, err := cloudrun.NewService(ctx, "service", &cloudrun.ServiceArgs{
			Location: pulumi.String("us-central1"),
			Template: &cloudrun.ServiceTemplateArgs{
				Spec: &cloudrun.ServiceTemplateSpecArgs{
					Containers: cloudrun.ServiceTemplateSpecContainerArray{
						&cloudrun.ServiceTemplateSpecContainerArgs{
							Image: image.ImageName, //pulumi.String("nullstring/buzz:dev"),
							Ports: cloudrun.ServiceTemplateSpecContainerPortArray{
								&cloudrun.ServiceTemplateSpecContainerPortArgs{
									ContainerPort: pulumi.Int(8000), // Replace with your desired port
								},
							},
						},
					},
				},
			},
		})
		// Check for errors
		if err != nil {
			return err
		}

		// Export the URL
		ctx.Export("url", service.Statuses.Index(pulumi.Int(0)).Url())
		return nil
	})
}
