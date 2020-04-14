// Code generated by MockGen. DO NOT EDIT.
// Source: debugger.go

// Package mock_debugger is a generated GoMock package.
package mock_debugger

import (
	binary "encoding/binary"
	debugger "github.com/AtomicMalloc/debugger/debugger"
	op "github.com/go-delve/delve/pkg/dwarf/op"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRegisters is a mock of Registers interface
type MockRegisters struct {
	ctrl     *gomock.Controller
	recorder *MockRegistersMockRecorder
}

// MockRegistersMockRecorder is the mock recorder for MockRegisters
type MockRegistersMockRecorder struct {
	mock *MockRegisters
}

// NewMockRegisters creates a new mock instance
func NewMockRegisters(ctrl *gomock.Controller) *MockRegisters {
	mock := &MockRegisters{ctrl: ctrl}
	mock.recorder = &MockRegistersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRegisters) EXPECT() *MockRegistersMockRecorder {
	return m.recorder
}

// GetRegister mocks base method
func (m *MockRegisters) GetRegister(arg0 string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRegister", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRegister indicates an expected call of GetRegister
func (mr *MockRegistersMockRecorder) GetRegister(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRegister", reflect.TypeOf((*MockRegisters)(nil).GetRegister), arg0)
}

// SetRegister mocks base method
func (m *MockRegisters) SetRegister(arg0 string, arg1 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRegister", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRegister indicates an expected call of SetRegister
func (mr *MockRegistersMockRecorder) SetRegister(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRegister", reflect.TypeOf((*MockRegisters)(nil).SetRegister), arg0, arg1)
}

// DwarfRegisters mocks base method
func (m *MockRegisters) DwarfRegisters() *op.DwarfRegisters {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DwarfRegisters")
	ret0, _ := ret[0].(*op.DwarfRegisters)
	return ret0
}

// DwarfRegisters indicates an expected call of DwarfRegisters
func (mr *MockRegistersMockRecorder) DwarfRegisters() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DwarfRegisters", reflect.TypeOf((*MockRegisters)(nil).DwarfRegisters))
}

// MockRegisterHandler is a mock of RegisterHandler interface
type MockRegisterHandler struct {
	ctrl     *gomock.Controller
	recorder *MockRegisterHandlerMockRecorder
}

// MockRegisterHandlerMockRecorder is the mock recorder for MockRegisterHandler
type MockRegisterHandlerMockRecorder struct {
	mock *MockRegisterHandler
}

// NewMockRegisterHandler creates a new mock instance
func NewMockRegisterHandler(ctrl *gomock.Controller) *MockRegisterHandler {
	mock := &MockRegisterHandler{ctrl: ctrl}
	mock.recorder = &MockRegisterHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRegisterHandler) EXPECT() *MockRegisterHandlerMockRecorder {
	return m.recorder
}

// GetRegisters mocks base method
func (m *MockRegisterHandler) GetRegisters(arg0 uint32) (debugger.Registers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRegisters", arg0)
	ret0, _ := ret[0].(debugger.Registers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRegisters indicates an expected call of GetRegisters
func (mr *MockRegisterHandlerMockRecorder) GetRegisters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRegisters", reflect.TypeOf((*MockRegisterHandler)(nil).GetRegisters), arg0)
}

// SetRegisters mocks base method
func (m *MockRegisterHandler) SetRegisters(arg0 uint32, arg1 debugger.Registers) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRegisters", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRegisters indicates an expected call of SetRegisters
func (mr *MockRegisterHandlerMockRecorder) SetRegisters(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRegisters", reflect.TypeOf((*MockRegisterHandler)(nil).SetRegisters), arg0, arg1)
}

// MockControl is a mock of Control interface
type MockControl struct {
	ctrl     *gomock.Controller
	recorder *MockControlMockRecorder
}

// MockControlMockRecorder is the mock recorder for MockControl
type MockControlMockRecorder struct {
	mock *MockControl
}

// NewMockControl creates a new mock instance
func NewMockControl(ctrl *gomock.Controller) *MockControl {
	mock := &MockControl{ctrl: ctrl}
	mock.recorder = &MockControlMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockControl) EXPECT() *MockControlMockRecorder {
	return m.recorder
}

// IsPaused mocks base method
func (m *MockControl) IsPaused() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPaused")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPaused indicates an expected call of IsPaused
func (mr *MockControlMockRecorder) IsPaused() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPaused", reflect.TypeOf((*MockControl)(nil).IsPaused))
}

// Pause mocks base method
func (m *MockControl) Pause() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pause")
	ret0, _ := ret[0].(error)
	return ret0
}

// Pause indicates an expected call of Pause
func (mr *MockControlMockRecorder) Pause() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pause", reflect.TypeOf((*MockControl)(nil).Pause))
}

