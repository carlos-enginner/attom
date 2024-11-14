FROM golang:1.23

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the workspace
COPY . .

RUN go install github.com/air-verse/air@latest && \
    CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

RUN echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | tee /etc/apt/sources.list.d/goreleaser.list && apt update && apt install goreleaser

EXPOSE 8080
EXPOSE 2345

# Command to run Delve server
#CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient"]
ENTRYPOINT ["air"]