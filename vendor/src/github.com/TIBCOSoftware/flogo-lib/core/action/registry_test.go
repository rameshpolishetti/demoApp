package action

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockFactory struct {
}

func (f *MockFactory) New(config *Config) Action {
	return &MockAction{}
}

type MockAction struct {
}

func (a *MockAction) Run(context context.Context, uri string, options interface{}, handler ResultHandler) error {
	return nil
}

func clearFactory (){
	factories = make(map[string]Factory)
}

func clearActions (){
	actions = make(map[string]Action)
}

//TestRegisterFactoryEmptyRef
func TestRegisterFactoryEmptyRef(t *testing.T) {
	defer clearFactory()

	// Register factory
	err := RegisterFactory("", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "RegisterFactory: ref is empty", err.Error())
}

//TestRegisterFactoryNilFactory
func TestRegisterFactoryNilFactory(t *testing.T) {
	defer clearFactory()

	// Register factory
	err := RegisterFactory("github.com/mock", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "RegisterFactory: factory is nil", err.Error())
}

//TestRegisterFactoryDuplicated
func TestRegisterFactoryDuplicated(t *testing.T) {

	defer clearFactory()
	f := &MockFactory{}

	// Add factory: this time should pass
	err := RegisterFactory("github.com/mock", f)
	assert.Nil(t, err)
	// Add factory: this time should fail, duplicated
	err = RegisterFactory("github.com/mock", f)
	assert.NotNil(t, err)
	assert.Equal(t, "RegisterFactory: already registered factory for ref 'github.com/mock'", err.Error())
}

//TestRegisterFactoryOk
func TestRegisterFactoryOk(t *testing.T) {
	defer clearFactory()

	f := &MockFactory{}

	// Add factory
	err := RegisterFactory("github.com/mock", f)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(Factories()))
}

//TestGetFactoriesOk
func TestGetFactoriesOk(t *testing.T) {
	defer clearFactory()

	f := &MockFactory{}

	// Add factory
	err := RegisterFactory("github.com/mock", f)
	assert.Nil(t, err)

	// Get factory
	fs := Factories()
	assert.Equal(t, 1, len(fs))
}

//TestRegisterEmptyId
func TestRegisterEmptyId(t *testing.T) {
	defer clearActions()

	// Register
	err := Register("", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "error registering action, id is empty", err.Error())
}

//TestRegisterNilInstance
func TestRegisterNilInstance(t *testing.T) {
	defer clearActions()

	// Add instance
	err := Register("myInstanceId", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "error registering action for id 'myInstanceId', action is nil", err.Error())
}

//TestRegisterDuplicated
func TestRegisterDuplicated(t *testing.T) {
	defer clearActions()

	a := &MockAction{}

	// Add instance: this time should pass
	err := Register("myinstanceId", a)
	assert.Nil(t, err)
	// Add instance: this time should fail, duplicated
	err = Register("myinstanceId", a)
	assert.NotNil(t, err)
	assert.Equal(t, "Error registering action, action already registered for id 'myinstanceId'", err.Error())
}

//TestRegisterOk
func TestRegisterOk(t *testing.T) {
	defer clearActions()

	a := &MockAction{}

	// Add instance
	err := Register("myinstanceId", a)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(Actions()))
}

//TestGetActionOk
func TestGetActionOk(t *testing.T) {
	defer clearActions()

	a := &MockAction{}

	// Add instance
	err := Register("myinstanceId", a)
	assert.Nil(t, err)

	myInstance :=  Get("myinstanceId")
	assert.NotNil(t, myInstance)

	myUnknownInstance := Get("myunknowninstanceId")
	assert.Nil(t, myUnknownInstance)
}
