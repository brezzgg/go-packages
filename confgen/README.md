# confgen

**confgen** is a Go code generator for configuration structs based on HCL schemas.

Instead of manually writing structs, default constructors, and validation logic — define your schema once in an `.hcl` file and `confgen` generates everything else.

## How It Works

The CLI handles **parsing** `.hcl` files, **generating** Go source via `text/template`, and **writing** output to disk (or stdout). At runtime, `confgen.Unmarshal` ties together decoding, default merging, and validation in a single call.

## Installation

```bash
go get -u github.com/brezzgg/go-packages/confgen@latest
go install github.com/brezzgg/go-packages/confgen/cmd/confgen@latest
```

## Quick Start

### 1. Define the output block

Every `.hcl` file must contain an `output` block:

```hcl
output {
  package = "config"          # generated Go package name
  output  = "internal/config" # output path for the generated file
  formats = ["yaml", "json"]  # struct tag formats
}
```

> Set `output = "stdout"` to print generated code to the console instead of writing to a file.

### 2. Describe your configuration schema

```hcl
# schema.hcl

output {
  package = "config"
  output  = "internal/config"
  formats = ["yaml", "json"]
}

generate "config" {
  field "address" {
    type    = "string"
    default = "0.0.0.0"
  }

  field "port" {
    type    = "uint16"
    default = "8080"
    desc    = "listening port"
  }

  field "admin_secret" {
    type     = "string"
    required = true
    desc     = "JWT secret for admin access"
  }

  field "server" {
    type   = "object"
    object = "server_options"
  }
}

object "server_options" {
  field "timeout" {
    type    = "int"
    default = "30"
    desc    = "timeout in seconds"
  }

  field "max_conns" {
    type    = "int"
    default = "100"
  }
}
```

### 3. Run the generator

```bash
confgen schema.hcl
```

Output — `internal/config/config.gen.go`:

```go
package config

type Config struct {
    Address     string        `yaml:"address" json:"address"`
    Port        uint16        `yaml:"port" json:"port"`
    AdminSecret string        `yaml:"admin_secret" json:"admin_secret" confgen_required:"true"`
    Server      ServerOptions `yaml:"server" json:"server"`
}

type ServerOptions struct {
    Timeout  int `yaml:"timeout" json:"timeout"`
    MaxConns int `yaml:"max_conns" json:"max_conns"`
}

func DefaultConfig() Config {
    return Config{
        Address:     "0.0.0.0",
        Port:        8080,
        AdminSecret: "",
        Server:      DefaultServerOptions(),
    }
}

func (c Config) Defaults() any { return DefaultConfig() }

func DefaultServerOptions() ServerOptions {
    return ServerOptions{Timeout: 30, MaxConns: 100}
}

func (c ServerOptions) Defaults() any { return DefaultServerOptions() }
```

### 4. Use the generated code

```go
package main

import (
    "fmt"
    "github.com/brezzgg/go-packages/confgen"
    schema "my-module/internal/config"
    "gopkg.in/yaml.v3"
)

func main() {
    var cfg schema.Config
    decoder := confgen.NewDecoder("config.yaml", yaml.Unmarshal)
    if err := decoder.Decode(&cfg); err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", cfg)
}
```

Or use `confgen.Unmarshal` directly:

```go
data, _ := os.ReadFile("config.yaml")

var cfg schema.Config
if err := confgen.Unmarshal(data, &cfg, yaml.Unmarshal); err != nil {
    panic(err) // e.g. field "AdminSecret" is required
}
```

`confgen.Unmarshal` accepts any compatible unmarshaler — `yaml.Unmarshal`, `json.Unmarshal`, or any third-party decoder.

## CLI Reference

```bash
confgen [command]
```

#### `generate`

Generate config form hcl schema

```bash
confgen generate <path> [path2 ...] [flags]
```

`path` can be either an `.hcl` file or a directory. When a directory is given, all `.hcl` files inside it are processed.

| Flag                       | Description                                                 |
| -------------------------- | ----------------------------------------------------------- |
| `-p`, `--filename-pattern` | specify filenames regular expression (default "^.+\\.hcl$") |
| `-r`, `--recursive`        | recursively walk directories                                |
| `-d`, `--working-dir`      | specify working directory (default ".")                     |

#### `validate`

Validate hcl schema syntax

```bash
confgen validate <path> [path2 ...] [flags]
```

| Flag                       | Description                                                 |
| -------------------------- | ----------------------------------------------------------- |
| `-p`, `--filename-pattern` | specify filenames regular expression (default "^.+\\.hcl$") |
| `-r`, `--recursive`        | recursively walk directories                                |
| `-d`, `--working-dir`      | specify working directory (default ".")                     |

#### `sample`

Generate example config file based on hcl schema

```bash
confgen validate <schema.hcl> <json|yaml> [flags]
```

| Flag                  | Description                             |
| --------------------- | --------------------------------------- |
| `-d`, `--working-dir` | specify working directory (default ".") |

### Usage samples

```bash
# single file
confgen generate schemas/schema.hcl

# multiple files
confgen generate schemas/service.hcl /schemas/v2/service.hcl

# directory (non-recursive)
confgen generate schemas

# directory (recursive)
confgen generate -r schemas

# validate
confgen validate schemas/v2/schema.hcl

# sample in stdout
confgen sample schemas/schema.hcl yaml

# sample in file
confgen sample schemas/schema.hcl yaml > sample.yaml
```

