# jsparser-go

A golang implementation of an application that retrieves a list of invocations of a specific imported node module
that has been imported using `requires()` syntax.


There's no CLI at the moment, so update the variables in main as needed to find the import/invocations you are looking for.

### Default Values:
```go
	moduleName := "fs/promises"
	propName := "readFile"
	srcPath := "./resources/app.js"
```


## Running jsparser-go

### With Docker
Use [make rules](./Makefile) to execute the jpsarser-go application ([main.go](./main.go))

To build the docker image and run the application with Docker Compose, run the following command:
```shell
$ make docker
```

### Without Docker
To run jsparser-go in your local golang environment, run the following command:
```shell
$ make run
```
