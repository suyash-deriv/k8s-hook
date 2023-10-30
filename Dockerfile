FROM golang:1.21 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /k8s-hook

FROM gcr.io/distroless/base-debian12 AS build-release-stage
WORKDIR /
COPY --from=build-stage /k8s-hook /k8s-hook
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/k8s-hook"]