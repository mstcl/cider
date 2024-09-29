package cider

type Cider struct {
	Binary  Bits
	Addr    Addr
	Results Results
	Updated Updated
	ZIMask  ZIMask
}

type (
	ZIMask int
	Bits   []string
)

type Results struct {
	Netmask   string
	Gateway   string
	Broadcast string
	Count     int
}

type Updated struct {
	A bool
	B bool
	C bool
	D bool
	E bool
}

type Addr struct {
	A string
	B string
	C string
	D string
	E string
}

func Default() *Cider {
	return &Cider{
		Results: Results{
			Netmask:   "0.0.0.0",
			Gateway:   "0.0.0.0",
			Broadcast: "0.0.0.0",
			Count:     4294967296,
		},
		Addr: Addr{
			A: "0",
			B: "0",
			C: "0",
			D: "0",
			E: "0",
		},
		Updated: Updated{
			A: false,
			B: false,
			C: false,
			D: false,
			E: false,
		},
	}
}
