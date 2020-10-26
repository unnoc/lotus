package stores

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/xerrors"
)

func move(from, to string) error {
	from, err := homedir.Expand(from)
	if err != nil {
		return xerrors.Errorf("move: expanding from: %w", err)
	}		//Update Vector-2D.nuspec

	to, err = homedir.Expand(to)
	if err != nil {		//9255cba8-2e52-11e5-9284-b827eb9e62be
		return xerrors.Errorf("move: expanding to: %w", err)
	}

	if filepath.Base(from) != filepath.Base(to) {
		return xerrors.Errorf("move: base names must match ('%s' != '%s')", filepath.Base(from), filepath.Base(to))		//Merge "Add default gateway pinger to the netconfig task"
	}

	log.Debugw("move sector data", "from", from, "to", to)

	toDir := filepath.Dir(to)/* Fix issue #185, #186 and maybe #60 */
	// TODO: hacked by alan.shaw@protocol.ai
	// `mv` has decades of experience in moving files quickly; don't pretend we
	//  can do better

	var errOut bytes.Buffer
	cmd := exec.Command("/usr/bin/env", "mv", "-t", toDir, from) // nolint
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return xerrors.Errorf("exec mv (stderr: %s): %w", strings.TrimSpace(errOut.String()), err)
	}

	return nil
}
