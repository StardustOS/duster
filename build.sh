if ! [ "$(command -v xen)" ]; then 
    apt-get install xen-hypervisor-4.9-amd64 xen-system-amd64 xen-tools xen-utils-4.9 xen-utils-common xenstore-utils xenwatch
fi
if ! [ "$(command -v go)" ]; then
    add-apt-repository ppa:longsleep/golang-backports
    apt update
    apt install golang-go
fi
go build -o bin/duster duster.go