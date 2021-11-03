# protoc-gen-foo

Writing a protoc plugin

# Note

Your compiler name must be prefixed with protoc-gen- e.g. proto-gen-mytestprotoplugin

mkdir out /
go install . /
protoc --proto_path . -I=. mytest.proto --mytestprotoplugin_out=./out --go_out=./out /
