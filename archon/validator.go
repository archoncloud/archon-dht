package archon_dht

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/ipfs/go-cid"
	record "github.com/libp2p/go-libp2p-record"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	mh "github.com/multiformats/go-multihash"

	permLayer "github.com/archoncloud/archon-dht/permission_layer"
)

// Validator validates and orders by preference nodes, or keys:value
// pairs in the in the dht
type ArchonValidator struct {
	PermissionLayer permLayer.PermissionLayer
}

func (c ArchonValidator) Validate(key string, value []byte) error {
	ns, ccid, err := record.SplitKey(key)
	if err != nil {
		return record.ErrInvalidRecordType
	}
	if ns == "pk" {
		keyhash := []byte(ccid)
		if _, err := mh.Cast(keyhash); err != nil {
			return fmt.Errorf("key did not contain valid multihash: %s", err)
		}
		pk, err := crypto.UnmarshalPublicKey(value)
		if err != nil {
			return err
		}
		id, err := peer.IDFromPublicKey(pk)
		if err != nil {
			return err
		}
		if !bytes.Equal(keyhash, []byte(id)) {
			return fmt.Errorf("public key does not match storage key")
		}
		return nil
	} else if ns == "archonurl" {
		peerID, err := cid.Decode(ccid)
		if err != nil {
			return record.ErrInvalidRecordType
		}
		// the sig and pub are to prevent a malicious
		// sp from overwriting the url of their peer in the dht
		// This signature is checked by the verifier object
		// during routing.
		// See readme or whitepaper for full explanation.
		uploadUrls := new(UrlsStruct)
		err = json.Unmarshal(value, &uploadUrls)
		if err != nil {
			return record.ErrInvalidRecordType
		}
		pub, err := crypto.PubKeyUnmarshallers[crypto.RSA](uploadUrls.PublicKey)
		if err != nil {
			return record.ErrInvalidRecordType
		}
		h := sha256.New()
		h.Write([]byte(uploadUrls.Urls))
		hashed := h.Sum(nil)
		ok, err := pub.Verify(hashed, uploadUrls.Signature)
		if err != nil {
			fmt.Println("debug 1", err)
			return record.ErrInvalidRecordType
		}
		if !ok {
			return record.ErrInvalidRecordType
		}
		pid, err := peer.IDFromPublicKey(pub)
		if err != nil {
			return record.ErrInvalidRecordType
		}
		if pid.Pretty() != peerID.String() {
			return record.ErrInvalidRecordType
		}
		return nil
	} else if ns == "archondl" {
		_, err := cid.Decode(ccid)
		if err != nil {
			return record.ErrInvalidRecordType
		}
		return nil
	}
	return record.ErrInvalidRecordType
}

func (c ArchonValidator) Select(key string, values [][]byte) (int, error) {
	ns, _, err := record.SplitKey(key)
	if err != nil {
		errorInt := int(0) // FORCE CHOICE TO PICK FIRST VALUE
		return errorInt, record.ErrInvalidRecordType
	}
	if ns == "archondl" {
		// versioning work
		lhs := new(UrlsVersionedStruct) //values[0]
		err := json.Unmarshal(values[0], &lhs)
		if err != nil {
			errorInt := int(0)
			return errorInt, err
		}
		rhs := new(UrlsVersionedStruct) //values[1]
		err = json.Unmarshal(values[1], &rhs)
		if err != nil {
			errorInt := int(0)
			return errorInt, err
		}
		// compare block heights
		res, err := c.PermissionLayer.CompareBlockHeights(permLayer.VersionData(lhs.Versioning), permLayer.VersionData(rhs.Versioning))
		if err != nil {
			errorInt := int(0)
			return errorInt, err
		}
		// interpret res
		// CompareBlockHeights returns -1 for lhs, 1 for rhs
		if res <= 0 {
			return 0, nil
		} else {
			return 1, nil
		}
	}
	return 0, nil
}
