FROM golang:alpine3.10 as builder
WORKDIR /go/src/github.com/alekssaul/kube-tags2iaas
COPY . .
RUN mkdir -p /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/kube-tags2iaas .


FROM alpine:latest
RUN apk update ;  apk add --no-cache ca-certificates ; update-ca-certificates ; mkdir /app
WORKDIR /app
COPY --from=builder /app .
CMD /app/kube-tags2iaas
EXPOSE 8080