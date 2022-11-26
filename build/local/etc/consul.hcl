// data_dir = "./data/consul"
node_name = "consul-server"
server    = true
bootstrap = true

ui_config {
  enabled = true
}

datacenter = "dc1"
// data_dir   = "consul/data"
log_level  = "INFO"

addresses {
  http = "0.0.0.0"
}

connect {
  enabled = true
}

// ports {
//   grpc = 8502
// }