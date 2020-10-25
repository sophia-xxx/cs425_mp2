#! /bin/bash

if [ ! -d "go" ]; then
	mkdir -p /home/$USER/go/src
fi

if [ ! -d "src" ]; then
	mkdir -p /home/$USER/go/pkg
fi

export GOROOT=/usr/lib/golang
export GOPATH=/home/$USER/go

go get github.com/c-bata/go-prompt
go get github.com/jinzhu/copier
go get github.com/golang/protobuf/proto
go get github.com/golang/protobuf/ptypes/timestamp
go get google.golang.org/protobuf/reflect/protoreflect
go get google.golang.org/protobuf/runtime/protoimpl
go get github.com/gogf/greuse