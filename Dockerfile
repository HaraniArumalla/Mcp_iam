# Use the official Golang image to build the app

FROM registry.gitlab.com/tmobile/citadel/containers/tmo-go-base:latest AS certs

FROM golang:latest as builder

RUN apt-get update && apt-get install -y ca-certificates openssl

ARG cert_location=/usr/local/share/ca-certificates


# Get certificate from "github.com"
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
# Get certificate from "proxy.golang.org"
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt
# Update certificates
RUN update-ca-certificates

# # # Set the working directory
WORKDIR /app

#RUN go env -w GOPROXY=direct GOFLAGS="-insecure" GOINSECURE="go.googlesource.com/*,github.com/*,golang.org/*"

# # Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Copy the source code
COPY . .

# Download dependencies
RUN go mod download

# # # Generate Swagger docs
# RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init

# Build the Go app
RUN go run -mod=mod github.com/99designs/gqlgen generate
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server/main.go

FROM registry.gitlab.com/tmobile/citadel/containers/tmo-go-base:latest
#FROM ubuntu:latest

ARG USERNAME=nonroot
ARG USER_UID=1000
ARG GROUP_NAME=nonroot
ARG USER_GID=1000

# RUN whoami
# # Set the working directory
WORKDIR /home/nonroot

# # Copy the built application from the builder
COPY --from=builder --chown=$USERNAME:$GROUP_NAME /app/app /home/$USERNAME/

RUN chmod +x /home/$USERNAME/app

# Switch to nonroot user
USER $USERNAME

ENV PORT=8080

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./app"]
