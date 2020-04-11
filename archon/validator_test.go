package archon_dht

import (
	"crypto/sha256"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/libp2p/go-libp2p-core/peer"
)

func TestUrlValidation(t *testing.T) {

	seed := int64(123)
	p, err := GetRSAKey(seed)
	if err != nil {
		assert.Equal(t, 1, 0, err.Error())
	}
	nodeID, err := peer.IDFromPrivateKey(p)
	if err != nil {
		assert.Equal(t, 1, 0, err.Error())
	}
	archonUlKey := "/archonurl/" + nodeID.Pretty()
	exampleUrl := "helloExampleUrl.com"
	pub := p.GetPublic()
	bPub, err := pub.Raw()
	if err != nil {
		assert.Equal(t, 1, 0, err.Error())
	}
	h := sha256.New()
	h.Write([]byte(exampleUrl))
	hashed := h.Sum(nil)
	sig, err := p.Sign(hashed)
	if err != nil {
		assert.Equal(t, 1, 0, err.Error())
	}

	var ULUs UrlsStruct = UrlsStruct{Urls: exampleUrl,
		Signature: sig,
		PublicKey: bPub}
	uploadUrls, err := json.Marshal(ULUs)
	if err != nil {
		assert.Equal(t, 1, 0, err)
	}
	archonValidator := new(ArchonValidator)
	err = archonValidator.Validate(archonUlKey, uploadUrls)
	if err != nil {
		assert.Equal(t, 1, 0, err)
	} else {
		assert.Equal(t, 1, 1, "/archonurl/ validates")
	}

}
