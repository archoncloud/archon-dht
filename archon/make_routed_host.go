package archon_dht

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"

	ma "github.com/multiformats/go-multiaddr"

	dht "github.com/archoncloud/archon-dht/mods/kad-dht-mod"

	"github.com/archoncloud/archon-dht/mods/libp2p-mod"

	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
)

var (
	AutoNATServiceDialTimeout   = 15 * time.Second
	AutoNATServiceResetInterval = 1 * time.Minute

	AutoNATServiceThrottle = 3
)

// makeRoutedHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true. It will bootstrap using the
// provided PeerInfo
func makeRoutedHost(config DHTConnectionConfig, bootstrapPeers []peer.AddrInfo) (*rhost.RoutedHost, *dht.IpfsDHT, error) {

	// If the seed is 1, will use real cryptographic randomness.
	// Otherwise, use a deterministic source of randomness to make
	// generated keys stay the same across multiple runs
	priv, err := GetRSAKey(config.Seed)
	if err != nil {
		return nil, nil, err
	}

	var myPartialMultiAddress string
	myPartialMultiAddress = config.MyPartialMultiAddress

	preTcp := regexp.MustCompile("^.*/tcp/")
	port := preTcp.ReplaceAllString(myPartialMultiAddress, "")

	// address factory
	publicPeer, err := ma.NewMultiaddr(myPartialMultiAddress)
	if err != nil {
		return nil, nil, err
	}

	addressFactory := func(addrs []ma.Multiaddr) []ma.Multiaddr {
		if publicPeer != nil {
			addrs = append(addrs, publicPeer)
		}
		return addrs
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/" + port),
		libp2p.Identity(priv),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.AddrsFactory(addressFactory),
		libp2p.NATPortMap(),
	}

	ctx := context.Background()
	// note: we are using archon specific modified libp2p libraries
	basicHost, err := libp2p.NewPermissioned(ctx, config.PermissionLayer, opts...)
	if err != nil {
		return nil, nil, err
	}

	// Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
	dstore := dsync.MutexWrap(ds.NewMapDatastore())

	// Make the DHT
	dht := dht.NewDHT(ctx, basicHost, dstore, config.PermissionLayer)

	archonValidator := new(ArchonValidator)
	// see validator.go
	archonValidator.PermissionLayer = config.PermissionLayer
	dht.Validator = archonValidator

	// Make the routed host
	routedHost := rhost.Wrap(basicHost, dht)
	// self checks if registered with smart contract in init,
	// which calls this function
	if !config.IAmBootstrap { // ask george@archon.cloud if curious
		// connect to the chosen dht nodes
		// first curate bootstrap peer list to be only those
		// registered with smart contract
		var validBootstrapPeers []peer.AddrInfo
		if config.PermissionLayer.Permissioned() {
			to := 15 * time.Second
			validatedBootstrapPeers, _ := config.PermissionLayer.ValidatePeers(bootstrapPeers, to) // w timeout
			validBootstrapPeers = validatedBootstrapPeers
		} else {
			validBootstrapPeers = bootstrapPeers
		}

		if (len(validBootstrapPeers) == 1) && (len(validBootstrapPeers[0].Addrs) == 0) {
			return nil, nil, fmt.Errorf("error makeRoutedHost, the bootstrap set has 0 SC validated nodes")
		}

		err = bootstrapConnect(ctx, routedHost, validBootstrapPeers)
		if err != nil {
			return nil, nil, err
		}

		// Bootstrap the host
		err = dht.Bootstrap(ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))
	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addrs := routedHost.Addrs()
	log.Println("I can be reached at:")
	for _, addr := range addrs {
		log.Println(addr.Encapsulate(hostAddr))
	}

	//log.Printf("Now run \"./routed-echo -l %d -d %s%s\" on a different terminal\n", routedHost.ID().Pretty(), config.Global) // FOR DEEP DEBUGGING

	return routedHost, dht, nil
}
