package module

import link "github.com/behavioral-ai/resiliency/link"

func Resolve(name string) (bool, any) {
	switch name {
	case link.NamespaceNameAuth:
		return true, link.Authorization
	case link.NamespaceNameLog:
		return true, link.Logger
	default:
		return false, nil
	}
}
