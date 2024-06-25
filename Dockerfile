FROM golang:1.22 AS builder

WORKDIR /go/src/pwd

COPY go.mod go.sum ./
COPY . .
RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -o ./pwdapp .

FROM scratch

COPY --from=builder /go/src/pwd/pwdapp /pwdapp

EXPOSE 4000

CMD ["/pwdapp"]
