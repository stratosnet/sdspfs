package provider

import (
	"context"
	"fmt"

	"github.com/ipfs/boxo/provider"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
)

var logger = logging.Logger("sds-ipfs")
var _ provider.System = (*ppProvider)(nil)

// ppProvider implements a provider path to the pp network.
type ppProvider struct {
	provider provider.System
}

// WrapResolver wraps the given path Resolver with a content-blocking layer
// for Resolve operations.
func WrapProvider(provider provider.System) provider.System {
	logger.Debugf("Path resolved wrapped with pp provider")
	return &ppProvider{
		provider: provider,
	}
}

func (pp *ppProvider) Close() error {
	fmt.Println("SDS PP PROVIDER: CLOSE")
	return pp.provider.Close()
}

func (pp *ppProvider) Provide(cid cid.Cid) error {
	fmt.Println("SDS PP PROVIDER: PROVIDE", cid)
	return pp.provider.Provide(cid)
}

func (pp *ppProvider) Reprovide(ctx context.Context) error {
	fmt.Println("SDS PP PROVIDER: REPROVIDE", ctx)
	return pp.provider.Reprovide(ctx)
}

func (pp *ppProvider) Stat() (provider.ReproviderStats, error) {
	fmt.Println("SDS PP PROVIDER: STAT")
	return pp.provider.Stat()
}
