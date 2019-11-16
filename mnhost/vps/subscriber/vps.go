package subscriber

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/util/log"

	vps "vircle.io/mnhost/interface/proto/vps"
)

type Vps struct{}

func (e *Vps) Handle(ctx context.Context, msg *vps.Message) error {
	fmt.Println("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *vps.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}

func NewNode(ctx context.Context, msg *vps.Message) error {
	fmt.Println("Function new node: ", msg.Say)
	return nil
}

func DelNode(ctx context.Context, msg *vps.Message) error {
	fmt.Println("Function del node: ", msg.Say)
	return nil
}
