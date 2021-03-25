package pimclient

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
)

type testPageItem struct {
	Code string `json:"code"`
	Data int    `json:"data"`
}

func TestParsing(t *testing.T) {
	res := `{
		"_links": {
			"self": {
				"href": "https://pim.e-store.net/api/rest/v1/resources?page=2&limt=2"
			},
			"first": {
				"href": "https://pim.e-store.net/api/rest/v1/resources?page=1&limt=2"
			},
			"previous": {
				"href": "https://pim.e-store.net/api/rest/v1/resources?page=1&limt=2"
			},
			"next": {
				"href": "https://pim.e-store.net/api/rest/v1/resources?page=3&limt=2"
			}
		},
		"current_page": 2,
		"_embedded": {
			"items": [
				{
					"_links": {
						"self": {
							"href": "https://pim.e-store.net/api/rest/v1/resources/one"
						}
					},
					"code": "one",
					"data": 17
				},
				{
					"_links": {
						"self": {
							"href": "https://pim.e-store.net/api/rest/v1/resources/two"
						}
					},
					"code": "two",
					"data": 18
				}
			]
		}
	}`
	body := io.NopCloser(strings.NewReader(res))

	p, err := newPage(body, func(d *pageItemDecoder) interface{} {
		items := []testPageItem{}

		for d.more() {
			var item testPageItem
			d.decode(&item)
			items = append(items, item)
		}

		return items
	})

	if err != nil {
		t.Fatal("newPage() failed with an error: ", err)
	}

	if p.next != "https://pim.e-store.net/api/rest/v1/resources?page=3&limt=2" {
		t.Error("Expected no next page link but received: ", p.next)
	}

	items := p.Items.([]testPageItem)

	if len(items) != 2 {
		t.Fatalf("Expected 2 items but received %d.", len(items))
	}

	one, _ := json.Marshal(items[0])
	if string(one) != `{"code":"one","data":17}` {
		t.Errorf("First item not as expected:\n\n%s", string(one))
	}

	two, _ := json.Marshal(items[1])
	if string(two) != `{"code":"two","data":18}` {
		t.Errorf("Second item not as expected:\n\n%s", string(two))
	}
}
