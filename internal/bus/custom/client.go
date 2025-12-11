package custom

import (
	"b3e/internal/core/command"
	"context"
	"encoding/json"
	"fmt"
	"net"
)

type TCPBus struct {
	address string
}

func NewTCPBus(addr string) *TCPBus {
	return &TCPBus{address: addr}
}

func (b *TCPBus) Publish(ctx context.Context, cmd command.CastVote) error {
	conn, err := net.Dial("tcp", b.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, _ := json.Marshal(cmd)
	fmt.Fprintf(conn, "PUSH %s\n", string(data))

	return nil
}
