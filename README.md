# Aegis

Proxy in a can with an Api to control a ban list.

--- 

### Compile and run
```shell
git clone <repository-url>.git
cd ./aegis
go build
./aegis
```

### Run with docker

#### From container registry
```shell
docker run -p 8765:8765 -p 25565:25565 ghcr.io/gdwr/aegis:main --destination=127.0.0.1
```

#### Build from source
```shell
docker build . -t aegis
docker run -p 8765:8765 -p 25565:25565 aegis --destination=127.0.0.1
```


### API
Runs on `:8765` by default.

#### Get ban list
```shell
curl -X GET 127.0.0.1:8765
```

#### Add to ban list
```shell
curl -X POST -d '"192.168.1.1"' 127.0.0.1:8765
```

#### Delete from the ban list
```shell
curl -X DELETE -d '"192.168.1.1"' 127.0.0.1:8765
```
