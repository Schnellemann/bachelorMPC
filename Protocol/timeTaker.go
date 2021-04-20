package protocol

import (
	config "MPC/Config"
	"log"
	"time"
)

type TimeMeasuring struct {
	prot   Prot
	log    *log.Logger
	config *config.Config
}

func mkTimeMeasuringProt(prot Prot, config *config.Config, log *log.Logger) *TimeMeasuring {
	tm := new(TimeMeasuring)
	tm.prot = prot
	tm.log = log
	tm.config = config
	return tm
}

func (tM *TimeMeasuring) startNetwork() {
	startTime := time.Now()
	tM.prot.startNetwork()
	endTime := time.Now()
	log.Printf("Starting the network for party %v took %v.\n", tM.config.VariableConfig.PartyNr, endTime.Sub(startTime))
}
func (tM *TimeMeasuring) calculate() int64 {
	startTime := time.Now()
	res := tM.prot.calculate()
	endTime := time.Now()
	log.Printf("Calculating instructions for party %v took %v.\n", tM.config.VariableConfig.PartyNr, endTime.Sub(startTime))
	return res
}

func (tM *TimeMeasuring) setupTree() {
	startTime := time.Now()
	tM.prot.setupTree()
	endTime := time.Now()
	log.Printf("Parsing the instructions tree for party %v took %v.\n", tM.config.VariableConfig.PartyNr, endTime.Sub(startTime))

}
func (tM *TimeMeasuring) runPreprocess() {
	startTime := time.Now()
	tM.prot.runPreprocess()
	endTime := time.Now()
	log.Printf("Running preprocess for party %v took %v.\n", tM.config.VariableConfig.PartyNr, endTime.Sub(startTime))

}
func (tM *TimeMeasuring) Run() int64 {
	startTime := time.Now()
	tM.startNetwork()
	tM.setupTree()
	tM.runPreprocess()
	res := tM.calculate()
	endTime := time.Now()
	log.Printf("Running the full protocol for party %v took %v.\n", tM.config.VariableConfig.PartyNr, endTime.Sub(startTime))
	return res
}
func (tM *TimeMeasuring) Destroy() {
	tM.prot.Destroy()
}
