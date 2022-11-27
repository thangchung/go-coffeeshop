job "grpc-gw" {
  datacenters = ["dc1"]

  constraint {
    attribute = "${attr.kernel.name}"
    value     = "linux"
  }

  group "svc" {
    count = 1

    network {
      mode = "bridge"

      port "http" {
        to = 5000
      }
    }

    service {
      name = "grpc-gw-http"
      port = "5000"

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "product-api-grpc"
              local_bind_port  = 5001
            }
            upstreams {
              destination_name = "counter-api-grpc"
              local_bind_port  = 5002
            }
          }
        }
      }

      tags = [
        "traefik.enable=true",
        "traefik.consulcatalog.connect=true",
        "traefik.port=5000",
        "traefik.http.routers.api.entryPoints=web",
        "traefik.http.routers.api.rule=PathPrefix(`/api`)",
        "traefik.http.routers.api.middlewares=api-stripprefix",
        "traefik.http.middlewares.api-stripprefix.stripprefix.prefixes=/api",
      ]
    }

    task "grpc-gw" {
      driver = "raw_exec"

      artifact {
        source = "git::https://github.com/thangchung/go-coffeeshop"
        destination = "local/repo"
      }

      config {
        command = "bash"
        args = [
          "-c",
          "cd local/repo/cmd/proxy && go mod tidy && go mod download && CGO_ENABLED=0 go run github.com/thangchung/go-coffeeshop/cmd/proxy"
        ]
      }

      env {
        APP_NAME          = "proxy-service in docker"
        GRPC_PRODUCT_HOST = "${NOMAD_UPSTREAM_IP_product_api_grpc}"
        GRPC_PRODUCT_PORT = "${NOMAD_UPSTREAM_PORT_product_api_grpc}"
        GRPC_COUNTER_HOST = "${NOMAD_UPSTREAM_IP_counter_api_grpc}"
        GRPC_COUNTER_PORT = "${NOMAD_UPSTREAM_PORT_counter_api_grpc}"
      }

      resources {
        cpu    = 100
        memory = 128
      }
    }
  }
}