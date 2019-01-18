package mregistry

import (
	"fmt"
	"testing"

	"golang.org/x/sys/windows/registry"
)

func TestSetMultipleStringValues(t *testing.T) {
	names := []string{
		"TestStringValue", "TestStringValue2", "TestStringValue3", "TestStringValue4Kitty",
	}
	values := []string{
		"tested string value 1", "tested string value 2", "tested string value 3", "kitty",
	}
	if err := SetMultipleExpandStringValues(registry.CURRENT_USER, `Software\TestTest`,
		registry.QUERY_VALUE|registry.SET_VALUE, names, values); err != nil {
		fmt.Println(err)
	}
}

func TestSetMultipleExpandStringValues(t *testing.T) {
	names := []string{
		"TestExpandStringValue", "TestExpandStringValue2", "TestExpandStringValue3", "TestExpandStringValue4Kitty",
	}
	values := []string{
		"tested expand string value 1", "tested expand string value 2", "tested expand string value 3", "kitty",
	}
	if err := SetMultipleExpandStringValues(registry.CURRENT_USER, `Software\TestTest`,
		registry.QUERY_VALUE|registry.SET_VALUE, names, values); err != nil {
		fmt.Println(err)
	}
}

func TestSetMultipleBinaryValues(t *testing.T) {
	names := []string{
		"TestBinaryValue", "TestBinaryValue2", "TestBinaryValue3", "TestBinaryValue4Kitty",
	}
	values := [][]byte{
		[]byte("TestBinaryValue"), []byte("TestBinaryValue2"), []byte("TestBinaryValue3"), []byte("kitty"),
	}
	if err := SetMultipleBinaryValues(registry.CURRENT_USER, `Software\TestTest`,
		registry.QUERY_VALUE|registry.SET_VALUE, names, values); err != nil {
		fmt.Println(err)
	}
}

func TestSetMultipleDWordValues(t *testing.T) {
	names := []string{
		"TestDwordValue", "TestDwordValue2", "TestDwordValue3", "TestDwordValue4Kitty",
	}
	values := []uint32{
		0x275cbe, 337, 1337, 0x972bce2,
	}
	if err := SetMultipleDWordValues(registry.CURRENT_USER, `Software\TestTest`,
		registry.QUERY_VALUE|registry.SET_VALUE, names, values); err != nil {
		fmt.Println(err)
	}
}

func TestSetMultipleQWordValues(t *testing.T) {
	names := []string{
		"TestQwordValue", "TestQwordValue2", "TestQwordValue3", "TestQwordValue4Kitty",
	}
	values := []uint64{
		0xCCCCCCCC, 337, 1337, 0xDEADBEEF,
	}
	if err := SetMultipleQWordValues(registry.CURRENT_USER, `Software\TestTest`,
		registry.QUERY_VALUE|registry.SET_VALUE, names, values); err != nil {
		fmt.Println(err)
	}
}
