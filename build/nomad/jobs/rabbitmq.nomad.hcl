job "rabbitmq" {
  datacenters = ["dc1"]

  group "rabbitmq" {
    network {
      mode = "bridge"

      port "broker" {
        static = 5672
      }

      port "broker_ui" {
        static = 15672
      }
    }

    service {
      name = "rabbitmq"
      port = "broker"

      connect {
        sidecar_service { }
      }
    }

    task "rabbitmq" {
      driver = "docker"

      config {
        image = "rabbitmq:3.11-management-alpine"
        ports = ["broker", "broker_ui"]
      }

      env {}
    }
  }
}