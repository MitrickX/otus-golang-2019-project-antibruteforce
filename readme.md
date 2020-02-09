<h1>Otus golang 2019 project - Anti brute force microservice</h1>

[![Build Status](https://travis-ci.com/MitrickX/otus-golang-2019-project-antibruteforce.svg?branch=master)](https://travis-ci.com/MitrickX/otus-golang-2019-project-antibruteforce)

[![Go Report Card](https://goreportcard.com/badge/github.com/MitrickX/otus-golang-2019-project-antibruteforce)](https://goreportcard.com/report/github.com/MitrickX/otus-golang-2019-project-antibruteforce)

Specification https://github.com/OtusTeam/Go/blob/master/projects/anti-bruteforce.md <br>

<h2>Algorithm</h2>
For this anti brute force service Token bucket algorithm (https://en.wikipedia.org/wiki/Token_bucket) is used
<br>

<h2>GRPC server</h2>
To run server execute command grpc, port specify by env var GRPC_PORT

> GRPC_PORT=port ./antibruteforce grpc

If GRPC_PORT is missed port server will be started on port 50051

To run server in container do this.

>
> make run<br>
> 

Port 50051 will be available for communicate with GRPC server.<br>
Now you can build CLI client and run commands  

<h2>GRPC CLI client</h2>
Before execute client commands you could specified server host by set up env var GRPC_SERVER_HOST, if this var is missed it will be assumed that server is 127.0.0.1 host<br>
Also you can specify port of GRPC server by env var GRPC_SERVER_HOST, if this var is missed it will be assumed that port is 50051<br>
<br>
<strong>Example (default)</strong>:

>
> make build<br>
> ./antibruteforce auth test 1234 193.192.170.13<br>
>

<strong>Example (set vars)</strong>:

>
> make build<br>
> export GRPC_SERVER_HOST=localhost<br>
> export GRPC_SERVER_PORT=50052<br>
> ./antibruteforce auth test 1234 193.192.170.13<br>
>

<h3>List of CLI client commands</h3>

>
> add <kind> <ip> [flags] - Add IP into black or white list<br>
> delete <kind> <ip> [flags] - Delete ip from black of white list<br>
> clear --login=<login> --password=<password> --ip=<ip> [flags] - Clear bucket(s) for login, password or ip<br>
> auth <login> <password> <ip> [flags] - Check that auth is allowed for login, password and ip<br>
>
  
Each command support --help option. Use it for explore details of commands

<h2>Integration and unit tests</h2>
For run integration tests go to build/package and run 

> make test<br>

For run unit tests run from the root of project and run 

> go test -v --race -tags unit ./...<br>