## Schema Reference

### Structure overview

```
.hcl file
├── output                        (required, exactly one)
│   ├── package  string           generated Go package name
│   ├── output   string           output path or "stdout"
│   └── formats  []string         struct tag formats, e.g. ["yaml", "json"]
│
├── generate "<name>"             (one or more)
│   └── field "<name>"
│       ├── type      string      field type
│       ├── default   string      default value as string literal
│       ├── required  bool        fail if missing in config file
│       ├── desc      string      added as a Go comment
│       └── object    string      ref to object block (for type = "object"|"list")
│
└── object "<name>"               (zero or more, reusable nested structs)
    └── field "<name>"
        └── ...                   same attributes as above
```

### Field Attributes

| Attribute  | Type   | Description                                                                       |
| ---------- | ------ | --------------------------------------------------------------------------------- |
| `type`     | string | Field type (see table below)                                                      |
| `default`  | string | Default value as a string literal                                                 |
| `required` | bool   | If `true`, returns an error when the field is missing from the config             |
| `desc`     | string | Description, added as a comment in the generated code                             |
| `object`   | string | Name of the referenced `object` block (when `type = "object"` or `type = "list"`) |

### Supported Types

| HCL type                     | Go type                    | Zero value        |
| ---------------------------- | -------------------------- | ----------------- |
| `string` / `str`             | `string`                   | `""`              |
| `bool` / `boolean`           | `bool`                     | `false`           |
| `int`, `int8` ... `int64`    | corresponding int type     | `0`               |
| `uint`, `uint8` ... `uint64` | corresponding uint type    | `0`               |
| `float32`, `float64`         | `float32` / `float64`      | `0`               |
| `byte`                       | `byte`                     | `0`               |
| `object` / `obj`             | nested struct              | `Default<Name>()` |
| `list [elements type]`       | `[]string` or `[]<Object>` | `nil`             |

### Blocks

#### `generate`

The top-level configuration block. Each `generate` block produces a Go struct, a `Default*()` constructor, and a `Defaults()` method.

```hcl
generate "app" {
  field "debug" {
    type    = "bool"
    default = "false"
  }
}
```

#### `object`

A reusable nested struct. Referenced from `generate` or other `object` blocks via `type = "object"` and `object = "<name>"`.

```hcl
object "database" {
  field "dsn" {
    type     = "string"
    required = true
  }
  field "pool_size" {
    type    = "int"
    default = "10"
  }
}
```

#### Slice of objects

To define a field as a slice of nested structs:

```hcl
object "replica" {
  field "host" { type = "string" }
  field "port" { type = "uint16"; default = "5432" }
}

generate "db_config" {
  field "replicas" {
    type   = "list object"
    object = "replica"
  }
}
```

Result: `Replicas []Replica`

## Runtime API

### `confgen.Unmarshal`

```go
func Unmarshal(data []byte, v any, unmarshaler Unmarshaler) error
```

1. Decodes `data` into `v` using the provided `unmarshaler`
2. Calls `Defaults()` if `v` implements `Defaultable` — fills zero fields with defaults
3. Checks `confgen_required:"true"` struct tags — returns an error if the field value is zero

### `confgen.NewDecoder`

```go
decoder := confgen.NewDecoder("config.yaml", yaml.Unmarshal)
err := decoder.Decode(&cfg)
```

A convenience wrapper around `Unmarshal` that reads the file itself.

### `Defaultable` interface

Generated structs automatically implement this interface:

```go
type Defaultable interface {
    Defaults() any
}
```

You can also implement it manually for custom structs — just add a `Defaults() any` method.

## Programmatic API

Generation can be invoked directly from Go code using the `generator` package:

```go
import (
    "github.com/brezzgg/go-packages/confgen/internal/generator"
    "github.com/brezzgg/go-packages/confgen/internal/parser"
)

schemas, _ := parser.Parse([]string{"schema.hcl"}, false, ".+")
for _, s := range schemas {
    code, err := generator.Generate(s,
        generator.WithCustomTemplate(myTemplate),
        generator.WithTemplateFunctions(map[string]any{
            "MyHelper": myHelperFunc,
        }),
    )
}
```

### Generator Options

| Option                                        | Description                                                         |
| --------------------------------------------- | ------------------------------------------------------------------- |
| `WithCustomTemplate(tmpl string)`             | Replaces the built-in template                                      |
| `WithTemplate(tmplName string)`               | Pick built-in template                                              |
| `WithTemplateFunctions(funcs map[string]any)` | Adds custom functions to the template (does not override built-ins) |
| `WithTemplateData(data map[string]any)`       | Adds custom data to the template context                            |

### Built-in Template Functions

| Function    | Description                                              |
| ----------- | -------------------------------------------------------- |
| `ToPascal`  | `my_field` → `MyField`                                   |
| `ToCamel`   | `my-field` → `MyField`                                   |
| `ToSnake`   | `MyField` → `my_field`                                   |
| `ToType`    | `schema.Field` → Go type string                          |
| `ToDefault` | `schema.Field` → Go default expression                   |
| `ToBool`    | Converts `*bool`, `bool`, or `string` → `bool`           |
| `ToTag`     | Generates struct tags from formats and the required flag |
| `Deref`     | Dereference a pointer                                    |

## Examples

More examples are available in [examples/confgen](https://github.com/brezzgg/go-packages/blob/main/examples/confgen).
