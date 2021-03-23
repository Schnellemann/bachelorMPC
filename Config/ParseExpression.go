package config

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func parseExpression(exp string) ast.Expr {
	aexp, err := parser.ParseExpr(exp)
	if err != nil {
		fmt.Printf("parsing failed: %s\n", err)
		return nil
	}
	return aexp

}

func isValid(instructions []instruction) error {
	//TODO
	return nil
}

//Returns result identifier and a list of instructions which are the operations that needs to be computed in the given order
//returns convertErr if conversion failed
func convertAstToExpressionList(exp ast.Expr) (string, []instruction, error) {
	resNum, instructions, err := convertAstAux(1, exp)
	isValidErr := isValid(instructions)
	if isValidErr != nil {
		return "", nil, isValidErr
	}
	return "r" + strconv.Itoa(resNum), instructions, err
}

func convertAstAux(unique int, exp ast.Expr) (int, []instruction, error) {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return convertBinaryExpr(unique, exp)
	case *ast.BasicLit, *ast.Ident:
		return 0, nil, nil
	case *ast.ParenExpr:
		return convertAstAux(unique, exp.X)
	default:
		return 0, nil, fmt.Errorf("Failed to convert ast to instructions list")
	}
}

func convertBinaryExpr(unique int, exp *ast.BinaryExpr) (int, []instruction, error) {
	leftUnique, leftInstructions, err := convertAstAux(unique, exp.X)
	if err != nil {
		return 0, nil, err
	}
	rightUnique, rightInstructions, err := convertAstAux(leftUnique, exp.Y)
	if err != nil {
		return 0, nil, err
	}
	unique = rightUnique + 1
	op := exp.Op
	var binaryIns instruction
	switch op {
	case token.ADD:
		binaryIns = createAddition(getExprName(exp.X, leftUnique), getExprName(exp.Y, rightUnique), "r"+strconv.Itoa(unique))
	case token.MUL:
		//Need to see if its scale by constant or multiply
		switch x := exp.X.(type) {
		case *ast.BasicLit:
			binaryIns = createScalar(x.Value, getExprName(exp.Y, rightUnique), "r"+strconv.Itoa(unique))
		default:
			switch y := exp.Y.(type) {
			case *ast.BasicLit:
				binaryIns = createScalar(y.Value, getExprName(exp.X, rightUnique), "r"+strconv.Itoa(unique))
			default:
				//None were basic lits so mul
				binaryIns = createMultiplication(getExprName(exp.X, leftUnique), getExprName(exp.Y, rightUnique), "r"+strconv.Itoa(unique))
			}
		}

	}
	instructions := append(leftInstructions, binaryIns)
	instructions = append(instructions, rightInstructions...)
	return unique, instructions, nil
}

func getExprName(node ast.Expr, lastUnique int) string {
	switch t := node.(type) {
	case *ast.Ident:
		return t.Name
	default:
		return "r" + strconv.Itoa(lastUnique)
	}
}
