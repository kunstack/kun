package auth_test

import (
	"bytes"
	"encoding/json"
	"github.com/aapelismith/kuntunnel/pkg/auth"
	"sigs.k8s.io/yaml"
	"testing"
)

type opts struct {
	Name  string `json:"name,omitempty"`
	Title string `json:"title,omitempty"`
	Foo   string `json:"foo,omitempty"`
}

func TestPluginOptions(t *testing.T) {
	data := opts{Name: "hello", Title: "World", Foo: "666"}

	tgt := opts{}

	b, err := yaml.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	opts := auth.NewPluginOptions()

	if err := yaml.Unmarshal(b, opts); err != nil {
		t.Fatal(err)
	}

	if err := opts.Unmarshal(&tgt); err != nil {
		t.Fatal(err)
	}

	if tgt.Name != data.Name {
		t.Fail()
	}

	d, err := yaml.Marshal(opts)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, d) {
		t.Fail()
	}
}

func TestPluginOptions_UnmarshalJSON(t *testing.T) {
	data := []byte(`
 		{
			"foo":"666",
			"name":"hello",
			"title":"World"
		}
	`)

	opts := auth.NewPluginOptions()

	if err := json.Unmarshal(data, opts); err != nil {
		t.Fatal(err)
	}

	if opts.String() != `{"foo":"666","name":"hello","title":"World"}` {
		t.Fail()
	}
}
