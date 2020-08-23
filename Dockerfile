FROM ubuntu:18.04

WORKDIR /build

RUN apt-get update 

RUN apt-get install -y software-properties-common

RUN add-apt-repository ppa:longsleep/golang-backports 

RUN apt-get update 

RUN apt-get install -y xen-system-amd64

RUN apt install -y golang-go

RUN apt-get install -y libxen-dev xen-tools

RUN apt-get install -y xenwatch

CMD go build duster.go 
