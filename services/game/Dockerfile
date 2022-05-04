FROM golang:1.17 as builder
WORKDIR /go/src/github.com/zuri03/game/

COPY . .

RUN CGO_ENABLED=0 go build -o  build/Game ./main.go
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/zuri03/TriviaApp/
COPY --from=builder /go/src/github.com/zuri03/game/build/Game ./main
EXPOSE 8080
ENTRYPOINT [ "./main" ]