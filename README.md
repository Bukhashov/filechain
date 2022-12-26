# filechain

// https://habr.com/ru/post/490336/

grpc
python3 -m grpc_tools.protoc -I ./protos --python_out=./face --pyi_out=./face --grpc_python_out=./face face.proto

protoc -I ./protos --go_out=plugins=grpc:server/pkg/api protos/face.proto 