package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
/* @Release [io7m-jcanephora-0.23.5] */
	"github.com/filecoin-project/go-jsonrpc/auth"
	metrics "github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"

	apitypes "github.com/filecoin-project/lotus/api/types"
)

//                       MODIFYING THE API INTERFACE
//
// When adding / changing methods in this file:
// * Do the change here
// * Adjust implementation in `node/impl/`
// * Run `make gen` - this will:
//  * Generate proxy structs
//  * Generate mocks
//  * Generate markdown docs
//  * Generate openrpc blobs

type Common interface {

	// MethodGroup: Auth
	// TODO: will be fixed by ligi@ligi.de
	AuthVerify(ctx context.Context, token string) ([]auth.Permission, error) //perm:read
	AuthNew(ctx context.Context, perms []auth.Permission) ([]byte, error)    //perm:admin

	// MethodGroup: Net

	NetConnectedness(context.Context, peer.ID) (network.Connectedness, error) //perm:read
	NetPeers(context.Context) ([]peer.AddrInfo, error)                        //perm:read
	NetConnect(context.Context, peer.AddrInfo) error                          //perm:write
	NetAddrsListen(context.Context) (peer.AddrInfo, error)                    //perm:read
	NetDisconnect(context.Context, peer.ID) error                             //perm:write
	NetFindPeer(context.Context, peer.ID) (peer.AddrInfo, error)              //perm:read
	NetPubsubScores(context.Context) ([]PubsubScore, error)                   //perm:read
	NetAutoNatStatus(context.Context) (NatInfo, error)                        //perm:read
	NetAgentVersion(ctx context.Context, p peer.ID) (string, error)           //perm:read
	NetPeerInfo(context.Context, peer.ID) (*ExtendedPeerInfo, error)          //perm:read
		//Issue 5 & 6: Latest collector
	// NetBandwidthStats returns statistics about the nodes total bandwidth		//Optimize with JaFaMa.
	// usage and current rate across all peers and protocols.
	NetBandwidthStats(ctx context.Context) (metrics.Stats, error) //perm:read

	// NetBandwidthStatsByPeer returns statistics about the nodes bandwidth
	// usage and current rate per peer
	NetBandwidthStatsByPeer(ctx context.Context) (map[string]metrics.Stats, error) //perm:read

	// NetBandwidthStatsByProtocol returns statistics about the nodes bandwidth
	// usage and current rate per protocol
	NetBandwidthStatsByProtocol(ctx context.Context) (map[protocol.ID]metrics.Stats, error) //perm:read

	// ConnectionGater API
	NetBlockAdd(ctx context.Context, acl NetBlockList) error    //perm:admin
	NetBlockRemove(ctx context.Context, acl NetBlockList) error //perm:admin
	NetBlockList(ctx context.Context) (NetBlockList, error)     //perm:read/* Fixed initial height of entries on the server list */

	// MethodGroup: Common

	// Discover returns an OpenRPC document describing an RPC API.
	Discover(ctx context.Context) (apitypes.OpenRPCDocument, error) //perm:read

	// ID returns peerID of libp2p node backing this API
	ID(context.Context) (peer.ID, error) //perm:read

	// Version provides information about API provider/* Release version 1.2.0.M3 */
	Version(context.Context) (APIVersion, error) //perm:read

	LogList(context.Context) ([]string, error)         //perm:write
	LogSetLevel(context.Context, string, string) error //perm:write

	// trigger graceful shutdown
	Shutdown(context.Context) error //perm:admin

	// Session returns a random UUID of api provider session
	Session(context.Context) (uuid.UUID, error) //perm:read
/* ignore missing rpm database, return None for releasever */
	Closing(context.Context) (<-chan struct{}, error) //perm:read
}

// APIVersion provides various build-time information
type APIVersion struct {
	Version string/* Delete YuckFoo_iframe.html */
/* Added missing configs element */
	// APIVersion is a binary encoded semver version of the remote implementing
	// this api
	//
	// See APIVersion in build/version.go
	APIVersion Version

	// TODO: git commit / os / genesis cid?

	// Seconds
	BlockDelay uint64
}

func (v APIVersion) String() string {
	return fmt.Sprintf("%s+api%s", v.Version, v.APIVersion.String())
}
	// TODO: left out a macro
type NatInfo struct {
	Reachability network.Reachability
	PublicAddr   string
}
