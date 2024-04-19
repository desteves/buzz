package main

import (
	"os"

	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/cloudrun"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const APPNAME = "buzz"

func main() {
	// Create a new Pulumi project
	pulumi.Run(func(ctx *pulumi.Context) error {
		// conf := config.New(ctx, "buzz")

		username := os.Getenv("DOCKER_USR") //conf.Require("DOCKER_USR")
		password := os.Getenv("DOCKER_PAT") //conf.RequireSecret("DOCKER_PAT")

		currentStackName := ctx.Stack()

		// Create a Docker image from a Dockerfile and push it to Docker Hub.
		image, err := docker.NewImage(ctx, APPNAME, &docker.ImageArgs{
			Build: &docker.DockerBuildArgs{
				Platform: pulumi.String("linux/amd64"),
				Context:  pulumi.String("../app"), // Path to the directory with Dockerfile and source.
			},
			ImageName: pulumi.Sprintf(username + "/" + APPNAME + ":" + currentStackName),
			SkipPush:  pulumi.Bool(false),
			Registry: &docker.RegistryArgs{
				Server:   pulumi.String("docker.io"), // Docker Hub server.
				Username: pulumi.String(username),
				Password: pulumi.String(password),
			},
		})
		if err != nil {
			return err
		}

		// Create a new Cloud Run Service using the image
		service, err := cloudrun.NewService(ctx, APPNAME, &cloudrun.ServiceArgs{
			Location: pulumi.String("us-central1"),
			Template: &cloudrun.ServiceTemplateArgs{
				Spec: &cloudrun.ServiceTemplateSpecArgs{
					Containers: cloudrun.ServiceTemplateSpecContainerArray{
						&cloudrun.ServiceTemplateSpecContainerArgs{
							Image: image.ImageName,
							Ports: cloudrun.ServiceTemplateSpecContainerPortArray{
								&cloudrun.ServiceTemplateSpecContainerPortArgs{
									ContainerPort: pulumi.Int(8000),
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Export the URL
		ctx.Export("url", service.Statuses.Index(pulumi.Int(0)).Url())

		// Create Cloudflare CDN
		if currentStackName == "prod" {

			conf := config.New(ctx, "cloudflare")

			// Create a new Cloudflare CDN
			// Configure your Cloudflare Zone ID here.
			zoneID := conf.Require("zoneID")
			domain := conf.Require("domain")
			sub := conf.Require("sub")
			gcp := config.New(ctx, "gcp")
			//  For example, if you are using Cloudflare CDN, you should
			// turn off the "Always use https" option in the "Edge Certificates" tab of the SSL/TLS tab.

			_, err = cloudrun.NewDomainMapping(ctx, "buzz", &cloudrun.DomainMappingArgs{
				Location: pulumi.String("us-central1"),
				Name:     pulumi.String(sub + "." + domain),
				Metadata: &cloudrun.DomainMappingMetadataArgs{
					Namespace: pulumi.String(gcp.Require("project")),
				},
				Spec: &cloudrun.DomainMappingSpecArgs{
					RouteName: service.Name,
				},
			})
			if err != nil {
				return err
			}

			// Set the DNS record for the CDN.
			_, err := cloudflare.NewRecord(ctx, "cdnRecord", &cloudflare.RecordArgs{
				ZoneId:  pulumi.String(zoneID),  // Replace with your actual Zone ID
				Name:    pulumi.String(sub),     // The subdomain or record name
				Type:    pulumi.String("CNAME"), // Typically a CNAME for CDN usage
				Value:   pulumi.String(domain),  // The value of the record, like a CDN endpoint
				Proxied: pulumi.Bool(true),      // Set to true to proxy traffic through Cloudflare (provides CDN and DDoS protection)
			})
			if err != nil {
				return err
			}

			// Enable Argo Smart Routing on the zone.
			_, err = cloudflare.NewArgo(ctx, "cdnArgo", &cloudflare.ArgoArgs{
				ZoneId:       pulumi.String(zoneID),
				SmartRouting: pulumi.String("on"),
			})
			if err != nil {
				return err
			}

		}

		return nil
	})
}
