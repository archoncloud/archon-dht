package dht_permission_layers

import (
	"time"

	permLayer "github.com/archoncloud/archon-dht/permission_layer"

	"github.com/libp2p/go-libp2p-core/peer"
)

const updateSPProfileInterval time.Duration = 30 * time.Minute

// NONPERMISSIONED
type NonPermissioned struct {
}

func (n NonPermissioned) ID() permLayer.PermissionLayerID {
	return permLayer.NotPermissionId
}

func (n NonPermissioned) Permissioned() bool {
	return false
}

func (n NonPermissioned) ValidatePeersPtrArr(bootstrapPeers []*peer.AddrInfo, timeout time.Duration) (res []*peer.AddrInfo, err error) {
	return nil, nil
}

func (n NonPermissioned) ValidatePeer(pid peer.ID, timout time.Duration) (bool, error) {
	return true, nil
}

func (n NonPermissioned) ValidatePeers(bootstrapPeers []peer.AddrInfo, timeout time.Duration) (res []peer.AddrInfo, err error) {
	return nil, nil
}

func (n NonPermissioned) UpdateIndividualSPProfileCache(pid peer.ID) {
	// TODO
}

func (n NonPermissioned) UpdateSPProfileCache(pids []peer.ID) {
	// do nothing
}

func (n NonPermissioned) CompareBlockHeights(lhs, rhs permLayer.VersionData) (int, error) {
	//  TODO
	return 0, nil
}

func (n NonPermissioned) GetBlockHeight() (string, error) {
	// TODO
	return "", nil
}

func (n NonPermissioned) GetBlockHash(blockHeight string) (string, error) {
	// TODO
	return "", nil
}

func (n NonPermissioned) NewVersionData() (*permLayer.VersionData, error) {
	// Not needed, as non-permissioned never stores data
	return nil, nil
}

// END NONPERMISSIONED
