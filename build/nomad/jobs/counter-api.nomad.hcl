job "counter-api" {
  datacenters = ["dc1"]

  constraint {
    attribute = "${attr.kernel.name}"
    value     = "linux"
  }

  group "counter-api" {
    count = 1

    network {
      mode = "bridge"

      port "grpc" {
        to = 5002
      }
    }

    service {
      name = "counter-api-grpc"
      port = "5002"

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "postgres-db"
              local_bind_port  = 5432
            }
            upstreams {
              destination_name = "rabbitmq"
              local_bind_port  = 5672
            }
            upstreams {
              destination_name = "product-api-grpc"
              local_bind_port  = 5001
            }
          }
        }
      }
    }

    task "counter-api" {
      driver = "raw_exec"

      artifact {
        source      = "git::https://github.com/thangchung/go-coffeeshop"
        destination = "local/repo"
      }

      config {
        command = "bash"
        args = [
          "-c",
          "cd local/repo/cmd/counter && go mod tidy && go mod download && CGO_ENABLED=0 go run -tags migrate github.com/thangchung/go-coffeeshop/cmd/counter"
        ]
      }

      env {
        APP_NAME           = "counter-service in docker"
        IN_DOCKER          = "false"
        PG_URL             = "postgres://postgres:P@ssw0rd@${attr.unique.network.ip-address}:5432/postgres"
        RABBITMQ_URL       = "amqp://guest:guest@${attr.unique.network.ip-address}:5672/"
        PRODUCT_CLIENT_URL = "${NOMAD_UPSTREAM_ADDR_product_api_grpc}"
      }

      resources {
        cpu    = 150
        memory = 200
      }
    }
  }
}