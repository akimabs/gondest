# Gondest CLI

Gondest CLI is a tool to simplify the creation of application structure using GoFiber with a modular approach like NestJS.

## Installation

To install Gondest CLI:

```bash
to be launch
```

## Basic Commands

### 1. `gondest init [app-name]`

Initializes a new project with the standard structure for GoFiber.

#### Example:

```bash
gondest init myApp
```

This will create a directory `myApp` with the following structure:

```
myApp/
|-- domains/
|   |-- app.controller.go
|   |-- app.service.go
|   |-- app.module.go
|-- main.go
|-- go.mod
```

### 2. `gondest generate [type] [name]`

Generates files for a new controller, service, or module.

#### Arguments:

- `type`: The type of file to generate (`controller`, `service`, or `module`).
- `name`: The name of the module or controller to create.

#### Example:

```bash
gondest generate module user
```

This command will create files like `user.controller.go`, `user.service.go`, and `user.module.go` in the `domains/` directory.

## Flags

- `-h, --help`: Displays help for each command.
