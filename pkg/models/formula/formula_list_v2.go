package formula

import (
	"fmt"
	"sync"

	col "github.com/hamza72x/go-color"
)

type FormulaListV2 struct {
	count int

	root *FormulaNodeV2

	// if parent is ready to be installed or not
	readyParents []string

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

func NewFormulaListV2(mainFormula *Formula) *FormulaListV2 {
	return &FormulaListV2{
		hasDataMap: make(map[string]bool),
		lock:       &sync.RWMutex{},
		root:       NewFormulaNodeV2(mainFormula),
	}
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

func (list *FormulaListV2) AddChild(parent *FormulaNodeV2, child *Formula) bool {
	if list.hasFormula(child) {
		return false
	}

	list.lock.Lock()
	defer list.lock.Unlock()

	newNode := &FormulaNodeV2{
		formula: child,
	}

	parent.children = append(parent.children, newNode)
	list.count++
	list.hasDataMap[child.Name] = true

	return true
}

func (list *FormulaListV2) removeChild(parent *FormulaNodeV2, child *Formula) {
	list.lock.Lock()
	defer list.lock.Unlock()

	for i, c := range parent.children {
		if c.formula.Name == child.Name {
			// remove the child from the parent
			parent.children = append(parent.children[:i], parent.children[i+1:]...)
			break
		}
	}

	list.count--
	delete(list.hasDataMap, child.Name)
}

func (list *FormulaListV2) Iterate(threads int, fn func(*Formula)) {
	list.lock.RLock()
	defer list.lock.RUnlock()

	wg := &sync.WaitGroup{}
	ch := make(chan int, threads)

	list.nodeIterator(list.root, wg, ch, fn)

	wg.Wait()
}

func (list *FormulaListV2) nodeIterator(node *FormulaNodeV2, wg *sync.WaitGroup, ch chan int, fn func(*Formula)) {
	wg.Add(1)
	ch <- 1

	defer wg.Done()
	defer func() { <-ch }()

	fmt.Printf("Childrens of %s: ", col.Red(node.formula.Name))

	for _, child := range node.children {
		fmt.Printf("%s ", col.Red(child.formula.Name))
	}

	fmt.Println("")

	for _, child := range node.children {
		list.nodeIterator(child, wg, ch, fn)
	}

	// If there is no child, then we can call the callback
	// afterwards remove that child from the parent
	if len(node.children) == 0 {
		// fmt.Println(col.Red("???"), "interator", node.formula.Name)
		// fn(node.formula)
		// list.removeChild(node, node.formula)
		return
	}

	// var wg2 sync.WaitGroup
	// var ch2 = make(chan int, len(node.children))

	// If there is a child, then we need to wait for all of the children to finish
	// for _, child := range node.children {
	// wg2.Add(1)
	// fmt.Println(col.Black("hmm"), "queued", child.formula.Name, "parent", node.formula.Name)
	// go func(child *FormulaNodeV2) {
	// 	ch2 <- 1
	// 	list.nodeIterator(child, wg, ch, fn)
	// 	wg2.Done()
	// 	<-ch2
	// }(child)
	// }

	// wg2.Wait()

	// After all the children are done, we can call the callback
	// fn(node.formula)
	// list.removeChild(node, node.formula)
}
