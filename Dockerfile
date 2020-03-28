FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y libssl1.0.0 wget apt-utils lsb-release curl
 
RUN wget https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-$(lsb_release -rs).deb
RUN dpkg -i seabolt-1.7.4-Linux-ubuntu-$(lsb_release -rs).deb
RUN rm seabolt-1.7.4-Linux-ubuntu-$(lsb_release -rs).deb

ENV GOLANG_VERSION 1.14.1

RUN curl -sSL https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV PATH /usr/local/go/bin:$PATH

RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -0 main .