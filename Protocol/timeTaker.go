package protocol

import (
	config "MPC/Config"
	"time"
)

type Times struct {
	Network    time.Duration
	Calculate  time.Duration
	SetupTree  time.Duration
	Preprocess time.Duration
	Run        time.Duration
}

func mkTimes() *Times {
	t := new(Times)
	t.Network = 0
	t.Calculate = 0
	t.SetupTree = 0
	t.Preprocess = 0
	t.Run = 0
	return t
}

type TimeMeasuring struct {
	prot   Prot
	Timer  Times
	config *config.Config
}

func MkTimeMeasuringProt(prot Prot, config *config.Config) *TimeMeasuring {
	tm := new(TimeMeasuring)
	tm.prot = prot
	tm.Timer = *mkTimes()
	tm.config = config
	return tm
}

func (tM *TimeMeasuring) startNetwork() {
	startTime := time.Now()
	tM.prot.startNetwork()
	endTime := time.Now()
	tM.Timer.Network = endTime.Sub(startTime)
	//log.Printf("Starting the network for party %v took %v.\n", tM.config.VariableConfig.PartyNr, endTime.Sub(startTime))
}
func (tM *TimeMeasuring) calculate() int64 {
	startTime := time.Now()
	res := tM.prot.calculate()
	endTime := time.Now()
	tM.Timer.Calculate = endTime.Sub(startTime)
	//log.Printf("Calculating instructions for party %v took %v.\n", tM.config.VariableConfig.PartyNr, endTime.Sub(startTime))
	return res
}

func (tM *TimeMeasuring) setupTree() {
	startTime := time.Now()
	tM.prot.setupTree()
	endTime := time.Now()
	tM.Timer.SetupTree = endTime.Sub(startTime)
	//log.Printf("Parsing the instructions tree for party %v took %v.\n", tM.config.VariableConfig.PartyNr, endTime.Sub(startTime))

}
func (tM *TimeMeasuring) runPreprocess() {
	startTime := time.Now()
	tM.prot.runPreprocess()
	endTime := time.Now()
	tM.Timer.Preprocess = endTime.Sub(startTime)
	//log.Printf("Running preprocess for party %v took %v.\n", tM.config.VariableConfig.PartyNr, endTime.Sub(startTime))

}
func (tM *TimeMeasuring) Run() int64 {
	startTime := time.Now()
	tM.startNetwork()
	tM.setupTree()
	tM.runPreprocess()
	res := tM.calculate()
	endTime := time.Now()
	tM.Timer.Run = endTime.Sub(startTime)
	//log.Printf("Running the full protocol for party %v took %v.\n", tM.config.VariableConfig.PartyNr, endTime.Sub(startTime))
	return res
}
func (tM *TimeMeasuring) Destroy() {
	tM.prot.Destroy()
}
