package parsing

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

type parse struct {
	multCounter int
}

func ParseExpression(exp string) ast.Expr {
	aexp, err := parser.ParseExpr(exp)
	if err != nil {
		fmt.Printf("parsing failed: %s\n", err)
		return nil
	}
	return aexp

}

func (tree *InstructionTree) CountMults() int {
	if tree == nil {
		return 0
	} else {
		switch tree.Instruction.(type) {
		case *MultInstruction:
			return 1 + tree.Left.CountMults() + tree.Right.CountMults()
		default:
			return tree.Left.CountMults() + tree.Right.CountMults()
		}

	}
}

func ConvertAstToTree(exp ast.Expr) (*InstructionTree, error) {
	p := &parse{1}
	_, tree, err := p.convertAstAux(0, exp)
	return tree, err
}

func (tree *InstructionTree) GetResultName() string {
	switch node := tree.Instruction.(type) {
	case *AddInstruction:
		return node.Result
	case *ScalarInstruction:
		return node.Result
	case *MultInstruction:
		return node.Result
	default:
		return "Error"
	}
}

func (p *parse) convertAstAux(unique int, exp ast.Expr) (int, *InstructionTree, error) {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return p.convertBinaryExpr(unique, exp)
	case *ast.BasicLit, *ast.Ident:
		return unique, nil, nil
	case *ast.ParenExpr:
		return p.convertAstAux(unique, exp.X)
	default:
		return 0, nil, fmt.Errorf("Failed to convert ast to InstructionTree - Unsupported ast-token")
	}
}

func (p *parse) convertBinaryExpr(unique int, exp *ast.BinaryExpr) (int, *InstructionTree, error) {
	leftUnique, leftSubTree, err := p.convertAstAux(unique, exp.X)
	if err != nil {
		return 0, nil, err
	}
	rightUnique, rightSubTree, err := p.convertAstAux(leftUnique, exp.Y)
	if err != nil {
		return 0, nil, err
	}
	unique = rightUnique + 1
	op := exp.Op
	var binaryIns Instruction
	switch op {
	case token.ADD:
		binaryIns = createAddition(getExprName(exp.X, leftUnique), getExprName(exp.Y, rightUnique), "r"+strconv.Itoa(unique))
	case token.MUL:
		//Need to see if its scale by constant or multiply
		switch x := exp.X.(type) {
		case *ast.BasicLit:
			scalar, err := strconv.Atoi(x.Value)
			if err != nil {
				return 0, nil, err
			}
			binaryIns = createScalar(scalar, getExprName(exp.Y, rightUnique), "r"+strconv.Itoa(unique))
		default:
			switch y := exp.Y.(type) {
			case *ast.BasicLit:
				scalar, err := strconv.Atoi(y.Value)
				if err != nil {
					return 0, nil, err
				}
				binaryIns = createScalar(scalar, getExprName(exp.X, rightUnique), "r"+strconv.Itoa(unique))
			default:
				//None were basic lits so mul
				binaryIns = createMultiplication(getExprName(exp.X, leftUnique), getExprName(exp.Y, rightUnique), "r"+strconv.Itoa(unique), p.multCounter)
				p.multCounter += 1
			}
		}
	}
	tree := &InstructionTree{Left: leftSubTree, Right: rightSubTree, Instruction: binaryIns}
	return unique, tree, nil
}

func getExprName(node ast.Expr, lastUnique int) string {
	switch t := node.(type) {
	case *ast.Ident:
		return t.Name
	default:
		return "r" + strconv.Itoa(lastUnique)
	}
}
