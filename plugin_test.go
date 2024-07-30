package main

import (
	"errors"
	"fmt"
	"plugin"
	"testing"

	iplugin "github.com/ipfs/kubo/plugin"
)

// TestPluginLoad smoke tsting if plugin successfully compiled and work on current OS
func TestPluginLoad(t *testing.T) {
	pl, err := plugin.Open("sdspfs.so")
	if err != nil {
		fmt.Println("Open FAIL")
		panic(err)
	}
	fmt.Println("Open OK")
	pls, err := pl.Lookup("Plugins")
	if err != nil {
		fmt.Println("Lookup FAIL")
		panic(err)
	}
	fmt.Println("Lookup OK")

	_, ok := pls.(*[]iplugin.Plugin)
	if !ok {
		fmt.Println("Type assertion FAIL")
		panic(errors.New("filed 'Plugins' didn't contain correct type"))
	}
	fmt.Println("Type assertion OK")
}
