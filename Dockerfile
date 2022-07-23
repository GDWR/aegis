FROM golang:1.18 as builder

WORKDIR /build

ADD . .
RUN go build -o /output/aegis

FROM alpine:3.16

RUN apk add --no-cache libc6-compat

WORKDIR /aegis
COPY --from=builder /output .

EXPOSE 8765
ENTRYPOINT ["./aegis", "--port=25565"]
