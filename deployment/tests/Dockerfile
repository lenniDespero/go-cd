FROM golang:1.14
WORKDIR /app
## Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./
## Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
## Copy the source from the current directory to the Working Directory inside the container
COPY . .
#RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
#Build binary for test
RUN CGO_ENABLED=0 go build -o ./bin/deployer ./cmd/main.go
WORKDIR /app/tests
#RUN go test -v
CMD ["go", "test", "-v", "--tags=integration"]