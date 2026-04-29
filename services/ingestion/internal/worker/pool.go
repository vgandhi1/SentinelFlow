package worker

import (
	"context"
	"sync"
)

// RawMessage is the unit of work: raw bytes (Protobuf) from MQTT.
type RawMessage struct {
	Payload []byte
}

// Handler processes a batch of raw messages (e.g. unmarshal Protobuf, send to Kafka).
type Handler interface {
	HandleBatch(ctx context.Context, batch []RawMessage) error
}

// Pool runs a fixed number of workers that pull from ch and call Handler in batches.
func Pool(ctx context.Context, ch <-chan RawMessage, h Handler, numWorkers int, batchSize int) {
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := make([]RawMessage, 0, batchSize)
			for {
				select {
				case <-ctx.Done():
					if len(batch) > 0 {
						_ = h.HandleBatch(ctx, batch)
					}
					return
				case m, ok := <-ch:
					if !ok {
						if len(batch) > 0 {
							_ = h.HandleBatch(ctx, batch)
						}
						return
					}
					batch = append(batch, m)
					for len(batch) < batchSize {
						select {
						case m2, ok2 := <-ch:
							if !ok2 {
								_ = h.HandleBatch(ctx, batch)
								return
							}
							batch = append(batch, m2)
						default:
							goto send
						}
					}
				send:
					if len(batch) > 0 {
						_ = h.HandleBatch(ctx, batch)
						batch = batch[:0]
					}
				}
			}
		}()
	}
	wg.Wait()
}
