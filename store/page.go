package store

import (
	"fmt"
	"math"

	"github.com/ttacon/chalk"
)

// Page represents a page in the document
type Page struct {
	Page  int
	Limit float64
	Count float64
	Skip  int
	Max   int
}

// Init does some necessary calculations
func (p *Page) Init() {
	Limit := int(p.Limit)
	p.Skip = (p.Page * Limit) - Limit
	p.Max = int(math.Ceil(p.Count / p.Limit))
}

// NextPage changes the page to imitate going to the next page
func (p *Page) NextPage() {
	p.Page++
	p.Init()
}

// PrevPage changes the page to imitate going to the prev page
func (p *Page) PrevPage() {
	if p.Page == 1 {
		return
	}
	p.Page--
	p.Init()
}

// End check if page is it's end
func (p *Page) End() bool {
	return p.Page >= p.Max
}

// More checks if there is more data to display
func (p *Page) More() bool {
	return p.Page < p.Max
}

func (p *Page) String() string {
	return fmt.Sprintf("Page: %d\t Limit: %f\tCount:%f\t Skip: %d Max: %d More: %t End: %t\n", p.Page, p.Limit, p.Count, p.Skip, p.Max, p.More(), p.End())
}

// Pretty prettifies the Page struct
func (p *Page) Pretty() string {
	return fmt.Sprintf("%sPage %d of %d%s", chalk.Blue, p.Page, p.Max, chalk.Reset)
}

// Commands returns the commands to be entered
func (p *Page) Commands() string {
	return fmt.Sprintf("Please enter \nlast: to view last page\nnext: to view next page\nprev: to view previous page\nexit: to exit the shell")
}
