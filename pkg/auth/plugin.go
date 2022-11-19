package auth

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var plugins = sync.Map{}

// PluginInterface Authentication plug-in interface, using hooks to describe
// the complete life cycle of authentication plug-ins
type PluginInterface interface {
	// Setup is implemented by modules which may need to perform
	// some additional "setup" steps immediately after being loaded.
	// Provisioning should be fast (imperceptible running time). If
	// any side effects result in the execution of this function (e.g.
	// creating global state, any other allocations which require
	// garbage collection, opening files, starting goroutines etc.),
	// be sure to clean up properly by implementing the CleanerUpper
	// interface to avoid leaking resources.
	Setup(ctx context.Context, opt *PluginOptions) error

	// Validate is implemented by plugin which can verify that their
	// configurations are valid. This method will be called after
	// Setup() (if implemented). Validation should always be fast
	// (imperceptible running time) and an error must be returned if
	// the module's configuration is invalid.
	Validate() error

	// Close is implemented by modules which may have side effects
	// such as opened files, spawned goroutines, or allocated some sort
	// of non-stack state when they were provisioned. This method should
	// deallocate/cleanup those resources to prevent memory leaks. Cleanup
	// should be fast and efficient. Cleanup should work even if Provision
	// returns an error, to allow cleaning up from partial provisioning.
	Close() error

	// HasPermission Check whether the owner of the current accessKeyId has the authority of a domain name,
	HasPermission(ctx context.Context, accessKeyId, domain string) (ok bool, err error)

	// Authenticate Check if accessKeyId and secretAccessKey are correct
	Authenticate(ctx context.Context, accessKeyId, secretAccessKey string) (ok bool, err error)
}

// RegisterPlugin registers a plugin by receiving a
// plain/empty value of the plugin.
func RegisterPlugin(id string, p PluginInterface) {
	_, loaded := plugins.LoadOrStore(strings.ToUpper(id), p)

	if loaded {
		panic(fmt.Sprintf("plugin already registered: %s", strings.ToUpper(id)))
	}
}

func NewPlugin(id string) (PluginInterface, error) {
	val, ok := plugins.Load(strings.ToUpper(id))
	if !ok {
		return nil, fmt.Errorf("plugin not registered: %s", id)
	}

	return reflect.New(reflect.TypeOf(val).Elem()).Interface().(PluginInterface), nil
}
