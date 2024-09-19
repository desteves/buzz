package main

import (
	"fmt"
	"os"

	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/cloudrun"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	// Create a new Pulumi project
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a Docker image from a Dockerfile and push it to Docker Hub.
		username := os.Getenv("DOCKER_USR")
		var APPNAME = "buzz"
		var TAG = fmt.Sprintf("docker.io/%s/%s:latest", username, "buzz")
		if ctx.Stack() == "dev" {
			fmt.Println("Building and pushing Docker image for dev stack...")
			// Build and push an image to ECR with inline caching.
			_, err := dockerbuild.NewImage(ctx, APPNAME, &dockerbuild.ImageArgs{
				// Tag our image with our ECR repository's address.
				Tags: pulumi.StringArray{
					pulumi.String(TAG),
				},
				Context: &dockerbuild.BuildContextArgs{
					Location: pulumi.String("../app"),
				},
				Dockerfile: &dockerbuild.DockerfileArgs{
					Location: pulumi.String("../app/Dockerfile"),
				},
				// Build the image for AMD.
				Platforms: dockerbuild.PlatformArray{
					dockerbuild.Platform_Linux_amd64,
				},
				// Push the final result to the registries
				Push: pulumi.Bool(true),
				// Provide Registry credentials.
				Registries: dockerbuild.RegistryArray{
					&dockerbuild.RegistryArgs{
						Address:  pulumi.String("docker.io"), // Docker Hub server.
						Password: pulumi.String(os.Getenv("DOCKER_PAT")),
						Username: pulumi.String(username),
					},
				},
			})
			if err != nil {
				return err
			}
		}

		if ctx.Stack() == "dev" {
			fmt.Println("Skipping cloud deployment for dev stack...")
			return nil
		} else {
			fmt.Println("Deploying to Google Cloud Run...")
		}

		googleProject := os.Getenv("GOOGLE_PROJECT")
		region := os.Getenv("GOOGLE_REGION")
		deterministicURL := "https://" + APPNAME + "-" + googleProject + "." + region + ".run.app"
		ctx.Export("URL", pulumi.String(deterministicURL))

		// Create a new Cloud Run Service using the image
		service, err := cloudrun.NewService(ctx, APPNAME, &cloudrun.ServiceArgs{
			// https://www.pulumi.com/docs/concepts/resources/names/
			Name:     pulumi.String(APPNAME), // overrides autonaming
			Location: pulumi.String(region),
			Template: &cloudrun.ServiceTemplateArgs{
				Spec: &cloudrun.ServiceTemplateSpecArgs{
					Containers: cloudrun.ServiceTemplateSpecContainerArray{
						&cloudrun.ServiceTemplateSpecContainerArgs{
							Image: pulumi.String(TAG),
							Ports: cloudrun.ServiceTemplateSpecContainerPortArray{
								&cloudrun.ServiceTemplateSpecContainerPortArgs{
									ContainerPort: pulumi.Int(8000),
								},
							},
							Envs: cloudrun.ServiceTemplateSpecContainerEnvArray{
								&cloudrun.ServiceTemplateSpecContainerEnvArgs{
									Name:  pulumi.String("GOOGLE_OAUTH_CLIENT_ID"),
									Value: pulumi.String(os.Getenv("GOOGLE_OAUTH_CLIENT_ID")),
								},
								&cloudrun.ServiceTemplateSpecContainerEnvArgs{
									Name:  pulumi.String("GOOGLE_OAUTH_CLIENT_SECRET"),
									Value: pulumi.String(os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")),
								},
								&cloudrun.ServiceTemplateSpecContainerEnvArgs{
									Name:  pulumi.String("GEMINI_API_KEY"),
									Value: pulumi.String(os.Getenv("GEMINI_API_KEY")),
								},
								&cloudrun.ServiceTemplateSpecContainerEnvArgs{
									Name:  pulumi.String("REDIR"),
									Value: pulumi.Sprintf("%s/%s", deterministicURL, "app"),
								},
								&cloudrun.ServiceTemplateSpecContainerEnvArgs{
									Name:  pulumi.String("SERVER_ADDR"),
									Value: pulumi.String(os.Getenv("SERVER_ADDR")),
								},
							},
						},
					},
				},
			},
		},
			pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "5m"}),
			pulumi.DeleteBeforeReplace(true))
		if err != nil {
			return err
		}

		// Create an IAM member to make the service publicly accessible.
		_, err = cloudrun.NewIamMember(ctx, "invoker", &cloudrun.IamMemberArgs{
			Service:  service.Name,
			Location: service.Location,
			Role:     pulumi.String("roles/run.invoker"),
			Member:   pulumi.String("allUsers"),
		})
		if err != nil {
			return err
		}
		// TODO add domain mapping resource
		// https://www.pulumi.com/docs/reference/pkg/gcp/cloudrun/domainmapping/

		// Create a new Cloudflare CDN
		zoneID := os.Getenv("CLOUDFLARE_ZONE")
		// Set the DNS record for the CDN.
		_, err = cloudflare.NewRecord(ctx, APPNAME, &cloudflare.RecordArgs{
			ZoneId:  pulumi.String(zoneID),                 // Replace with your actual Zone ID
			Name:    pulumi.String(APPNAME),                // The subdomain or record name
			Type:    pulumi.String("CNAME"),                // Typically a CNAME for CDN usage
			Content: pulumi.String("ghs.googlehosted.com"), // The value of the record, like a CDN endpoint
			Proxied: pulumi.Bool(false),                    // Set to true to proxy traffic through Cloudflare
		})
		if err != nil {
			return err
		}
		return nil
	})
}
