package config

import (
	"testing"
)

func TestParseExpressionAdd(t *testing.T) {
	toConvert := "p1+p2"
	exp := ParseExpression(toConvert)
	insTree, _ := ConvertAstToTree(exp)

	if insTree.Left != nil || insTree.Right != nil {
		t.Error("Incorret instruction tree, has children but shouldnt")
	}
	if insTree.Instruction.Op != Add {
		t.Errorf("Wrong instruction got %v expected %v", insTree.Instruction.Op, Add)
	}
	if insTree.Instruction.Left != "p1" {
		t.Errorf("Wrong instruction left side got %v expected %v", insTree.Instruction.Left, "p1")
	}
	if insTree.Instruction.Right != "p2" {
		t.Errorf("Wrong instruction right side got %v expected %v", insTree.Instruction.Right, "p2")
	}
	if insTree.Instruction.Result == "" {
		t.Error("Result was not set")
	}
}

func TestParseExpressionMul(t *testing.T) {
	toConvert := "p1*p2"
	exp := ParseExpression(toConvert)
	insTree, _ := ConvertAstToTree(exp)
	if insTree.Left != nil || insTree.Right != nil {
		t.Error("Incorret instruction tree, has children but shouldnt")
	}
	if insTree.Instruction.Op != Multiply {
		t.Errorf("Wrong instruction got %v expected %v", insTree.Instruction.Op, Multiply)
	}
	if insTree.Instruction.Left != "p1" {
		t.Errorf("Wrong instruction left side got %v expected %v", insTree.Instruction.Left, "p1")
	}
	if insTree.Instruction.Right != "p2" {
		t.Errorf("Wrong instruction right side got %v expected %v", insTree.Instruction.Right, "p2")
	}
	if insTree.Instruction.Result == "" {
		t.Error("Result was not set")
	}
}

func TestParseExpressionScalar(t *testing.T) {
	toConvert := "1*p1"
	exp := ParseExpression(toConvert)
	insTree, _ := ConvertAstToTree(exp)

	if insTree.Left != nil || insTree.Right != nil {
		t.Error("Incorret instruction tree, has children but shouldnt")
	}
	if insTree.Instruction.Op != Scalar {
		t.Errorf("Wrong instruction got %v expected %v", insTree.Instruction.Op, Scalar)
	}
	if insTree.Instruction.Left != "1" {
		t.Errorf("Wrong instruction left side got %v expected %v", insTree.Instruction.Left, "1")
	}
	if insTree.Instruction.Right != "p1" {
		t.Errorf("Wrong instruction right side got %v expected %v", insTree.Instruction.Right, "p1")
	}
	if insTree.Instruction.Result == "" {
		t.Error("Result was not set")
	}

	toConvert = "p1*2"
	exp = ParseExpression(toConvert)
	insTree, _ = ConvertAstToTree(exp)

	if insTree.Left != nil || insTree.Right != nil {
		t.Error("Incorret instruction tree, has children but shouldnt")
	}
	if insTree.Instruction.Op != Scalar {
		t.Errorf("Wrong instruction got %v expected %v", insTree.Instruction.Op, Scalar)
	}
	if insTree.Instruction.Left != "2" {
		t.Errorf("Wrong instruction left side got %v expected %v", insTree.Instruction.Left, "2")
	}
	if insTree.Instruction.Right != "p1" {
		t.Errorf("Wrong instruction right side got %v expected %v", insTree.Instruction.Right, "p1")
	}
	if insTree.Instruction.Result == "" {
		t.Error("Result was not set")
	}
}
