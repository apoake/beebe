package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func Test_StrMock(t *testing.T) {
	str := `@str(1, 10, lower)`
	result, err := mockManager.MockData(&str)
	assert.Empty(t, err)
	t.Log(result)

	str = "hello_@str(@num(2,3), @num(5,10), lower)_@str(3, 3)"
	result, err = mockManager.MockData(&str)
	assert.Empty(t, err)
	t.Log(result)

	str = "hello_@str(1, 10, lower)_@str(3,3)"
	result, err = mockManager.MockData(&str)
	assert.Empty(t, err)
	t.Log(result)

	str = "hello_@stre(@num(1,1), @num(3,3), <index@num(1,100)>)"
	result, err = mockManager.MockData(&str)
	assert.Empty(t, err)
	t.Log(result)

	str = "hello_@stre(@num(1,1), @num(3,3), <index@num(1,100)>)"
	result, err = mockManager.MockData(&str)
	assert.Empty(t, err)
	t.Log(result)
}

func Test_PickMock(t *testing.T) {
	str := "@pick([\"a\", \"b\", \"c\"])"
	result, err := mockManager.MockData(&str)
	assert.Empty(t, err)
	t.Log(result)
}

func Test_AddressMock(t *testing.T) {
	str := "@address()"
	result, err := mockManager.MockData(&str)
	assert.Empty(t, err)
	t.Log(result)
}

func Test_CityMock(t *testing.T) {
	str := "@city()"
	result, err := mockManager.MockData(&str)
	assert.Empty(t, err)
	t.Log(result)
}

func Test_ProvinceMock(t *testing.T) {
	str := "@province()"
	result, err := mockManager.MockData(&str)
	assert.Empty(t, err)
	t.Log(result)
}

func Test_RegionMock(t *testing.T) {
	str := "@city()"
	result, err := mockManager.MockData(&str)
	assert.Empty(t, err)
	t.Log(result)
}