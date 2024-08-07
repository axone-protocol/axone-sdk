package dataverse

import "context"

type Client interface {
	GetGovAddr(context.Context, string) (string, error)
	ExecGov(context.Context, string, string) (interface{}, error)
}
