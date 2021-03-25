package config

import (
	"fmt"
	"testing"
)

func TestParseExpressionAdd(t *testing.T) {
	toConvert := "p1+p2"
	exp := ParseExpression(toConvert)
	res, instructions, _ := ConvertAstToExpressionList(exp)
	if len(instructions) != 1 {
		t.Errorf("Number of instructions incorrect got %v expected %v", len(instructions), 1)
	}
	if instructions[0].Op != Add {
		t.Errorf("Wrong instruction got %v expected %v", instructions[0].Op, Add)
	}
	if instructions[0].Left != "p1" {
		t.Errorf("Wrong instruction left side got %v expected %v", instructions[0].Left, "p1")
	}
	if instructions[0].Right != "p2" {
		t.Errorf("Wrong instruction right side got %v expected %v", instructions[0].Right, "p2")
	}
	if instructions[0].Result != res {
		t.Errorf("Result did not match returned %v and created %v", res, instructions[0].Result)
	}
}

func TestParseExpressionMul(t *testing.T) {
	toConvert := "p1*p2"
	exp := ParseExpression(toConvert)
	res, instructions, _ := ConvertAstToExpressionList(exp)
	if len(instructions) != 1 {
		t.Errorf("Number of instructions incorrect got %v expected %v", len(instructions), 1)
	}
	if instructions[0].Op != Multiply {
		t.Errorf("Wrong instruction got %v expected %v", instructions[0].Op, Multiply)
	}
	if instructions[0].Left != "p1" {
		t.Errorf("Wrong instruction left side got %v expected %v", instructions[0].Left, "p1")
	}
	if instructions[0].Right != "p2" {
		t.Errorf("Wrong instruction right side got %v expected %v", instructions[0].Right, "p2")
	}
	if instructions[0].Result != res {
		t.Errorf("Result did not match returned %v and created %v", res, instructions[0].Result)
	}
}

func TestParseExpressionScalar(t *testing.T) {
	toConvert := "1*p1"
	exp := ParseExpression(toConvert)
	res, instructions, _ := ConvertAstToExpressionList(exp)
	if len(instructions) != 1 {
		t.Errorf("Number of instructions incorrect got %v expected %v", len(instructions), 1)
	}
	if instructions[0].Op != Scalar {
		t.Errorf("Wrong instruction got %v expected %v", instructions[0].Op, Scalar)
	}
	if instructions[0].Left != "1" {
		t.Errorf("Wrong instruction left side got %v expected %v", instructions[0].Left, "1")
	}
	if instructions[0].Right != "p1" {
		t.Errorf("Wrong instruction right side got %v expected %v", instructions[0].Right, "p1")
	}
	if instructions[0].Result != res {
		t.Errorf("Result did not match returned %v and created %v", res, instructions[0].Result)
	}

	toConvert = "p1*2"
	exp = ParseExpression(toConvert)
	res, instructions, _ = ConvertAstToExpressionList(exp)
	if len(instructions) != 1 {
		t.Errorf("Number of instructions incorrect got %v expected %v", len(instructions), 1)
	}
	if instructions[0].Op != Scalar {
		t.Errorf("Wrong instruction got %v expected %v", instructions[0].Op, Scalar)
	}
	if instructions[0].Left != "2" {
		t.Errorf("Wrong instruction left side got %v expected %v", instructions[0].Left, "2")
	}
	if instructions[0].Right != "p1" {
		t.Errorf("Wrong instruction right side got %v expected %v", instructions[0].Right, "p1")
	}
	if instructions[0].Result != res {
		t.Errorf("Result did not match returned %v and created %v", res, instructions[0].Result)
	}
}

func TestParseExpressionCombined(t *testing.T) {
	toConvert := "(1*p1+p2)*p3"
	exp := ParseExpression(toConvert)
	res, instructions, _ := ConvertAstToExpressionList(exp)
	if len(instructions) != 3 {
		t.Errorf("Number of instructions incorrect got %v expected %v", len(instructions), 3)
	}
	if instructions[0].Op != Scalar {
		t.Errorf("Wrong third instruction got %v expected %v", instructions[0].Op, Scalar)
	}
	if instructions[1].Op != Add {
		t.Errorf("Wrong third instruction got %v expected %v", instructions[1].Op, Add)
	}
	if instructions[2].Op != Multiply {
		t.Errorf("Wrong third instruction got %v expected %v", instructions[2].Op, Multiply)
	}

	if instructions[2].Result != res {
		t.Errorf("Results did not match returned %v and created %v", res, instructions[0].Result)
	}
	if instructions[0].Result == instructions[1].Result || instructions[0].Result == instructions[2].Result || instructions[1].Result == instructions[2].Result {
		t.Errorf("Overlapping results")
	}
	if instructions[2].Result != "r3" {
		t.Errorf("Result from multiply is %v expected %v", instructions[2].Result, "r3")
	}
	if instructions[1].Result != "r2" {
		t.Errorf("Result from add is %v expected %v", instructions[1].Result, "r2")
	}
	if instructions[0].Result != "r1" {
		t.Errorf("Result from scalar is %v expected %v", instructions[0].Result, "r1")
	}
}

func TestParseExpressionBigCombined(t *testing.T) {
	toConvert := "(1*p1+p2)*p3+(p1*p2+4*p1)"
	exp := ParseExpression(toConvert)
	_, instructions, _ := ConvertAstToExpressionList(exp)
	if len(instructions) != 7 {
		t.Errorf("Number of instructions incorrect got %v expected %v", len(instructions), 7)
	}
	if instructions[0].Op != Scalar {
		t.Errorf("Wrong first instruction got %v expected %v", instructions[0].Op, Scalar)
	}
	if instructions[1].Op != Add {
		t.Errorf("Wrong second instruction got %v expected %v", instructions[1].Op, Add)
	}
	if instructions[2].Op != Multiply {
		t.Errorf("Wrong third instruction got %v expected %v", instructions[2].Op, Multiply)
	}
	if instructions[3].Op != Multiply {
		t.Errorf("Wrong fourth instruction got %v expected %v", instructions[3].Op, Multiply)
	}
	if instructions[4].Op != Scalar {
		t.Errorf("Wrong fifth instruction got %v expected %v", instructions[4].Op, Scalar)
	}
	if instructions[5].Op != Add {
		t.Errorf("Wrong sixth instruction got %v expected %v", instructions[5].Op, Add)
	}
	if instructions[6].Op != Add {
		t.Errorf("Wrong seventh instruction got %v expected %v", instructions[6].Op, Add)
	}
	fmt.Print(instructions)
}
