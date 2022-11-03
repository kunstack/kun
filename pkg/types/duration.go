package types

import (
	"encoding/json"
	"github.com/spf13/pflag"
	"io"
	"strings"
	"time"
)

var (
	_ pflag.Value      = (*Duration)(nil)
	_ json.Unmarshaler = (*Duration)(nil)
)

type Duration time.Duration

func (d *Duration) Set(s string) error {
	t, err := time.ParseDuration(s)
	*d = Duration(t)
	return err
}

func (d *Duration) Type() string {
	return "duration"
}

func (d *Duration) String() string {
	return (*time.Duration)(d).String()
}

func (d *Duration) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return io.EOF
	}

	var (
		err error
		dur time.Duration
	)

	if bytes[0] == byte('"') && bytes[len(bytes)-1] == byte('"') {
		dur, err = time.ParseDuration(strings.Trim(string(bytes), `"`))
	} else {
		err = json.Unmarshal(bytes, &dur)
	}

	*d = Duration(dur)
	return err
}

// MarshalJSON implements the json.Marshaler interface.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
