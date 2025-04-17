package states

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"server/internal/server"
	"server/internal/server/objects"
	"server/pkg/packets"
	"time"
)

type InGame struct {
	client                 server.ClientInterfacer
	player                 *objects.Player
	logger                 *log.Logger
	cancelPlayerUpdateLoop context.CancelFunc
}

func (g *InGame) Name() string {
	return "InGame"
}

func (g *InGame) SetClient(client server.ClientInterfacer) {
	g.client = client
	loggingPrefix := fmt.Sprintf("Client %d [%s]: ", client.Id(), g.Name())
	g.logger = log.New(log.Writer(), loggingPrefix, log.LstdFlags)
}

func (g *InGame) OnEnter() {
	g.logger.Printf("Adding Player to the %s shared collection", g.player.Name)
	go g.client.SharedGameObjects().Players.Add(g.player, g.client.Id())

	// Set the initial properties of the player
	g.player.X = rand.Float64() * 1000
	g.player.Y = rand.Float64() * 1000
	g.player.Speed = 150.0
	g.player.Radius = 20.0

	// Send the player's initial state to the client
	g.client.SocketSend(packets.NewPlayer(g.client.Id(), g.player))

	// Send the spores to the client in the background
	// Added logic send in batches
	go g.sendInitialSpores(20, 50*time.Millisecond)

}

func (g *InGame) HandlerMessage(senderId uint64, message packets.Msg) {
	switch message := message.(type) {
	case *packets.Packet_Player:
		g.handlePlayer(senderId, message)
	case *packets.Packet_PlayerDirection:
		g.handlePlayerDirection(senderId, message)
	case *packets.Packet_Chat:
		g.handleChatMessage(senderId, message)
	case *packets.Packet_SporeConsumed:
		g.logger.Printf("Spore %d consumed by client %d", message.SporeConsumed.SporeId, senderId)
	}
}

func (g *InGame) handlePlayerDirection(senderId uint64, message *packets.Packet_PlayerDirection) {
	if senderId == g.client.Id() {
		g.player.Direction = message.PlayerDirection.Direction

		// If this is the first time receiving a player direction message from our client, start the player update loop
		if g.cancelPlayerUpdateLoop == nil {
			ctx, cancel := context.WithCancel(context.Background())
			g.cancelPlayerUpdateLoop = cancel
			go g.playerUpdateLoop(ctx)
		}
	}
}

func (g *InGame) handleChatMessage(senderId uint64, message *packets.Packet_Chat) {
	if senderId == g.client.Id() {
		g.client.Broadcast(message)
	} else {
		g.client.SocketSendAs(message, senderId)
	}
}
func (g *InGame) handlePlayer(senderId uint64, message *packets.Packet_Player) {
	if senderId == g.client.Id() {
		g.logger.Println("Received player message from our own client, ignoring")
		return
	}
	g.client.SocketSendAs(message, senderId)
}

func (g *InGame) playerUpdateLoop(ctx context.Context) {
	const delta float64 = 0.05
	ticker := time.NewTicker(time.Duration(delta*1000) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			g.syncPlayer(delta)
		case <-ctx.Done():
			return
		}
	}
}

func (g *InGame) OnExit() {
	if g.cancelPlayerUpdateLoop != nil {
		g.cancelPlayerUpdateLoop()
	}
	g.client.SharedGameObjects().Players.Remove(g.client.Id())
}

func (g *InGame) syncPlayer(delta float64) {
	newX := g.player.X + g.player.Speed*math.Cos(g.player.Direction)*delta
	newY := g.player.Y + g.player.Speed*math.Sin(g.player.Direction)*delta

	g.player.X = newX
	g.player.Y = newY

	updatePacket := packets.NewPlayer(g.client.Id(), g.player)
	g.client.Broadcast(updatePacket)
	go g.client.SocketSend(updatePacket)
}

func (g *InGame) sendInitialSpores(batchSize int, delayBetweenBatches time.Duration) {
	sporesBatch := make(map[uint64]*objects.Spore, batchSize)

	g.client.SharedGameObjects().Spores.ForEach(func(sporeId uint64, spore *objects.Spore) {
		sporesBatch[sporeId] = spore

		if len(sporesBatch) >= batchSize {
			g.client.SocketSend(packets.NewSporesBatch(sporesBatch))
			sporesBatch = make(map[uint64]*objects.Spore, batchSize)
			time.Sleep(delayBetweenBatches)
		}
	})
	// Send any remaining spores
	if len(sporesBatch) > 0 {
		g.client.SocketSend(packets.NewSporesBatch(sporesBatch))
	}
}
