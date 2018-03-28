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
	Title       string
	Description string
	Address     string
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

	var output []listing_struct
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

func scrapeListings(ctxt context.Context, c *chromedp.CDP, output []listing_struct) []string {
	// force max timeout of 15 seconds for retrieving and processing the data
	var cancel func()
	ctxt, cancel = context.WithTimeout(ctxt, 25*time.Second)
	defer cancel()

	// run task list
	var url = "http://www.rightmove.co.uk/property-for-sale/find.html?locationIdentifier=USERDEFINEDAREA%5E%7B%22id%22%3A4773322%7D&minBedrooms=3&maxPrice=900000&sortType=6&propertyTypes=detached%2Csemi-detached%2Cterraced&primaryDisplayPropertyType=houses"

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

	var addresses []cdptypes.NodeID
	if err := c.Run(ctxt, chromedp.NodeIDs(`div.l-searchResult address.propertyCard-address`, &addresses)); err != nil {
		fmt.Errorf("could not get projects: %v", err)
	}

	var descriptions []cdptypes.NodeID
	if err := c.Run(ctxt, chromedp.NodeIDs(`div.l-searchResult .propertyCard-description`, &descriptions)); err != nil {
		fmt.Errorf("could not get projects: %v", err)
	}

	fmt.Println("%v", len(listings))
	fmt.Println("%v", len(titles))
	fmt.Println("%v", len(addresses))
	fmt.Println("%v", len(descriptions))

	var title_strings []string
	var address_strings []string
	var description_strings []string

	var title_resu string
	var address_resu string
	var description_resu string

	for i := 0; i < (len(listings) - 1); i++ {
		if strings.Contains(listings[i].Attributes[1], "is-hidden") == false {
			var title_temp_ids []cdptypes.NodeID
			var address_temp_ids []cdptypes.NodeID
			var description_temp_ids []cdptypes.NodeID

			title_temp_ids = append(title_temp_ids, titles[i])
			address_temp_ids = append(address_temp_ids, addresses[i])
			description_temp_ids = append(description_temp_ids, descriptions[i])

			if err := c.Run(ctxt, chromedp.Text(title_temp_ids, &title_resu, chromedp.ByNodeID)); err != nil {
				fmt.Errorf("could not get projects: %v", err)
			}
			if err := c.Run(ctxt, chromedp.Text(address_temp_ids, &address_resu, chromedp.ByNodeID)); err != nil {
				fmt.Errorf("could not get projects: %v", err)
			}
			if err := c.Run(ctxt, chromedp.Text(description_temp_ids, &description_resu, chromedp.ByNodeID)); err != nil {
				fmt.Errorf("could not get projects: %v", err)
			}
			title_strings = append(title_strings, strings.TrimSpace(title_resu))
			address_strings = append(address_strings, strings.TrimSpace(address_resu))
			description_strings = append(description_strings, strings.TrimSpace(description_resu))

			var new_struct listing_struct
			new_struct = listing_struct{Title: strings.TrimSpace(title_resu), Address: strings.TrimSpace(address_resu), Description: strings.TrimSpace(description_resu)}
			output = append(output, new_struct)
		}
	}

	fmt.Printf("HI MUM %#v", output)
	// fmt.Printf("HI MUM %#v", title_strings)
	// fmt.Printf("HI MUM %#v", address_strings)
	// fmt.Printf("HI MUM %#v", description_strings)

	return title_strings
}
