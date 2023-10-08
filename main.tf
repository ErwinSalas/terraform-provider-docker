terraform {
  required_providers {
    docker = {
      source  = "opentofu-docker"
      version = "1.0.0"
    }

  }
}


provider "docker" {
  alias = "mydocker"
  host = "unix:///var/run/docker.sock"
}

resource "docker_container" "example" {
  name         = "my-container"
  image        = "nginx:latest"
  ports {
    internal = 80
    external = 8080
  }
}