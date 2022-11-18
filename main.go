package main

import (
	"fmt"

	"github.com/fulviodenza/docker_rest/internal/docker_client"
)

/*
&{POST http://%2Fvar%2Frun%2Fdocker.sock/v1.41/containers/1d3c35a41a5ad5720021cdf15847decc0234266de67b5ceacb11cd55f43e0b4a/start HTTP/1.1 1 1 map[Content-Type:[text/plain]] {} 0x14aa1e0 0 [] false docker map[] map[] <nil> map[]   <nil> <nil> <nil> 0xc00019c000}

*/

const UBUNTU_IMAGE = "ubuntu"

func main() {
	c, err := docker_client.NewDockerClient()
	if err != nil {
		panic(err)
	}

	err = c.Pull(UBUNTU_IMAGE + ":latest")
	if err != nil {
		panic(err)
	}

	idContainer, err := c.Create(UBUNTU_IMAGE, []string{"cat", "/proc/loadavg"})
	if err != nil {
		panic(err)
	}

	fmt.Println(idContainer)

	err = c.Start(idContainer, UBUNTU_IMAGE)
	if err != nil {
		panic(err)
	}
	// ctx := context.Background()
	// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// if err != nil {
	// 	panic(err)
	// }
	// defer cli.Close()

	// reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	// defer reader.Close()
	// io.Copy(os.Stdout, reader)

	// resp, err := cli.ContainerCreate(ctx, &container.Config{
	// 	Image: "alpine",
	// 	Cmd:   []string{"cat", "/proc/loadavg"},
	// 	Tty:   false,
	// }, nil, nil, nil, "")
	// if err != nil {
	// 	panic(err)
	// }

	// for {
	// 	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
	// 		panic(err)
	// 	}

	// 	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	// 	select {
	// 	case err := <-errCh:
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	case <-statusCh:
	// 	}

	// 	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	// }

}
