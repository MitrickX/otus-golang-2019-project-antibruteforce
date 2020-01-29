<h1>Otus golang 2019 project - Anti brute force microservice</h1>
Specification [here](https://github.com/OtusTeam/Go/blob/master/projects/anti-bruteforce.md) <br>

<h2>Algorithms</h2>
Leaky Bucket algorithm specification is [here](http://www.idc-online.com/technical_references/pdfs/data_communications/Leaky_Bucket_Algorithm.pdf)
<br>
Token Bucket algorithm specification is [here](https://en.wikipedia.org/wiki/Token_bucket)
<br><br>

For this anti brute force service [Token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket) is more suitable and it is implemented
<br>

<h2>GRPC server</h2>
To run server execute command grpc, port specify by env var GRPC_PORT
```
GRPC_PORT=port ./antibruteforce grpc
```
If GRPC_PORT is missed port server will be started on port 50051

<h2>GRPC CLI client</h2>
Before execute client commands you must specified server host by set up env var GRPC_SERVER_HOST<br>
Also you can specify port of GRPC server by env var GRPC_SERVER_HOST, if this var is missed it will be assumed that port is 50051<br>
<br>
**Example**:
```
export GRPC_SERVER_HOST=localhost
export GRPC_SERVER_PORT=50052
./antibruteforce auth test 1234 193.192.170.13
```

<h3>List of CLI client commands</h3>
```
add <kind> <ip> [flags] - Add IP into black or white list
delete <kind> <ip> [flags] - Delete ip from black of white list
clear --login=<login> --password=<password> --ip=<ip> [flags] - Clear bucket(s) for login, password or ip
auth <login> <password> <ip> [flags] - Check that auth is allowed for login, password and ip
```
Each command support --help option. Use it for explore details of commands

<h2>Integration and unit tests</h2>
For run integration tests go to build/package and run 
```make test```
For run unit tests run from the root of project and run ```go test -v --race -tags unit ./...```