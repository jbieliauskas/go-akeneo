package pimclient

import (
	"encoding/json"
	"io"
)

// Page is a page
type Page struct {
	Len   int
	Items interface{}

	next        string
	decodeItems decodePageItemsFunc
}

type decodePageItemsFunc func(*pageItemDecoder) interface{}

type pageItemDecoder struct {
	d   *json.Decoder
	err error
}

func (p *Page) IsLast() bool {
	return p.next == ""
}

func (p *Page) Next(c PIMClient) (Page, error) {
	if p.IsLast() {
		return *p, nil
	}

	return c.getPage(p.next, p.decodeItems)
}

func newPage(body io.ReadCloser, decodeItems decodePageItemsFunc) (Page, error) {
	defer body.Close()

	var p Page
	var t json.Token
	var err error = nil
	d := json.NewDecoder(body)
	links, items := false, false

	for !(links && items) && err == nil {
		t, err = d.Token()

		if err == nil {
			if t == "_links" {
				err = decodePageLinks(d, &p)
				links = true
			} else if t == "items" {
				err = decodePageItems(d, &p, decodeItems)
				items = true
			}
		}
	}

	p.decodeItems = decodeItems

	return p, err
}

func decodePageLinks(d *json.Decoder, p *Page) error {
	type link struct {
		Href string `json:"href"`
	}
	var links struct {
		First *link `json:"first"`
		Prev  *link `json:"previous"`
		Next  *link `json:"next"`
	}

	err := d.Decode(&links)
	if err != nil {
		return err
	}

	if links.Next != nil {
		p.next = links.Next.Href
	}

	return nil
}

func decodePageItems(d *json.Decoder, p *Page, decodeItems decodePageItemsFunc) error {
	var dec pageItemDecoder
	dec.d = d
	_, dec.err = d.Token()

	p.Items = decodeItems(&dec)

	return dec.err
}

func (dec *pageItemDecoder) more() bool {
	return dec.err == nil && dec.d.More()
}

func (dec *pageItemDecoder) decode(v interface{}) {
	dec.err = dec.d.Decode(v)
}
