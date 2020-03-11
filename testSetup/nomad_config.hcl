# Increase log verbosity
log_level = "DEBUG"
bind_addr= "0.0.0.0"

# Setup data dir
data_dir = "/tmp/nomadJobUpdate/nomadData/"

# Enable the server
server {
  enabled = true

  bootstrap_expect = 1
}

client {
  enabled       = true
}

consul {
  address = "127.0.0.1:8500"
}

tls {
  http = true
  rpc  = true

  ca_file   = "certificates/nomad-ca.pem"
  cert_file = "certificates/server.pem"
  key_file  = "certificates/server-key.pem"

  verify_server_hostname = true
  verify_https_client    = true
}