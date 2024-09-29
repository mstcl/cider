package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"net/netip"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mstcl/cider/internal/calculation"
	"github.com/mstcl/cider/internal/cider"
)

func Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleRequests(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleRequests(w http.ResponseWriter, r *http.Request) {
	var c *cider.Cider

	if len(r.URL.Query()) != 0 {
		addr, binary, mask, updated := parseParams(r.URL.Query())
		input := formAddr(addr)

		results, err := getResults(input)
		if err != nil {
			http.Error(w, fmt.Sprintf("bad request: %v", err), 400)
		}

		Logger.Debug("calculation", "results", results)

		c = &cider.Cider{
			Results: cider.Results{
				Netmask:   results.Netmask.String(),
				Gateway:   results.Gateway.String(),
				Broadcast: results.Broadcast.String(),
				Count:     results.Count,
			},
			ZIMask:  mask,
			Addr:    addr,
			Binary:  binary,
			Updated: updated,
		}
	} else {
		c = cider.Default()
	}

	getIndex(w, r, c)
}

func getIndex(w http.ResponseWriter, _ *http.Request, data *cider.Cider) {
	index := filepath.Join("web", "templates", "index.tmpl")
	t, err := template.ParseFiles(index)
	if err != nil {
		http.Error(w, "error parsing templates", 500)
	}

	t.ExecuteTemplate(w, "index", data)
}

func formAddr(a cider.Addr) string {
	var addr strings.Builder

	addr.WriteString(strings.Join([]string{a.A, a.B, a.C, a.D}, "."))
	addr.WriteString("/")
	addr.WriteString(a.E)

	return addr.String()
}

func parseParams(v url.Values) (cider.Addr, cider.Bits, cider.ZIMask, cider.Updated) {
	binarySlice := make([]string, 0, 32)
	addrSlice := make([]string, 4)
	addrSlice[0] = v.Get("a")
	addrSlice[1] = v.Get("b")
	addrSlice[2] = v.Get("c")
	addrSlice[3] = v.Get("d")

	changedSlice := make([]bool, 4)

	for i, v := range addrSlice {
		n, err := strconv.Atoi(v)
		if err != nil {
			addrSlice[i] = "0"
			changedSlice[i] = true
		}

		if n < 0 {
			addrSlice[i] = "0"
			changedSlice[i] = true
		} else if n > 255 {
			addrSlice[i] = "255"
			changedSlice[i] = true
		} else {
			addrSlice[i] = strconv.Itoa(n)
		}

		binarySlice = append(binarySlice, strings.Split(
			fmt.Sprintf("%08b", n), "")...)
	}

	changedPrefix := false
	prefix := v.Get("e")

	n, err := strconv.Atoi(prefix)
	if err != nil {
		prefix = "0"
		changedPrefix = true
	}

	if n < 0 {
		prefix = "0"
		changedPrefix = true
	} else if n > 32 {
		prefix = "32"
		changedPrefix = true
	} else {
		prefix = strconv.Itoa(n)
	}

	n, _ = strconv.Atoi(prefix)

	return cider.Addr{
			A: addrSlice[0],
			B: addrSlice[1],
			C: addrSlice[2],
			D: addrSlice[3],
			E: prefix,
		}, binarySlice,
		cider.ZIMask(n - 1),
		cider.Updated{
			A: changedSlice[0],
			B: changedSlice[1],
			C: changedSlice[2],
			D: changedSlice[3],
			E: changedPrefix,
		}
}

func getResults(input string) (*calculation.Results, error) {
	prefix, err := netip.ParsePrefix(input)
	if err != nil {
		return nil, fmt.Errorf("input: %v", err)
	}

	if prefix.Bits() == -1 {
		return nil, fmt.Errorf("input: prefix bits is -1 - something went wrong")
	}

	results, err := calculation.GetResults(prefix)
	if err != nil {
		return nil, fmt.Errorf("calculation: %v", err)
	}

	return results, nil
}
