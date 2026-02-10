## Build and run

```bash
GOOS=linux GOARCH=amd64 go build -o bot cmd/bot/main.go
```

Transfer to server using scp

```
scp -i key bot Dockerfile docker-compose.yml .env user@177.77.77.777:/root/spends
```

Build


```bash
docker-compose build
docker-compose up -d
```


## Check
```bash
docker-compose logs
docker-compose logs -f   
```

## Stop
```bash
docker-compose down
```
