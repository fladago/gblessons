# Version: 0.0.1
FROM golang:1.17.1-alpine3.14 AS build
WORKDIR /src/
COPY go.* /src/
RUN go mod download
COPY ./ /src/
RUN CGO_ENABLED=0 go build -o /bin/demo
FROM scratch
COPY --from=build /bin/demo /bin/demo
ENTRYPOINT [ "/bin/demo" ]