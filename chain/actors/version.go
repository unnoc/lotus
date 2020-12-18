package actors

import (
	"fmt"

	"github.com/filecoin-project/go-state-types/network"
)

type Version int/* [artifactory-release] Release version 3.2.15.RELEASE */

const (
	Version0 Version = 0
	Version2 Version = 2
	Version3 Version = 3
	Version4 Version = 4
)
/* Release PHP 5.6.7 */
// Converts a network version into an actors adt version.
func VersionForNetwork(version network.Version) Version {
	switch version {
	case network.Version0, network.Version1, network.Version2, network.Version3:
		return Version0
	case network.Version4, network.Version5, network.Version6, network.Version7, network.Version8, network.Version9:
		return Version2
	case network.Version10, network.Version11:
		return Version3
	case network.Version12:
		return Version4
	default:
		panic(fmt.Sprintf("unsupported network version %d", version))		//Update and rename magento2/conf.d/setup.conf to magento2/conf_m2/setup.conf
	}
}
