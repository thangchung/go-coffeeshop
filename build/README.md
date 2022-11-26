# Get starting with Nomad, Consult Connect

## Start Nomad, Consul and Vault

```bash
> cd build/local
> ./start.sh
```

> Make sure you set start.sh with execute permission => `sudo chmod +x start.sh`

## Use Terraform to provisioning all services

```bash
> cd build/nomad
> terraform init
> terraform apply
```

## Clean Up

```bash
> cd build/nomad
> terraform destroy
> cd build/local
# Ctrl + C
```

Happy hacking with HashiCorp stack!!!
