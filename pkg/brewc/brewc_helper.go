package brewc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/hamza72x/brewc/pkg/constant"
	"github.com/hamza72x/brewc/pkg/models/formula"
	"github.com/hamza72x/brewc/pkg/util"
	col "github.com/hamza72x/go-color"
)

// GetAllFormulas returns all of the formulas.
// this skips the main formula.
func (b *BrewC) GetFormulaListV2(name string, includeAll bool) (*formula.FormulaListV2, error) {
	mainFormula, err := b.getFormulaJSON(name)

	if err != nil {
		return nil, err
	}

	list := formula.NewFormulaListV2(mainFormula)

	err = b.setAllFormulasRecursiveV2(name, list, list.Root(), includeAll)

	if err != nil {
		return nil, err
	}

	return list, nil
}

// setAllFormulasRecursive returns all of the formulas recursively.
// this skips the main formula.
func (b *BrewC) setAllFormulasRecursiveV2(name string, list *formula.FormulaListV2, parentNode *formula.FormulaNodeV2, includeAll bool) error {

	var wg sync.WaitGroup
	var conn = make(chan int, b.threads)

	var mainFormula, err = b.getFormulaJSON(name)

	if err != nil {
		return err
	}

	// if the formula is already installed, then we don't need to install it again.
	// that's also means that all of its dependencies are already installed too.
	if mainFormula.IsInstalled() && !includeAll {
		return nil
	}

	list.AddChild(parentNode, formula.NewFormulaNodeV2(mainFormula))

	for _, dep := range mainFormula.Dependencies {
		wg.Add(1)

		go func(dep string) {
			conn <- 1

			defer wg.Done()
			defer func() { <-conn }()

			f, err := b.getFormulaJSON(dep)

			if err != nil {
				fmt.Printf("%s Error getting formula %s: %v", constant.RedArrow, dep, err)
				return
			}

			if f.IsInstalled() && !includeAll {
				return
			}

			newNode := formula.NewFormulaNodeV2(f)

			if list.AddChild(parentNode, newNode) {
				fmt.Printf("%s ", col.Info(dep))
			}

			if err := b.setAllFormulasRecursiveV2(dep, list, newNode, includeAll); err != nil {
				fmt.Printf("%s Error getting nested formula %s: %v", constant.RedArrow, dep, err)
				return
			}

		}(dep)
	}

	wg.Wait()

	return nil
}

// DECIDE: should we use the github API to get the list of formulas?
// Or check local installation folder
func (b *BrewC) getFormulaJSON(name string) (*formula.Formula, error) {
	var f formula.Formula

	var url = util.GetFormulaURL(name)

	resp, err := b.doGET(url)

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
func (b *BrewC) doGET(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
