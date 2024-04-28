package mock

import (
	"web_server/internal/entity"

	"github.com/golang/mock/gomock"
)

type MockProductInterface struct {
	ctrl     *gomock.Controller
	recorder *MockProductInterfaceMockRecorder
}

type MockProductInterfaceMockRecorder struct {
	mock *MockProductInterface
}

func NewMockProductInterface(ctrl *gomock.Controller) *MockProductInterface {
	mock := &MockProductInterface{ctrl: ctrl}
	mock.recorder = &MockProductInterfaceMockRecorder{mock}
	return mock
}

func (m *MockProductInterface) EXPECT() *MockProductInterfaceMockRecorder {
	return m.recorder
}

func (m *MockProductInterface) Create(product *entity.Product) error {
	ret := m.ctrl.Call(m, "Create", product)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete implements database.ProductInterface.
func (m *MockProductInterface) Delete(id string) error {
	panic("unimplemented")
}

// GetAll implements database.ProductInterface.
func (m *MockProductInterface) GetAll(page int, limit int, sort string) ([]entity.Product, error) {
	panic("unimplemented")
}

// GetByID implements database.ProductInterface.
func (m *MockProductInterface) GetByID(id string) (*entity.Product, error) {
	panic("unimplemented")
}

// Update implements database.ProductInterface.
func (m *MockProductInterface) Update(product *entity.Product) error {
	panic("unimplemented")
}
