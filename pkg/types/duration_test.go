package types_test

import (
	"encoding/json"
	"github.com/aapelismith/kun/pkg/types"
	"testing"
	"time"
)

func TestDuration_Set(t *testing.T) {
	d := types.Duration(0)

	if err := d.Set("1s"); err != nil {
		t.Fatal(err)
	}

	if d != types.Duration(time.Second) {
		t.Fail()
	}
}

func TestDuration_String(t *testing.T) {
	d := types.Duration(time.Second)

	if d.String() != "1s" {
		t.Fail()
	}
}

func TestDuration_UnmarshalJSON2(t *testing.T) {
	result := struct {
		Time types.Duration `json:"time"`
	}{}

	if err := json.Unmarshal([]byte(`{"time":"1s"}`), &result); err != nil {
		t.Fatal(err)
	}

	if result.Time != types.Duration(time.Second) {
		t.Fail()
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", data)
}

func TestDuration_MarshalJSON(t *testing.T) {
	d := types.Duration(time.Second)

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != `"1s"` {
		t.Fail()
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	d := types.Duration(0)

	if err := json.Unmarshal([]byte(`"1s"`), &d); err != nil {
		t.Fatal(err)
	}

	if d != types.Duration(time.Second) {
		t.Fail()
	}

	if err := json.Unmarshal([]byte("2000000000"), &d); err != nil {
		t.Fatal(err)
	}

	if d != types.Duration(time.Second*2) {
		t.Fail()
	}
}
