package game

import (
	"errors"
	"fmt"
	"github.com/devleesch001/Quantum-go/frames"
	"github.com/devleesch001/Quantum-go/maps"
	"io"
	"log/slog"
	"net"
)

const (
	MaxClients = 128 // Maximum number of clients allowed
)

// Game represents the game server, managing clients and maps.
type Game struct {
	clients        [MaxClients]*client
	maps           maps.Maps
	serverLn       net.Listener
	boxMessageChan chan dataMessage
}

func New() *Game {
	return &Game{
		clients:        [MaxClients]*client{},
		boxMessageChan: make(chan dataMessage, 100),
	}
}

// Close closes the server listener, stopping the game server.
func (g *Game) Close() error {
	if g.serverLn == nil {
		return nil
	}
	return g.serverLn.Close()
}

// handleConnection listens for incoming client connections and adds them to the game.
func (g *Game) handleConnection() {
	for {
		conn, err := g.serverLn.Accept()
		if err != nil {
			slog.Error("Error accepting connection", "error", err)
			continue
		}

		g.addClient(conn)
	}
}

func (g *Game) findFreeSlot() (uint8, error) {
	for i := 1; i < MaxClients; i++ {
		if g.clients[i] == nil {
			return uint8(i), nil
		}
	}
	return 0, errors.New("no free slots available")
}

// addClient creates a new client instance and starts handling its messages.
func (g *Game) addClient(conn net.Conn) {
	slog.Info("New connection established", "address", conn.RemoteAddr().String())

	slotID, err := g.findFreeSlot()
	if err != nil {
		slog.Error("No free slots available for new client", "error", err)
		_ = conn.Close()
		return
	}

	c := newClient(slotID, conn)

	for _, existingClient := range g.clients {
		if existingClient == nil {
			continue
		}

		frame, err := frames.New(existingClient.id, frames.NewClientInfo(
			existingClient.color,
			existingClient.x,
			existingClient.y,
			existingClient.faceID,
			existingClient.bodyID,
			existingClient.legsID,
			existingClient.name,
			true,
		)).MarshalBinary()
		if err != nil {
			continue
		}

		slog.Debug("Sending existing client info", "clientID", existingClient.id, "to", c.conn.RemoteAddr().String())
		_, _ = conn.Write(frame)
	}

	g.clients[c.id] = c

	go g.handleClientMessages(c)
}

func (g *Game) initMap() error {
	var m maps.Maps
	if err := m.Init(); err != nil {
		return err
	}

	g.maps = m

	return nil
}

func (g *Game) Start(addr string) error {
	if err := g.initMap(); err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	slog.Info("Server started on " + addr)

	g.serverLn = ln

	go g.handleConnection()
	go g.handleSendMessages()

	return nil
}

func (g Game) String() string {
	return fmt.Sprintf("maps: %+v", g.maps)
}

func (g *Game) handleClientMessages(c *client) {
	for {
		buf := make([]byte, 4096)
		n, err := c.conn.Read(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				slog.Debug(c.Label()+" Closed connection", "name", c.name, "id", c.id, "address", c.conn.RemoteAddr().String())
			} else if errors.Is(err, io.EOF) {
				slog.Debug(c.Label()+" EOF reached, closing connection", "name", c.name, "id", c.id, "address", c.conn.RemoteAddr().String())
			} else {
				slog.Error(c.Label()+" Error reading from connection", "name", c.name, "id", c.id, "address", c.conn.RemoteAddr().String(), "error", err)
			}

			g.disconnectClient(c)
			c.conn.Close()
			return
		}

		if n == 0 {
			continue
		}

		var raw = buf[:n]
		f, err := frames.UnmarshalFrame(raw)
		if err != nil {
			slog.Error("Impossible de decoder le message", "message", fmt.Sprintf("% 02x", raw), "error", err)
			continue
		}

		slog.Debug(c.Label(), "fd", len(raw), "message", f.String(), "from", c.conn.RemoteAddr().String(), "raw", fmt.Sprintf("% 02x", raw))
		switch payload := f.IPayload.(type) {
		case *frames.ClientHello:
			slog.Info(c.Label(), "id", c.id, "joined", payload.String())

			c.y = payload.Y()
			c.x = payload.X()
			c.color = payload.Color()
			c.faceID = payload.FaceID()
			c.bodyID = payload.BodyID()
			c.legsID = payload.LegsID()
			c.name = payload.Name()

			newResponse, err := frames.New(c.id, frames.NewClientInfo(c.color, c.x, c.y, c.faceID, c.bodyID, c.legsID, c.name, true)).MarshalBinary()

			if err == nil {
				g.boxMessageChan <- dataMessage{
					clientID: c.id,
					data:     newResponse,
				}

			}

		case *frames.ClientPosition:
			slog.Info(c.Label(), "id", c.id, "name", c.name, "pos", "("+payload.String()+")")

			c.x = payload.X()
			c.y = payload.Y()

			frame, err := frames.New(c.id, payload).MarshalBinary()
			if err != nil {
				continue
			}

			g.boxMessageChan <- dataMessage{
				clientID: c.id,
				data:     frame,
			}

		case *frames.ClientMessage:
			slog.Info(c.Label(), "id", c.id, "name", c.name, "said", payload.String())

			frame, err := frames.New(c.id, payload).MarshalBinary()
			if err != nil {
				continue
			}

			g.boxMessageChan <- dataMessage{
				clientID: c.id,
				data:     frame,
			}

		default:
			slog.Warn(c.Label(), "id", c.id, "name", c.name, "unknown", f.String(), "payload", raw)
		}

	}
}

func (g *Game) handleSendMessages() {
	for message := range g.boxMessageChan {
		for _, c := range g.clients {
			if c == nil || c.id == message.clientID {
				continue
			}

			if c.conn == nil {
				slog.Warn("Client connection is nil, skipping", "clientID", c.id)
				go g.disconnectClient(c)
				continue
			}

			write, err := c.conn.Write(message.data)
			if err != nil {
				slog.Error("Error writing to client connection", "clientID", c.id, "error", err)

				go g.disconnectClient(c)
				continue
			}

			slog.Debug("Message sent to client", "clientID", c.id, "bytesWritten", write, "data", fmt.Sprintf("%x", message.data))
		}
	}
}

func (g *Game) disconnectClient(c *client) {
	slog.Info(c.Label()+" disconnected", "id", c.id, "name", c.name, "address", c.conn.RemoteAddr().String())
	binary, err := c.DisconnectFrame()
	if err == nil {
		g.boxMessageChan <- dataMessage{
			clientID: c.id,
			data:     binary,
		}
	}

	g.clients[c.id] = nil
}
