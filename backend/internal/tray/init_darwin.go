//go:build darwin

package tray

import (
	"github.com/go-webgpu/goffi/ffi"
	"github.com/go-webgpu/goffi/types"
)

func initPlatform() {
	handle, err := ffi.LoadLibrary("/System/Library/Frameworks/AppKit.framework/AppKit")
	if err != nil {
		return
	}
	sym, err := ffi.GetSymbol(handle, "NSApplicationLoad")
	if err != nil {
		return
	}
	cif := &types.CallInterface{}
	if err := ffi.PrepareCallInterface(cif, types.DefaultCall, types.VoidTypeDescriptor, nil); err != nil {
		return
	}
	_, _ = ffi.CallFunction(cif, sym, nil, nil)
}
