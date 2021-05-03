package graph

import (
	prot "MPC/Protocol"

	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Excel struct {
	file       *excel.File
	rowCounter int
}

func MkExcel() *Excel {
	e := new(Excel)
	e.file = excel.NewFile()
	e.file.SetCellValue(e.file.GetSheetName(0), "B1", "Network (ms)")
	e.file.SetCellValue(e.file.GetSheetName(0), "C1", "Calculate (ms)")
	e.file.SetCellValue(e.file.GetSheetName(0), "D1", "SetupTree (ms)")
	e.file.SetCellValue(e.file.GetSheetName(0), "E1", "Preprocess (ms)")
	e.rowCounter = 2
	return e
}

func (e *Excel) Plot(title string, variableName string) error {
	e.file.SetCellValue(e.file.GetSheetName(0), "A1", variableName)
	err := e.file.SaveAs(title + ".xlsx")
	return err

}

func (e *Excel) AddData(variable int, data *prot.Times) {
	pos, _ := excel.CoordinatesToCellName(1, e.rowCounter)
	e.file.SetCellValue(e.file.GetSheetName(0), pos, variable)
	pos, _ = excel.CoordinatesToCellName(2, e.rowCounter)
	e.file.SetCellValue(e.file.GetSheetName(0), pos, data.Network.Milliseconds())
	pos, _ = excel.CoordinatesToCellName(3, e.rowCounter)
	e.file.SetCellValue(e.file.GetSheetName(0), pos, data.Calculate.Milliseconds())
	pos, _ = excel.CoordinatesToCellName(4, e.rowCounter)
	e.file.SetCellValue(e.file.GetSheetName(0), pos, data.SetupTree.Milliseconds())
	pos, _ = excel.CoordinatesToCellName(5, e.rowCounter)
	e.file.SetCellValue(e.file.GetSheetName(0), pos, data.Preprocess.Milliseconds())
	e.rowCounter += 1
}
