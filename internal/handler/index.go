package handler

import (
	"fmt"
	"net/http"
	"net/netip"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mstcl/cider/internal/calculation"
	"github.com/mstcl/cider/internal/cider"
)

type Request struct {
	A string `query:"a"`
	B string `query:"b"`
	C string `query:"c"`
	D string `query:"d"`
	E string `query:"e"`
}

func Index(e echo.Context) error {
	var c *cider.Cider

	if len(e.QueryParams()) != 0 {
		var r Request
		err := e.Bind(&r)
		if err != nil {
			return e.String(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		}

		addr, binary, mask, updated := parseParams(&r)
		input := formAddr(addr)

		results, err := getResults(input)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("calculation: %v", err))
		}

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

	if err := e.Render(http.StatusOK, "index", c); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("templating: %v", err))
	}

	return nil
}

func formAddr(a cider.Addr) string {
	var addr strings.Builder

	addr.WriteString(strings.Join([]string{a.A, a.B, a.C, a.D}, "."))
	addr.WriteString("/")
	addr.WriteString(a.E)

	return addr.String()
}

func parseParams(r *Request) (cider.Addr, cider.Bits, cider.ZIMask, cider.Updated) {
	binarySlice := make([]string, 0, 32)
	addrSlice := make([]string, 4)
	addrSlice[0] = r.A
	addrSlice[1] = r.B
	addrSlice[2] = r.C
	addrSlice[3] = r.D

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
	prefix := r.E

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
