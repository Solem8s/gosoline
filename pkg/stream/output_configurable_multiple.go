package stream

import (
	"context"
	"fmt"
	"github.com/applike/gosoline/pkg/cfg"
	"github.com/applike/gosoline/pkg/mon"
	"github.com/hashicorp/go-multierror"
)

type multiOutput struct {
	outputs []Output
}

func (m *multiOutput) WriteOne(ctx context.Context, msg WritableMessage) error {
	err := &multierror.Error{}

	for _, output := range m.outputs {
		err = multierror.Append(err, output.WriteOne(ctx, msg))
	}

	return err.ErrorOrNil()
}

func (m *multiOutput) Write(ctx context.Context, batch []WritableMessage) error {
	err := &multierror.Error{}

	for _, output := range m.outputs {
		err = multierror.Append(err, output.Write(ctx, batch))
	}

	return err.ErrorOrNil()
}

func NewConfigurableMultiOutput(config cfg.Config, logger mon.Logger, base string) Output {
	key := fmt.Sprintf("%s.types", ConfigurableOutputKey(base))
	ts := config.Get(key).(map[string]interface{})
	output := &multiOutput{
		outputs: make([]Output, 0),
	}

	for outputName := range ts {
		name := fmt.Sprintf("%s.types.%s", base, outputName)
		o := NewConfigurableOutput(config, logger, name)
		output.outputs = append(output.outputs, o)
	}

	return output
}
