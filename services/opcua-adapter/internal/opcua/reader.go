package opcua

import (
	"context"
	"log"
	"time"

	"sentinelflow/opcua-adapter/internal/config"

	"github.com/gopcua/opcua/ua"
)

// Reader polls OPC UA nodes and sends values to the provided channel.
type Reader struct {
	endpoint string
	mappings []config.NodeMapping
	interval time.Duration
}

// NewReader creates a reader for the given config.
func NewReader(cfg *config.Config) *Reader {
	return &Reader{
		endpoint: cfg.Endpoint,
		mappings: cfg.NodeMappings,
		interval: time.Duration(cfg.SamplingInterval) * time.Millisecond,
	}
}

// Value holds a single node value for forwarding.
type Value struct {
	SensorID   string
	AssetID    string
	MetricName string
	Value      float64
	Critical   bool
}

// Run connects to OPC UA, polls each node at the configured interval, and sends Value to out until ctx is done.
func (r *Reader) Run(ctx context.Context, out chan<- Value) error {
	if len(r.mappings) == 0 {
		log.Print("opcua: no node mappings configured (set OPCUA_NODES); exiting")
		return nil
	}
	client, err := opcua.NewClient(r.endpoint)
	if err != nil {
		return err
	}
	if err := client.Connect(ctx); err != nil {
		return err
	}
	defer client.Close(ctx)

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			for _, m := range r.mappings {
				nodeID, err := ua.ParseNodeID(m.NodeID)
				if err != nil {
					log.Printf("opcua: invalid node id %q: %v", m.NodeID, err)
					continue
				}
				v, err := client.Node(nodeID).Value(ctx)
				if err != nil {
					log.Printf("opcua: read %q: %v", m.NodeID, err)
					continue
				}
				if v == nil {
					continue
				}
				f := variantToFloat64(v)
				out <- Value{
					SensorID:   m.SensorID,
					AssetID:    m.AssetID,
					MetricName: m.MetricName,
					Value:      f,
					Critical:   false,
				}
			}
		}
	}
}

func variantToFloat64(v *ua.Variant) float64 {
	if v == nil {
		return 0
	}
	switch v.Type() {
	case ua.TypeIDDouble:
		return v.Float()
	case ua.TypeIDFloat:
		return float64(v.Float())
	case ua.TypeIDInt32:
		return float64(v.Int())
	case ua.TypeIDInt64:
		return float64(v.Int())
	case ua.TypeIDUint32:
		return float64(v.Uint())
	case ua.TypeIDUint64:
		return float64(v.Uint())
	default:
		return 0
	}
}
efault:
		return 0
	}
}
.Uint())
	default:
		return 0
	}
}
efault:
		return 0
	}
}
