FROM golang:1.23.6-alpine3.21 as build
WORKDIR /llcm
COPY go.mod go.sum ./
COPY main.go .
RUN go build -o main main.go
FROM alpine:3.21
COPY --from=build /llcm/main /main
ENTRYPOINT [ "/main" ]
