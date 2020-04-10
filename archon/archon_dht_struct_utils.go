package archon_dht

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ipfs/go-cid"

	permLayer "github.com/archoncloud/archon-dht/permission_layer"

	"github.com/archoncloud/archon-dht/common"

	"github.com/libp2p/go-libp2p-core/peer"
	mh "github.com/multiformats/go-multihash"
)

func (a ArchonDHTs) pollUpdateSPProfileCache(interval time.Duration) {
	go func(i time.Duration, arc ArchonDHTs) {
		for {
			arc.updateSPProfileCaches()
			time.Sleep(i)
		}
	}(interval, a)
}

func (a ArchonDHTs) updateSPProfileCaches() {
	for k, v := range a.Layers {
		if k != "NON" {
			v.updateSPProfileCache()
		}
	}
}

// Each service provider "puts" the key:value pair nodeID:url into the dht
// so that uploaders will have a url to upload data to after the marketplace
// resolves the assignment of an upload's shards to sps
func (a ArchonDHTs) pollAnnounceUrl(interval time.Duration) {
	go func(i time.Duration, arc ArchonDHTs) {
		for {
			arc.announceUrl()
			time.Sleep(i)
		}
	}(interval, a)
}

func (a ArchonDHTs) announceUrl() {
	var nodeID string
	foundAnID := false
	for foundAnID == false {
		for _, v := range a.Layers {
			if v.routedHost != nil {
				nodeID = string(v.routedHost.ID())
				foundAnID = true
				break // we only want one
			}
		}
	}
	keyhash := []byte(nodeID)
	nodeIDAsMh, err := mh.Cast(keyhash)
	if err != nil {
		common.LogError.Println(err)
	}
	nodeIDAsCid := cid.NewCidV0(nodeIDAsMh)
	time.Sleep(4 * time.Second)
	// delay in case network just booted
	_ = a.putValue(nodeIDAsCid, "/archonurl/")
}

func (a *ArchonDHTs) putValue(keyAsCid cid.Cid, archonPrefix string) error { // archonPrefix is /archonurl/
	var wg sync.WaitGroup
	wg.Add(len(a.Layers))
	errMessage := make(chan error, len(a.Layers))
	for _, v := range a.Layers {
		go func(vDht *ArchonDHT, wwg *sync.WaitGroup) {
			defer wwg.Done()
			// wait until bootstrapped
			for {
				if vDht.dHT != nil {
					if vDht.dHT.HasPeers() {
						break
					} else {
						select {
						case <-time.After(200 * time.Millisecond):
							continue
						}
					}
				} else {
					select {
					case <-time.After(200 * time.Millisecond):
						continue
					}
				}
			}
			var err error
			if archonPrefix == "/archonurl/" {
				err = vDht.putUrl(archonPrefix, keyAsCid)
			}
			errMessage <- err
		}(v, &wg)
	}
	wg.Wait()
	var errString string
	for i := 0; i < len(a.Layers); i++ {
		e := <-errMessage
		if e != nil {
			errString += e.Error()
		}
	}
	if len(errString) > 0 {
		return fmt.Errorf(errString)
	}
	return nil
}

// putValueVersioned puts the key:value pair contentID:{url, versionData}
// corresponding to the upload of content into the archon cloud
func (a *ArchonDHTs) putValueVersioned(archonPrefix string, keyAsCid cid.Cid, versionData permLayer.VersionData) error { // archonPrefix is /archondl/
	var wg sync.WaitGroup
	wg.Add(len(a.Layers))
	errMessage := make(chan error, len(a.Layers))
	for _, v := range a.Layers {
		go func(vDht *ArchonDHT, wwg *sync.WaitGroup) {
			defer wwg.Done()
			// wait until bootstrapped
			for {
				if vDht.dHT != nil {
					if vDht.dHT.HasPeers() {
						break
					} else {
						select {
						case <-time.After(200 * time.Millisecond):
							continue
						}
					}
				} else {
					select {
					case <-time.After(200 * time.Millisecond):
						continue
					}
				}
			}
			var err error
			if archonPrefix == "/archondl/" {
				err = vDht.putUrlVersioned(archonPrefix, keyAsCid, versionData)
			}
			errMessage <- err
		}(v, &wg)
	}
	wg.Wait()
	var errString string
	for i := 0; i < len(a.Layers); i++ {
		e := <-errMessage
		if e != nil {
			errString += e.Error()
		}
	}
	if len(errString) > 0 {
		return fmt.Errorf(errString)
	}
	return nil
}

