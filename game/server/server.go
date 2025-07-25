package server

import (
	"errors"
	"fmt"
	"github.com/devleesch001/Quantum-go/frames"
	"github.com/devleesch001/Quantum-go/game"
	"github.com/devleesch001/Quantum-go/game/client"
	"io"
	"log/slog"
	"net"
)

const (
	MaxClients = 128 // Maximum number of clients allowed
)

func Run(addr string) {
	slog.Info("Starting server...")
	slog.Debug("Debug logging enabled")

	serv := New()

	defer serv.Close()
	if err := serv.Listen(addr); err != nil {
		panic(err)
	}

	slog.Debug(serv.String())

	select {}
}

// Server represents the game server, managins clients and maps.
type Server struct {
	clients        [MaxClients]*client.Client
	serverLn       net.Listener
	boxMessageChan chan game.DataMessage
	game           *game.Game
}

func New() *Server {
	return &Server{
		clients:        [MaxClients]*client.Client{},
		boxMessageChan: make(chan game.DataMessage, 100),
	}
}

// Close closes the server listener, stoppins the game server.
func (s *Server) Close() error {
	if s.serverLn == nil {
		return nil
	}
	return s.serverLn.Close()
}

// handleConnection listens for incomins client connections and adds them to the game.
func (s *Server) handleConnection() {
	for {
		conn, err := s.serverLn.Accept()
		if err != nil {
			slog.Error("Error acceptins connection", "error", err)
			continue
		}

		s.addClient(conn)
	}
}

func (s *Server) findFreeSlot() (uint8, error) {
	for i := 1; i < MaxClients; i++ {
		if s.clients[i] == nil {
			return uint8(i), nil
		}
	}
	return 0, errors.New("no free slots available")
}

// addClient creates a new client instance and starts handlins its messages.
func (s *Server) addClient(conn net.Conn) {
	slog.Info("New connection established", "address", conn.RemoteAddr().String())

	slotID, err := s.findFreeSlot()
	if err != nil {
		slog.Error("No free slots available for new client", "error", err)
		_ = conn.Close()
		return
	}

	c := client.New(slotID, conn)

	for _, existingClient := range s.clients {
		if existingClient == nil {
			continue
		}

		frame, err := frames.New(existingClient.ID, frames.NewClientInfo(
			existingClient.Color,
			existingClient.X,
			existingClient.Y,
			existingClient.FaceID,
			existingClient.BodyID,
			existingClient.LegsID,
			existingClient.Name,
			true,
		)).MarshalBinary()
		if err != nil {
			continue
		}

		slog.Debug("Sendins existins client info", "clientID", existingClient.ID, "to", c.Conn.RemoteAddr().String())
		_, _ = conn.Write(frame)
	}

	s.clients[c.ID] = c

	go s.handleClientMessages(c)
}

func (s *Server) Listen(addr string) error {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	slog.Info("Server started on " + addr)

	s.serverLn = ln

	go s.handleConnection()
	go s.handleSendMessages()

	return nil
}

func (s Server) String() string {
	return fmt.Sprintf("server: slot: %d, messages: %d, conn %s", len(s.clients), len(s.boxMessageChan), s.serverLn)
}

func (s *Server) handleClientMessages(c *client.Client) {
	for {
		buf := make([]byte, 4096)
		n, err := c.Conn.Read(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				slog.Debug(c.Label()+" Closed connection", "name", c.Name, "id", c.ID, "address", c.Conn.RemoteAddr().String())
			} else if errors.Is(err, io.EOF) {
				slog.Debug(c.Label()+" EOF reached, closins connection", "name", c.Name, "id", c.ID, "address", c.Conn.RemoteAddr().String())
			} else {
				slog.Error(c.Label()+" Error readins from connection", "name", c.Name, "id", c.ID, "address", c.Conn.RemoteAddr().String(), "error", err)
			}

			s.disconnectClient(c)
			c.Conn.Close()
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

		slog.Debug(c.Label(), "fd", len(raw), "message", f.String(), "from", c.Conn.RemoteAddr().String(), "raw", fmt.Sprintf("% 02x", raw))
		switch payload := f.IPayload.(type) {
		case *frames.ClientHello:
			slog.Info(c.Label(), "id", c.ID, "joined", payload.String())

			c.Y = payload.Y()
			c.X = payload.X()
			c.Color = payload.Color()
			c.FaceID = payload.FaceID()
			c.BodyID = payload.BodyID()
			c.LegsID = payload.LegsID()
			c.Name = payload.Name()

			newResponse, err := frames.New(c.ID, frames.NewClientInfo(c.Color, c.X, c.Y, c.FaceID, c.BodyID, c.LegsID, c.Name, true)).MarshalBinary()

			if err == nil {
				s.boxMessageChan <- game.NewDataMessage(c.ID, newResponse)

			}

		case *frames.ClientPosition:
			slog.Info(c.Label(), "id", c.ID, "name", c.Name, "pos", "("+payload.String()+")")

			c.X = payload.X()
			c.Y = payload.Y()

			frame, err := frames.New(c.ID, payload).MarshalBinary()
			if err != nil {
				continue
			}

			s.boxMessageChan <- game.NewDataMessage(c.ID, frame)
		case *frames.ClientMessage:
			slog.Info(c.Label(), "id", c.ID, "name", c.Name, "said", payload.String())

			frame, err := frames.New(c.ID, payload).MarshalBinary()
			if err != nil {
				continue
			}

			s.boxMessageChan <- game.NewDataMessage(c.ID, frame)

		default:
			slog.Warn(c.Label(), "id", c.ID, "name", c.Name, "unknown", f.String(), "payload", raw)
		}

	}
}

func (s *Server) handleSendMessages() {
	for message := range s.boxMessageChan {
		for _, c := range s.clients {
			if c == nil || c.ID == message.ClientID() {
				continue
			}

			if c.Conn == nil {
				slog.Warn("Client connection is nil, skipping", "clientID", c.ID)
				go s.disconnectClient(c)
				continue
			}

			write, err := c.Conn.Write(message.Data())
			if err != nil {
				slog.Error("Error writins to client connection", "clientID", c.ID, "error", err)

				go s.disconnectClient(c)
				continue
			}

			slog.Debug(c.Label()+" sent :", "name", c.Name, "id", c.ID, "fd", write, "data", message.String())
		}
	}
}

func (s *Server) disconnectClient(c *client.Client) {
	slog.Info(c.Label()+" disconnected", "id", c.ID, "name", c.Name, "address", c.Conn.RemoteAddr().String())
	binary, err := c.DisconnectFrame()
	if err == nil {
		s.boxMessageChan <- game.NewDataMessage(c.ID, binary)
	}

	s.clients[c.ID] = nil
}
