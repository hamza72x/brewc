package formula

import (
	"fmt"
	"sync"
)

type FormulaListV2 struct {
	count int

	root *FormulaNodeV2

	// key string: formula name
	hasDataMap map[string]bool

	// lock is used to make the list thread-safe
	lock *sync.RWMutex
}

// The client will receive a callback for each formula.
type FormulaNodeV2 struct {
	formula *Formula
	// And the callback will be triggered only if there is no child
	children []*FormulaNodeV2
}

func (list *FormulaListV2) Count() int {
	return list.count
}

func (list *FormulaListV2) Root() *FormulaNodeV2 {
	return list.root
}

func (node *FormulaNodeV2) FormulaName() string {
	return node.formula.Name
}

func NewFormulaListV2(mainFormula *Formula) *FormulaListV2 {
	list := &FormulaListV2{
		hasDataMap: make(map[string]bool),
		lock:       &sync.RWMutex{},
		root:       NewFormulaNodeV2(mainFormula),
	}

	list.hasDataMap[mainFormula.Name] = true
	list.count++

	return list
}

func NewFormulaNodeV2(formula *Formula) *FormulaNodeV2 {
	return &FormulaNodeV2{
		formula: formula,
	}
}

func (list *FormulaListV2) hasFormula(f *Formula) bool {
	list.lock.RLock()
	defer list.lock.RUnlock()
	if _, ok := list.hasDataMap[f.Name]; ok {
		return true
	}
	return false
}

func (list *FormulaListV2) AddChild(parent *FormulaNodeV2, child *FormulaNodeV2) bool {
	if list.hasFormula(child.formula) {
		return false
	}

	list.lock.Lock()
	defer list.lock.Unlock()

	parent.children = append(parent.children, child)
	list.count++
	list.hasDataMap[child.FormulaName()] = true

	// fmt.Printf("Adding %s as a child of %s\n", col.Red(child.FormulaName), col.Purple(parent.FormulaName()))
	// fmt.Printf("Childrens of %s:", col.Red(parent.formula.Name))

	// for _, child := range parent.children {
	// 	fmt.Printf(" %s", col.Purple(child.formula.Name))
	// }

	// fmt.Println("")

	return true
}

// IterateChildFirst iterates over the list in a child-first manner.
// This means that the callback will be called only if there is no child of the given node.
// Otherwise, the callback will be called after all of the children have been processed.
func (list *FormulaListV2) IterateChildFirst(threads int, fn func(*Formula)) {
	list.nodeIterator(list.root, fn)
}

func (list *FormulaListV2) nodeIterator(node *FormulaNodeV2, fn func(*Formula)) {

	// If there is no child, then we can call the callback
	if len(node.children) == 0 {
		fn(node.formula)
		return
	}

	var wg sync.WaitGroup
	var ch = make(chan int, 5)

	fmt.Printf("🛠  Resolving dependencies for %s 🛠\n", node.formula.Name)

	// If there is a child, then we need to wait for all of the children to finish
	for _, child := range node.children {
		wg.Add(1)
		go func(child *FormulaNodeV2) {
			ch <- 1
			list.nodeIterator(child, fn)
			wg.Done()
			<-ch
		}(child)
	}

	wg.Wait()

	// After all the children are done, we can call the callback
	fmt.Printf("🎉 Completed all dependencies of %s 🎉\n", node.formula.Name)
	fn(node.formula)
}
