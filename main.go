package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
	"github.com/fulviodenza/docker_rest/internal/docker_client"
)

const interrupt_task_file = "./tmp/interrupt_task.txt"

func watch(watcher fsnotify.Watcher, ch_interrupt, ch_exit chan struct{}) chan struct{} {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
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

	image := docker_client.UBUNTU_IMAGE
	if len(os.Args) > 1 {
		image = os.Args[1]
	}

	ctx := context.Background()

	// instantiate watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("[fsnotify.NewWatcher() error]: ", err)
	}
	defer watcher.Close()

	if err := watcher.Add(interrupt_task_file); err != nil {
		log.Fatal("[watcher.Add error]: ", err)
	}

	var (
		ch_kill      = make(chan struct{}, 1)
		ch_interrupt = make(chan struct{}, 1)
	)

	go watch(*watcher, ch_interrupt, ch_kill)

	// instantiate docker client
	c := docker_client.NewDockerClient()

	containers, err := c.List(ctx)
	if err != nil {
		log.Fatal("[List error]: ", err)
	}

	foundImage := false
	for _, ct := range containers {
		if ct.Image == image {
			foundImage = true
		}
	}

	if !foundImage {
		log.Println("Image not found, pulling...")
		if err := c.Pull(image); err != nil {
			log.Fatal("[Pull error]: ", err)
		}
	}

start:

	// The first three fields in this file are load average figures giving
	// the number of jobs in the run queue (state R) or waiting for disk I/O (state D)
	// averaged over 1, 5, and 15 minutes. They are the same as the load average
	// numbers given by uptime(1) and other programs. The fourth field consists of
	// two numbers separated by a slash (/). The first of these is the number of
	// currently runnable kernel scheduling entities (processes, threads).
	// The value after the slash is the number of kernel scheduling entities
	// that currently exist on the system. The fifth field is the PID of the process
	// that was most recently created on the system.
	idContainer, err := c.Create(image, []string{"cat", "/proc/loadavg"})
	if err != nil {
		log.Fatal("[Create]: error ", err)
	}
	defer c.Destroy(ctx, idContainer)

	fmt.Println(idContainer)

	// The for loop is used to loop endless and at each iteration,
	// it creates the desidred image and run it using its id.
	// The loop stops when the container has been deleted
	// e.g. using the `docker rm -f {id}` command
	for {
		select {
		// When I receive something on the ch_kill channel,
		// recover the main function
		case <-ch_kill:
			os.Exit(1)
		// When I receive something on the ch_interrupt channel,
		// recover the main function
		case <-ch_interrupt:
			goto start
		default:
			if err := c.Start(idContainer); err != nil {
				log.Fatal("[Start error]: ", err)
			}
			if err := c.Logs(idContainer); err != nil {
				log.Fatal("[Logs error]: ", err)
			}
		}
	}
}
