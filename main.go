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
	scrapeListings(ctxt, c, output)
	// err = c.Run(ctxt, scrapeListings(ctxt, c, output))
	// if err != nil {
	// 	log.Fatal(err)
	// }

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

func scrapeListings(ctxt context.Context, c *chromedp.CDP, output map[string]listing_struct) []string {
	// force max timeout of 15 seconds for retrieving and processing the data
	var cancel func()
	ctxt, cancel = context.WithTimeout(ctxt, 25*time.Second)
	defer cancel()

	// run task list
	var url = "http://www.rightmove.co.uk/property-for-sale/find.html?locationIdentifier=POSTCODE%5E377902&minBedrooms=3&maxPrice=900000&radius=0.25&sortType=6&propertyTypes=detached%2Csemi-detached%2Cterraced&primaryDisplayPropertyType=houses"

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

	var titles []cdptypes.NodeID
	if err := c.Run(ctxt, chromedp.NodeIDs(`div.l-searchResult h2.propertyCard-title`, &titles)); err != nil {
		fmt.Errorf("could not get projects: %v", err)
	}

	var title_strings []string
	var resu string
	for i := 0; i < (len(listings) - 1); i++ {
		if strings.Contains(listings[i].Attributes[1], "is-hidden") == false {
			var temp_ids []cdptypes.NodeID
			temp_ids = append(temp_ids, titles[i])

			if err := c.Run(ctxt, chromedp.Text(temp_ids, &resu, chromedp.ByNodeID)); err != nil {
				fmt.Errorf("could not get projects: %v", err)
			}
			title_strings = append(title_strings, resu)
			fmt.Printf("HI MUM %#v", title_strings)
			fmt.Printf("HI MUM %#v", resu)
		}
	}

	return title_strings

	// // var title_strings []*string
	// var tasks chromedp.Tasks
	// var res string
	// for i := 0; i < (len(listings) - 1); i++ {
	// 	if strings.Contains(listings[i].Attributes[1], "is-hidden") == false {
	// 		if err := c.Run(ctxt, chromedp.Text(titles[i], &res, chromedp.ByNodeID)); err != nil {
	// 			fmt.Errorf("could not get projects: %v", err)
	// 		}
	// 		var task chromedp.Tasks
	// 		task = chromedp.Tasks{
	// 			chromedp.Text(titles[i], &res, chromedp.ByNodeID),
	// 		}
	// 		tasks = append(tasks, task)
	// 	}
	// }
	// return tasks
}
