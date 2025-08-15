FROM --platform=arm64 golang:1.24.6-alpine AS build
WORKDIR /build
COPY go.mod go.sum ./
COPY main.go .
RUN go mod download && go build -trimpath -o main main.go

FROM --platform=arm64 alpine
COPY --from=build /build/main /var/task/main
ENTRYPOINT [ "/var/task/main" ]
