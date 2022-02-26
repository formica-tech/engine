package domain_test

import (
	"fmt"
	"github.com/hamzali/formica-engine/domain"
	"testing"
	"time"
)

func TestOEE(t *testing.T) {
	t.Run("should calculate OEE values", func(st *testing.T) {
		result := domain.OEEResult{
			Start:                 time.Now(),
			End:                   time.Now().Add(time.Hour * 24),
			TotalProduction:       230,
			TotalWorkDuration:     (time.Hour * 20).Nanoseconds(),
			NotGoodProduction:     10,
			UnplannedStopDuration: (time.Minute * 120).Nanoseconds(),
			PlannedStopDuration:   (time.Minute * 120).Nanoseconds(),
			IdealCycle:            time.Minute * 6,
			Counts:                nil,
			Durations:             nil,
		}
		expectedOEE := 0.869565
		actualOEE := result.OEE()

		if fmt.Sprintf("%f", actualOEE) != fmt.Sprintf("%f", expectedOEE) {
			st.Errorf("expected OEE is %f but received %f", actualOEE, expectedOEE)
			return
		}
	})
}
