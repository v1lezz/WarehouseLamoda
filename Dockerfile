FROM golang:1.22
RUN mkdir /server
ADD . /warehouse/
WORKDIR /warehouse
RUN go build -o srv ./cmd/warehouse
CMD ["/warehouse/srv"]