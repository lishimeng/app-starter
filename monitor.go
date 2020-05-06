package app

import (
	"sync"
	"time"
)

type M struct {
	appStartTime time.Time
	requestTimes int

	minCost   time.Duration
	maxCost   time.Duration
	totalCost time.Duration

	lock *sync.Mutex
}

type Monitor struct {
	AppStartTime  string `json:"startTime"`
	RequestAmount int    `json:"requests"`
	MaxCost       string `json:"maxCost,omitempty"`
	MinCost       string `json:"minCost,omitempty"`
	CostAverage   string `json:"costAverage,omitempty"`
}

var m *M

func init() {
	m = &M{
		appStartTime: time.Now(),
		lock:         new(sync.Mutex),
	}
}

func RecordRequest(costTime time.Duration) {
	go m.record(costTime)
}

func (m *M) record(costTime time.Duration) {
	m.lock.Lock()

	if m.requestTimes == 0 {
		m.minCost = costTime
		m.maxCost = costTime
	} else {
		if m.maxCost < costTime {
			m.maxCost = costTime
		}
		if m.minCost > costTime {
			m.minCost = costTime
		}
	}

	m.totalCost += costTime
	m.requestTimes++

	m.lock.Unlock()
}

func GetStatus() Monitor {

	state := Monitor{
		AppStartTime:  m.appStartTime.Format(time.ANSIC),
		RequestAmount: m.requestTimes,
	}
	if m.requestTimes > 0 {
		ave := time.Duration(uint64(m.totalCost) / uint64(m.requestTimes))
		state.CostAverage = ave.String()
		state.MaxCost = m.maxCost.String()
		state.MinCost = m.minCost.String()
	}

	return state
}

