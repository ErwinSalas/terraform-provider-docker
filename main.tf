terraform {
  required_providers {
    docker = {
      source = "ErwinSalas/docker"
      version = "1.0.1"
    }
  }
}


provider "docker" {
  docker_host = "unix:///var/run/docker.sock"
}

resource "docker_container" "example" {
  name         = "my-container"
  image        = "nginx:latest"
  ports {
    internal = 80
    external = 8080
  }
}