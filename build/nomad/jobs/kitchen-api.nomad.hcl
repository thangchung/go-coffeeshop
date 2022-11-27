job "kitchen-api" {
  datacenters = ["dc1"]

  constraint {
    attribute = "${attr.kernel.name}"
    value     = "linux"
  }

  group "kitchen-api" {
    count = 1
    
    network {
      mode = "bridge"

      port "http" { 
        to = 5004 
      }
    }

    service {
      name = "kitchen-api"
      port = "5004"
    }

    task "kitchen-api" {
      driver = "raw_exec"

      artifact {
        source      = "git::https://github.com/thangchung/go-coffeeshop"
        destination = "local/repo"
      }

      config {
        command = "bash"
        args = [
          "-c",
          "cd local/repo/cmd/kitchen && go mod tidy && go mod download && CGO_ENABLED=0 go run -tags migrate github.com/thangchung/go-coffeeshop/cmd/kitchen"
        ]
      }

      env {
        APP_NAME     = "kitchen-service in docker"
        IN_DOCKER    = "false"
        PG_URL       = "postgres://postgres:P@ssw0rd@${attr.unique.network.ip-address}:5432/postgres"
        RABBITMQ_URL = "amqp://guest:guest@${attr.unique.network.ip-address}:5672/"
      }

      resources {
        cpu    = 100
        memory = 200
      }
    }
  }
}