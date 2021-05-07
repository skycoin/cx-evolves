### Summary

CX executes worker executes the cx program and send back the output. 

### Deployment of Worker
For more information, run
```
go run cmd/server.go help 
```

```
go run cmd/server.go -port=[port number for worker to use]
```

### Example

```
go run cmd/server.go -port=9090
```

### Notes
1. If no arguments are specified, the program will deploy 1 worker at port 9090.
