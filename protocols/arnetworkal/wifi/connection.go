package wifi

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/phayes/freeport"
)

const (
	// maxUDPDataBytes represents the practical maximum numbers of data bytes in
	// a UDP packet.
	maxUDPDataBytes = 65507
)

var (
	// These are vars instead of a consts so that they can be overridden by unit
	// tests. They're not exported, so there is no danger of anyone else
	// tampering with these.
	deviceIP      = "192.168.42.1"
	discoveryPort = 44444
)

type connectionNegotiationRequest struct {
	D2CPort        int    `json:"d2c_port"`
	ControllerType string `json:"controller_type"`
	ControllerName string `json:"controller_name"`
}

type connectionNegotiationResponse struct {
	Status  int `json:"status"`
	C2DPort int `json:"c2d_port"`
}

type connection struct {
	c2dPort    int
	c2dAddr    *net.UDPAddr
	c2dConn    *net.UDPConn
	d2cPort    int
	d2cConn    *net.UDPConn
	rcvFrameCh chan arnetworkal.Frame
	rcvErrCh   chan error
	rcvStopCh  chan struct{}
	rcvDoneCh  chan struct{}
	// This function is overridable by unit tests
	encodeFrame func(frame arnetworkal.Frame) []byte
	// This function is overridable by unit tests
	decodeData func(data []byte) ([]arnetworkal.Frame, error)
}

// NewConnection returns a UDP/IP based implementation of the
// arnetworkal.Connection interface.
func NewConnection() (arnetworkal.Connection, error) {
	// Select an available port
	d2cPort, err := freeport.GetFreePort()
	if err != nil {
		return nil,
			fmt.Errorf("error selecting available client-side port: %s", err)
	}

	// Negotiate the connection. This is how the client informs the device
	// of the UDP port it will listen on. In response, the device informs
	// the client of which UDP port it will listen on.
	negAddr, err := net.ResolveTCPAddr(
		"tcp",
		fmt.Sprintf("%s:%d", deviceIP, discoveryPort),
	)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error resolving address for connection negotiation: %s",
				err,
			)
	}
	negConn, err := net.DialTCP("tcp", nil, negAddr)
	if err != nil {
		return nil, fmt.Errorf("error negotiating connection: %s", err)
	}
	defer negConn.Close()
	jsonBytes, err := json.Marshal(
		connectionNegotiationRequest{
			D2CPort:        d2cPort,
			ControllerType: "computer",
			ControllerName: "go-parrot",
		},
	)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error marshaling connection negotiation request: %s",
				err,
			)
	}
	jsonBytes = append(jsonBytes, 0x00)
	if _, err := negConn.Write(jsonBytes); err != nil {
		return nil,
			fmt.Errorf(
				"error sending connection negotiation request: %s",
				err,
			)
	}
	data, err := bufio.NewReader(negConn).ReadBytes(0x00)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error receiving connection negotiation response: %s",
				err,
			)
	}
	var negRes connectionNegotiationResponse
	if err := json.Unmarshal(data[:len(data)-1], &negRes); err != nil {
		return nil,
			fmt.Errorf(
				"error unmarshaling connection negotiation response: %s",
				err,
			)
	}
	// Any non-zero status is a refused connection.
	if negRes.Status != 0 {
		return nil,
			errors.New(
				"connection negotiation failed; connection refused by device",
			)
	}

	// Establish an outbound connection...
	c2dAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf("%s:%d", deviceIP, negRes.C2DPort),
	)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error resolving address for outbound connection: %s",
				err,
			)
	}
	c2dConn, err := net.DialUDP("udp", nil, c2dAddr)
	if err != nil {
		return nil,
			fmt.Errorf("error establishing outbound connection: %s", err)
	}

	// Establish an inbound connection...
	d2cAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf(":%d", d2cPort),
	)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error resolving address for inbound connection: %s",
				err,
			)
	}
	d2cConn, err := net.ListenUDP("udp", d2cAddr)
	if err != nil {
		return nil, fmt.Errorf("error establishing inbound connection: %s", err)
	}

	conn := &connection{
		c2dPort:     negRes.C2DPort,
		c2dAddr:     c2dAddr,
		c2dConn:     c2dConn,
		d2cPort:     d2cPort,
		d2cConn:     d2cConn,
		rcvFrameCh:  make(chan arnetworkal.Frame),
		rcvErrCh:    make(chan error),
		rcvStopCh:   make(chan struct{}),
		rcvDoneCh:   make(chan struct{}),
		encodeFrame: defaultEncodeFrame,
		decodeData:  defaultDecodeData,
	}

	go conn.receivePackets()

	return conn, nil
}

func (c *connection) Send(frame arnetworkal.Frame) error {
	if _, err := c.c2dConn.Write(c.encodeFrame(frame)); err != nil {
		return fmt.Errorf("error writing frame to outbound connection: %s", err)
	}
	return nil
}

func (c *connection) receivePackets() {
	defer close(c.rcvDoneCh)
	data := make([]byte, maxUDPDataBytes)
	for {
		select {
		case <-c.rcvStopCh:
			return
		default:
			if err :=
				c.d2cConn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
				select {
				case c.rcvErrCh <- fmt.Errorf("error setting read deadline: %s", err):
					continue
				case <-c.rcvStopCh:
					return
				}
			}
			bytesRead, _, err := c.d2cConn.ReadFromUDP(data)
			if err != nil {
				// Timeouts are ok. We deliberately timeout every three seconds to
				// give ourselves a chance to be interrupted. Handle all other errors.
				if opErr, ok := err.(*net.OpError); !ok || !opErr.Timeout() {
					select {
					case c.rcvErrCh <- fmt.Errorf(
						"error receiving data from inbound connection: %s",
						err,
					):
						continue
					case <-c.rcvStopCh:
						return
					}
				}
				continue
			}
			frames, err := c.decodeData(data[0:bytesRead])
			if err != nil {
				select {
				case c.rcvErrCh <- fmt.Errorf("error decoding inbound data: %s", err):
					continue
				case <-c.rcvStopCh:
					return
				}
			}
			for _, frame := range frames {
				select {
				case c.rcvFrameCh <- frame:
				case <-c.rcvStopCh:
				}
			}
		}
	}
}

func (c *connection) Receive() (arnetworkal.Frame, bool, error) {
	select {
	case frame := <-c.rcvFrameCh:
		return frame, true, nil
	case err := <-c.rcvErrCh:
		return arnetworkal.Frame{}, false, err
	case <-c.rcvDoneCh:
		return arnetworkal.Frame{}, false, nil
	}
}

func (c *connection) Close() error {
	// Signal the goroutine that is receiving packets from the device to stop
	// listening-- we don't want to close the d2c connection while it is.
	close(c.rcvStopCh)
	// Wait for confirmation that the goroutine has stopped listening.
	<-c.rcvDoneCh
	// Now it is safe to close connections.
	errStrs := []string{}
	if err := c.c2dConn.Close(); err != nil {
		errStrs = append(
			errStrs,
			fmt.Sprintf("error closing outbound connection: %s\n", err),
		)
	}
	if err := c.d2cConn.Close(); err != nil {
		errStrs = append(
			errStrs,
			fmt.Sprintf("error closing inbound connection: %s\n", err),
		)
	}
	if len(errStrs) > 0 {
		return fmt.Errorf(
			"error(s) closing connection: %s",
			strings.Join(errStrs, "; "),
		)
	}
	return nil
}