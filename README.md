# docker_rest

docker_rest is a simple REST Client which implements some simple docker operations
to run and test it locally you can run the following command from the root of the repo:
`go run main.go`

To interrupt the execution, you need to grab the container id and interrupt if by another terminal
launching the following command:
`docker rm -f {id}`

The documentation can be found in the `/docs` folder.
