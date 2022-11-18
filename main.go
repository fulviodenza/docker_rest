package main

import "github.com/fulviodenza/docker_rest/internal/docker_client"

func main() {
	c, err := docker_client.NewDockerClient()
	if err != nil {
		panic(err)
	}

	err = c.Pull("ubuntu:latest")
	if err != nil {
		panic(err)
	}

	err = c.Run("ubuntu", "top")
	if err != nil {
		panic(err)
	}

}
