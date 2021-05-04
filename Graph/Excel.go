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
	fileExisted  bool
	currentSheet string
}

func MkExcel(title string, variableName string) *Excel {
	e := new(Excel)
	e.fileName = title
	e.variableName = variableName
	e.currentSheet = "For Graphs"
	file, err := excel.OpenFile(title + ".xlsx")
	if err == nil {
		e.file = file
		e.fileExisted = true
		for _, sheetName := range e.file.GetSheetList() {
			if sheetName != "For graphs" {
				e.file.DeleteSheet(sheetName)
			}
		}

	} else {
		e.file = excel.NewFile()
		e.file.SetSheetName(e.file.GetSheetName(0), "For graphs")
		e.fileExisted = false
	}
	e.rowCounter = 2
	return e
}

func (e *Excel) setTemplate(sheet string) {
	e.file.SetCellValue(sheet, "A1", e.variableName)
	e.file.SetCellValue(sheet, "B1", "Network (ms)")
	e.file.SetCellValue(sheet, "C1", "Calculate (ms)")
	e.file.SetCellValue(sheet, "D1", "SetupTree (ms)")
	e.file.SetCellValue(sheet, "E1", "Preprocess (ms)")
}

func (e *Excel) Plot() error {
	if e.fileExisted {
		return e.file.Save()
	} else {
		return e.file.SaveAs(e.fileName + ".xlsx")
	}

}

func (e *Excel) NewSeries(name string) {
	e.file.NewSheet(name)
	e.rowCounter = 2
	e.setTemplate(name)
	e.currentSheet = name
}

func (e *Excel) AddData(variable int, data *prot.Times) {
	if e.currentSheet == "For Graphs" {
		fmt.Println("No series to add data to, add a new series before adding data.")
		return
	}
	pos, _ := excel.CoordinatesToCellName(1, e.rowCounter)
	e.file.SetCellValue(e.currentSheet, pos, variable)
	pos, _ = excel.CoordinatesToCellName(2, e.rowCounter)
	e.file.SetCellValue(e.currentSheet, pos, data.Network.Milliseconds())
	pos, _ = excel.CoordinatesToCellName(3, e.rowCounter)
	e.file.SetCellValue(e.currentSheet, pos, data.Calculate.Milliseconds())
	pos, _ = excel.CoordinatesToCellName(4, e.rowCounter)
	e.file.SetCellValue(e.currentSheet, pos, data.SetupTree.Milliseconds())
	pos, _ = excel.CoordinatesToCellName(5, e.rowCounter)
	e.file.SetCellValue(e.currentSheet, pos, data.Preprocess.Milliseconds())

	//Insert SUM() in column 7
	pos, _ = excel.CoordinatesToCellName(7, e.rowCounter)
	from, _ := excel.CoordinatesToCellName(3, e.rowCounter)
	to, _ := excel.CoordinatesToCellName(5, e.rowCounter)
	e.file.SetCellFormula(e.currentSheet, pos, "=SUM("+from+":"+to+")")
	e.rowCounter += 1
}
