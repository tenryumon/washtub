package repositories

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-memdb"
	"github.com/nsqsink/washtub/internal/models"
	httpLib "github.com/nsqsink/washtub/pkg/http"
)

type WorkerStore struct {
	db     *memdb.MemDB
	txn    *memdb.Txn
	schema *memdb.TableSchema
}

// Init New Worker WorkerStore from in-memory DB
func NewWorkerStore(db *memdb.MemDB) models.WorkerStore {
	return &WorkerStore{
		db:     db,
		txn:    db.Txn(false),
		schema: models.WorkerSchema,
	}
}

// Fetch to get all worker data from WorkerStore
func (w *WorkerStore) Fetch(ctx context.Context, request httpLib.FetchRequest) (res []models.Worker, err error) {
	// Get all data
	list, err := w.txn.Get(w.schema.Name, "id")
	if err != nil {
		return res, err
	}

	// Parse to object
	for obj := list.Next(); obj != nil; obj = list.Next() {
		res = append(res, obj.(models.Worker))
	}

	// Return data
	return res, err
}

// GetByID to get data by worker id from WorkerStore.
func (w *WorkerStore) GetByID(ctx context.Context, id string) (res models.Worker, err error) {
	// Lookup by id
	raw, err := w.txn.First("worker", "id", id)
	if err != nil || raw == nil {
		return res, err
	}

	// Parse to object
	res, ok := raw.(models.Worker)
	if !ok {
		return res, fmt.Errorf("failed to parse data: %v", raw)
	}

	// Return data
	return res, err
}

// Stash to clear worker data from WorkerStore.
func (w *WorkerStore) Stash(ctx context.Context, worker models.Worker) error {
	panic("unimplemented")
}

// Store to insert worker data to WorkerStore.
func (w *WorkerStore) Store(ctx context.Context, worker models.Worker) (models.Worker, error) {
	// Create a write transaction
	w.txn = w.db.Txn(true)

	// Insert worker
	err := w.txn.Insert(w.schema.Name, worker)

	// Commit the transaction
	w.txn.Commit()

	// Create read-only transaction
	w.txn = w.db.Txn(false)
	defer w.txn.Abort()

	return worker, err
}
