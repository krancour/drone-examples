package arnetwork

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

type c2dBuffer struct {
	C2DBufferConfig
	*buffer
	conn  arnetworkal.Connection
	seq   uint8
	ackCh chan Frame
}

func newC2DBuffer(
	bufCfg C2DBufferConfig,
	conn arnetworkal.Connection,
) *c2dBuffer {
	buf := &c2dBuffer{
		C2DBufferConfig: bufCfg,
		buffer:          newBuffer(bufCfg.Size, bufCfg.IsOverwriting),
		conn:            conn,
	}

	go buf.writeFrames()

	return buf
}

func (c *c2dBuffer) writeFrames() {
	for frame := range c.outCh {
		c.writeFrame(frame)
	}
}

func (c *c2dBuffer) writeFrame(frame Frame) {
	for attempts := 0; attempts <= c.MaxRetries || c.MaxRetries == -1; attempts++ {
		netFrame := arnetworkal.Frame{
			ID:   c.ID,
			Type: c.FrameType,
			Seq:  c.seq,
			Data: frame.Data,
		}
		c.seq++
		if err := c.conn.Send(netFrame); err != nil {
			log.Printf("error sending frame: %s\n", err)
		}
		if netFrame.Type == arnetworkal.FrameTypeDataWithAck {
			select {
			case ack := <-c.ackCh:
				if bytes.Equal(
					[]byte(fmt.Sprintf("%d", netFrame.Seq)),
					ack.Data,
				) {
					return
				}
			case <-time.After(c.AckTimeout):
			}
		}
	}
}