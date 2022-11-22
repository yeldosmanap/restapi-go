FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

COPY ./ /github.com/yeldosmanap/restapi-go
WORKDIR /github.com/yeldosmanap/restapi-go

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/api/main.go
# ./.bin/api ./cmd/api/main.go
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /github.com/yeldosmanap/restapi-go/.bin/app .
COPY --from=builder /github.com/yeldosmanap/restapi-go/.env ./.env
COPY --from=builder /github.com/yeldosmanap/restapi-go/configs ./configs/
COPY --from=builder /github.com/yeldosmanap/restapi-go/prometheus ./prometheus/

EXPOSE 8080

CMD ["./app"]