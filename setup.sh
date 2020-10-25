#! /bin/bash

if [ ! -d "go" ]; then
	mkdir -p /home/chunhao3/go/src
fi

if [ ! -d "src" ]; then
	mkdir -p /home/chunhao3/go/pkg
fi

export GOROOT=/usr/lib/golang
export GOPATH=/home/chunhao3/go

go get github.com/jinzhu/copier
go get github.com/golang/protobuf/proto
go get github.com/golang/protobuf/ptypes/timestamp
go get google.golang.org/protobuf/reflect/protoreflect
go get google.golang.org/protobuf/runtime/protoimpl
go get github.com/gogf/greuse
