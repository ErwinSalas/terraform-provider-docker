package docker

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"docker_host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DOCKER_HOST"}, "unix:///var/run/docker.sock"),
				Description: "Docker daemon endpoint.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"docker_container": resourceDockerContainer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Read the Docker host configuration from the provider schema.
	dockerHost := d.Get("docker_host").(string)

	// Initialize the Docker client.
	cli, err := client.NewClientWithOpts(client.WithHost(dockerHost))
	if err != nil {
		return nil, fmt.Errorf("Error creating Docker client: %w", err)
	}

	// Return the Docker client as the provider configuration.
	return cli, nil
}
