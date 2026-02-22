package quality

import (
	"biometrics-cli/internal/opencode"
	"context"
)

// EnforceQualityGate startet einen Folge-Agenten, der die Arbeit verifiziert.
func EnforceQualityGate(ctx context.Context, exec *opencode.Executor, req opencode.AgentRequest) error {
	verifyPrompt := `Sicher? Führe eine vollständige Selbstreflexion durch. 
Prüfe jede deiner Aussagen, verifiziere, ob ALLE Restriktionen des Initial-Prompts 
exakt eingehalten wurden. Stelle alles Fehlende fertig. 
Führe 'go vet ./...' und 'go fmt ./...' aus.`

	verifyReq := opencode.AgentRequest{
		ProjectID: req.ProjectID,
		Model:     "minimax-m2.5", // Verifikation immer mit schnellem Modell
		Prompt:    verifyPrompt,
		Category:  "quick",
	}

	result := exec.RunAgent(ctx, verifyReq)
	return result.Error
}
