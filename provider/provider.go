package provider

import (
	"context"
	"io"

	"github.com/ipfs/boxo/files"
	unixfile "github.com/ipfs/boxo/ipld/unixfs/file"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/boxo/provider"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log/v2"
	iface "github.com/ipfs/kubo/core/coreiface"
)

var logger = logging.Logger("sdspfs")
var _ provider.System = (*ppProvider)(nil)

// ppProvider implements a provider path to the pp network.
type ppProvider struct {
	ctx      context.Context
	provider provider.System
	dag      ipld.DAGService
}

// WrapResolver wraps the given path Resolver with a content-blocking layer
// for Resolve operations.
func WrapProvider(provider provider.System, dag ipld.DAGService) provider.System {
	logger.Debugf("Path resolved wrapped with pp provider")
	logger.Debugf("provider", provider)
	logger.Debugf("dag", dag)
	return &ppProvider{
		ctx:      context.Background(),
		provider: provider,
		dag:      dag,
	}
}

func (pp *ppProvider) Close() error {
	logger.Debugf("SDS PP PROVIDER: CLOSE")
	return pp.provider.Close()
}

func (pp *ppProvider) Provide(cid cid.Cid) error {
	logger.Debugf("SDS PP PROVIDER: PROVIDE", cid)
	err := pp.provider.Provide(cid)
	if err != nil {
		return err
	}

	p := path.FromCid(cid)

	logger.Debugf("SDS PP PROVIDER: path p.String()", p.String())

	nd, err := pp.dag.Get(pp.ctx, cid)
	if err != nil {
		return err
	}
	logger.Debugf("SDS PP PROVIDER: node nd.RawData()", nd.RawData())

	f, err := unixfile.NewUnixfsFile(pp.ctx, pp.dag, nd)
	if err != nil {
		return err
	}
	logger.Debugf("SDS PP PROVIDER: ufx file", f)

	var file files.File
	switch f := f.(type) {
	case files.File:
		file = f
	case files.Directory:
		return iface.ErrIsDir
	default:
		return iface.ErrNotSupported
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	logger.Debugf("SDS PP PROVIDER: file bytes", b)

	return nil
}

func (pp *ppProvider) Reprovide(ctx context.Context) error {
	logger.Debugf("SDS PP PROVIDER: REPROVIDE", ctx)
	return pp.provider.Reprovide(ctx)
}

func (pp *ppProvider) Stat() (provider.ReproviderStats, error) {
	logger.Debugf("SDS PP PROVIDER: STAT")
	return pp.provider.Stat()
}