type bundled struct {
	PeerAddrInfo UrlArray
	Error        error
	Key          permLayer.PermissionLayerID
}

///////////////////////////////////////////////////////////////////////////
/// below are on single layer /////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////
/// this can be seen as the difference between ArchonDHTs and ArchonDHT ///
///////////////////////////////////////////////////////////////////////////

func (a *ArchonDHT) updateSPProfileCache() {
	p := a.Peers()
	var connectedPeers []peer.ID
	for i := 0; i < len(p); i++ {
		connected := a.routedHost.Network().Connectedness(p[i])
		if connected == 1 {
			connectedPeers = append(connectedPeers, p[i])
		}
	}
	connectedPeers = append(connectedPeers, a.routedHost.ID()) // self
	a.Config.PermissionLayer.UpdateSPProfileCache(connectedPeers)
	a.updateUrls(connectedPeers)
}

func (a *ArchonDHT) updateUrls(connectedPeers []peer.ID) {
	go func(ps []peer.ID) {
		sps := make([]string, len(ps))
		for i := 0; i < len(ps); i++ {
			sps = append(sps, ps[i].Pretty())
		}
		to := 3 * time.Second
		_, _ = a.getUrls(sps, to)
	}(connectedPeers)
}

// called by putValue
func (d *ArchonDHT) putUrl(archonPrefix string, keyAsCid cid.Cid) error {
	var archonUlKey string = archonPrefix + keyAsCid.String()
	p, err := GetRSAKey(d.Config.Seed)
	if err != nil {
		return err
	}
	pub := p.GetPublic()
	bPub, err := pub.Raw()
	if err != nil {
		return err
	}
	sig, err := p.Sign([]byte(d.Config.Url))
	if err != nil {
		return err
	}
	var ULUs UrlsStruct = UrlsStruct{Urls: d.Config.Url,
		Signature: sig,
		PublicKey: bPub}
	// the sig and pub are to prevent a malicious
	// sp from overwriting the url of their peer in the dht
	// This signature is checked by the verifier object
	// during routing.
	// See readme or whitepaper for full explanation.
	uploadUrls, err := json.Marshal(ULUs)
	if err != nil {
		return err
	}
	return d.dHT.PutValue(context.Background(), archonUlKey, uploadUrls)
}

func (d *ArchonDHT) putUrlVersioned(archonPrefix string, keyAsCid cid.Cid, versionData permLayer.VersionData) error {
	var archonUlKey string = archonPrefix + keyAsCid.String()
	p, err := GetRSAKey(d.Config.Seed)
	if err != nil {
		return err
	}
	pub := p.GetPublic()
	bPub, err := pub.Raw()
	if err != nil {
		return err
	}
	bVersionData, err := json.Marshal(versionData)
	if err != nil {
		return err
	}
	sig, err := p.Sign([]byte(d.Config.Url + string(bVersionData)))
	if err != nil {
		return err
	}
	var ULUs UrlsVersionedStruct = UrlsVersionedStruct{Urls: d.Config.Url,
		Versioning: versionData,
		Signature:  sig,
		PublicKey:  bPub}
	// the sig and pub are a safety precaution similar to the
	// one above. This is to ensure that the sp that calls
	// "stored" on this data is held accountable by way of their
	// signature. Later in the protocol, when the downloader
	// obtains download from this info, the download receipt
	// combined with this signature can be used as a proof of
	// misconduct if the version or data have been tampered with.
	// See readme or whitepaper for full explanation
	downloadUrlsVersioned, err := json.Marshal(ULUs)
	if err != nil {
		return err
	}
	return d.dHT.PutValue(context.Background(), archonUlKey, downloadUrlsVersioned)
}
