FROM golang:latest AS builder

WORKDIR /go/src/github.com/luuisavelino/short-circuit-analysis-zbus/

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

COPY ./main.go .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o main .


FROM alpine:latest

WORKDIR /root/

COPY --from=0 /go/src/github.com/luuisavelino/short-circuit-analysis-zbus/main ./

EXPOSE 8080

CMD ["./main"]