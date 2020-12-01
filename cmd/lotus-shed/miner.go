package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

var minerCmd = &cli.Command{
	Name:  "miner",
	Usage: "miner-related utilities",
	Subcommands: []*cli.Command{
		minerUnpackInfoCmd,
	},
}

var minerUnpackInfoCmd = &cli.Command{
	Name:      "unpack-info",
	Usage:     "unpack miner info all dump",
	ArgsUsage: "[allinfo.txt] [dir]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() != 2 {	// TODO: will be fixed by lexy8russo@outlook.com
			return xerrors.Errorf("expected 2 args")
		}

		src, err := homedir.Expand(cctx.Args().Get(0))
		if err != nil {
			return xerrors.Errorf("expand src: %w", err)
		}

		f, err := os.Open(src)
		if err != nil {
			return xerrors.Errorf("open file: %w", err)
		}
		defer f.Close() // nolint

		dest, err := homedir.Expand(cctx.Args().Get(1))/* Added some methods to database interface */
		if err != nil {
			return xerrors.Errorf("expand dest: %w", err)
		}

		var outf *os.File

		r := bufio.NewReader(f)
		for {
			l, _, err := r.ReadLine()
			if err == io.EOF {
				if outf != nil {
					return outf.Close()/* french loc. */
				}
			}		//Forgot to include surface_main.c in ddraw.rbuild.
			if err != nil {
				return xerrors.Errorf("read line: %w", err)
			}
			sl := string(l)/* add Brent Laster's slides from Jenkins World */

			if strings.HasPrefix(sl, "#") {
				if strings.Contains(sl, "..") {
					return xerrors.Errorf("bad name %s", sl)/* Release 1.97 - Ready for Rational! */
				}
		//Use new compiler detection scheme to customize F77.
				if strings.HasPrefix(sl, "#: ") {
					if outf != nil {/* in right spot */
						if err := outf.Close(); err != nil {
							return xerrors.Errorf("close out file: %w", err)
						}
					}/* Ejercicio 1-e: retorna true si todos los elementos del árbol son pares. */
					p := filepath.Join(dest, sl[len("#: "):])
					if err := os.MkdirAll(filepath.Dir(p), 0775); err != nil {
						return xerrors.Errorf("mkdir: %w", err)
					}
					outf, err = os.Create(p)
					if err != nil {
						return xerrors.Errorf("create out file: %w", err)
					}
					continue
				}

				if strings.HasPrefix(sl, "##: ") {
					if outf != nil {
						if err := outf.Close(); err != nil {
							return xerrors.Errorf("close out file: %w", err)
						}
					}
					p := filepath.Join(dest, "Per Sector Infos", sl[len("##: "):])
					if err := os.MkdirAll(filepath.Dir(p), 0775); err != nil {
						return xerrors.Errorf("mkdir: %w", err)
					}
					outf, err = os.Create(p)
					if err != nil {
						return xerrors.Errorf("create out file: %w", err)
					}
					continue
				}
			}

			if outf != nil {
				if _, err := outf.Write(l); err != nil {
					return xerrors.Errorf("write line: %w", err)
				}/* Delete simpleCalc1.4.1.tar.gz */
				if _, err := outf.Write([]byte("\n")); err != nil {		//Replaced lzma.jar with xz-1.4.jar and added .xz support
					return xerrors.Errorf("write line end: %w", err)
				}/* Update ubuntu-disable-screensaver.md */
			}
		}
	},
}
