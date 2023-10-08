package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// List containers connected to a specific network.
func listContainersConnectedToNetwork(cli *client.Client, networkID string) ([]types.Container, error) {
	containers := []types.Container{}

	// List all containers.
	containerListOpts := types.ContainerListOptions{}
	allContainers, err := cli.ContainerList(context.Background(), containerListOpts)
	if err != nil {
		return nil, fmt.Errorf("Error listing containers: %w", err)
	}

	// Iterate through all containers and check if they are connected to the specified network.
	for _, container := range allContainers {
		networks := container.NetworkSettings.Networks
		if _, exists := networks[networkID]; exists {
			containers = append(containers, container)
		}
	}

	return containers, nil
}

// Update a container's network settings to use the new network.
func updateContainerNetworkSettings(cli *client.Client, containerID, oldNetworkID, newNetworkID string) error {
	// Disconnect the container from the old network.
	disconnectErr := cli.NetworkDisconnect(context.Background(), oldNetworkID, containerID, true)
	if disconnectErr != nil {
		return fmt.Errorf("Error disconnecting container from the old network: %w", disconnectErr)
	}

	// Connect the container to the new network.
	connectErr := cli.NetworkConnect(context.Background(), newNetworkID, containerID, nil)
	if connectErr != nil {
		return fmt.Errorf("Error connecting container to the new network: %w", connectErr)
	}

	return nil
}

// Implement logic to generate a unique network name.
func generateUniqueNetworkName() string {
	// You can implement your logic here to generate a unique network name.
	// This could involve a combination of a prefix, timestamp, or a random string.
	// Ensure that the generated name is unique to avoid conflicts.
	return "my_unique_network_name"
}

// Create a new Docker network with the specified driver.
func createDockerNetwork(cli *client.Client, driver string) (types.NetworkCreateResponse, error) {
	networkName := generateUniqueNetworkName() // Generate a unique network name.
	networkCreateOpts := types.NetworkCreate{
		Driver: driver,
	}

	return cli.NetworkCreate(context.Background(), networkName, networkCreateOpts)
}

// Update containers to use the new network.
func updateContainersForNewNetwork(cli *client.Client, oldNetworkID string, newNetwork types.NetworkCreateResponse) error {
	// List containers that are connected to the old network.
	containers, err := listContainersConnectedToNetwork(cli, oldNetworkID)
	if err != nil {
		return err
	}

	// Update each container's network settings to use the new network.
	for _, container := range containers {
		err := updateContainerNetworkSettings(cli, container.ID, oldNetworkID, newNetwork.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
