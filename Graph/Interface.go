package graph

import (
	prot "MPC/Protocol"
)

type Interface interface {
	Plot() error
	AddData(variable int, data *prot.Times)
	NewSeries(name string)
}
