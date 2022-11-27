job "product-api" {
  datacenters = ["dc1"]

  constraint {
    attribute = "${attr.kernel.name}"
    value     = "linux"
  }

  group "svc" {
    count = 1

    network {
      mode = "bridge"

      port "grpc" {
        to = 5001
      }
    }

    service {
      name = "product-api-grpc"
      port = "5001"

      connect {
        sidecar_service {}
      }
    }

    task "product-api" {
      driver = "raw_exec"

      artifact {
        source      = "git::https://github.com/thangchung/go-coffeeshop"
        destination = "local/repo"
      }

      config {
        command = "bash"
        args = [
          "-c",
          "cd local/repo/cmd/product && go mod tidy && go mod download && CGO_ENABLED=0 go run github.com/thangchung/go-coffeeshop/cmd/product"
        ]
      }

      env {
        APP_NAME = "product-service in docker"
      }

      resources {
        cpu    = 100
        memory = 128
      }
    }
  }
}