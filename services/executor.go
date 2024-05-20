package services

import (
	"context"
	"log/slog"
	"os/exec"
)

type customOutput struct{ logger *slog.Logger }

func (cO *customOutput) Write(p []byte) (int, error) {
	cO.logger.Info("received output: ", "cmdOutput", string(p))
	return len(p), nil
}

type executor struct {
	logger *slog.Logger
	debug  bool
}

func (e *executor) execute(ctx context.Context, command string, args []string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	if e.debug {
		cmd.Stdout = &customOutput{logger: e.logger}
	}
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (e *executor) hasExecutable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
