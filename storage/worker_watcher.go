package storage

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/coreos/etcd/clientv3"
	"github.com/cybozu-go/neco"
)

// WorkerWatcher has callback handlers to handle status changes
type WorkerWatcher struct {
	handleStatus func(context.Context, int, *neco.UpdateStatus) (bool, error)
}

// NewWorkerWatcher creates a new WorkerWatcher
func NewWorkerWatcher(
	handleStatus func(context.Context, int, *neco.UpdateStatus) (bool, error),
) WorkerWatcher {
	return WorkerWatcher{
		handleStatus: handleStatus,
	}
}

// Watch watches worker changes until deadline is reached.
// If the handleStatus returns (true, nil) this returns nil
// Otherwise non-nil error is returned
func (w WorkerWatcher) Watch(ctx context.Context, rev int64, storage Storage) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := storage.etcd.Watch(
		ctx, KeyWorkerStatusPrefix,
		clientv3.WithRev(rev+1), clientv3.WithFilterDelete(), clientv3.WithPrefix(),
	)
	for resp := range ch {
		for _, ev := range resp.Events {
			st := new(neco.UpdateStatus)
			err := json.Unmarshal(ev.Kv.Value, st)
			if err != nil {
				return err
			}
			lrn, err := strconv.Atoi(string(ev.Kv.Key[len(KeyWorkerStatusPrefix):]))
			if err != nil {
				return err
			}
			completed, err := w.handleStatus(ctx, lrn, st)
			if err != nil {
				return err
			}
			if completed {
				return nil
			}
		}
	}
	return ErrTimedOut
}
