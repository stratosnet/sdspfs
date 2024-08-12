package provider

import (
	"context"
	"fmt"
	"io"

	"github.com/ipfs/boxo/files"
	unixfile "github.com/ipfs/boxo/ipld/unixfs/file"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/boxo/provider"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log/v2"
	iface "github.com/ipfs/kubo/core/coreiface"
	"github.com/sdspfs/rpc"
)

var logger = logging.Logger("sdspdf")
var _ provider.System = (*ppProvider)(nil)

// ppProvider implements a provider path to the pp network.
type ppProvider struct {
	ctx      context.Context
	provider provider.System
	dag      ipld.DAGService
	rpc      *rpc.Rpc
}

// WrapResolver wraps the given path Resolver with a content-blocking layer
// for Resolve operations.
func WrapProvider(provider provider.System, dag ipld.DAGService) provider.System {
	logger.Debugf("Path resolved wrapped with pp provider")
	fmt.Println("provider", provider)
	fmt.Println("dag", dag)
	// TODO: Add config
	rpc, _ := rpc.NewRpc("https://sds-gateway.thestratos.org/private/rpc/sjwRNkqXQ0MI5bPXjTs2Q2vXheA")
	return &ppProvider{
		ctx:      context.Background(),
		provider: provider,
		dag:      dag,
		rpc:      rpc,
	}
}

func (pp *ppProvider) Close() error {
	fmt.Println("SDS PP PROVIDER: CLOSE")
	return pp.provider.Close()
}

func (pp *ppProvider) Provide(cid cid.Cid) error {
	fmt.Println("SDS PP PROVIDER: PROVIDE", cid)
	err := pp.provider.Provide(cid)
	if err != nil {
		return err
	}

	p := path.FromCid(cid)

	fmt.Println("SDS PP PROVIDER: path p.String()", p.String())

	nd, err := pp.dag.Get(pp.ctx, cid)
	if err != nil {
		return err
	}
	fmt.Println("SDS PP PROVIDER: node nd.RawData()", nd.RawData())

	f, err := unixfile.NewUnixfsFile(pp.ctx, pp.dag, nd)
	if err != nil {
		return err
	}
	fmt.Println("SDS PP PROVIDER: ufx file", f)

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

	fmt.Println("SDS PP PROVIDER: file bytes", b)

	oz, err := pp.rpc.GetOzone()
	fmt.Println("oz", oz)
	fmt.Println("err", err)
	if err != nil {
		return err
	}

	return nil
}

func (pp *ppProvider) Reprovide(ctx context.Context) error {
	fmt.Println("SDS PP PROVIDER: REPROVIDE", ctx)
	return pp.provider.Reprovide(ctx)
}

func (pp *ppProvider) Stat() (provider.ReproviderStats, error) {
	fmt.Println("SDS PP PROVIDER: STAT")
	return pp.provider.Stat()
}
