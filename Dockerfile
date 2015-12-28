# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.5

RUN go get github.com/Masterminds/glide

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/ssoudan/frontend

# Build the command inside the container.
ENV GO15VENDOREXPERIMENT 1
WORKDIR /go/src/github.com/ssoudan/frontend
RUN glide install
RUN go install github.com/ssoudan/frontend/cmd/frontend

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/frontend -logtostderr -backend http://$BACKEND_SERVICE_HOST:$BACKEND_SERVICE_PORT

# Document that the service listens on port 8080.
EXPOSE 8080
