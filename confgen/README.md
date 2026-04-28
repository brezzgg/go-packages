# confgen

**confgen** is a modern code generation platform for Go configuration files.

Instead of manually writing config structs, default values, and validation logic — you define a single HCL schema and confgen generates everything for you: typed Go structs, default constructors, and required field validation.

## How It Works

```
schema.hcl  ──▶  confgen  ──▶  config.gen.go
                                 ├── typed structs
                                 ├── Default*() constructors
                                 └── Defaults() / validation hooks
```

At runtime, `confgen.Unmarshal` ties it all together:
1. Decodes your YAML / JSON / TOML file into the generated struct
2. Fills missing fields with schema-defined defaults
3. Returns an error if any required field is absent

## Installation

```bash
go get -u github.com/brezzgg/go-packages/confgen@latest
go install github.com/brezzgg/go-packages/confgen/cmd/confgen@latest
```

## Get Started

### 1. Define output parameters

Every `.hcl` schema file must contain an `output` block:

```hcl
output {
  package = "config"         # generated package name
  output  = "gen/config"     # relative output path
  formats = ["yaml", "json"] # struct tag formats
}
```

> **Tip:** Set `output = "stdout"` to print generated code to the console instead of writing to a file.

### 2. Describe your configuration schema

Create a `.hcl` file with any name and define your configuration using a `generate` block:

```hcl
# schema.hcl

generate "config" {
  field "address" {
    type = "string"
  }

  field "port" {
    type    = "uint16"
    default = "8080"
  }

  field "timeout" {
    type    = "int"
    default = "10"
    desc    = "timeout in seconds"
  }

  field "admin_secret" {
    type     = "string"
    required = true
    desc     = "admin secret JWT token"
  }
}
```

Fields fall into two categories:

| Type | Behavior |
|---|---|
| **Default field** | Uses the specified default value if not provided in the config file |
| **Required field** | Returns a validation error if not provided |

### 3. Generate


```
confgen schema.hcl
```

The generated file `gen/config/config.gen.go` will look something like this:

```go
package config

type Config struct {
    // Address generated from hcl field "address". Default value is "".
    Address string `yaml:"address" json:"address"`
    // Port generated from hcl field "port". Default value is 8080.
    Port uint16 `yaml:"port" json:"port"`
    // Timeout generated from hcl field "timeout". Default value is 10.
    // Description: timeout in seconds
    Timeout int `yaml:"timeout" json:"timeout"`
    // AdminSecret generated from hcl field "admin_secret". Default value is "". Field is required.
    // Description: admin secret JWT token
    AdminSecret string `yaml:"admin_secret" json:"admin_secret" confgen_required:"true"`
}

func DefaultConfig() Config {
    return Config{
        Address:     "",
        Port:        8080,
        Timeout:     10,
        AdminSecret: "",
    }
}

func (c Config) Defaults() any {
    return DefaultConfig()
}
```

### 4. Parse

**Create your config file:**

```yaml
# schema.yaml
address: "0.0.0.0"
port: 9090
timeout: 30
admin_secret: "supersecret"
```

**Parse it in Go:**

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/brezzgg/go-packages/confgen"
    schema "my-module/gen/config"
    "gopkg.in/yaml.v3"
)

func main() {
    data, err := os.ReadFile("schema.yaml")
    if err != nil {
        panic(err)
    }

    var cfg schema.Config
    if err := confgen.Unmarshal(data, &cfg, yaml.Unmarshal); err != nil {
        panic(err)
    }

    b, _ := json.MarshalIndent(&cfg, "", "  ")
    fmt.Println(string(b))
}
```

`confgen.Unmarshal` accepts any unmarshaler function — use `yaml.Unmarshal`, `json.Unmarshal`, or any compatible third-party decoder.

All fields provided — output reflects the config file:

```json
{
  "address": "0.0.0.0",
  "port": 9090,
  "timeout": 30,
  "admin_secret": "supersecret"
}
```

`port` and `timeout` omitted — defaults from the schema are applied automatically:

```json
{
  "address": "0.0.0.0",
  "port": 8080,
  "timeout": 10,
  "admin_secret": "supersecret"
}
```

`admin_secret` omitted — `confgen.Unmarshal` returns an error:

```
field "AdminSecret" is required
```

> **Tip:** To see more examples click [here](https://github.com/brezzgg/go-packages/blob/main/examples/confgen).

## Schema Reference

### Field Attributes

| Attribute | Description |
|---|---|
| `type` | Field type (see table above) |
| `default` | Default value as string literal |
| `required` | If `true`, returns error when field is missing |
| `desc` | Description, added as a comment in generated code |
| `object` | When `type = "object"` or `type = "list"`. Name of the referenced object block |

### Blocks

#### `generate`

The top-level configuration struct. Every `generate` block produces a Go struct along with a `Default*()` constructor and a `Defaults()` method used by `confgen.Unmarshal`.

```hcl
generate "app" {
  field "debug" {
    type    = "bool"
    default = "false"
  }

  field "server" {
    type   = "object"
    object = "server"
  }
}
```

#### `object`

A reusable nested struct. Objects are referenced from `generate` or other `object` blocks via `type = "object"` and `object = "block_name"`.

```hcl
object "server" {
  field "host" {
    type    = "string"
    default = "0.0.0.0"
  }

  field "port" {
    type    = "uint16"
    default = "8080"
  }
}
```

The output will look like this:

```go
type App struct {
    Debug  bool   `yaml:"debug" json:"debug"`
    Server Server `yaml:"server" json:"server"`
}

type Server struct {
    Host string `yaml:"host" json:"host"`
    Port uint16 `yaml:"port" json:"port"`
}
```
