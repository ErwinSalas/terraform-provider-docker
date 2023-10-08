package main

import (
	"context"

	"log"
	"os"

	tf5server "github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	tf5muxserver "github.com/hashicorp/terraform-plugin-mux/tf5muxserver"

	"github.com/ErwinSalas/terraform-provider-docker/docker"
)

// const providerName = "registry.terraform.io/hashicorp/kubernetes"
const providerName = "dummy-provider"

// Generate docs for website
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	mainProvider := docker.Provider().GRPCProvider
	//manifestProvider := manifest.Provider()

	ctx := context.Background()
	muxer, err := tf5muxserver.NewMuxServer(ctx, mainProvider)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	opts := []tf5server.ServeOpt{}

	tf5server.Serve(providerName, muxer.ProviderServer, opts...)
}

// func containerCliTest() {
// 	cli, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		panic(err)
// 	}

// 	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
// 	fmt.Println(containers)
// 	if err != nil {
// 		panic(err)
// 	}

// 	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
// 	fmt.Println(images)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, container := range containers {
// 		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
// 	}
// }
