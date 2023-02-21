package formula

import (
	"net/http"
	"sync"
	"time"
)

type formulaWorkStatus int

const (
	notStarted formulaWorkStatus = iota
	working
	worked
)

type FormulaList struct {
	count int

	root *FormulaNode

	// key string: formula name
	workStatuses map[string]formulaWorkStatus

	// lock is used to make the list thread-safe
	lock *sync.RWMutex

	httpClient *http.Client

	threads int
}

func newFormulaList(mainFormula *Formula) *FormulaList {
	list := &FormulaList{
		workStatuses: make(map[string]formulaWorkStatus),
		lock:         &sync.RWMutex{},
		root:         newFormulaNode(mainFormula),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	list.workStatuses[mainFormula.Name] = notStarted
	list.count++

	return list
}
