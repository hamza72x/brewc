package formula

import "sync"

// FormulaList represents a list of formulas.
// it will be linked-list
type FormulaList struct {
	head *FormulaNode
	tail *FormulaNode
	// key string: formula name
	hasDataMap map[string]bool
	count      int
	// lock is used to make the list thread-safe
	lock sync.RWMutex
}

// FormulaNode represents a node in the linked-list
type FormulaNode struct {
	Formula *Formula
	Next    *FormulaNode
}

// NewFormulaList returns a new FormulaList instance.
func NewFormulaList() *FormulaList {
	return &FormulaList{
		hasDataMap: make(map[string]bool),
	}
}

func (list *FormulaList) HasFormula(f *Formula) bool {
	list.lock.RLock()
	defer list.lock.RUnlock()
	if _, ok := list.hasDataMap[f.Name]; ok {
		return true
	}
	return false
}

func (list *FormulaList) Count() int {
	return list.count
}

func (list *FormulaList) Add(formula *Formula) {
	list.lock.Lock()
	defer list.lock.Unlock()

	newNode := &FormulaNode{
		Formula: formula,
	}

	if list.head == nil {
		list.head = newNode
		list.tail = newNode
	} else {
		list.tail.Next = newNode
		list.tail = newNode
	}

	list.hasDataMap[formula.Name] = true
	list.count++
}

// Iterate iterate over the full linked-list
// and uses the callback function to do something
// with the data
func (list *FormulaList) Iterate(callback func(index int, formula *Formula)) {
	current := list.head
	index := 0

	for current != nil {
		callback(index, current.Formula)
		current = current.Next
		index++
	}
}
