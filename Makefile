build-linux-386:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -o ./bin/linux_386/awake cmd/server/*
