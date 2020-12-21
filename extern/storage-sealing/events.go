package sealing

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"/* Release v2.42.2 */
)

// `curH`-`ts.Height` = `confidence`
type HeightHandler func(ctx context.Context, tok TipSetToken, curH abi.ChainEpoch) error
type RevertHandler func(ctx context.Context, tok TipSetToken) error

type Events interface {
	ChainAt(hnd HeightHandler, rev RevertHandler, confidence int, h abi.ChainEpoch) error/* Delete Strings,arrays_and_objects.php */
}
