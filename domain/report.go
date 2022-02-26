package domain

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type EntityStateTransition struct {
	open  EntitySignal
	close *EntitySignal
}

type EntityStateInfo struct {
	State    string
	Duration time.Duration
	Open     bool
}

func (s *EntityStateTransition) Info() EntityStateInfo {
	return EntityStateInfo{
		State:    s.State(),
		Duration: s.Duration(),
		Open:     s.Open(),
	}
}

func (s *EntityStateTransition) State() string {
	return s.open.Event
}

func (s *EntityStateTransition) Duration() time.Duration {
	if s.close == nil {
		return time.Now().Sub(s.open.Timestamp)
	}
	return s.close.Timestamp.Sub(s.open.Timestamp)
}

func (s *EntityStateTransition) Open() bool {
	return s.close == nil
}

func (s *EntityStateTransition) OpenSignal() EntitySignal {
	return s.open
}

func (s *EntityStateTransition) CloseSignal() *EntitySignal {
	if s.close == nil {
		return nil
	}
	temp := &EntitySignal{}
	*temp = *s.close
	return temp
}

func NewStateTransition(open EntitySignal, close *EntitySignal) (EntityStateTransition, error) {
	result := EntityStateTransition{}
	if close != nil && open.Timestamp.After(close.Timestamp) {
		return result, errors.New("invalid signal, open signal is after close signal")
	}

	return EntityStateTransition{
		open:  open,
		close: close,
	}, nil
}

func UpdateEntityStateTransitions(stateTransitions []EntityStateTransition, signal EntitySignal) ([]EntityStateTransition, error) {
	newTransitions := make([]EntityStateTransition, len(stateTransitions))
	copy(newTransitions, stateTransitions)

	if len(newTransitions) == 0 {
		transition, err := NewStateTransition(signal, nil)
		if err != nil {
			return newTransitions, err
		}
		return append(newTransitions, transition), nil
	}

	lastTransition := newTransitions[len(newTransitions)-1]

	if lastTransition.close != nil {
		newTransition, err := NewStateTransition(*lastTransition.close, &signal)
		if err != nil {
			return newTransitions, err
		}
		newTransitions[len(newTransitions)-1] = newTransition
	} else {
		newTransition, err := NewStateTransition(lastTransition.OpenSignal(), &signal)
		if err != nil {
			return newTransitions, err
		}
		newTransitions[len(newTransitions)-1] = newTransition
	}

	newTransition, err := NewStateTransition(signal, nil)
	if err != nil {
		return newTransitions, err
	}
	newTransitions = append(newTransitions, newTransition)
	return newTransitions, nil
}

type ProcessDataPoint struct {
	SignalID  string
	EntityID  string
	Label     string
	Value     float64
	Timestamp time.Time
}

func SignalToProcessData(s EntitySignal) (p ProcessDataPoint, err error) {
	value, err := strconv.ParseFloat(fmt.Sprintf("%v", s.Payload["value"]), 64)
	if err != nil {
		return
	}
	p.Value = value
	p.Label = fmt.Sprintf("%v", s.Payload["label"])
	p.SignalID = s.ID
	p.EntityID = s.EntityID
	p.Timestamp = s.Timestamp
	return
}

type OEEResult struct {
	Start                 time.Time        `json:"start"`
	End                   time.Time        `json:"end"`
	TotalProduction       float64          `json:"total_production"`
	TotalWorkDuration     int64            `json:"total_work_duration"`
	NotGoodProduction     float64          `json:"not_good_production"`
	UnplannedStopDuration int64            `json:"unplanned_stop_duration"`
	PlannedStopDuration   int64            `json:"planned_stop_duration"`
	IdealCycle            time.Duration    `json:"ideal_cycle"`
	Counts                map[string]int64 `json:"counts"`
	Durations             map[string]int64 `json:"durations"`
	LastSignal            *EntitySignal    `json:"last_signal"`
}

type OEECalculations struct {
	PlannedDuration   int64   `json:"planned_duration"`
	PlannedProduction int64   `json:"planned_production"`
	TotalDuration     int64   `json:"total_duration"`
	GoodProduction    float64 `json:"good_production"`
	Availability      float64 `json:"availability"`
	Performance       float64 `json:"performance"`
	Quality           float64 `json:"quality"`
	OEE               float64 `json:"oee"`
}

func (r *OEEResult) Calculations() OEECalculations {
	return OEECalculations{
		PlannedDuration:   r.PlannedDuration(),
		PlannedProduction: r.PlannedProduction(),
		TotalDuration:     r.TotalDuration(),
		GoodProduction:    r.GoodProduction(),
		Availability:      r.Availability(),
		Performance:       r.Performance(),
		Quality:           r.Quality(),
		OEE:               r.OEE(),
	}
}

