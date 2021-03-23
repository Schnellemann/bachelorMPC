package config

import (
	"testing"
)

func TestParseExpressionAdd(t *testing.T) {
	toConvert := "p1+p2"
	exp := parseExpression(toConvert)
	res, instructions, _ := convertAstToExpressionList(exp)
	if len(instructions) != 1 {
		t.Errorf("Number of instructions incorrect got %v expected %v", len(instructions), 1)
	}
	if instructions[0].op != add {
		t.Errorf("Wrong instruction got %v expected %v", instructions[0].op, add)
	}
	if instructions[0].left != "p1" {
		t.Errorf("Wrong instruction left side got %v expected %v", instructions[0].left, "p1")
	}
	if instructions[0].right != "p2" {
		t.Errorf("Wrong instruction right side got %v expected %v", instructions[0].right, "p2")
	}
	if instructions[0].result != res {
		t.Errorf("Result did not match returned %v and created %v", res, instructions[0].result)
	}
}

func TestParseExpressionMul(t *testing.T) {
	toConvert := "p1*p2"
	exp := parseExpression(toConvert)
	res, instructions, _ := convertAstToExpressionList(exp)
	if len(instructions) != 1 {
		t.Errorf("Number of instructions incorrect got %v expected %v", len(instructions), 1)
	}
	if instructions[0].op != multiply {
		t.Errorf("Wrong instruction got %v expected %v", instructions[0].op, multiply)
	}
	if instructions[0].left != "p1" {
		t.Errorf("Wrong instruction left side got %v expected %v", instructions[0].left, "p1")
	}
	if instructions[0].right != "p2" {
		t.Errorf("Wrong instruction right side got %v expected %v", instructions[0].right, "p2")
	}
	if instructions[0].result != res {
		t.Errorf("Result did not match returned %v and created %v", res, instructions[0].result)
	}
}

func TestParseExpressionScalar(t *testing.T) {
	toConvert := "1*p1"
	exp := parseExpression(toConvert)
	res, instructions, _ := convertAstToExpressionList(exp)
	if len(instructions) != 1 {
		t.Errorf("Number of instructions incorrect got %v expected %v", len(instructions), 1)
	}
	if instructions[0].op != scalar {
		t.Errorf("Wrong instruction got %v expected %v", instructions[0].op, scalar)
	}
	if instructions[0].left != "1" {
		t.Errorf("Wrong instruction left side got %v expected %v", instructions[0].left, "1")
	}
	if instructions[0].right != "p1" {
		t.Errorf("Wrong instruction right side got %v expected %v", instructions[0].right, "p1")
	}
	if instructions[0].result != res {
		t.Errorf("Result did not match returned %v and created %v", res, instructions[0].result)
	}
}

func TestParseExpressionCombined(t *testing.T) {
	toConvert := "(1*p1+p2)*p3"
	exp := parseExpression(toConvert)
	res, instructions, _ := convertAstToExpressionList(exp)
	if len(instructions) != 3 {
		t.Errorf("Number of instructions incorrect got %v expected %v", len(instructions), 3)
	}
	if instructions[0].op != scalar {
		t.Errorf("Wrong third instruction got %v expected %v", instructions[0].op, scalar)
	}
	if instructions[1].op != add {
		t.Errorf("Wrong third instruction got %v expected %v", instructions[1].op, add)
	}
	if instructions[2].op != multiply {
		t.Errorf("Wrong third instruction got %v expected %v", instructions[2].op, multiply)
	}

	if instructions[2].result != res {
		t.Errorf("Results did not match returned %v and created %v", res, instructions[0].result)
	}
	if instructions[0].result == instructions[1].result || instructions[0].result == instructions[2].result || instructions[1].result == instructions[2].result {
		t.Errorf("Overlapping results")
	}
	if instructions[2].result != "r3" {
		t.Errorf("Result from multiply is %v expected %v", instructions[2].result, "r3")
	}
	if instructions[1].result != "r2" {
		t.Errorf("Result from add is %v expected %v", instructions[1].result, "r2")
	}
	if instructions[0].result != "r1" {
		t.Errorf("Result from scalar is %v expected %v", instructions[0].result, "r1")
	}
}

func TestParseExpressionBigCombined(t *testing.T) {
	toConvert := "(1*p1+p2)*p3+(p1*p2+4*p1)"
	exp := parseExpression(toConvert)
	_, instructions, _ := convertAstToExpressionList(exp)
	if len(instructions) != 7 {
		t.Errorf("Number of instructions incorrect got %v expected %v", len(instructions), 7)
	}
	if instructions[0].op != scalar {
		t.Errorf("Wrong first instruction got %v expected %v", instructions[0].op, scalar)
	}
	if instructions[1].op != add {
		t.Errorf("Wrong second instruction got %v expected %v", instructions[1].op, add)
	}
	if instructions[2].op != multiply {
		t.Errorf("Wrong third instruction got %v expected %v", instructions[2].op, multiply)
	}
	if instructions[3].op != multiply {
		t.Errorf("Wrong fourth instruction got %v expected %v", instructions[3].op, multiply)
	}
	if instructions[4].op != scalar {
		t.Errorf("Wrong fifth instruction got %v expected %v", instructions[4].op, scalar)
	}
	if instructions[5].op != add {
		t.Errorf("Wrong sixth instruction got %v expected %v", instructions[5].op, add)
	}
	if instructions[6].op != add {
		t.Errorf("Wrong seventh instruction got %v expected %v", instructions[6].op, add)
	}
	//fmt.Print(instructions)
}
