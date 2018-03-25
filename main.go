// Command simple is a chromedp example demonstrating how to do a simple google
// search.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	cdptypes "github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type listing_struct struct {
	Title string
}

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	var output map[string]listing_struct
	var res string
	err = c.Run(ctxt, scrapeListings(ctxt, c, output, &res))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HI MUM %v", res)
	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range output {
		log.Printf("project %s (%s):", k, v.Title)
	}
}

func scrapeListings(ctxt context.Context, c *chromedp.CDP, output map[string]listing_struct, res *string) chromedp.Tasks {
	// force max timeout of 15 seconds for retrieving and processing the data
	var cancel func()
	ctxt, cancel = context.WithTimeout(ctxt, 25*time.Second)
	defer cancel()

	// run task list
	var url = "http://www.rightmove.co.uk/property-for-sale/find.html?locationIdentifier=USERDEFINEDAREA%5E%7B%22id%22%3A4773322%7D&minBedrooms=3&maxPrice=900000&minPrice=500000&sortType=6&propertyTypes=detached%2Csemi-detached%2Cterraced&primaryDisplayPropertyType=houses&includeSSTC=true"

	if err := c.Run(ctxt, chromedp.Navigate(url)); err != nil {
		fmt.Errorf("could not navigate to github: %v", err)
	}

	if err := c.Run(ctxt, chromedp.WaitVisible(`select.pagination-dropdown`)); err != nil {
		fmt.Errorf("could not get section: %v", err)
	}

	var listings []*cdptypes.Node
	if err := c.Run(ctxt, chromedp.Nodes(`div.l-searchResult`, &listings)); err != nil {
		fmt.Errorf("could not get projects: %v", err)
	}

	var titles []*cdptypes.Node
	if err := c.Run(ctxt, chromedp.Nodes(`div.l-searchResult h2.propertyCard-title`, &titles)); err != nil {
		fmt.Errorf("could not get projects: %v", err)
	}

	var tasks chromedp.Tasks
	for i := 0; i < (len(listings) - 1); i++ {
		if strings.Contains(listings[i].Attributes[1], "is-hidden") == false {
			var task chromedp.Tasks
			task = chromedp.Tasks{
				chromedp.Text(titles[i].FullXPath(), res, chromedp.BySearch),
			}
			tasks = append(tasks, task)
			fmt.Printf("HI MUM %T : %s : %s", res, res, *res)
		}
	}
	return tasks
}
