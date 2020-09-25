package debugger_test

import (
	"errors"
	"testing"

	"github.com/StardustOS/duster/debugger"
	mocks "github.com/StardustOS/duster/mock_debugger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	address  = uint64(0xff)
	size     = uint(1)
	breakInt = byte(0xcc)
)

var (
	content = []byte{0x25}
)

func TestAdd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	memoryAccess := mocks.NewMockMemoryAccess(mockCtrl)
	address := uint64(0x1)
	size := uint(1)

	gomock.InOrder(
		memoryAccess.EXPECT().Read(address, size).Return(content, nil),
		memoryAccess.EXPECT().Write(address, []byte{breakInt}, uint(1)).Return(nil),
	)

	//Test the breakpoint manager will successful write to memory
	manager := debugger.NewBreakpointManager(memoryAccess)
	err := manager.Add(address)
	assert.Nil(t, err)

	// Test that the case where breakpoint already exists at that address.
	err = manager.Add(address)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Error: breakpoint already at 1")

}

// Tests Add will return a read error from memory
func TestAddReadError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	memoryAccess := mocks.NewMockMemoryAccess(mockCtrl)
	errExpected := errors.New("Error: Something went wrong with memory")
	memoryAccess.EXPECT().Read(address, size).Return(content, errExpected)
	manager := debugger.NewBreakpointManager(memoryAccess)
	err := manager.Add(address)
	assert.Equal(t, err, errExpected)
}

// Tests Add will return a write error from memory
func TestAddWriteError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	errExpected := errors.New("Error: Something went wrong with memory")
	memoryAccess := mocks.NewMockMemoryAccess(mockCtrl)
	memoryAccess.EXPECT().Read(address, size).Return(content, nil)
	memoryAccess.EXPECT().Write(address, []byte{breakInt}, uint(1)).Return(errExpected)
	manager := debugger.NewBreakpointManager(memoryAccess)
	err := manager.Add(address)
	assert.Equal(t, err, errExpected)
}

// // Tests Remove will successful write back the original byte
// func TestRemove(t *testing.T) {
// 	mem := new(mockMem)
// 	mem.On("Read", address, size).Return(content, nil)
// 	mem.On("Write", address, []byte{breakInt}, uint(1)).Return(nil)
// 	manager := NewBreakpointManager(mem)
// 	err := manager.Add(address)
// 	assert.Nil(t, err)
// 	mem.On("Write", address, content, uint(1)).Return(nil)
// 	err = manager.Remove(address)
// 	assert.Nil(t, err)

// 	// Tests doesn't allow for double remove
// 	err = manager.Remove(address)
// 	assert.Equal(t, BreakPointError{address, NotFound}, err)
// }

// // Tests that remove will pass up the write error if it occurs
// func TestRemoveWriteError(t *testing.T) {
// 	mem := new(mockMem)
// 	e := errors.New("An error occured")
// 	mem.On("Read", address, size).Return(content, nil)
// 	mem.On("Write", address, []byte{breakInt}, uint(1)).Return(nil)
// 	manager := NewBreakpointManager(mem)
// 	err := manager.Add(address)
// 	assert.Nil(t, err)
// 	mem.On("Write", address, content, uint(1)).Return(e)
// 	err = manager.Remove(address)
// 	assert.Equal(t, e, err)
// }

// //Tests addresses returns a valid set of addresses
// func TestAddresses(t *testing.T) {
// 	mem := new(mockMem)
// 	manager := NewBreakpointManager(mem)

// 	for i := uint64(0); i < 3; i += 1 {
// 		mem.On("Read", address + i, size).Return(content, nil)
// 		mem.On("Write", address + i, []byte{breakInt}, uint(1)).Return(nil)
// 		err := manager.Add(address + i)
// 		assert.Nil(t, err)
// 	}
// 	list := manager.Addresses()
// 	for index, current := range list {
// 		assert.Equal(t, address + uint64(index), current)
// 	}
// }
