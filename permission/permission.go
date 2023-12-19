package permission

import "fmt"

const DEFAULT_ROLE = "guest"

type permMethod map[string]bool
type permController map[string]permMethod
type PermRules map[string]permController

//controller=no _Controller postfix!!!
func (pr *PermRules) IsAllowed(role, controller, method string) bool{
	if role == "" {
		role = DEFAULT_ROLE
	}
	if pr_contr, ok := map[string]permController(*pr)[role]; ok {
		if pr_meth, ok := pr_contr[controller]; ok {
			if pr_allowed, ok := pr_meth[method]; ok {
				return pr_allowed
			}
		}
	}
	return false
}


type Provider interface {
	InitManager(manParams []interface{}) error
	Reload() error
	IsAllowed(role, controller, method string) bool
}

var provides = make(map[string]Provider)

// Register makes a permission provider available by name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, provide Provider) {
	if provide == nil {
		panic("permission: Register provide is nil")
	}
	if _, dup := provides[name]; dup {
		panic("permission: Register called twice for provide " + name)
	}
	provides[name] = provide
}

func NewManager(manName string, manParams ...interface{}) (Provider, error) {
	manager, ok := provides[manName]
	if !ok {
		return nil, fmt.Errorf("permission: unknown manager %q (forgotten import?)", manName)
	}
	manager.InitManager(manParams)
	return manager, nil
}

