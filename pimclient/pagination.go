package pimclient

import (
	"encoding/json"
	"io"
)

// Page is a page
type Page struct {
	Len   int
	Items interface{}

	first, prev, next string
}

type pageItemDecoder struct {
	d   *json.Decoder
	err error
}

func newPage(body io.ReadCloser, decodeItems func(d pageItemDecoder) interface{}) (Page, error) {
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

	if links.First != nil {
		p.first = links.First.Href
	}
	if links.Prev != nil {
		p.prev = links.Prev.Href
	}
	if links.Next != nil {
		p.next = links.Next.Href
	}

	return nil
}

func decodePageItems(d *json.Decoder, p *Page, decodeItems func(d pageItemDecoder) interface{}) error {
	var dec pageItemDecoder
	dec.d = d
	_, dec.err = d.Token()

	p.Items = decodeItems(dec)

	return dec.err
}

func (dec *pageItemDecoder) more() bool {
	return dec.err == nil && dec.d.More()
}

func (dec *pageItemDecoder) decode(v interface{}) {
	dec.err = dec.d.Decode(v)
}
