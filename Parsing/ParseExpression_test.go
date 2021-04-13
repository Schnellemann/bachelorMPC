package parsing

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

	switch node := insTree.Instruction.(type) {
	case *AddInstruction:
		if node.Left != "p1" {
			t.Errorf("Wrong instruction left side got %v expected %v", node.Left, "p1")
		}
		if node.Right != "p2" {
			t.Errorf("Wrong instruction right side got %v expected %v", node.Right, "p2")
		}
		if node.Result == "" {
			t.Error("Result was not set")
		}
	default:
		t.Error("Was not a Add instruction")
	}
}

func TestParseExpressionMul(t *testing.T) {
	toConvert := "p1*p2"
	exp := ParseExpression(toConvert)
	insTree, _ := ConvertAstToTree(exp)
	if insTree.Left != nil || insTree.Right != nil {
		t.Error("Incorret instruction tree, has children but shouldnt")
	}

	switch node := insTree.Instruction.(type) {
	case *MultInstruction:
		if node.Left != "p1" {
			t.Errorf("Wrong instruction left side got %v expected %v", node.Left, "p1")
		}
		if node.Right != "p2" {
			t.Errorf("Wrong instruction right side got %v expected %v", node.Right, "p2")
		}
		if node.Result == "" {
			t.Error("Result was not set")
		}
		if node.Num != 1 {
			t.Errorf("Number was %v but expect 1", node.Num)
		}
	default:
		t.Error("Was not a Mult instruction")
	}
}

func TestCombined(t *testing.T) {
	toConvert := "(p1*p2)+((p1*p3)*2)"
	exp := ParseExpression(toConvert)
	insTree, _ := ConvertAstToTree(exp)

	//Check outermost addition
	switch node := insTree.Instruction.(type) {
	case *AddInstruction:
		if node.Left != "r1" {
			t.Errorf("Wrong instruction left side got %v expected %v", node.Left, "r1")
		}
		if node.Right != "r3" {
			t.Errorf("Wrong instruction right side got %v expected %v", node.Right, "r3")
		}
		if node.Result == "" {
			t.Error("Result was not set")
		}
	default:
		t.Error("Was not an Add instruction")
	}

	//Check p1*p2
	switch node := insTree.Left.Instruction.(type) {
	case *MultInstruction:
		if node.Left != "p1" {
			t.Errorf("Wrong instruction left side got %v expected %v", node.Left, "p1")
		}
		if node.Right != "p2" {
			t.Errorf("Wrong instruction right side got %v expected %v", node.Right, "p2")
		}
		if node.Result == "" {
			t.Error("Result was not set")
		}
		if node.Num != 1 {
			t.Errorf("Number was %v but expect 1", node.Num)
		}
	default:
		t.Error("Was not a Mult instruction")
	}

	//Check (p1*p3)*2
	switch node := insTree.Right.Instruction.(type) {
	case *ScalarInstruction:
		if node.Scalar != 2 {
			t.Errorf("Wrong instruction left side got %v expected %v", node.Scalar, 2)
		}
		if node.Variable != "r2" {
			t.Errorf("Wrong instruction right side got %v expected %v", node.Variable, "r2")
		}
		if node.Result != "r3" {
			t.Error("Result was not set")
		}

	default:
		t.Error("Was not a Scalar instruction")
	}

	//Check (p1*p3)
	switch node := insTree.Right.Left.Instruction.(type) {
	case *MultInstruction:
		if node.Left != "p1" {
			t.Errorf("Wrong instruction left side got %v expected %v", node.Left, "p1")
		}
		if node.Right != "p3" {
			t.Errorf("Wrong instruction right side got %v expected %v", node.Right, "p3")
		}
		if node.Result != "r2" {
			t.Error("Result was not set")
		}
		if node.Num != 2 {
			t.Errorf("Number was %v but expect 2", node.Num)
		}
	default:
		t.Error("Was not a Mult instruction")
	}

}

func TestParseExpressionScalar(t *testing.T) {
	toConvert := "1*p1"
	exp := ParseExpression(toConvert)
	insTree, _ := ConvertAstToTree(exp)

	if insTree.Left != nil || insTree.Right != nil {
		t.Error("Incorret instruction tree, has children but shouldnt")
	}

	switch node := insTree.Instruction.(type) {
	case *ScalarInstruction:
		if node.Scalar != 1 {
			t.Errorf("Wrong instruction left side got %v expected %v", node.Scalar, 1)
		}
		if node.Variable != "p1" {
			t.Errorf("Wrong instruction right side got %v expected %v", node.Variable, "p1")
		}
		if node.Result == "" {
			t.Error("Result was not set")
		}
	default:
		t.Error("Was not a Scalar instruction")
	}

	toConvert = "p1*2"
	exp = ParseExpression(toConvert)
	insTree, _ = ConvertAstToTree(exp)

	switch node := insTree.Instruction.(type) {
	case *ScalarInstruction:
		if node.Scalar != 2 {
			t.Errorf("Wrong instruction left side got %v expected %v", node.Scalar, 2)
		}
		if node.Variable != "p1" {
			t.Errorf("Wrong instruction right side got %v expected %v", node.Variable, "p1")
		}
		if node.Result == "" {
			t.Error("Result was not set")
		}
	default:
		t.Error("Was not a Scalar instruction")
	}
}
