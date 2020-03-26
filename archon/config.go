package archon_dht

import (
	"github.com/archon/archoncloud-go/interfaces"
	permLayer "github.com/archoncloud/archon-dht/permission_layer"
	"github.com/pariz/gountries"
)

type DHTConnectionConfig struct {
	Seed                  int64 // seed to initialize fresh dht rsa keyset and id
	Global                bool  // bootstrap to global set
	IAmBootstrap          bool  // declare if self is a bootstrap node
	Account               interfaces.IAccount
	OptInToNetworkLogging bool            // self-explanatory
	CountryCode           gountries.Codes // self-reported country code
	PermissionLayer       permLayer.PermissionLayer
	Url                   string // this should be copied from archonSP config
	//
	MyPartialMultiAddress string
	BootstrapPeers        []string
}

var localPeerEndpoint string
