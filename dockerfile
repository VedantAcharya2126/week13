FROM golang:1.23 AS builder
WORKDIR /app
COPY  go.mod go.sum ./ 
RUN go mod download 
COPY  . .
# RUN go build -o main .
RUN CSG_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM alpine:3.20
# RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder  /app/app ./app
COPY --from=builder  /app/templates ./templates
EXPOSE 80
CMD [ "./app" ]