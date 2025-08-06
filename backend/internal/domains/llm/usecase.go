package llm

import "context"

type UseCase struct {
	svc Service
}

func NewUseCase(svc Service) UseCase {
	return UseCase{
		svc: svc,
	}
}

func (u UseCase) AskLLM(ctx context.Context, prompt string) (string, error) {
	return u.svc.CallOllama(ctx, prompt)
}