func (r *OEEResult) PlannedDuration() int64 {
	return r.TotalDuration() - r.PlannedStopDuration
}

func (r *OEEResult) PlannedProduction() int64 {
	if r.IdealCycle.Nanoseconds() == 0 {
		return 0
	}
	return r.PlannedDuration() / r.IdealCycle.Nanoseconds()
}

func (r *OEEResult) TotalDuration() int64 {
	return r.End.UnixNano() - r.Start.UnixNano()
}

func (r *OEEResult) GoodProduction() float64 {
	return r.TotalProduction - r.NotGoodProduction
}

func (r *OEEResult) Availability() float64 {
	plannedDuration := r.PlannedDuration()
	if plannedDuration == 0 {
		return 0
	}
	return float64(plannedDuration-r.UnplannedStopDuration) / float64(plannedDuration)
}

func (r *OEEResult) Performance() float64 {
	plannedProduction := r.PlannedProduction()
	if plannedProduction == 0 {
		return 0
	}
	return r.GoodProduction() / float64(plannedProduction)
}

func (r *OEEResult) Quality() float64 {
	if r.TotalProduction == 0 {
		return 0
	}
	return r.GoodProduction() / r.TotalProduction
}

func (r *OEEResult) OEE() float64 {
	return r.Availability() * r.Performance() * r.Quality()
}

type SignalFilter struct {
	Max time.Duration
	Min time.Duration
}

type OEEInput struct {
	Machines          []string
	PlannedEvents     []string
	UnplannedEvents   []string
	CountableEvents   []string
	IdealCycle        int64
	SignalFilter      SignalFilter
	plannedEventMap   map[string]bool
	unplannedEventMap map[string]bool
	countableEventMap map[string]bool
}

func (r *OEEInput) IsPlannedEvent(event string) bool {
	if r.plannedEventMap != nil {
		return r.plannedEventMap[event]
	}
	r.plannedEventMap = map[string]bool{}
	for _, e := range r.PlannedEvents {
		r.plannedEventMap[e] = true
	}
	return r.IsPlannedEvent(event)
}

func (r *OEEInput) IsUnplannedEvent(event string) bool {
	if r.unplannedEventMap != nil {
		return r.unplannedEventMap[event]
	}
	r.unplannedEventMap = map[string]bool{}
	for _, e := range r.UnplannedEvents {
		r.unplannedEventMap[e] = true
	}
	return r.IsUnplannedEvent(event)
}

func (r *OEEInput) IsCountableEvent(event string) bool {
	if r.countableEventMap != nil {
		return r.countableEventMap[event]
	}

	r.countableEventMap = map[string]bool{}
	for _, e := range r.CountableEvents {
		r.countableEventMap[e] = true
	}
	return r.IsCountableEvent(event)
}

const (
	ProductionEvent  string = "PRODUCTION"
	FailEvent               = "FAIL"
	BreakEvent              = "BREAK"
	NotGoodEvent            = "NOT_GOOD"
	ProcessDataEvent        = "PROCESS_DATA"
)

var ErrInvalidSignalOrder = errors.New("last signal is after current signal, check signal order")

func UpdateOEECalculation(result OEEResult, params OEEInput, s EntitySignal) (OEEResult, error) {
	result.IdealCycle = time.Duration(params.IdealCycle)
	if s.Event == NotGoodEvent {
		result.NotGoodProduction += 1
		result.Counts[s.Event] += 1
		return result, nil
	}

	if params.IsCountableEvent(s.Event) {
		result.Counts[s.Event] += 1
	}

	if result.LastSignal == nil {
		result.LastSignal = &EntitySignal{}
		*result.LastSignal = s
		return result, nil
	}

	d := s.Timestamp.Sub(result.LastSignal.Timestamp)

	if d < 0 {
		return result, ErrInvalidSignalOrder
	}

	if d < params.SignalFilter.Min || (params.SignalFilter.Max != time.Duration(0) && d > params.SignalFilter.Max) {
		return result, nil
	}

	if result.LastSignal.Event == ProductionEvent {
		result.TotalProduction += 1
		result.TotalWorkDuration += d.Nanoseconds()
	}

	if result.LastSignal.Event == FailEvent || params.IsUnplannedEvent(result.LastSignal.Event) {
		result.UnplannedStopDuration += d.Nanoseconds()
	}

	if result.LastSignal.Event == BreakEvent || params.IsPlannedEvent(result.LastSignal.Event) {
		result.PlannedStopDuration += d.Nanoseconds()
	}

	result.Durations[result.LastSignal.Event] += d.Nanoseconds()

	*result.LastSignal = s
	return result, nil
}
