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


### Как добавлять категории трат?

В файле internal/bot/keyboards.go 
1) 
```go
btnNewCategory := categoryMenu.Text("New") 
```

2)
```go
categoryMenu.Reply(
	categoryMenu.Row(btnSnacks, btnProducts),
	categoryMeny.Row(btnNewCategori), // Добавлили сюда кнопку
```

3) 
```go
func GetCategoryHandlers() map[string]string {
	return map[string]string{
		"🍿 Снеки":       "снеки",
		"🛒 Продукты":    "продукты",
		"Новая кнопка": "Новая категория",
	}
```
