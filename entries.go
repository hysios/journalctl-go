package journalctl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	_GET  = "GET"
	_POST = "POST"
)

type Options struct {
	Cursor string
	Skip   int
	Limit  int
	Follow bool
	Key    string
	Boot   bool
}

// Entries will return a list of entries.
// If a non-nil filter is passed, any non-empty entries will be used
// to filter the query.
// For example, Entry{SYSTEMD_UNIT: "my-service.service"} will return only entries
// matching this unit name.
// Note that the systemd journal does not allow filtering on all journal field names.
func (c *Client) Entries(filter *Entry, opt *Options) (entries chan Entry, err error) {
	entries = make(chan Entry)

	values, header := c.buildOptions(opt)

	if filter != nil {
		// TODO allow filtering of other fields
		if unit := filter.SYSTEMDUNIT; unit != "" {
			values.Set("_SYSTEMD_UNIT", unit)
		}
	}

	const endpoint = "/entries"
	host := c.host()

	req, err := http.NewRequest(_GET, host+"/entries?"+values.Encode(), nil)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")
	for key, value := range header {
		req.Header[key] = value
	}
	// mergo.Merge(&req.Header, header)

	resp, err := c.http.Do(req)
	if err != nil {
		return
	}

	go func() error {
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			var entry Entry
			bts := scanner.Bytes()
			if len(bts) == 0 {
				continue
			}
			err = json.Unmarshal(bts, &entry)
			if err != nil {
				log.Print(err)
				close(entries)
				return err
			}
			entries <- entry
		}

		if err := scanner.Err(); err != nil {
			log.Print(err)
			close(entries)
			return err
		}
		close(entries)
		return nil
	}()

	return entries, nil
}

func (c *Client) buildOptions(opt *Options) (url.Values, http.Header) {
	var (
		values = url.Values{}
		header = make(http.Header)
	)

	if opt != nil {

		header.Add("Range", fmt.Sprintf("entries=%s:%s:%s", opt.Cursor, Itos(opt.Skip), Itos(opt.Limit)))

		if follow := opt.Follow; follow {
			values.Set("follow", "1")
		}
		if key := opt.Key; key != "" {
			values.Set("key", key)
		}
		if boot := opt.Boot; boot {
			values.Set("boot", "1")
		}
	}
	return values, header
}

func Itos(i int) string {
	if i == 0 {
		return ""
	}
	return fmt.Sprintf("%d", i)
}
