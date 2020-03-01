FROM 1.13-alpine AS dependencies 
WORKDIR /app
RUN go env -w GO111MODULE="on"
COPY go.sum go.mod ./
RUN go mod tidy

FROM dependencies AS build 
WORKDIR /app
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bitly main.go 
EXPOSE 8001
ENTRYPOINT [ "bitly", "-token", "" ]
