# docker_rest

docker_rest is a simple REST Client which implements some simple docker operations. </br>
To run and test it locally you can run the following command from the root of the repo: </br>
`go run main.go` 
</br>
To interrupt the execution, you need to grab the container id and interrupt if by another terminal
launching the following command: </br>
`docker rm -f {id}`

## documentation
The documentation can be found in the `/docs` folder.

## output
If you'll run the program, you'll see strange numbers poppin' out, that numbers come from the command: </br>
`cat /proc/loadavg`
</br>

As said in a stackoverflow question:

- The first three fields in this file are load average figures giving the number of jobs in the run queue (state R) or waiting for disk I/O (state D) averaged over 1, 5, and 15 minutes. They are the same as the load average numbers given by uptime(1) and other programs.

- The fourth field consists of two numbers separated by a slash (/). The first of these is the number of currently executing kernel scheduling entities (processes, threads); this will be less than or equal to the number of CPUs. The value after the slash is the number of kernel scheduling entities that currently exist on the system.

- The fifth field is the PID of the process that was most recently created on the system.

(source [stackoverflow](https://stackoverflow.com/questions/11987495/what-do-the-numbers-in-proc-loadavg-mean-on-linux))
