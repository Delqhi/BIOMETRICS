package chaos

import (
	"biometrics-cli/internal/metrics"
	"biometrics-cli/internal/state"
	"math/rand"
	"time"
)

func RunChaosMonkey() {
	for {
		if state.GlobalState.ChaosEnabled {
			time.Sleep(time.Duration(rand.Intn(120)+60) * time.Second)
			event := rand.Intn(3)
			state.GlobalState.SetChaos(true)

			switch event {
			case 0:
				state.GlobalState.Log("CHAOS", "Simulating Serena process failure")
				metrics.ChaosEventsTotal.WithLabelValues("serena_fail").Inc()
			case 1:
				state.GlobalState.Log("CHAOS", "Simulating Boulder file corruption")
				metrics.ChaosEventsTotal.WithLabelValues("boulder_fail").Inc()
			case 2:
				state.GlobalState.Log("CHAOS", "Injecting artificial model latency")
				metrics.ChaosEventsTotal.WithLabelValues("latency").Inc()
			}

			time.Sleep(30 * time.Second)
			state.GlobalState.SetChaos(false)
			state.GlobalState.Log("CHAOS", "Monkey returned to cage. Systems restoring...")
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}
