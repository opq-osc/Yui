//go:build tinygo.wasm

// Code generated by protoc-gen-go-plugin. DO NOT EDIT.
// versions:
// 	protoc-gen-go-plugin v0.1.0
// 	protoc               v3.19.4
// source: proto/library/systemInfo/export/systemInfo.proto

package export

import (
	context "context"
	emptypb "github.com/knqyf263/go-plugin/types/known/emptypb"
	wasm "github.com/knqyf263/go-plugin/wasm"
	_ "unsafe"
)

type systemInfo struct{}

func NewSystemInfo() SystemInfo {
	return systemInfo{}
}

//go:wasm-module systemInfo
//export cpu_info
//go:linkname _cpu_info
func _cpu_info(ptr uint32, size uint32) uint64

func (h systemInfo) CpuInfo(ctx context.Context, request *emptypb.Empty) (*SystemInfoReply, error) {
	buf, err := request.MarshalVT()
	if err != nil {
		return nil, err
	}
	ptr, size := wasm.ByteToPtr(buf)
	ptrSize := _cpu_info(ptr, size)
	wasm.FreePtr(ptr)

	ptr = uint32(ptrSize >> 32)
	size = uint32(ptrSize)
	buf = wasm.PtrToByte(ptr, size)

	response := new(SystemInfoReply)
	if err = response.UnmarshalVT(buf); err != nil {
		return nil, err
	}
	return response, nil
}

//go:wasm-module systemInfo
//export mem_info
//go:linkname _mem_info
func _mem_info(ptr uint32, size uint32) uint64

func (h systemInfo) MemInfo(ctx context.Context, request *emptypb.Empty) (*SystemInfoReply, error) {
	buf, err := request.MarshalVT()
	if err != nil {
		return nil, err
	}
	ptr, size := wasm.ByteToPtr(buf)
	ptrSize := _mem_info(ptr, size)
	wasm.FreePtr(ptr)

	ptr = uint32(ptrSize >> 32)
	size = uint32(ptrSize)
	buf = wasm.PtrToByte(ptr, size)

	response := new(SystemInfoReply)
	if err = response.UnmarshalVT(buf); err != nil {
		return nil, err
	}
	return response, nil
}
