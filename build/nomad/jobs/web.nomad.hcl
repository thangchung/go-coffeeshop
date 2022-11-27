job "web" {
  datacenters = ["dc1"]

  constraint {
    attribute = "${attr.kernel.name}"
    value     = "linux"
  }

  group "web" {
    count = 1

    network {
      mode = "bridge"

      port "http" {
        to = 8888
      }
    }

    service {
      name = "web-http"
      port = "8888"

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "grpc-gw-http"
              local_bind_port  = 5000
            }
          }
        }
      }

      tags = [
        "traefik.enable=true",
        "traefik.consulcatalog.connect=true",
        "traefik.port=8888",
        "traefik.http.routers.web.entryPoints=web",
        "traefik.http.routers.web.rule=PathPrefix(`/`)",
      ]
    }

    task "web" {
      driver = "raw_exec"

      artifact {
        source      = "git::https://github.com/thangchung/go-coffeeshop"
        destination = "local/repo"
      }

      config {
        command = "bash"
        args = [
          "-c",
          "cd local/repo/cmd/web && go mod tidy && go mod download && CGO_ENABLED=0 go run github.com/thangchung/go-coffeeshop/cmd/web"
        ]
      }

      env {
        REVERSE_PROXY_URL = "http://localhost/api"
        WEB_PORT          = 8888
      }

      resources {
        cpu    = 100
        memory = 128
      }
    }
  }
}