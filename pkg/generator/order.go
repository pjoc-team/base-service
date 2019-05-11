package generator

import (
	"fmt"
	"github.com/pjoc-team/base-service/pkg/util"
	"hash/fnv"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

type OrderGenerator struct {
	ClusterId   *string
	MachineId   *string
	Concurrency int
	maxIndex    uint32
	indexWidth  int
	index       uint32
}

// yyyyMMddHHmmSS
const TIME_LAYOUT = "20060102150405"
const ZERO_BYTE = byte('0')

func New(clusterId *string, concurrency int) *OrderGenerator {
	g := &OrderGenerator{}
	id := fmt.Sprint(getIdentityId())
	g.MachineId = &id
	g.ClusterId = clusterId
	g.Concurrency = concurrency
	g.maxIndex = uint32(concurrency)
	for g.indexWidth = 0; concurrency > 0; g.indexWidth++ {
		concurrency = concurrency / 10
	}
	return g
}

func getIdentityId() uint32 {
	if name, err := os.Hostname(); err == nil {
		h := fnv.New32()
		h.Write([]byte(name))
		sum32 := h.Sum32()
		return sum32
	} else {
		ip := util.GetHostIP()
		h := fnv.New32()
		h.Write([]byte(ip))
		sum32 := h.Sum32()
		return sum32
	}
}

func dateStr() string {
	date := time.Now().Format(TIME_LAYOUT)
	nanosecond := time.Now().Nanosecond() / 1000
	return fmt.Sprintf("%s%d", date, nanosecond)
}

func (g *OrderGenerator) GenerateIndex() string {
	index := atomic.AddUint32(&g.index, 1) % g.maxIndex
	s := fmt.Sprint(index)
	leastSize := g.indexWidth - len(s)

	builder := strings.Builder{}
	if leastSize > 0 {
		bytes := make([]byte, leastSize)
		for i := 0; i < leastSize; i++ {
			bytes[i] = ZERO_BYTE
		}
		builder.WriteString(string(bytes))
	}
	builder.WriteString(s)
	return builder.String()
}

func (g *OrderGenerator) GenerateOrderId() string {
	builder := strings.Builder{}
	dateStr := dateStr()
	builder.WriteString(dateStr)
	builder.WriteString(*g.ClusterId)
	builder.WriteString(*g.MachineId)
	builder.WriteString(g.GenerateIndex())
	return builder.String()
}
