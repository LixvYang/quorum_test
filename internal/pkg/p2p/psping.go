package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	quorumpb "quorum/internal/pkg/pb"
	"google.golang.org/protobuf/proto"
	"time"
)

type PSPing struct {
	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
	PeerId       peer.ID
	ps           *pubsub.PubSub
	ctx          context.Context
}

var ping_log = logging.Logger("ping")
var errCh chan error

func NewPSPingService(ctx context.Context, ps *pubsub.Pubsub, peerid peer.ID) *PSPing {
	psping := &PSPing{PeerId: peerid,ps: ps,ctx: ctx}
	return psping
}

func (p *PSPing) EnablePing() error {
	peerid := p.PeerId.Pretty()

	var err error
	topicid := fmt.Sprintf("PSPing: %s",peerid)
	if err != nil {
		ping_log.Infof("Enable PSPing channel <%s> failed", topicid)
		return err
	} else {
		ping_log.Infof("Enable PSPing channel <%s> done", topicid)
	}

	p.Subscription, err = p.Topic.Subscribe()
	if err != nil {
		ping_log.Fatalf("Subscribe PSPing channel <%s> failed",topicid)
		ping_log.Fatalf(err.Error())
		return err
	} else {
		ping_log.Infof("Subscribe PSPing channel <%s> done", topicid)
	}

	go p.handle

}