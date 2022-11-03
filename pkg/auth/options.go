package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/pflag"
)

var (
	_ json.Marshaler   = (*PluginOptions)(nil)
	_ json.Unmarshaler = (*PluginOptions)(nil)
	_ pflag.Value      = (*PluginOptions)(nil)
)

// PluginOptions like json.RawMessage but implements pflag.Value
type PluginOptions []byte

// String implements fmt.Stringer
func (p *PluginOptions) String() string {
	buf := bytes.NewBuffer(nil)
	_ = json.Compact(buf, *p)

	return buf.String()
}

// Set implements pflag.Value
func (p *PluginOptions) Set(s string) error {
	return p.UnmarshalJSON([]byte(s))
}

// Type implements pflag.Value
func (p *PluginOptions) Type() string {
	return "string"
}

// UnmarshalJSON implements json.Unmarshaler
func (p *PluginOptions) UnmarshalJSON(val []byte) error {
	if val == nil {
		return errors.New("auth.PluginOptions: UnmarshalJSON on nil pointer")
	}

	if !bytes.Equal(val, []byte("null")) {
		*p = append((*p)[0:0], val...)
	}
	return nil
}

// MarshalJSON may get called on pointers or values, so implement MarshalJSON on value.
func (p PluginOptions) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("null"), nil
	}
	return p, nil
}

// Unmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by obj
func (p *PluginOptions) Unmarshal(obj any) error {
	return json.Unmarshal(*p, obj)
}

type Options struct {
	// Method The name of the authentication plugin to be used
	Method string `yaml:"method,omitempty" json:"method,omitempty"`

	// Options for the authentication plugin to be used
	Options *PluginOptions `yaml:"options,omitempty" json:"options,omitempty"`
}

func (o *Options) SetDefaults() {
	o.Method = "STATIC"
	if o.Options == nil {
		o.Options = NewPluginOptions()
	}
}

func (o *Options) Validate() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if o.Method == "" {
		return errors.New("method is required filed")
	}

	mod, err := NewPlugin(o.Method)
	if err != nil {
		return fmt.Errorf("unable create plugin %s, got: %v", o.Method, err)
	}

	defer func() {
		_ = mod.Close()
	}()

	if err := mod.Setup(ctx, o.Options); err != nil {
		return fmt.Errorf("unable setup plugin %s, got: %v", o.Method, err)
	}

	if err := mod.Validate(); err != nil {
		return fmt.Errorf("validate plugin %s failed, got: %v", o.Method, err)
	}
	return nil
}

// AddFlags add auth related command line parameters
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Method, "auth.method", o.Method, "The name of the authentication plugin to be used")
	fs.Var(o.Options, "auth.options", "Options for the authentication plugin to be used (json format)")
}

// NewPluginOptions  create `zero` plugin options
func NewPluginOptions() *PluginOptions {
	return new(PluginOptions)
}

// NewOptions create `zero` auth options
func NewOptions() *Options {
	return &Options{Options: NewPluginOptions()}
}
