/*
 * Copyright 2021 Aapeli.Smith<aapeli.nian@gmail.com>.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slices"
)

var (
	encodings = []string{"json", "console"}
)

var (
	_ pflag.Value = (*atomicLevel)(nil)
)

type atomicLevel struct {
	lvl *zap.AtomicLevel
}

func (a *atomicLevel) String() string {
	return a.lvl.String()
}

func (a *atomicLevel) Set(s string) error {
	lvl, err := zapcore.ParseLevel(s)
	if err != nil {
		return err
	}
	a.lvl.SetLevel(lvl)
	return nil
}

func (a *atomicLevel) Type() string {
	return "string"
}

type Options struct {
	// Options the zap options
	Options []zap.Option `yaml:"-" json:"-"`

	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	Level zap.AtomicLevel `json:"level" yaml:"level"`

	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktrace more liberally.
	Development bool `json:"development" yaml:"development"`

	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`

	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktrace are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`

	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	Sampling *zap.SamplingConfig `json:"sampling" yaml:"sampling"`

	// Encoding sets the logger's encoding. Valid values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	Encoding string `json:"encoding" yaml:"encoding"`

	// EncoderConfig sets options for the chosen encoder. See
	// zapcore.EncoderConfig for details.
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`

	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`

	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`

	// InitialFields is a collection of fields to add to the root logger.
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}

// SetDefaults sets the default values.
func (o *Options) SetDefaults() {
	o.Options = make([]zap.Option, 0)

	o.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	o.Encoding = "console"

	if o.Sampling == nil {
		o.Sampling = &zap.SamplingConfig{}
	}

	o.Sampling.Initial = 100
	o.Sampling.Thereafter = 100

	o.EncoderConfig.TimeKey = "ts"
	o.EncoderConfig.LevelKey = "level"
	o.EncoderConfig.NameKey = "logger"
	o.EncoderConfig.CallerKey = "caller"
	o.EncoderConfig.FunctionKey = zapcore.OmitKey
	o.EncoderConfig.MessageKey = "msg"
	o.EncoderConfig.StacktraceKey = "stacktrace"
	o.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
	o.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	o.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	o.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	o.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	o.OutputPaths = []string{"stderr"}
	o.ErrorOutputPaths = []string{"stderr"}
}

// AddFlags add related command line parameters
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.Var(&atomicLevel{lvl: &o.Level}, "log.level", "Log level to configure the "+
		"verbosity of logging. Can be one of 'debug', 'info', 'warn', 'error', 'dpanic', 'panic', 'fatal'")

	fs.BoolVar(&o.Development, "log.development", o.Development, "puts the logger in development mode")

	fs.BoolVar(&o.DisableCaller, "log.disable-caller", o.DisableCaller,
		"Completely disables automatic stacktrace capturing. By default, stacktrace are captured"+
			" for WarnLevel and above logs in development and ErrorLevel and above in production.")

	fs.BoolVar(&o.DisableStacktrace, "log.disable-stacktrace", o.DisableStacktrace,
		"completely disables automatic stacktrace capturing. By  default, stacktrace are captured"+
			" for WarnLevel and above logs in  development and ErrorLevel and above in production")

	if o.Sampling != nil {
		fs.IntVar(&o.Sampling.Initial, "log.sampling.initial", o.Sampling.Initial, "Set the initial value for log sampling.")

		fs.IntVar(&o.Sampling.Thereafter, "log.sampling.thereafter", o.Sampling.Thereafter, "Set the thereafter value for log sampling.")
	}

	fs.StringVar(&o.Encoding, "log.encoding", o.Encoding, "sets the logger's encoding. "+
		"Valid values are 'json' and 'console', as well as any third-party encodings registered via  RegisterEncoder.")

	fs.StringArrayVar(&o.OutputPaths, "log.output-paths", o.OutputPaths, "the list of URLs or file paths to write logging output to.")
}

// Validate verify the configuration and return an error if correct
func (o *Options) Validate() error {
	if o.Encoding == "" {
		return errors.New("no encoder name specified")
	}

	if !slices.Contains(encodings, o.Encoding) {
		return fmt.Errorf("no encoder registered for name '%s'", o.Encoding)
	}
	return nil
}

// NewOptions returns a `zero` instance
func NewOptions() *Options {
	return &Options{
		Sampling:      &zap.SamplingConfig{},
		EncoderConfig: zapcore.EncoderConfig{},
	}
}
