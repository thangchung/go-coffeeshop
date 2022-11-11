# go-coffeeshop

The coffeeshop with golang stack

## Services

No. | Service | URI
--- | --- | ---
1 | grpc-gateway | [http://localhost:5000](http://localhost:5000)
2 | product service | [http://localhost:5001](http://localhost:5001)
3 | counter service | [http://localhost:5002](http://localhost:5002)
4 | barista service | [http://localhost:5003](http://localhost:5003)
5 | kitchen service | [http://localhost:5004](http://localhost:5004)

## Package

```go
go get -u github.com/ilyakaznacheev/cleanenv
```

## Debug Apps

[Debug golang app in monorepo](https://github.com/thangchung/go-coffeeshop/wiki/Golang#debug-app-in-monorepo)

## Connect to app running on host machine in `docker-from-docker` feature

From host machine, Ubuntu 22.04 at here, type:

```bash
> ifconfig
```

Then find `docker0` IP address in the output, just like

```bash
docker0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
    inet 172.17.0.1  netmask 255.255.0.0  broadcast 172.17.255.255
    inet6 fe80::42:6dff:fe39:ba02  prefixlen 64  scopeid 0x20<link>
    ether 02:42:6d:39:ba:02  txqueuelen 0  (Ethernet)
    RX packets 994  bytes 122399 (122.3 KB)
    RX errors 0  dropped 0  overruns 0  frame 0
    TX packets 759  bytes 772987 (772.9 KB)
    TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
```

In here, we can find `172.17.0.1`, then use it for all accesses in `docker-from-docker`

## Generate Protobuf/gRPC

```bash
> buf mod update proto
> buf generate
```

## Env

Create `.env` at your root with content as

```bash
PG_URL=postgres://postgres:P@ssw0rd@<your devcontainer IP>:5432/postgres
```

## Migrations

Add migration for `counter-api`

```bash
> migrate create -seq -dir cmd/counter/db/migrations -ext sql init_db
```

## Sample

``` mermaid
graph TB
  S[myproject]
  M[email service]
  S -.->|sends emails| M
  A[Anonymous User<br/>Person] -.->|Can register a new organisation| S
  B[System Admin<br/>Person] -.->|Manages organisations and accounts| S
  C[Organisation Admin<br/>Person] -.->|Manages organisation, can do X| S
  classDef person fill:#08427b,color:#fff,stroke:none;
  class A,B,C person;
  classDef currentSystem fill:#1168bd,color:#fff,stroke:none;
  class S currentSystem;
  classDef otherSystem fill:#999999,color:#fff,stroke:none;
  class M otherSystem;
```

## Project structure

- [project-layout](https://github.com/golang-standards/project-layout)
- [repository-structure](https://peter.bourgon.org/go-best-practices-2016/#repository-structure)
- [go-build-template](https://github.com/thockin/go-build-template)
- [go-wiki](https://github.com/golang/go/wiki/Articles#general)
