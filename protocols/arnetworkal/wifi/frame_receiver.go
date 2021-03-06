package wifi

import (
	"net"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/pkg/errors"
)

type frameReceiver struct {
	conn *net.UDPConn
	// This function is overridable by unit tests
	decodeDatagram     func(data []byte) ([]arnetworkal.Frame, error)
	datagramBuffer     []byte
	datagramBufferLock sync.Mutex
}

func (f *frameReceiver) Receive() ([]arnetworkal.Frame, error) {
	f.datagramBufferLock.Lock()
	defer f.datagramBufferLock.Unlock()
	log.Debug("reading / waiting for datagram from d2c connection")
	if err :=
		f.conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return nil, errors.Wrap(err, "error setting read deadline for datagram")
	}
	bytesRead, _, err := f.conn.ReadFromUDP(f.datagramBuffer) // nolint: errcheck
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			// TODO: Fix this-- handle more elegantly and reconnect, if possible!
			log.Fatal("detected a probable disconnection")
		}
		return nil,
			errors.Wrap(err, "error receiving datagram from d2c connection")
	}
	log.WithField(
		"bytesRead", bytesRead,
	).Debug("got datagram from d2c connection")
	// Very imporant-- make a COPY of the data since the datagramBuffer is reused
	// and slices are REFERENCES to a subset of an array or another slice.
	data := make([]byte, bytesRead)
	copy(data, f.datagramBuffer[0:bytesRead])
	return f.decodeDatagram(data)
}

func (f *frameReceiver) Close() {
	if f.conn != nil {
		log.Debug("closing d2c connection")
		if err := f.conn.Close(); err != nil {
			log.Errorf("error closing d2c connection: %s", err)
		}
		log.Debug("closed d2c connection")
	}
}
