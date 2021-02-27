package pimclient

import "encoding/json"

// Page is a page
type Page struct {
	res pageResponse
}

type pageResponse struct {
	Links map[string]struct {
		URL string `json:"href"`
	} `json:"_links"`

	Items struct {
		Items []map[string]interface{} `json:"items"`
	} `json:"_embedded"`
}

// IsFirst checks if this is the first page
func (page *Page) IsFirst() bool {
	return page.linkMissing("previous")
}

// IsLast checks if this is the last page
func (page *Page) IsLast() bool {
	return page.linkMissing("next")
}

// Next fetches next page
func (page *Page) Next(client *PIMClient) (Page, error) {
	return page.getNewPage("next", client)
}

// Previous fetches previous page
func (page *Page) Previous(client *PIMClient) (Page, error) {
	return page.getNewPage("previous", client)
}

// First fetches the first page.
func (page *Page) First(client *PIMClient) (Page, error) {
	return page.getNewPage("first", client)
}

// ItemCount returns the amount of items in page
func (page *Page) ItemCount() int {
	return len(page.res.Items.Items)
}

// At takes an item at given index and uses it to initialize given PIM entity.
func (page *Page) At(i int, entity interface{}) {
	raw, _ := json.Marshal(page.res.Items.Items[i])
	json.Unmarshal(raw, entity)
}

func (page *Page) linkMissing(link string) bool {
	_, found := page.res.Links[link]
	return !found
}

func (page *Page) getNewPage(link string, client *PIMClient) (Page, error) {
	var p Page

	req := client.newGetRequest(page.res.Links[link].URL)
	err := client.send(req, &p.res)

	return p, err
}
