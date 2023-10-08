package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDockerNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceDockerNetworkCreate,
		Read:   resourceDockerNetworkRead,
		Update: resourceDockerNetworkUpdate,
		Delete: resourceDockerNetworkDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Docker network.",
			},
			"driver": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "bridge",
				Description: "The driver to use for the Docker network.",
			},
			// Add more attributes as needed.
		},
	}
}

func resourceDockerNetworkCreate(d *schema.ResourceData, m interface{}) error {
	cli := m.(*client.Client)

	// Define the network configuration based on the Terraform resource attributes.
	networkName := d.Get("name").(string)
	networkDriver := d.Get("driver").(string)

	// Create the network using the Docker client.
	networkCreateOpts := types.NetworkCreate{
		Driver: networkDriver,
	}
	network, err := cli.NetworkCreate(context.Background(), networkName, networkCreateOpts)
	if err != nil {
		return fmt.Errorf("Error creating Docker network: %w", err)
	}

	// Set the ID of the created network as the resource ID in Terraform state.
	d.SetId(network.ID)

	return nil
}

// Define the Read function for Docker network resource.
func resourceDockerNetworkRead(d *schema.ResourceData, m interface{}) error {
	cli := m.(*client.Client)

	// Get the network ID from the resource data.
	networkID := d.Id()

	// Check if the Network still exists (may have been removed externally).
	networkInspect, err := cli.NetworkInspect(context.Background(), networkID, types.NetworkInspectOptions{})
	if err != nil {
		// If the network is not found, mark it as destroyed and remove it from Terraform state.
		return fmt.Errorf("Error inspecting Docker network: %w", err)
	}

	// Update resource data with network information.
	d.Set("name", networkInspect.Name)
	// Set other attributes as needed.

	return nil
}

// Define the Update function for Docker network resource.
func resourceDockerNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	cli := m.(*client.Client)

	// Get the network ID from the resource data.
	networkID := d.Id()

	// Get the desired driver from the Terraform configuration.
	newDriver := d.Get("driver").(string)

	// Create a new network with the desired driver.
	newNetwork, err := createDockerNetwork(cli, newDriver)
	if err != nil {
		return fmt.Errorf("Error creating a new Docker network: %w", err)
	}

	// Update containers to use the new network.
	err = updateContainersForNewNetwork(cli, networkID, newNetwork)
	if err != nil {
		return fmt.Errorf("Error updating containers: %w", err)
	}

	// Optionally, remove the old network if it's no longer needed.
	// This step is optional and depends on your requirements.

	// Mark the resource as "updated" in Terraform state.
	d.SetId(newNetwork.ID)

	return nil
}

func resourceDockerNetworkDelete(d *schema.ResourceData, m interface{}) error {
	cli := m.(*client.Client)

	// Get the network ID from the resource data.
	networkID := d.Id()

	// Remove the Docker network.
	err := cli.NetworkRemove(context.Background(), networkID)
	if err != nil {
		// If the network is not found, consider it as successfully deleted.
		if client.IsErrNotFound(err) {
			return nil
		}
		return fmt.Errorf("Error deleting Docker network: %w", err)
	}

	// Mark the resource as destroyed and remove it from Terraform state.
	d.SetId("")

	return nil
}
