# Gondest CLI

Gondest CLI is a tool to simplify the creation of application structure using GoFiber with a modular approach like NestJS.

## Requirements

- Go must be installed on your machine. You can download and install it from [here](https://golang.org/dl/).

## Installation

To install Gondest CLI via:

### curl

```bash
curl -sL https://raw.githubusercontent.com/akimabs/gondest/main/install.sh | bash
```

### homebrew

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
|-- utils/
|   |-- response.go
|-- main.go
|-- go.mod
```

### 2. `gondest generate [type] [name]`

Generates files for a new controller, service, module, or model.

#### Arguments:

- `type`: The type of file to generate (`controller`, `service`, `module`, or `model`).
- `name`: The name of the module or model to create.

#### Supported types:

- `module`
- `model`

#### Example:

```bash
gondest generate module user
```

### 3. `gondest config --[type]`

Configures integrations such as databases or other services.

#### Arguments:

##### Type for Config

The type of integration to configure (e.g., `db`, `redis`, more to come).

- `db`: The type of integration to configure database (e.g., mysql, postgres, sqlserver).
- `redis`: The type of integration to configure redis.

##### Flags for Databases:

- `--mysql`: Use MySQL as the database.
- `--postgres`: Use PostgreSQL as the database.
- `--sqlserver`: Use SQL Server as the database.

#### Example:

```bash
gondest config db --postgres
```

## Flags

- `-h, --help`: Displays help for each command.

## Contributors

Special thanks to [ChatGPT](https://openai.com/chatgpt) for assistance in code generation and documentation.

List generated with [contributors-img](https://contrib.rocks). [Updates every 24 hrs]

<a href="https://github.com/akimabs/gondest/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=akimabs/gondest" />
</a>
