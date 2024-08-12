package main

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/plugin"

	"github.com/stratosnet/sdspfs/provider"
	"go.uber.org/fx"
)

var logger = logging.Logger("sdspfs")

// Plugins sets the list of plugins to be loaded.
var Plugins = []plugin.Plugin{
	&sdsPlugin{},
}

// sdsPlugin is used for enabling sds pp node.
type sdsPlugin struct{}

var _ plugin.PluginFx = (*sdsPlugin)(nil)

func (p *sdsPlugin) Name() string {
	return "sdspfs"
}

func (p *sdsPlugin) Version() string {
	return "0.0.1"
}

func (p *sdsPlugin) Init(env *plugin.Environment) error {
	return nil
}

func (p *sdsPlugin) Options(info core.FXNodeInfo) ([]fx.Option, error) {
	logging.SetLogLevel("sdspfs", "INFO")
	logger.Info("Loading Sds plugin: another brand new ds")

	opts := append(
		info.FXOptions,
		fx.Decorate(provider.WrapProvider),
	)
	return opts, nil
}
