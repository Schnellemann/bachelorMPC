package graph

import (
	prot "MPC/Protocol"
)

type Interface interface {
	Plot(title string, variableName string)
	AddData(data prot.Times)
}
