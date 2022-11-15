package osbe

type Permissioner interface {	
	Reload() error
	IsAllowed(role, controller, method string) bool
}
