package formula

import (
	"net/http"
	"sync"
	"time"
)

type FormulaList struct {
	count int

	root *FormulaNode

	// key string: formula name
	uniques map[string]bool

	// lock is used to make the list thread-safe
	lock *sync.RWMutex

	httpClient *http.Client

	threads int

	iteratorChannelCount int
}

func newFormulaList(mainFormula *Formula, threads int) *FormulaList {
	list := &FormulaList{
		uniques: make(map[string]bool),
		lock:    &sync.RWMutex{},
		root:    newFormulaNode(mainFormula),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		threads: threads,
	}

	if threads == 0 {
		list.threads = 5
	}

	list.count++

	return list
}
