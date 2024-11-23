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
			Netmask:   "255.255.255.254",
			Gateway:   "192.168.1.0",
			Broadcast: "192.168.1.1",
			Count:     2,
		},
		Addr: Addr{
			A: "192",
			B: "168",
			C: "1",
			D: "0",
			E: "31",
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
