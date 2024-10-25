package tiles

import (
	"context"
	"log/slog"
	"os/exec"
)

type customOutput struct {
	logger *slog.Logger
}

func (cO *customOutput) Write(p []byte) (int, error) {
	cO.logger.Info("received output: ", "cmdOutput", string(p))
	return len(p), nil
}

type executor struct {
	logger *slog.Logger
	debug  bool
}

func (e *executor) execute(ctx context.Context, command string, params []string) error {
	cmd := exec.CommandContext(ctx, command, params...)
	if e.debug {
		cmd.Stderr = &customOutput{logger: e.logger}
		cmd.Stdout = &customOutput{logger: e.logger}
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (e *executor) executeWithOutput(
	ctx context.Context,
	command string,
	params []string,
) ([]byte, error) {
	cmd := exec.CommandContext(ctx, command, params...)

	output, err := cmd.Output()
	if err != nil {
		return []byte{}, err
	}

	return output, nil
}

func (e *executor) hasExecutable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
