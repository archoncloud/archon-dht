// Package secio is used to encrypt `go-libp2p-conn` connections. Connections wrapped by secio use secure sessions provided by this package to encrypt all traffic. A TLS-like handshake is used to setup the communication channel.
package secio

import (
	"context"
	"net"
	"time"

	dhtPermLayers "github.com/archoncloud/archon-dht/dht_permission_layers"
	permLayer "github.com/archoncloud/archon-dht/permission_layer"

	ci "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/sec"
)

// ID is secio's protocol ID (used when negotiating with multistream)
const ID = "/secio/1.0.0"

// SessionGenerator constructs secure communication sessions for a peer.
type Transport struct {
	LocalID         peer.ID
	PrivateKey      ci.PrivKey
	PermissionLayer permLayer.PermissionLayer
}

func New(sk ci.PrivKey) (*Transport, error) {
	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	return &Transport{
		LocalID:    id,
		PrivateKey: sk,
	}, nil
}

func NewNONSecureTransport(sk ci.PrivKey) (*Transport, error) {
	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	return &Transport{
		LocalID:         id,
		PrivateKey:      sk,
		PermissionLayer: dhtPermLayers.NonPermissioned{},
	}, nil
}

func NewETHSecureTransport(sk ci.PrivKey) (*Transport, error) {
	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	return &Transport{
		LocalID:         id,
		PrivateKey:      sk,
		PermissionLayer: dhtPermLayers.Ethereum{},
	}, nil
}

func NewNEOSecureTransport(sk ci.PrivKey) (*Transport, error) {
	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	return &Transport{
		LocalID:         id,
		PrivateKey:      sk,
		PermissionLayer: dhtPermLayers.Neo{},
	}, nil
}

func SetPermissionedSecureTransport(pid permLayer.PermissionLayerID) (s interface{}) {
	switch pid {
	case permLayer.EthPermissionId:
		return NewETHSecureTransport
	case permLayer.NeoPermissionId:
		return NewNEOSecureTransport
	default:
		return NewNONSecureTransport
	}
}

var _ sec.SecureTransport = (*Transport)(nil)

func (sg *Transport) SecureInbound(ctx context.Context, insecure net.Conn) (sec.SecureConn, error) {
	return sg.newSecureSession(ctx, sg.LocalID, sg.PrivateKey, insecure, "")
}
func (sg *Transport) SecureOutbound(ctx context.Context, insecure net.Conn, p peer.ID) (sec.SecureConn, error) {
	return sg.newSecureSession(ctx, sg.LocalID, sg.PrivateKey, insecure, p)
}

func (s *secureSession) SetReadDeadline(t time.Time) error {
	return s.insecure.SetReadDeadline(t)
}

func (s *secureSession) SetWriteDeadline(t time.Time) error {
	return s.insecure.SetWriteDeadline(t)
}

func (s *secureSession) SetDeadline(t time.Time) error {
	return s.insecure.SetDeadline(t)
}

func (s *secureSession) RemoteAddr() net.Addr {
	return s.insecure.RemoteAddr()
}

func (s *secureSession) LocalAddr() net.Addr {
	return s.insecure.LocalAddr()
}

// LocalPeer retrieves the local peer.
func (s *secureSession) LocalPeer() peer.ID {
	return s.localPeer
}

// LocalPrivateKey retrieves the local peer's PrivateKey
func (s *secureSession) LocalPrivateKey() ci.PrivKey {
	return s.localKey
}

// RemotePeer retrieves the remote peer.
func (s *secureSession) RemotePeer() peer.ID {
	return s.remotePeer
}

// RemotePublicKey retrieves the remote public key.
func (s *secureSession) RemotePublicKey() ci.PubKey {
	return s.remote.permanentPubKey
}
