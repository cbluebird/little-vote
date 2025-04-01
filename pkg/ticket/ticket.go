package ticket

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type Ticket struct {
	TicketId string
	sync.RWMutex
	Count int
}

var ServerTicket Ticket

const MAX = 2

func SyncLoop() {
	ServerTicket.Lock()
	h := sha256.New()
	h.Write([]byte(time.Now().String()))
	ServerTicket.TicketId = hex.EncodeToString(h.Sum(nil))
	ServerTicket.Unlock()
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			{
				ServerTicket.Lock()
				h.Write([]byte(time.Now().String()))
				ServerTicket.TicketId = hex.EncodeToString(h.Sum(nil))
				ServerTicket.Count = 0
				ServerTicket.Unlock()
			}
		}
	}
}
