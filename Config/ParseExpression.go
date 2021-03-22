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

func isValid(holdings []holding) error {
	//TODO
	return nil
}

//Returns result identifier and a list of holdings which are the operations that needs to be computed in the given order
//returns convertErr if conversion failed
func convertAstToExpressionList(exp ast.Expr) (string, []holding, error) {
	resNum, holdings, err := convertAstAux(1, exp)
	isValidErr := isValid(holdings)
	if isValidErr != nil {
		return "", nil, isValidErr
	}
	return "r" + strconv.Itoa(resNum), holdings, err
}

func convertAstAux(unique int, exp ast.Expr) (int, []holding, error) {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return convertBinaryExpr(unique, exp)
	default:
		return 0, nil, fmt.Errorf("Failed to convert ast to holdings list")
	}
}

func convertBinaryExpr(unique int, exp *ast.BinaryExpr) (int, []holding, error) {
	leftUnique, leftHoldings, err := convertAstAux(unique, exp.X)
	if err != nil {
		return 0, nil, err
	}
	rightUnique, rightHoldings, err := convertAstAux(leftUnique, exp.Y)
	if err != nil {
		return 0, nil, err
	}
	unique = rightUnique + 1
	op := exp.Op
	var binaryIns holding
	switch op {
	case token.ADD:
		binaryIns = createAddition(getExprName(exp.X, leftUnique), getExprName(exp.X, rightUnique), "r"+strconv.Itoa(unique))
	case token.MUL:
		binaryIns = createMultiplication(getExprName(exp.X, leftUnique), getExprName(exp.X, rightUnique), "r"+strconv.Itoa(unique))
	}
	holdings := append(leftHoldings, binaryIns)
	holdings = append(holdings, rightHoldings...)
	return unique, holdings, nil
}

func getExprName(node ast.Expr, lastUnique int) string {
	switch t := node.(type) {
	case *ast.Ident:
		return t.Name
	default:
		return "r" + strconv.Itoa(lastUnique)
	}
}
