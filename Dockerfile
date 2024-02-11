FROM golang:latest AS build

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go build -o /server

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /server /server

EXPOSE 8080

USER nonroot:noonroot

ENTRYPOINT [ "/server" ]


