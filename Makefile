build-386:
	export GOOS=linux; \
	export GOARCH=386; \
	CGO_ENABLED=1 go build -o "./bin/${GOOS}_${GOARCH}/awake" cmd/server/*
