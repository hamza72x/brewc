package formula

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hamza72x/brewc/pkg/constant"
	"github.com/hamza72x/brewc/pkg/util"
	col "github.com/hamza72x/go-color"
)

type GetFormulaListOpts struct {
	// default is false
	IncludeInstalled bool
	// default is 1
	// -1 means all the dependencies
	DependencyLevel int
}

// GetFormulaList returns a list of all the formulae
func GetFormulaList(name string, opts *GetFormulaListOpts) (*FormulaList, error) {

	if opts.DependencyLevel == 0 {
		opts.DependencyLevel = 1
	}

	mainFormula, err := GetFormulaJSON(name)

	if err != nil {
		return nil, err
	}

	list := newFormulaList(mainFormula)

	err = list.setNodesRecursive(name, list.root, opts, 1)

	if err != nil {
		return nil, err
	}

	return list, nil
}

// setNodes sets the nodes of the formula list.
func (list *FormulaList) setNodesRecursive(name string, parentNode *FormulaNode, opts *GetFormulaListOpts, level int) error {

	var wg sync.WaitGroup
	var conn = make(chan int, list.threads)

	var mainFormula, err = GetFormulaJSON(name)

	if err != nil {
		return err
	}

	// if the formula is already installed, then we don't need to install it again.
	// that's also means that all of its dependencies are already installed too.
	if mainFormula.IsInstalled() && !opts.IncludeInstalled {
		return nil
	}

	// if it's the main formula, we don't need to add it as a child of itself.
	if parentNode != list.root {
		list.AddChild(parentNode, newFormulaNode(mainFormula))
	}

	for _, dep := range mainFormula.Dependencies {
		wg.Add(1)

		go func(dep string) {
			conn <- 1

			defer wg.Done()
			defer func() { <-conn }()

			f, err := GetFormulaJSON(dep)

			if err != nil {
				fmt.Printf("%s Error getting formula %s: %v", constant.RedArrow, dep, err)
				return
			}

			if f.IsInstalled() && !opts.IncludeInstalled {
				return
			}

			newNode := newFormulaNode(f)

			if list.AddChild(parentNode, newNode) {
				fmt.Printf("%s ", col.Info(dep))
			}

			if level >= opts.DependencyLevel && opts.DependencyLevel != -1 {
				return
			}

			if err := list.setNodesRecursive(dep, newNode, opts, level+1); err != nil {
				fmt.Printf("%s Error getting nested formula %s: %v", constant.RedArrow, dep, err)
				return
			}

		}(dep)
	}

	wg.Wait()

	return nil
}

func (list *FormulaList) AddChild(parent *FormulaNode, child *FormulaNode) bool {
	list.lock.Lock()
	defer list.lock.Unlock()

	list.workStatuses[child.formula.Name] = notStarted
	parent.children = append(parent.children, child)
	list.count++

	return true
}

func (list *FormulaList) updateWorkStatus(node *FormulaNode, status formulaWorkStatus) {
	list.lock.Lock()
	defer list.lock.Unlock()

	list.workStatuses[node.formula.Name] = status
}

func (list *FormulaList) getWorkStatus(node *FormulaNode) formulaWorkStatus {
	list.lock.RLock()
	defer list.lock.RUnlock()

	return list.workStatuses[node.formula.Name]
}

func (list *FormulaList) waitForAllChildrenToBeWorked(node *FormulaNode) {
	list.lock.RLock()
	defer list.lock.RUnlock()

	for _, child := range node.children {
		for list.workStatuses[child.formula.Name] != worked {
			fmt.Printf("Waiting for %s to be worked", child.formula.Name)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// IterateChildFirst iterates over the list in a child-first manner.
// This means that the callback will be called only if there is no child of the given node.
// Otherwise, the callback will be called after all of the children have been processed.
func (list *FormulaList) IterateChildFirst(threads int, fn func(*Formula)) {
	list.childFirstIterator(list.root, fn)
}

func (list *FormulaList) childFirstIterator(node *FormulaNode, fn func(*Formula)) {

	// If there is no child, then we can call the callback
	if len(node.children) == 0 {
		if list.getWorkStatus(node) != notStarted {
			return
		}
		list.updateWorkStatus(node, working)
		fn(node.formula)
		list.updateWorkStatus(node, worked)
		return
	}

	var wg sync.WaitGroup
	var ch = make(chan int, list.threads/2)

	fmt.Printf("🛠  Resolving dependencies for %s 🛠\n", node.formula.Name)

	// If there is a child, then we need to wait for all of the children to finish
	for _, child := range node.children {
		wg.Add(1)
		go func(child *FormulaNode) {
			ch <- 1
			list.childFirstIterator(child, fn)
			wg.Done()
			<-ch
		}(child)
	}

	wg.Wait()

	// some other dependencies may triggered the same formula to be installed.
	list.waitForAllChildrenToBeWorked(node)

	// After all the children are done, we can call the callback
	fmt.Printf("🎉 Completed all dependencies of %s 🎉\n", node.formula.Name)
	fn(node.formula)
}

// IterateParentFirst iterates over the list in a parent-first manner.
func (list *FormulaList) IterateParentFirst(threads int, fn func(*Formula)) {
	list.parentFirstIterator(list.root, fn)
}

func (list *FormulaList) parentFirstIterator(node *FormulaNode, fn func(*Formula)) {
	fn(node.formula)

	var wg sync.WaitGroup
	var ch = make(chan int, list.threads/2)

	for _, child := range node.children {
		wg.Add(1)
		go func(child *FormulaNode) {
			ch <- 1
			list.parentFirstIterator(child, fn)
			wg.Done()
			<-ch
		}(child)
	}

	wg.Wait()
}

// DECIDE: should we use the github API to get the list of formulas?
// Or check local installation folder
func GetFormulaJSON(name string) (*Formula, error) {
	var f Formula

	var url = util.GetFormulaURL(name)

	resp, err := doGET(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&f)

	if err != nil {
		return nil, err
	}

	return &f, nil
}

// doGET makes a GET request to the given URL.
func doGET(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
