output {
  package = "config"
  output  = "examples/confgen/config"
  formats = ["yaml", "json"]
}

object "tls" {
  field "cert_file" {
    type     = "string"
    required = true
    desc     = "example cert.pem"
  }
  field "key_file" {
    type     = "string"
    required = true
    desc     = "example key.pem"
  }
}

object "pool" {
  field "max_open" {
    type    = "int"
    default = "10"
  }
  field "max_idle" {
    type    = "int"
    default = "5"
  }
  field "max_lifetime" {
    type    = "int64"
    default = "300"
  }
}

object "postgres" {
  field "host" {
    type    = "string"
    default = "localhost"
  }
  field "port" {
    type    = "int"
    default = "5432"
  }
  field "user" {
    type     = "string"
    required = true
  }
  field "password" {
    type     = "string"
    required = true
  }
  field "dbname" {
    type    = "string"
    default = "app"
  }
  field "pool" {
    type   = "object"
    object = "pool"
  }
  field "tls" {
    type   = "object"
    object = "tls"
  }
}

object "redis" {
  field "host" {
    type    = "string"
    default = "localhost"
  }
  field "port" {
    type    = "int"
    default = "6379"
  }
  field "password" {
    type = "string"
  }
  field "db" {
    type    = "int"
    default = "0"
  }
  field "pool" {
    type   = "object"
    object = "pool"
  }
}

object "server" {
  field "host" {
    type    = "string"
    default = "0.0.0.0"
    desc    = "listen host"
  }
  field "admins" {
    type = "list"
    object = "server_admins"
  }
  field "port" {
    type    = "int"
    default = "8080"
  }
  field "read_timeout" {
    type    = "int64"
    default = "30"
  }
  field "write_timeout" {
    type    = "int64"
    default = "30"
  }
  field "tls" {
    type   = "object"
    object = "tls"
  }
}

object "kafka_broker" {
  field "address" {
    type    = "string"
    default = "localhost:9092"
    desc    = "remote host of kafka broker instance"
  }
}

object "server_admins" {
  field "token" {
    required = true
    desc     = "base64 128 symbols"
  }
  field "name" {
    required = true
    desc     = ">= 3 symbols"
  }
}

object "kafka" {
  field "brokers" {
    type   = "list"
    object = "kafka_broker"
  }
  field "topic" {
    type    = "string"
    default = "events"
  }
  field "group_id" {
    type    = "string"
    default = "app-group"
  }
  field "max_retries" {
    type    = "int"
    default = "3"
  }
}

object "log" {
  field "level" {
    type    = "string"
    default = "info"
  }
  field "format" {
    type    = "string"
    default = "json"
  }
  field "output" {
    type     = "string"
    required = true
  }
}

generate "app" {
  field "env" {
    type    = "string"
    default = "production"
  }
  field "debug" {
    type    = "bool"
    default = "false"
  }
  field "server" {
    type   = "object"
    object = "server"
  }
  field "postgres" {
    type   = "object"
    object = "postgres"
  }
  field "redis" {
    type   = "object"
    object = "redis"
  }
  field "kafka" {
    type   = "object"
    object = "kafka"
  }
  field "log" {
    type   = "object"
    object = "log"
  }
}
