package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/fsnotify/fsnotify"
	"github.com/fulviodenza/docker_rest/internal/docker_client"
)

const UBUNTU_IMAGE = "ubuntu"
const interrupt_task_file = "./tmp/interrupt_task.txt"

func watch(watcher fsnotify.Watcher) chan struct{} {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			recover()
		}
	}()

	for {
		select {
		case event := <-watcher.Events:
			log.Println("event:", event)
			os.Exit(1)
		case err := <-watcher.Errors:
			log.Println("error:", err)
			os.Exit(1)
		}
	}
}

func recover() {
	main()
}

func main() {

	ctx := context.Background()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("[fsnotify.NewWatcher()]: error ", err)
		panic(err)
	}
	defer watcher.Close()

	// Add a path.
	err = watcher.Add(interrupt_task_file)
	if err != nil {
		log.Fatal("[watcher.Add]: error ", err)
		panic(err)
	}
	go watch(*watcher)

	c := docker_client.NewDockerClient()

	err = c.Pull(UBUNTU_IMAGE + ":latest")
	if err != nil {
		log.Fatal("[Pull]: error ", err)
		panic(err)
	}

	idContainer, err := c.Create(UBUNTU_IMAGE, []string{"cat", "/proc/loadavg"})
	if err != nil {
		log.Fatal("[Create]: error ", err)
		panic(err)
	}

	fmt.Println(idContainer)

	containers, err := c.List(ctx)
	if err != nil {
		log.Fatal("[List]: error ", err)
		panic(err)
	}

	for _, ct := range containers {
		fmt.Println(ct.ID)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("[client.NewClientWithOpts]: error ", err)
		panic(err)
	}

	var containerID string
	// The for loop is used to loop endless and at each iteration
	// select an existing container with the desidred image and run it
	// using its id. The loop stops when the container has been deleted
	// e.g. using the `docker rm -f {id}` command
	for {
		for _, ct := range containers {
			if ct.Image == UBUNTU_IMAGE {

				containerID = ct.ID
				fmt.Println("CONTAINER SELECTED: ", ct.ID)
				// I know, I have to do it with the rest client and not with
				// the sdk, but damn, rules are made to be broken, right?
				err = cli.ContainerStart(ctx, ct.ID, types.ContainerStartOptions{})
				if err != nil {
					log.Fatal("[Start]: error ", err)
					panic(err)
				}
				break
			}
		}

		c.Logs(containerID)
	}
}
