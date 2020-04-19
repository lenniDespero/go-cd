FROM golang as builder
WORKDIR /app
## Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./
## Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
## Copy the source from the current directory to the Working Directory inside the container
COPY . .
## Build the Go app
RUN GCO_ENABLED=0 GOOS=linux go build -o deployer ./cmd/main.go
#ENTRYPOINT ./deployer

### Get small image
FROM alpine
WORKDIR /opt/api
RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
## Copy build from sourse
COPY --from=builder /app/deployer .
## Add the wait script to the image
CMD ./deployer