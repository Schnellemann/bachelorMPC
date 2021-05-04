package graph

import (
	prot "MPC/Protocol"
	"fmt"

	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Excel struct {
	file         *excel.File
	rowCounter   int
	variableName string
	fileName     string
	currentSheet string
}

func MkExcel(title string, variableName string) *Excel {
	e := new(Excel)
	e.fileName = title
	e.variableName = variableName
	e.currentSheet = "For graphs"
	e.rowCounter = 2
	file, err := excel.OpenFile(title + ".xlsx")
	if err == nil {
		e.file = file
		e.file.SetActiveSheet(e.file.GetSheetIndex("For graphs"))
		for _, sheetName := range e.file.GetSheetList() {
			if sheetName != "For graphs" {
				e.file.DeleteSheet(sheetName)
			}
		}

	} else {
		e.file = excel.NewFile()
		e.file.SetSheetName(e.file.GetSheetName(e.file.GetActiveSheetIndex()), "For graphs")
	}
	return e
}

func (e *Excel) setTemplate(sheet string) {
	check(e.file.SetCellValue(sheet, "A1", e.variableName))
	check(e.file.SetCellValue(sheet, "B1", "Network (ms)"))
	check(e.file.SetCellValue(sheet, "C1", "Calculate (ms)"))
	check(e.file.SetCellValue(sheet, "D1", "SetupTree (ms)"))
	check(e.file.SetCellValue(sheet, "E1", "Preprocess (ms)"))
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (e *Excel) Plot() error {
	return e.file.SaveAs(e.fileName + ".xlsx")
}

func (e *Excel) NewSeries(name string) {
	e.file.NewSheet(name)
	e.rowCounter = 2
	e.setTemplate(name)
	e.currentSheet = name
}

func (e *Excel) AddData(variable int, data *prot.Times) {
	if e.currentSheet == "For graphs" {
		fmt.Println("No series to add data to, add a new series before adding data.")
		return
	}
	pos, _ := excel.CoordinatesToCellName(1, e.rowCounter)
	check(e.file.SetCellValue(e.currentSheet, pos, variable))
	pos, err := excel.CoordinatesToCellName(2, e.rowCounter)
	check(err)
	check(e.file.SetCellValue(e.currentSheet, pos, data.Network.Milliseconds()))
	pos, err = excel.CoordinatesToCellName(3, e.rowCounter)
	check(err)
	check(e.file.SetCellValue(e.currentSheet, pos, data.Calculate.Milliseconds()))
	pos, err = excel.CoordinatesToCellName(4, e.rowCounter)
	check(err)
	check(e.file.SetCellValue(e.currentSheet, pos, data.SetupTree.Milliseconds()))
	pos, err = excel.CoordinatesToCellName(5, e.rowCounter)
	check(err)
	check(e.file.SetCellValue(e.currentSheet, pos, data.Preprocess.Milliseconds()))

	//Insert SUM() in column 7
	pos, err = excel.CoordinatesToCellName(7, e.rowCounter)
	check(err)
	from, err := excel.CoordinatesToCellName(3, e.rowCounter)
	check(err)
	to, err := excel.CoordinatesToCellName(5, e.rowCounter)
	check(err)
	check(e.file.SetCellFormula(e.currentSheet, pos, "SUM("+from+":"+to+")"))
	e.rowCounter += 1
}
