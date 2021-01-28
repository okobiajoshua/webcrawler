# build stage
FROM golang as builder

# librdkafka Build from source
RUN git clone https://github.com/edenhill/librdkafka.git

WORKDIR librdkafka

RUN ./configure --prefix /usr

RUN make

RUN make install

# Build go binary

WORKDIR /app/
COPY ./ webcrawler

ENV GO111MODULE=on
WORKDIR /app/webcrawler
RUN go mod download


RUN go build -o main
RUN ls
# final stage
FROM ubuntu

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/lib/pkgconfig /usr/lib/pkgconfig
COPY --from=builder /usr/lib/librdkafka* /usr/lib/
COPY --from=builder /app/webcrawler/* /webcrawler/
WORKDIR /webcrawler
CMD ["./main"]