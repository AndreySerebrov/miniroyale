FROM  golang:1.19 as builder

WORKDIR /go/src
COPY . .
RUN  go mod tidy
RUN  go build  -o bin/app ./main.go

FROM  golang:1.19
WORKDIR /application
COPY --from=builder /go/src/bin/app /application
RUN ls -l
COPY --from=builder /go/src/assets/Words-of-Wisdom.txt /application/assets/Words-of-Wisdom.txt
RUN ls -l

RUN  chmod 777 /application/app

CMD ["/application/app"]
