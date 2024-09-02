package marshal

import (
	"encoding"
	"fmt"
	"net/netip"

	"github.com/trwk76/goweb/openapi/spec"
)

func ParseIPAddr(str string) (IPAddr, error) {
	tmp, err := netip.ParseAddr(str)
	if err != nil {
		return IPAddr{}, fmt.Errorf("not a valid IP address (ex: 10.147.0.0)")
	}

	return IPAddr(tmp), nil
}

func (v IPAddr) String() string {
	return netip.Addr(v).String()
}

type (
	IPAddr netip.Addr
)

func (IPAddr) Schema() spec.Schema {
	return spec.Schema{Type: spec.TypeString, Format: spec.Format("ipv4"), Pattern: "^((25[0-5]|(2[0-4]|1\\d|[1-9]|)\\d)\\.?\\b){4}$"}
}

func (v IPAddr) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *IPAddr) UnmarshalText(raw []byte) error {
	tmp, err := ParseIPAddr(string(raw))
	if err != nil {
		return err
	}

	*v = tmp
	return nil
}

var (
	_ encoding.TextMarshaler   = IPAddr{}
	_ encoding.TextUnmarshaler = (*IPAddr)(nil)
)
