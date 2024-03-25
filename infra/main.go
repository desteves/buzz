package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/cloudrun"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	// Create a new Pulumi project
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a new Cloud Run Service
		service, err := cloudrun.NewService(ctx, "service", &cloudrun.ServiceArgs{
			Location: pulumi.String("us-central1"),
			Template: &cloudrun.ServiceTemplateArgs{
				Spec: &cloudrun.ServiceTemplateSpecArgs{
					Containers: cloudrun.ServiceTemplateSpecContainerArray{
						&cloudrun.ServiceTemplateSpecContainerArgs{
							Image: pulumi.String("nullstring/buzz:dev"),
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
