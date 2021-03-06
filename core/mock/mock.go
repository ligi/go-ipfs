package coremock

import (
	"net"

	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore"
	syncds "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore/sync"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"

	commands "github.com/ipfs/go-ipfs/commands"
	core "github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo"
	config "github.com/ipfs/go-ipfs/repo/config"
	ds2 "github.com/ipfs/go-ipfs/thirdparty/datastore2"
	testutil "github.com/ipfs/go-ipfs/thirdparty/testutil"
	peer "gx/ipfs/QmQGwpJy9P4yXZySmqkZEXCmbBpJUb8xntCv8Ca4taZwDC/go-libp2p-peer"
	host "gx/ipfs/QmQgQeBQxQmJdeUSaDagc8cr2ompDwGn13Cybjdtzfuaki/go-libp2p/p2p/host"
	metrics "gx/ipfs/QmQgQeBQxQmJdeUSaDagc8cr2ompDwGn13Cybjdtzfuaki/go-libp2p/p2p/metrics"
	mocknet "gx/ipfs/QmQgQeBQxQmJdeUSaDagc8cr2ompDwGn13Cybjdtzfuaki/go-libp2p/p2p/net/mock"
	pstore "gx/ipfs/QmZ62t46e9p7vMYqCmptwQC1RhRv5cpQ5cwoqYspedaXyq/go-libp2p-peerstore"
)

// NewMockNode constructs an IpfsNode for use in tests.
func NewMockNode() (*core.IpfsNode, error) {
	ctx := context.Background()

	// effectively offline, only peer in its network
	return core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Host:   MockHostOption(mocknet.New(ctx)),
	})
}

func MockHostOption(mn mocknet.Mocknet) core.HostOption {
	return func(ctx context.Context, id peer.ID, ps pstore.Peerstore, bwr metrics.Reporter, fs []*net.IPNet) (host.Host, error) {
		return mn.AddPeerWithPeerstore(id, ps)
	}
}

func MockCmdsCtx() (commands.Context, error) {
	// Generate Identity
	ident, err := testutil.RandIdentity()
	if err != nil {
		return commands.Context{}, err
	}
	p := ident.ID()

	conf := config.Config{
		Identity: config.Identity{
			PeerID: p.String(),
		},
	}

	r := &repo.Mock{
		D: ds2.CloserWrap(syncds.MutexWrap(datastore.NewMapDatastore())),
		C: conf,
	}

	node, err := core.NewNode(context.Background(), &core.BuildCfg{
		Repo: r,
	})
	if err != nil {
		return commands.Context{}, err
	}

	return commands.Context{
		Online:     true,
		ConfigRoot: "/tmp/.mockipfsconfig",
		LoadConfig: func(path string) (*config.Config, error) {
			return &conf, nil
		},
		ConstructNode: func() (*core.IpfsNode, error) {
			return node, nil
		},
	}, nil
}
