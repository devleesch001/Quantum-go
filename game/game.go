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

// Game represents the game server, managing clients and maps.
type Game struct {
	clientCount    uint8
	clients        []*client
	maps           maps.Maps
	serverLn       net.Listener
	boxMessageChan chan dataMessage
}

// Close closes the server listener, stopping the game server.
func (g *Game) Close() error {
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

// addClient creates a new client instance and starts handling its messages.
func (g *Game) addClient(conn net.Conn) {
	slog.Info("New connection established", "address", conn.RemoteAddr().String())

	g.clientCount++

	c := newClient(g.clientCount, conn)

	for _, existingClient := range g.clients {
		frame, err := frames.New(c.id, frames.NewClientInfo(
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

	g.clients = append(g.clients, c)

	go g.HandleClientMessages(c)
}

func (g *Game) initMap() error {
	var m maps.Maps
	if err := m.Init(); err != nil {
		return err
	}

	g.maps = m

	return nil
}

func (g *Game) Start(port uint16) error {
	if err := g.initMap(); err != nil {
		return err
	}

	var strAddr = fmt.Sprintf("0.0.0.0:%d", port)

	ln, err := net.Listen("tcp", strAddr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	slog.Info("Server started on " + strAddr)

	g.serverLn = ln
	g.boxMessageChan = make(chan dataMessage, 100)

	go g.handleConnection()
	go g.HandleSendMessages()

	return nil
}

func (g Game) String() string {
	return fmt.Sprintf("maps: %+v", g.maps)
}

func (g *Game) HandleClientMessages(c *client) {
	for {

		buf := make([]byte, 4096)
		n, err := c.conn.Read(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				slog.Info("connection closed", "clientID", c.id, "address", c.conn.RemoteAddr().String())
			} else if errors.Is(err, io.EOF) {
				slog.Info("EOF reached, closing connection", "clientID", c.id, "address", c.conn.RemoteAddr().String())
			} else {
				slog.Error("Error reading from connection", "clientID", c.id, "address", c.conn.RemoteAddr().String(), "error", err)
			}

			newResponse, err := frames.New(c.id, frames.NewClientInfo(c.color, c.x, c.y, c.faceID, c.bodyID, c.legsID, c.name, false)).MarshalBinary()

			if err == nil {
				g.boxMessageChan <- dataMessage{
					clientID: c.id,
					data:     newResponse,
				}

			}

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

		slog.Debug("Received message", "message", f.String(), "from", c.conn.RemoteAddr().String())
		slog.Debug(fmt.Sprintf("%x", raw))

		switch payload := f.IPayload.(type) {
		case *frames.ClientHello:
			slog.Info("Client Hello", "client", payload.String())

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
			slog.Info("Client Position", "position", payload.String())

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
			slog.Info("Client Message", "message", payload.String())

			frame, err := frames.New(c.id, payload).MarshalBinary()
			if err != nil {
				continue
			}

			g.boxMessageChan <- dataMessage{
				clientID: c.id,
				data:     frame,
			}

		default:
			slog.Warn("Unknown", "message", f.String())
		}

	}
}

func (g *Game) HandleSendMessages() {
	for message := range g.boxMessageChan {
		for _, c := range g.clients {
			if c.id == message.clientID {
				slog.Debug("Skipping sending message to self", "clientID", c.id)
				continue
			}

			if c.conn == nil {
				slog.Warn("Client connection is nil, skipping", "clientID", c.id)
				continue
			}

			write, err := c.conn.Write(message.data)
			if err != nil {
				slog.Error("Error writing to client connection", "clientID", c.id, "error", err)
				continue
			}
			slog.Debug("Message sent to client", "clientID", c.id, "bytesWritten", write)
		}
	}
}
