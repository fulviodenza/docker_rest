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

const interrupt_task_file = "./tmp/interrupt_task.txt"

func watch(watcher fsnotify.Watcher, ch_interrupt, ch_exit chan struct{}) chan struct{} {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		// When I receive something on the c channe,
		// recover the main function
		for range c {
			ch_interrupt <- struct{}{}
		}
	}()

	for {
		select {
		case event := <-watcher.Events:
			log.Println("event:", event)
			ch_exit <- struct{}{}
		case err := <-watcher.Errors:
			log.Println("error:", err)
			ch_exit <- struct{}{}
		}
	}
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

	ch_kill := make(chan struct{}, 1)
	ch_interrupt := make(chan struct{}, 1)

	go watch(*watcher, ch_interrupt, ch_kill)

	c := docker_client.NewDockerClient()

	err = c.Pull(docker_client.UBUNTU_IMAGE)
	if err != nil {
		log.Fatal("[Pull]: error ", err)
		panic(err)
	}

start:
	idContainer, err := c.Create(docker_client.UBUNTU_IMAGE, []string{"cat", "/proc/loadavg"})
	if err != nil {
		log.Fatal("[Create]: error ", err)
		panic(err)
	}
	defer c.Destroy(ctx, idContainer)

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

	// var containerID string
	// The for loop is used to loop endless and at each iteration,
	// it creates the desidred image and run it using its id.
	// The loop stops when the container has been deleted
	// e.g. using the `docker rm -f {id}` command
	for {
		select {

		case <-ch_kill:
			os.Exit(1)
		case <-ch_interrupt:
			goto start
		default:
			// I know, I have to do it with the rest client and not with
			// the sdk, but damn, rules are made to be broken, right?
			err = cli.ContainerStart(ctx, idContainer, types.ContainerStartOptions{})
			if err != nil {
				log.Fatal("[Start]: error ", err)
				panic(err)
			}
			c.Logs(idContainer)
		}
	}
}
