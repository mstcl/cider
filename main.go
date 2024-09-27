package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/netip"
	"os"

	"github.com/mstcl/cider/internal/calculation"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	inputPtr := flag.String("cidr", "0.0.0.0/0", "CIDR block")
	flag.Parse()

	input := *inputPtr

	prefix, err := netip.ParsePrefix(input)
	if err != nil {
		logger.Error("input", "error", err)
	}

	if prefix.Bits() == -1 {
		logger.Error("input", "error", "prefix bits is -1 - something went wrong")
	}

	results, err := calculation.GetResults(prefix)
	if err != nil {
		logger.Error("calculation", "error", err)
	}

	logger.Debug("calculation", "results", results)

	fmt.Println("netmask:", results.Netmask)
	fmt.Println("gateway:", results.Gateway)
	fmt.Println("broadcast:", results.Broadcast)
	fmt.Println("count:", results.Count)
}
