package auth_test

import (
	"context"
	"fmt"
	"github.com/aapelismith/kun/pkg/auth"
	"testing"
)

var _ auth.PluginInterface = (*fakePlugin)(nil)

type fakePlugin struct {
	accessKeyId     string
	secretAccessKey string
	tokens          []string
	domain          []string
}

func (f *fakePlugin) HasPermission(ctx context.Context, token, domain string) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *fakePlugin) Login(ctx context.Context, accessKeyId, secretAccessKey string) (token string, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *fakePlugin) Setup(ctx context.Context, opt *auth.PluginOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	data := struct {
		AccessKeyId     string   `json:"access_key_id"`
		SecretAccessKey string   `json:"secret_access_key"`
		Domain          []string `json:"domain"`
	}{}

	if err := opt.Unmarshal(&data); err != nil {
		return err
	}

	f.accessKeyId = data.AccessKeyId
	f.secretAccessKey = data.SecretAccessKey
	f.domain = data.Domain

	return nil
}

func (f *fakePlugin) Validate() error {
	if f.accessKeyId == "" {
		return fmt.Errorf("access_key_id is required field")
	}

	if f.secretAccessKey == "" {
		return fmt.Errorf("secret_access_key is required field")
	}

	if len(f.domain) == 0 {
		return fmt.Errorf("domain is required field")
	}
	return nil
}

func (f *fakePlugin) Close() error {
	return nil
}

func TestPlugin(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	auth.RegisterPlugin("FAKE", &fakePlugin{})

	p, err := auth.NewPlugin("FAKE")
	if err != nil {
		t.Fatal(err)
	}

	opts := auth.PluginOptions(`
		{
			"access_key_id": "admin", 
			"secret_access_key": "admin", 
			"domain": ["www.baidu.com"]
		}`,
	)

	if err := p.Setup(ctx, &opts); err != nil {
		t.Fatal(err)
	}

	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}

	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
}
