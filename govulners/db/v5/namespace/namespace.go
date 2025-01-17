package namespace

import (
	"github.com/nextlinux/govulners/govulners/db/v5/pkg/resolver"
)

type Namespace interface {
	Provider() string
	Resolver() resolver.Resolver
	String() string
}
