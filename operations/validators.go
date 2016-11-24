package operations

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var (
	loadTLDsOnce sync.Once
	tlds         []byte
)

type Validator func(string) error

func validateNothing(input string) error {
	return nil
}

func validateCustomDomain(input string) error {
	if !strings.Contains(input, ".") || len(input) <= 4 {
		return fmt.Errorf("%s is not a valid domain name", input)
	}

	parts := strings.SplitN(input, ".", -1)
	tldName := parts[len(parts)-1]

	if tldName == "" {
		return fmt.Errorf("%s is not a valid domain name", input)
	}

	loadTLDsOnce.Do(func() {
		resp, err := http.Get("https://publicsuffix.org/list/effective_tld_names.dat")
		if err != nil {
			return
		}

		tlds, _ = ioutil.ReadAll(resp.Body)
	})

	if tlds == nil {
		return fmt.Errorf("unable to load TLDs")
	}

	scanner := bufio.NewScanner(bytes.NewReader(tlds))
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 || text[0] == '/' || text[0] == ' ' {
			continue
		}
		if tldName == text {
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return fmt.Errorf(".%s is not a registered TLD", tldName)
}
