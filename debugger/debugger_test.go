package debugger_test

import (
	"testing"

	"github.com/AtomicMalloc/debugger/debugger"
	mocks "github.com/AtomicMalloc/debugger/mock_debugger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setup(mockCtrl *gomock.Controller) (*mocks.MockMemoryAccess, *mocks.MockControl, *mocks.MockLineInformation, *mocks.MockRegisterHandler, *mocks.MockSymbol, *debugger.Debugger) {
	mem := mocks.NewMockMemoryAccess(mockCtrl)
	control := mocks.NewMockControl(mockCtrl)
	lineInfo := mocks.NewMockLineInformation(mockCtrl)
	registers := mocks.NewMockRegisterHandler(mockCtrl)
	symbols := mocks.NewMockSymbol(mockCtrl)
	debugger := debugger.NewDebugger(mem, control, lineInfo, registers, symbols)
	return mem, control, lineInfo, registers, symbols, debugger
}

//Checks that we can step through a program without issues 
//(assumes we're not going through breakpoint)
func TestStep(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	_, cntrl, lineInfo, regs, _, dbg := setup(mockCtrl)
	dummyRegisters := mocks.NewMockRegisters(mockCtrl)

	vcpu := uint32(0)
	rip := uint64(0x1)
	gomock.InOrder(
		cntrl.EXPECT().IsPaused().Return(true),
		regs.EXPECT().GetRegisters(vcpu).Return(dummyRegisters, nil),
		dummyRegisters.EXPECT().GetRegister("rflags").Return(uint64(0x0), nil),
		dummyRegisters.EXPECT().SetRegister("rflags", uint64(256)).Return(nil),
		regs.EXPECT().SetRegisters(vcpu, dummyRegisters).Return(nil),
		cntrl.EXPECT().IsPaused().Return(true),
		regs.EXPECT().GetRegisters(vcpu).Return(dummyRegisters, nil),
		dummyRegisters.EXPECT().GetRegister("rip").Return(rip, nil),
		lineInfo.EXPECT().IsNewLine(rip).Return(true),
		cntrl.EXPECT().Unpause().Return(nil),
	)
	err := dbg.Step(vcpu)
	assert.Nil(t, err)
}

//Tests we can add a breakpoint correctly
func TestAddBreakpoint(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mem, cntrl, lineInfo, _, _, dbg := setup(mockCtrl)
	vcpu := uint32(0)
	address := uint64(0x13)
	filename := "start.c"
	line := 3
	gomock.InOrder(
		cntrl.EXPECT().IsPaused().Return(true),
		lineInfo.EXPECT().Address(filename, line).Return(address),
		mem.EXPECT().Read(address, uint(1)).Return([]byte{0x1}, nil),
		mem.EXPECT().Write(address, []byte{0xCC}, uint(1)).Return(nil),
	)
	err := dbg.SetBreakpoint(filename, line, vcpu)
	assert.Nil(t, err)
}

//Tests we can remove a breakpoint correctly 
func TestRemoveBreakpoint(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mem, cntrl, lineInfo, _, _, dbg := setup(mockCtrl)
	vcpu := uint32(0)
	address := uint64(0x13)
	filename := "start.c"
	line := 3
	gomock.InOrder(
		cntrl.EXPECT().IsPaused().Return(true),
		lineInfo.EXPECT().Address(filename, line).Return(address),
		mem.EXPECT().Read(address, uint(1)).Return([]byte{0x1}, nil),
		mem.EXPECT().Write(address, []byte{0xCC}, uint(1)).Return(nil),
		cntrl.EXPECT().IsPaused().Return(true),
		lineInfo.EXPECT().Address(filename, line).Return(address),
		mem.EXPECT().Write(address, []byte{0x1}, uint(1)).Return(nil),
	)
	err := dbg.SetBreakpoint(filename, line, vcpu)
	assert.Nil(t, err)
	err = dbg.RemoveBreakpoint(filename, line, vcpu)
	assert.Nil(t, err)
}

//Tests we can continue correctly to the next breakpoint
func TestContinue(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	_, cntrl, _, regs, _, dbg := setup(mockCtrl)
	dummyRegisters := mocks.NewMockRegisters(mockCtrl)

	vcpu := uint32(0)
	gomock.InOrder(
		cntrl.EXPECT().IsPaused().Return(true),
		regs.EXPECT().GetRegisters(vcpu).Return(dummyRegisters, nil),
		dummyRegisters.EXPECT().GetRegister("rflags").Return(uint64(256), nil),
		dummyRegisters.EXPECT().SetRegister("rflags", uint64(0)).Return(nil),
		regs.EXPECT().SetRegisters(vcpu, dummyRegisters).Return(nil),
		cntrl.EXPECT().Unpause().Return(nil),
		cntrl.EXPECT().IsPaused().Return(true),
	)

	err := dbg.Continue(vcpu)
	assert.Nil(t, err)
}