// Unpause mocks base method
func (m *MockControl) Unpause() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unpause")
	ret0, _ := ret[0].(error)
	return ret0
}

// Unpause indicates an expected call of Unpause
func (mr *MockControlMockRecorder) Unpause() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unpause", reflect.TypeOf((*MockControl)(nil).Unpause))
}

// MockLineInformation is a mock of LineInformation interface
type MockLineInformation struct {
	ctrl     *gomock.Controller
	recorder *MockLineInformationMockRecorder
}

// MockLineInformationMockRecorder is the mock recorder for MockLineInformation
type MockLineInformationMockRecorder struct {
	mock *MockLineInformation
}

// NewMockLineInformation creates a new mock instance
func NewMockLineInformation(ctrl *gomock.Controller) *MockLineInformation {
	mock := &MockLineInformation{ctrl: ctrl}
	mock.recorder = &MockLineInformationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLineInformation) EXPECT() *MockLineInformationMockRecorder {
	return m.recorder
}

// CurrentLine mocks base method
func (m *MockLineInformation) CurrentLine() (string, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentLine")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// CurrentLine indicates an expected call of CurrentLine
func (mr *MockLineInformationMockRecorder) CurrentLine() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentLine", reflect.TypeOf((*MockLineInformation)(nil).CurrentLine))
}

// IsNewLine mocks base method
func (m *MockLineInformation) IsNewLine(arg0 uint64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNewLine", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNewLine indicates an expected call of IsNewLine
func (mr *MockLineInformationMockRecorder) IsNewLine(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNewLine", reflect.TypeOf((*MockLineInformation)(nil).IsNewLine), arg0)
}

// Address mocks base method
func (m *MockLineInformation) Address(arg0 string, arg1 int) uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Address", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Address indicates an expected call of Address
func (mr *MockLineInformationMockRecorder) Address(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Address", reflect.TypeOf((*MockLineInformation)(nil).Address), arg0, arg1)
}

// MockMemoryAccess is a mock of MemoryAccess interface
type MockMemoryAccess struct {
	ctrl     *gomock.Controller
	recorder *MockMemoryAccessMockRecorder
}

// MockMemoryAccessMockRecorder is the mock recorder for MockMemoryAccess
type MockMemoryAccessMockRecorder struct {
	mock *MockMemoryAccess
}

// NewMockMemoryAccess creates a new mock instance
func NewMockMemoryAccess(ctrl *gomock.Controller) *MockMemoryAccess {
	mock := &MockMemoryAccess{ctrl: ctrl}
	mock.recorder = &MockMemoryAccessMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMemoryAccess) EXPECT() *MockMemoryAccessMockRecorder {
	return m.recorder
}

// Read mocks base method
func (m *MockMemoryAccess) Read(address uint64, size uint) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", address, size)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockMemoryAccessMockRecorder) Read(address, size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockMemoryAccess)(nil).Read), address, size)
}

// Write mocks base method
func (m *MockMemoryAccess) Write(address uint64, bytes []byte, size uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", address, bytes, size)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write
func (mr *MockMemoryAccessMockRecorder) Write(address, bytes, size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockMemoryAccess)(nil).Write), address, bytes, size)
}

