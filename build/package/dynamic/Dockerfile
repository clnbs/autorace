FROM golang:1.14 AS builder
WORKDIR /go/src/github.com/clnbs/autorace
RUN apt update && apt upgrade -y && apt-get install libgl1-mesa-dev libxrandr-dev libxcursor-dev libxi-dev libghc-x11-dev binutils-mingw-w64 mingw-w64 libx11-dev xorg-dev -y
COPY . .
RUN go get -u ./...
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/dynamic/main.go

FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/clnbs/autorace/app .
ENTRYPOINT ["./app"]