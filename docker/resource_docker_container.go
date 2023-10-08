package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDockerContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceDockerContainerCreate,
		Read:   resourceDockerContainerRead,
		Update: resourceDockerContainerUpdate,
		Delete: resourceDockerContainerDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Docker container.",
			},
			"image": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Docker image name and tag (e.g., nginx:latest).",
			},
			"ports": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"internal": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The internal port inside the container.",
						},
						"external": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The external port on the host.",
						},
					},
				},
				Required:    true,
				Description: "Port mappings from container to host.",
			},
			// Add more schema attributes as needed for your Docker container configuration.
		},
	}
}

func resourceDockerContainerCreate(d *schema.ResourceData, m interface{}) error {
	cli := m.(*client.Client)
	ctx := context.Background()
	// Parse and extract resource attributes from d.
	imageName := d.Get("image").(string)
	fmt.Printf("🚀 ~  imageName: %s", imageName)
	containerName := d.Get("name").(string)
	fmt.Printf("🚀 ~  containerName: %s", containerName)

	internalPort := d.Get("ports.0.internal").(int)
	externalPort := d.Get("ports.0.external").(int)

	// Convert exposedPorts to the appropriate format for the Docker API.
	portBindings := map[nat.Port][]nat.PortBinding{
		nat.Port(fmt.Sprintf("%d/tcp", internalPort)): {{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", externalPort)}},
	}

	// Create container configuration.
	config := &container.Config{
		Image: imageName,
	}

	// Host configuration with port mappings.
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	// Create the container.
	resp, err := cli.ContainerCreate(
		ctx,
		config,
		hostConfig,
		nil, // Networking configuration (optional)
		nil,
		containerName, // Container name (Docker will generate one if empty)
	)
	if err != nil {
		return fmt.Errorf("Error creating Docker container: %w", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	// Set the ID of the created container in d.
	d.SetId(resp.ID)

	return nil
}

func resourceDockerContainerRead(d *schema.ResourceData, m interface{}) error {
	cli := m.(*client.Client)

	// Get the ID of the container from the resource data.
	containerID := d.Id()

	// Use Docker client to inspect the container.
	containerJSON, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return fmt.Errorf("Error inspecting Docker container: %w", err)
	}

	d.Set("image_name", containerJSON.Config.Image)
	d.Set("container_name", containerJSON.Name)
	// Add more attributes as needed based on the Docker container inspection results.

	return nil
}

func resourceDockerContainerUpdate(d *schema.ResourceData, m interface{}) error {
	cli := m.(*client.Client)

	// Get the ID of the container from the resource data.
	containerID := d.Id()

	// Check if the container still exists (may have been removed externally).
	_, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		// Handle the case where the container no longer exists.
		// You can choose to recreate the container or return an error based on your use case.
		return fmt.Errorf("Docker container no longer exists: %w", err)
	}

	// Implement logic to compare the current state with the desired state.
	// You may need to read attributes from the resource data (d) and the current container state.

	// Example: Update the container name (assuming it's an editable attribute).
	newContainerName := d.Get("container_name").(string)
	err = cli.ContainerRename(context.Background(), containerID, newContainerName)
	if err != nil {
		return fmt.Errorf("Error renaming Docker container: %w", err)
	}

	// Implement similar logic for other attributes you want to update.

	return nil
}

func resourceDockerContainerDelete(d *schema.ResourceData, m interface{}) error {
	cli := m.(*client.Client)

	// Get the ID of the container from the resource data.
	containerID := d.Id()

	// Use Docker client to stop and remove the container.
	err := cli.ContainerStop(context.Background(), containerID, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("Error stopping Docker container: %w", err)
	}

	err = cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true, // Remove associated volumes if needed
		Force:         true, // Force removal even if the container is running
	})
	if err != nil {
		return fmt.Errorf("Error removing Docker container: %w", err)
	}

	// Clear the ID from the resource data.
	d.SetId("")

	return nil
}
