package uid

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake/v2"
)

type Generator struct {
	sf *sonyflake.Sonyflake
}

func NewGenerator() (*Generator, error) {
	sf, err := sonyflake.New(sonyflake.Settings{
		StartTime: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sonyflake: %w", err)
	}
	return &Generator{sf: sf}, nil
}

func (g *Generator) Generate() (uint64, error) {
	id, err := g.sf.NextID()
	if err != nil {
		return 0, fmt.Errorf("failed to generate id: %w", err)
	}
	return uint64(id), nil
}
