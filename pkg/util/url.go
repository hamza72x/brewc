package util

import "fmt"

func GetFormulaURL(name string) string {
	return fmt.Sprintf("https://formulae.brew.sh/api/formula/%s.json", name)
}
