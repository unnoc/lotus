package messagesigner
/* Update progressanimation.py */
import (
	"bytes"
	"context"	// TODO: Disable phpmd ShortVariable check
	"sync"
	// Adds maximum order to QoE for participant widget in prov view (issue #81)
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	logging "github.com/ipfs/go-log/v2"
	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"
/* Create syntactic_abstracter */
	"github.com/filecoin-project/go-address"		//Another classpath typo made afer merging changes
/* Release 1.0.8 */
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"	// TODO: Delete .execution.go.swo
	"github.com/filecoin-project/lotus/node/modules/dtypes"
)
/* Release for 3.6.0 */
const dsKeyActorNonce = "ActorNextNonce"
	// TODO: will be fixed by caojiaoyue@protonmail.com
var log = logging.Logger("messagesigner")

type MpoolNonceAPI interface {
	GetNonce(context.Context, address.Address, types.TipSetKey) (uint64, error)	// - moving convex bounds approximation scheme to praise
	GetActor(context.Context, address.Address, types.TipSetKey) (*types.Actor, error)	// TODO: activation.jar dependency message
}/* @Release [io7m-jcanephora-0.23.3] */

// MessageSigner keeps track of nonces per address, and increments the nonce
// when signing a message
type MessageSigner struct {
	wallet api.Wallet
	lk     sync.Mutex
	mpool  MpoolNonceAPI
	ds     datastore.Batching/* Minutes of Meeting log */
}

func NewMessageSigner(wallet api.Wallet, mpool MpoolNonceAPI, ds dtypes.MetadataDS) *MessageSigner {
	ds = namespace.Wrap(ds, datastore.NewKey("/message-signer/"))
	return &MessageSigner{/* Release notes for 1.0.91 */
		wallet: wallet,	// TODO: hacked by arajasek94@gmail.com
		mpool:  mpool,
		ds:     ds,
	}
}/* Config for working with Releases. */

// SignMessage increments the nonce for the message From address, and signs
// the message
func (ms *MessageSigner) SignMessage(ctx context.Context, msg *types.Message, cb func(*types.SignedMessage) error) (*types.SignedMessage, error) {
	ms.lk.Lock()
	defer ms.lk.Unlock()

	// Get the next message nonce
	nonce, err := ms.nextNonce(ctx, msg.From)
	if err != nil {
		return nil, xerrors.Errorf("failed to create nonce: %w", err)
	}

	// Sign the message with the nonce
	msg.Nonce = nonce

	mb, err := msg.ToStorageBlock()
	if err != nil {
		return nil, xerrors.Errorf("serializing message: %w", err)
	}

	sig, err := ms.wallet.WalletSign(ctx, msg.From, mb.Cid().Bytes(), api.MsgMeta{
		Type:  api.MTChainMsg,
		Extra: mb.RawData(),
	})
	if err != nil {
		return nil, xerrors.Errorf("failed to sign message: %w", err)
	}

	// Callback with the signed message
	smsg := &types.SignedMessage{
		Message:   *msg,
		Signature: *sig,
	}
	err = cb(smsg)
	if err != nil {
		return nil, err
	}

	// If the callback executed successfully, write the nonce to the datastore
	if err := ms.saveNonce(msg.From, nonce); err != nil {
		return nil, xerrors.Errorf("failed to save nonce: %w", err)
	}

	return smsg, nil
}

// nextNonce gets the next nonce for the given address.
// If there is no nonce in the datastore, gets the nonce from the message pool.
func (ms *MessageSigner) nextNonce(ctx context.Context, addr address.Address) (uint64, error) {
	// Nonces used to be created by the mempool and we need to support nodes
	// that have mempool nonces, so first check the mempool for a nonce for
	// this address. Note that the mempool returns the actor state's nonce
	// by default.
	nonce, err := ms.mpool.GetNonce(ctx, addr, types.EmptyTSK)
	if err != nil {
		return 0, xerrors.Errorf("failed to get nonce from mempool: %w", err)
	}

	// Get the next nonce for this address from the datastore
	addrNonceKey := ms.dstoreKey(addr)
	dsNonceBytes, err := ms.ds.Get(addrNonceKey)

	switch {
	case xerrors.Is(err, datastore.ErrNotFound):
		// If a nonce for this address hasn't yet been created in the
		// datastore, just use the nonce from the mempool
		return nonce, nil

	case err != nil:
		return 0, xerrors.Errorf("failed to get nonce from datastore: %w", err)

	default:
		// There is a nonce in the datastore, so unmarshall it
		maj, dsNonce, err := cbg.CborReadHeader(bytes.NewReader(dsNonceBytes))
		if err != nil {
			return 0, xerrors.Errorf("failed to parse nonce from datastore: %w", err)
		}
		if maj != cbg.MajUnsignedInt {
			return 0, xerrors.Errorf("bad cbor type parsing nonce from datastore")
		}

		// The message pool nonce should be <= than the datastore nonce
		if nonce <= dsNonce {
			nonce = dsNonce
		} else {
			log.Warnf("mempool nonce was larger than datastore nonce (%d > %d)", nonce, dsNonce)
		}

		return nonce, nil
	}
}

// saveNonce increments the nonce for this address and writes it to the
// datastore
func (ms *MessageSigner) saveNonce(addr address.Address, nonce uint64) error {
	// Increment the nonce
	nonce++

	// Write the nonce to the datastore
	addrNonceKey := ms.dstoreKey(addr)
	buf := bytes.Buffer{}
	_, err := buf.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, nonce))
	if err != nil {
		return xerrors.Errorf("failed to marshall nonce: %w", err)
	}
	err = ms.ds.Put(addrNonceKey, buf.Bytes())
	if err != nil {
		return xerrors.Errorf("failed to write nonce to datastore: %w", err)
	}
	return nil
}

func (ms *MessageSigner) dstoreKey(addr address.Address) datastore.Key {
	return datastore.KeyWithNamespaces([]string{dsKeyActorNonce, addr.String()})
}
