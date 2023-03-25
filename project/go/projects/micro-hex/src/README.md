# GO Microservice

## Getting started

This scaffolded project will build and run as is. All it will do is log "Hello" to the console.
To get started:

1. Navigate to `dispatcher/service.go` and define your core service.
2. Implement your service and core logic in `dispatcher/logic.go`.
3. Run your service.

## Project structure
This template is following the Hexagonal Architecture. The general idea is that there is a core application which for
this project has been called `dispatcher` and multiple ports that essentially are connections to the outside. Inside the
directory with the same name are multiple files that contain definitions like domain models and configuration objects in
`moels.go`, serializers for the models in `serializers.go`, interfaces for all the "ports" in `ports.go`, the core
service interface in `service.go` and the implementation of the core service in `logic.go`. All ports should be located
under the directory `ports/`.

There are two types of ports in essence: _driven_ and _driver_ ports. Driven ports are used by the core service to perform
some operation like saving data with a database port, or perhaps sending a message to a message broker. Driver ports are
ports that initiate some process and call a function in the core service. This can be a message consumer, a http api or
similar interfaces that will do something based on an event from the outside. The two types are not currently separated
and can both be found under the `ports/` directory.

The project structure looks something like this:

```
├ my-service/
│ ├ cicd/ <- DevOps stuff
│ │
│ ├ core/
│ │ ├ logic.go
│ │ ├ models.go
│ │ ├ ports.go
│ │ ├ serializers.go
│ │ └ service.go
│ │
│ ├ ports/
│ │ ├ config/
│ │ │ └ config-file-loader-port.go
│ │ ├ logger/
│ │ │ ├ zap-config-builder.go
│ │ │ └ zap-logger-port.go
│ │
│ ├ serializers/
│ │
│ ├ config.json
│ ├ config.local.json
│ ├ Dockerfile
│ ├ go.mod
│ ├ go.sum
│ ├ main.go
│ ├ README.md
│ ├ version
│ └ zap.json
```

## The main function

In the main function located in `main.go` there is a top level context called `mainCtx`. This context's done channel
will
be closed when the program receives a SIGINT or SIGTERM from the operating systems. Pass this context to any part of the
code that needs to be signaled when the application is shutting down.

All background services should implement the interface `BackgroundService` defined in `ports.go`, which means they will
have
a function with the signature `Start(ctx context.Context) error`. All services that need to run in the background should
be defined in the function `getBackgroundServices`, usually by simply calling their factory method. In the main function
an `errgroup` is created and each background service is started with `g.Go()` (g being the errgroup), which essentially
runs them in a go routine. A derived context from `mainCtx` is passed to all services at this point. Using `g.Wait()`
the
main go routine is blocked until all background services have stopped gracefully OR one or more service has returned an
error.
The first background service to return an error will stop the application and print the error (which is return
by `g.Wait()`).
