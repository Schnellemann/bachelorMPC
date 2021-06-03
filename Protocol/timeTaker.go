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

type TimeMeasuring struct {
	prot   Prot
	Timer  *Times
	config *config.Config
}

func AverageTimes(times []*Times) *Times {
	t := new(Times)
	for _, v := range times {
		t.Network += v.Network
		t.Calculate += v.Calculate
		t.SetupTree += v.SetupTree
		t.Preprocess += v.Preprocess
		t.Run += v.Run
	}
	t.Network = time.Duration(int64(t.Network) / int64(len(times)))
	t.Calculate = time.Duration(int64(t.Calculate) / int64(len(times)))
	t.SetupTree = time.Duration(int64(t.SetupTree) / int64(len(times)))
	t.Preprocess = time.Duration(int64(t.Preprocess) / int64(len(times)))
	t.Run = time.Duration(int64(t.Run) / int64(len(times)))
	return t

}

func MkTimeMeasuringProt(prot Prot, config *config.Config, timer *Times) *TimeMeasuring {
	tm := new(TimeMeasuring)
	tm.prot = prot
	tm.Timer = timer
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
