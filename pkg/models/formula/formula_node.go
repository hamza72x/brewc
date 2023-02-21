package formula

type FormulaNode struct {
	formula  *Formula
	children []*FormulaNode
}

func newFormulaNode(formula *Formula) *FormulaNode {
	return &FormulaNode{
		formula: formula,
	}
}

func (list *FormulaList) Count() int {
	return list.count
}
