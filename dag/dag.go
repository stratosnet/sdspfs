package dag

import (
	"context"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log/v2"
)

var logger = logging.Logger("sdspfs")
var _ ipld.DAGService = (*ppDag)(nil)

type ppDag struct {
	dag ipld.DAGService
}

// WrapDag wraps to get a controll to dag service
func WrapDag(dag ipld.DAGService) ipld.DAGService {
	logger.Debugf("SDS PP: Wrapped with pp dag")
	logger.Debugf("dag", dag)
	return &ppDag{
		dag: dag,
	}
}

func (pp *ppDag) Add(ctx context.Context, nd ipld.Node) error {
	logger.Debugf("SDS PP DAG: Add")
	return pp.dag.Add(ctx, nd)
}

func (pp *ppDag) AddMany(ctx context.Context, nds []ipld.Node) error {
	logger.Debugf("SDS PP DAG: AddMany")
	return pp.dag.AddMany(ctx, nds)
}

func (pp *ppDag) Get(ctx context.Context, cid cid.Cid) (ipld.Node, error) {
	logger.Debugf("SDS PP DAG: Get")
	return pp.dag.Get(ctx, cid)
}

func (pp *ppDag) GetMany(ctx context.Context, cids []cid.Cid) <-chan *ipld.NodeOption {
	logger.Debugf("SDS PP DAG: GetMany")
	return pp.dag.GetMany(ctx, cids)
}

func (pp *ppDag) Remove(ctx context.Context, cid cid.Cid) error {
	logger.Debugf("SDS PP DAG: Remove")
	return pp.dag.Remove(ctx, cid)
}

func (pp *ppDag) RemoveMany(ctx context.Context, cids []cid.Cid) error {
	logger.Debugf("SDS PP DAG: RemoveMany")
	return pp.dag.RemoveMany(ctx, cids)
}
