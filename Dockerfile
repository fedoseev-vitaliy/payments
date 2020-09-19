ARG GOLANG_IMG=golang:1.15
ARG APLINE_IMG=alpine:3.10.3

FROM ${GOLANG_IMG} as builder

WORKDIR $GOPATH/src/github.com/fedoseev-vitaliy/payments

COPY cmd cmd
COPY tools tools
COPY vendor vendor
COPY internal internal
COPY main.go .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -a -ldflags="-w -s" -o /usr/local/bin/payments .


FROM ${APLINE_IMG} as final
COPY --from=builder /usr/local/bin/payments /usr/local/bin/payments

ENV BIND 0.0.0.0:80

EXPOSE 80

ENTRYPOINT ["payments"]
