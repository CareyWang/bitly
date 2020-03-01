FROM golang:1.13-alpine AS dependencies 
WORKDIR /app
RUN go env -w GO111MODULE="on"
COPY go.sum go.mod ./
RUN go mod tidy 

FROM dependencies as build
WORKDIR /app
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bitly main.go 

FROM scratch
WORKDIR /app
COPY --from=build /app/bitly ./
EXPOSE 8001
ENTRYPOINT ["/app/bitly"]
