package model

import (
	"testing"
	"fmt"
)


func Test_StrMock(t *testing.T) {
	str := `@str(1, 10, lower)`
	result, err := mockManager.Mock(&str)
	fmt.Printf("%v    %v\n", result, err)
}