package main

//go:generate protoc --go-plugin_out=. --go-plugin_opt=paths=source_relative proto/opq.proto
//go:generate go run proto/generate/generate.go
//go:generate protoc --go-plugin_out=. --go-plugin_opt=paths=source_relative proto/library/systemInfo/export/systemInfo.proto
