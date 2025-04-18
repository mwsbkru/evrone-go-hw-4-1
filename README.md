```shell
cd cmd/app
INPUT=./../../urls.txt OUTPUT=./../../result.csv PARALLELS=3 go run .
go test ./... -coverprofile cover.out&& go tool cover -html=cover.out
```