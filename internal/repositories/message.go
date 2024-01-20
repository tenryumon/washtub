package repositories

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-memdb"
	"github.com/nsqsink/washtub/internal/models"
	httpLib "github.com/nsqsink/washtub/pkg/http"
)

type MessageStore struct {
	db     *memdb.MemDB
	txn    *memdb.Txn
	schema *memdb.TableSchema
}

func NewMessageStore(db *memdb.MemDB) models.MessageStore {
	return &MessageStore{
		db:     db,
		txn:    db.Txn(false),
		schema: models.MessageSchema,
	}
}

// Fetch implements models.MessageStore.
func (w *MessageStore) Fetch(ctx context.Context, request httpLib.FetchRequest, workerID string) (res []models.Message, err error) {
	// Get all data
	list, err := w.txn.Get(w.schema.Name, "workerid", workerID)
	if err != nil {
		return res, err
	}

	// Parse to object
	for obj := list.Next(); obj != nil; obj = list.Next() {
		res = append(res, obj.(models.Message))
	}

	// Return data
	return res, err
}

// GetByID implements models.MessageStore.
func (w *MessageStore) GetByID(ctx context.Context, id string) (res models.Message, err error) {
	// Lookup by id
	raw, err := w.txn.First("worker", "id", id)
	if err != nil {
		return res, err
	}

	// Parse to object
	res, ok := raw.(models.Message)
	if !ok {
		return res, fmt.Errorf("failed to parse data: %v", raw)
	}

	// Return data
	return res, err

}

// Stash implements models.MessageStore.
func (w *MessageStore) Stash(ctx context.Context, message models.Message) error {
	panic("unimplemented")
}

// Store implements models.MessageStore.
func (w *MessageStore) Store(ctx context.Context, message models.Message) (models.Message, error) {
	// Create a write transaction
	w.txn = w.db.Txn(true)

	// Insert worker
	err := w.txn.Insert(w.schema.Name, message)

	// Commit the transaction
	w.txn.Commit()

	// Create read-only transaction
	w.txn = w.db.Txn(false)
	defer w.txn.Abort()

	return message, err
}
