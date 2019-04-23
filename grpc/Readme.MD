# Learning notes
## Go modules

Running `go mod init` gives  this error `go: cannot determine module path for source directory`

Give it a module name instead, e.g. `go mod init learngogrpc`

When importing qualify with learngogrpc like this 
```
import (
	pb "learngogrpc/hello"
)
```

## Grpc code generation 

- To include grpc code include plugin options `--go_out=plugins=grpc:...`
- Generated code includes an server interface, this means that the server need to provide a new Server type and all the server interface methods 
