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

//Returns result identifier and a list of holdings which are the operations that needs to be computed in the given order
//returns convertErr if conversion failed
func convertAstToExpressionList(exp ast.Expr) (string, []holding, error) {
	resNum, holdings, name, err := convertAstAux(1, exp)
	return name, holdings, err
}

func convertAstAux(unique int, exp ast.Expr) (int, []holding, string, error) {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return convertBinaryExpr(unique, exp)
	default:
		fmt.Printf("Failed to convert ast to holdings list")
		return 0, nil
	}
	return 0, nil
}


func convertBinaryExpr(unique int, exp *ast.BinaryExpr) (int, []holding, string, error) {
	op := exp.Op
		switch op {
		case token.ADD:
			unique, leftHoldings, err := convertAstAux(unique, exp.X)
			if err != nil {
				return 0, nil, err
			}
			unique, rightHoldings, name, err := convertAstAux(unique, exp.Y)
			addition := createAddition(x.Name, "r"+strconv.Itoa(num))
			switch x := exp.X.(type) {
			case *ast.Ident:
				switch y := exp.Y.(type) {
				case *ast.Ident:
					addition = createAddition(x.Name, y.Name)
				default:
					//y is a general expression so use r[num]
					num, rightHoldings, err := 
					unique = num + 1
					
					return unique, append([]holding{addition}, rightHoldings...), err
				}
			default:
				//X is not Ident
				switch y := exp.Y.(type) {
				case *ast.Ident:
					num, leftHoldings, err := convertAstAux(unique, exp.X)
					unique = num + 1
					addition = createAddition("r"+strconv.Itoa(num), y.Name)
					return unique, append(leftHoldings, addition), err
				default:

					return 0, nil, fmt.Errorf("Failed to convert ast to holdings list in addition")
				}

			}
		case token.MUL:
			
			return 0, nil
		}
}
