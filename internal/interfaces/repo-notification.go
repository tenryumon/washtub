package interfaces

import (
	"context"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

type Notification interface {
	Send(ctx context.Context, channel entities.Channel, metadata map[string]interface{}) error
	SendMulti(ctx context.Context, channels []entities.Channel, metadata map[string]interface{}) error
}