// MockVariable is a mock of Variable interface
type MockVariable struct {
	ctrl     *gomock.Controller
	recorder *MockVariableMockRecorder
}

// MockVariableMockRecorder is the mock recorder for MockVariable
type MockVariableMockRecorder struct {
	mock *MockVariable
}

// NewMockVariable creates a new mock instance
func NewMockVariable(ctrl *gomock.Controller) *MockVariable {
	mock := &MockVariable{ctrl: ctrl}
	mock.recorder = &MockVariableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockVariable) EXPECT() *MockVariableMockRecorder {
	return m.recorder
}

// Name mocks base method
func (m *MockVariable) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockVariableMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockVariable)(nil).Name))
}

// Parse mocks base method
func (m *MockVariable) Parse(arg0 []byte, arg1 binary.ByteOrder) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parse indicates an expected call of Parse
func (mr *MockVariableMockRecorder) Parse(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockVariable)(nil).Parse), arg0, arg1)
}

// Location mocks base method
func (m *MockVariable) Location() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Location")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Location indicates an expected call of Location
func (mr *MockVariableMockRecorder) Location() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Location", reflect.TypeOf((*MockVariable)(nil).Location))
}

// Size mocks base method
func (m *MockVariable) Size() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size")
	ret0, _ := ret[0].(int)
	return ret0
}

// Size indicates an expected call of Size
func (mr *MockVariableMockRecorder) Size() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockVariable)(nil).Size))
}

// MockSymbol is a mock of Symbol interface
type MockSymbol struct {
	ctrl     *gomock.Controller
	recorder *MockSymbolMockRecorder
}

// MockSymbolMockRecorder is the mock recorder for MockSymbol
type MockSymbolMockRecorder struct {
	mock *MockSymbol
}

// NewMockSymbol creates a new mock instance
func NewMockSymbol(ctrl *gomock.Controller) *MockSymbol {
	mock := &MockSymbol{ctrl: ctrl}
	mock.recorder = &MockSymbolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSymbol) EXPECT() *MockSymbolMockRecorder {
	return m.recorder
}

// GetSymbol mocks base method
func (m *MockSymbol) GetSymbol(arg0 string, arg1 uint64) (debugger.Variable, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSymbol", arg0, arg1)
	ret0, _ := ret[0].(debugger.Variable)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSymbol indicates an expected call of GetSymbol
func (mr *MockSymbolMockRecorder) GetSymbol(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSymbol", reflect.TypeOf((*MockSymbol)(nil).GetSymbol), arg0, arg1)
}

// IsPointer mocks base method
func (m *MockSymbol) IsPointer(arg0 debugger.Variable) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPointer", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPointer indicates an expected call of IsPointer
func (mr *MockSymbolMockRecorder) IsPointer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPointer", reflect.TypeOf((*MockSymbol)(nil).IsPointer), arg0)
}

// GetPointContentSize mocks base method
func (m *MockSymbol) GetPointContentSize(arg0 debugger.Variable) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPointContentSize", arg0)
	ret0, _ := ret[0].(int)
	return ret0
}

// GetPointContentSize indicates an expected call of GetPointContentSize
func (mr *MockSymbolMockRecorder) GetPointContentSize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPointContentSize", reflect.TypeOf((*MockSymbol)(nil).GetPointContentSize), arg0)
}

// ParsePointer mocks base method
func (m *MockSymbol) ParsePointer(arg0 debugger.Variable, arg1 []byte, arg2 binary.ByteOrder) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParsePointer", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParsePointer indicates an expected call of ParsePointer
func (mr *MockSymbolMockRecorder) ParsePointer(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParsePointer", reflect.TypeOf((*MockSymbol)(nil).ParsePointer), arg0, arg1, arg2)
}
