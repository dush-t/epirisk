# Using ubuntu as base because I need to setup some C dependencies
# that I couldn't setup with the golang base.
FROM ubuntu:18.04 as builder

# Setting up dependencies
RUN apt-get update
RUN apt-get install -y libssl1.0.0 wget apt-utils lsb-release curl

RUN curl https://pkgconfig.freedesktop.org/releases/pkg-config-0.29.tar.gz | tar -v -C /usr/local -xz
RUN /usr/local/pkg-config-0.29/configure && make install
 
RUN wget https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-$(lsb_release -rs).deb
RUN dpkg -i seabolt-1.7.4-Linux-ubuntu-$(lsb_release -rs).deb
RUN rm seabolt-1.7.4-Linux-ubuntu-$(lsb_release -rs).deb

# Installing GoLang 1.14.1
ENV GOLANG_VERSION 1.14.1

RUN curl -sSL https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV PATH /usr/local/go/bin:$PATH

RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

# Set current working directory inside the container
WORKDIR /app

# Download all dependencies. Doing this before the rest of the code so that dependencies
# are cached by docker if there is no change in these two files, irrespective of changes in code
COPY go.mod go.sum ./
RUN go mod download


# Copying source code
COPY . .

# Building
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]