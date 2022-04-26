package p2p

import (
	"errors"
	"fmt"
	"github.com/ChainSafe/chainbridge-core/config/relayer"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peerstore"
	noise "github.com/libp2p/go-libp2p-noise"
	"github.com/rs/zerolog/log"
)

// NewHost creates new host.Host from private key and relayer configuration
func NewHost(privKey crypto.PrivKey, rconf relayer.MpcRelayerConfig) (host.Host, error) {
	if privKey == nil {
		return nil, errors.New("unable to create libp2p host: private key not defined")
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", rconf.Port)),
		libp2p.Identity(privKey),
		libp2p.DisableRelay(),
		libp2p.Security(noise.ID, noise.New),
	}

	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to create libp2p host: %v", err)
	}

	log.Info().Str("peerID", h.ID().Pretty()).Msgf(
		"new libp2p host created with address: %s", h.Addrs()[0].String(),
	)

	for _, p := range rconf.Peers {
		h.Peerstore().AddAddr(p.ID, p.Addrs[0], peerstore.PermanentAddrTTL)
	}
	return h, nil
}
