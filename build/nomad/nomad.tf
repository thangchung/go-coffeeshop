provider "nomad" {
  address = "http://localhost:4646"
  version = "~> 1.4"
}

resource "nomad_job" "traefik" {
  jobspec = file("${path.module}/jobs/traefik.nomad.hcl")
}

resource "nomad_job" "postgres_db" {
  jobspec    = file("${path.module}/jobs/postgresdb.nomad.hcl")
  depends_on = [nomad_job.traefik]
}

resource "nomad_job" "rabbitmq" {
  jobspec    = file("${path.module}/jobs/rabbitmq.nomad.hcl")
  depends_on = [nomad_job.traefik]
}

resource "nomad_job" "product_api" {
  jobspec    = file("${path.module}/jobs/product-api.nomad.hcl")
  depends_on = [nomad_job.traefik]
}

resource "nomad_job" "counter_api" {
  jobspec    = file("${path.module}/jobs/counter-api.nomad.hcl")
  depends_on = [nomad_job.postgres_db, nomad_job.rabbitmq, nomad_job.product_api]
}

resource "nomad_job" "barista_api" {
  jobspec    = file("${path.module}/jobs/barista-api.nomad.hcl")
  depends_on = [nomad_job.postgres_db, nomad_job.rabbitmq]
}

resource "nomad_job" "kitchen_api" {
  jobspec    = file("${path.module}/jobs/kitchen-api.nomad.hcl")
  depends_on = [nomad_job.postgres_db, nomad_job.rabbitmq]
}

resource "nomad_job" "grpc_gw" {
  jobspec    = file("${path.module}/jobs/grpc-gw.nomad.hcl")
  depends_on = [nomad_job.product_api, nomad_job.counter_api, nomad_job.barista_api, nomad_job.kitchen_api]
}

resource "nomad_job" "web" {
  jobspec    = file("${path.module}/jobs/web.nomad.hcl")
  depends_on = [nomad_job.grpc_gw]
}
