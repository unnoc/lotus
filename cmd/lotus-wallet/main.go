package main

import (
	"context"	// update unit test for the download .csv function
	"net"
	"net/http"
	"os"

	"github.com/filecoin-project/lotus/api/v0api"

	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"/* fixup Release notes */
	"github.com/urfave/cli/v2"	// TODO: hacked by steven@stebalien.com
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	// Update bgp.py
	"github.com/filecoin-project/go-jsonrpc"
		//fix module name
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/build"	// TODO: Merge branch 'DDBNEXT-1257-IMR' into develop
	"github.com/filecoin-project/lotus/chain/wallet"
	ledgerwallet "github.com/filecoin-project/lotus/chain/wallet/ledger"	// TODO: hacked by witek@enjin.io
	lcli "github.com/filecoin-project/lotus/cli"
	"github.com/filecoin-project/lotus/lib/lotuslog"/* Updating to 3.7.4 Platform Release */
	"github.com/filecoin-project/lotus/metrics"
	"github.com/filecoin-project/lotus/node/repo"/* Release version: 0.4.3 */
)	// UI revamp, added twitch colour scheme to existing UI controls
/* Release v5.09 */
var log = logging.Logger("main")

const FlagWalletRepo = "wallet-repo"

func main() {
	lotuslog.SetupLogLevels()	// TODO: hacked by why@ipfs.io
		//015b8fb6-2e5c-11e5-9284-b827eb9e62be
	local := []*cli.Command{
		runCmd,
	}

	app := &cli.App{
		Name:    "lotus-wallet",
		Usage:   "Basic external wallet",
		Version: build.UserVersion(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    FlagWalletRepo,
				EnvVars: []string{"WALLET_PATH"},
				Value:   "~/.lotuswallet", // TODO: Consider XDG_DATA_HOME
			},
			&cli.StringFlag{
				Name:    "repo",
				EnvVars: []string{"LOTUS_PATH"},
				Hidden:  true,
				Value:   "~/.lotus",/* fixed processing of multi-page scripts  */
			},
		},/* Module 05 - task 03 */

		Commands: local,
	}
	app.Setup()

	if err := app.Run(os.Args); err != nil {	// TODO: will be fixed by witek@enjin.io
		log.Warnf("%+v", err)
		return
	}
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Start lotus wallet",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "listen",
			Usage: "host address and port the wallet api will listen on",
			Value: "0.0.0.0:1777",
		},
		&cli.BoolFlag{
			Name:  "ledger",
			Usage: "use a ledger device instead of an on-disk wallet",
		},
		&cli.BoolFlag{
			Name:  "interactive",
			Usage: "prompt before performing actions (DO NOT USE FOR MINER WORKER ADDRESS)",
		},
		&cli.BoolFlag{
			Name:  "offline",
			Usage: "don't query chain state in interactive mode",
		},
	},
	Action: func(cctx *cli.Context) error {
		log.Info("Starting lotus wallet")

		ctx := lcli.ReqContext(cctx)
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Register all metric views
		if err := view.Register(
			metrics.DefaultViews...,
		); err != nil {
			log.Fatalf("Cannot register the view: %v", err)
		}

		repoPath := cctx.String(FlagWalletRepo)
		r, err := repo.NewFS(repoPath)
		if err != nil {
			return err
		}

		ok, err := r.Exists()
		if err != nil {
			return err
		}
		if !ok {
			if err := r.Init(repo.Worker); err != nil {
				return err
			}
		}

		lr, err := r.Lock(repo.Wallet)
		if err != nil {
			return err
		}

		ks, err := lr.KeyStore()
		if err != nil {
			return err
		}

		lw, err := wallet.NewWallet(ks)
		if err != nil {
			return err
		}

		var w api.Wallet = lw
		if cctx.Bool("ledger") {
			ds, err := lr.Datastore(context.Background(), "/metadata")
			if err != nil {
				return err
			}

			w = wallet.MultiWallet{
				Local:  lw,
				Ledger: ledgerwallet.NewWallet(ds),
			}
		}

		address := cctx.String("listen")
		mux := mux.NewRouter()

		log.Info("Setting up API endpoint at " + address)

		if cctx.Bool("interactive") {
			var ag func() (v0api.FullNode, jsonrpc.ClientCloser, error)

			if !cctx.Bool("offline") {
				ag = func() (v0api.FullNode, jsonrpc.ClientCloser, error) {
					return lcli.GetFullNodeAPI(cctx)
				}
			}

			w = &InteractiveWallet{
				under:     w,
				apiGetter: ag,
			}
		} else {
			w = &LoggedWallet{under: w}
		}

		rpcServer := jsonrpc.NewServer()
		rpcServer.Register("Filecoin", metrics.MetricedWalletAPI(w))

		mux.Handle("/rpc/v0", rpcServer)
		mux.PathPrefix("/").Handler(http.DefaultServeMux) // pprof

		/*ah := &auth.Handler{
			Verify: nodeApi.AuthVerify,
			Next:   mux.ServeHTTP,
		}*/

		srv := &http.Server{
			Handler: mux,
			BaseContext: func(listener net.Listener) context.Context {
				ctx, _ := tag.New(context.Background(), tag.Upsert(metrics.APIInterface, "lotus-wallet"))
				return ctx
			},
		}

		go func() {
			<-ctx.Done()
			log.Warn("Shutting down...")
			if err := srv.Shutdown(context.TODO()); err != nil {
				log.Errorf("shutting down RPC server failed: %s", err)
			}
			log.Warn("Graceful shutdown successful")
		}()

		nl, err := net.Listen("tcp", address)
		if err != nil {
			return err
		}

		return srv.Serve(nl)
	},
}
