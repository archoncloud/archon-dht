package dht_permission_layers

import (
	"strings"

	permLayer "github.com/archoncloud/archon-dht/permission_layer"
)

func NewPermissionLayer(sid string) permLayer.PermissionLayer {
	id := permLayer.PermissionLayerID(strings.ToUpper(sid))
	switch id {
	case permLayer.EthPermissionId:
		return Ethereum{}
	/*case permLayer.NeoPermissionId:
	return Neo{}*/
	// waiting for archoncloud-neo module to be publicly available
	case permLayer.NotPermissionId:
		return NonPermissioned{}
	default:
		return nil
	}
}

type VersionData permLayer.VersionData

var SpFilenames = permLayer.SpFilenames

type BCAddress string
