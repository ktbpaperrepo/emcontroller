package dependency

import (
	"fmt"
	"io"
	"net/http"
)

var DepUrls []string

func Exec() error {
	if err := callAllUrls(); err != nil {
		return fmt.Errorf("Call dependent urls error: %w", err)
	}
	return nil
}

func callAllUrls() error {
	for _, url := range DepUrls {
		fmt.Println("call url:", url)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("url: %s, make request error: %w", url, err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("url: %s, do request error: %w", url, err)
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode > 299 {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("url: %s, res.statusCode is %d, read res.Body error: %w", url, res.StatusCode, err)
			} else {
				return fmt.Errorf("url: %s, res.statusCode is %d, res.Body is %s", url, res.StatusCode, string(body))
			}
		}
	}
	return nil
}
