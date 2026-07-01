package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	maxCodeBytes   = 64 * 1024
	maxStdinBytes  = 8 * 1024
	maxOutputBytes = 32 * 1024
	runTimeout     = 5 * time.Second
)

var mainFuncPattern = regexp.MustCompile(`func\s+main\s*\(`)

// RunInput is the payload for executing practice code locally.
type RunInput struct {
	Language string
	Code     string
	Stdin    string
}

// RunResult is returned after executing user code.
type RunResult struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exit_code"`
	DurationMs int64  `json:"duration_ms"`
	TimedOut   bool   `json:"timed_out"`
}

// Executor runs user code in a local interpreter.
type Executor interface {
	Run(ctx context.Context, input RunInput) (RunResult, error)
}

type localExecutor struct{}

// NewLocalExecutor returns an executor that runs code via local CLI tools.
func NewLocalExecutor() Executor {
	return &localExecutor{}
}

func (e *localExecutor) Run(ctx context.Context, input RunInput) (RunResult, error) {
	lang := strings.ToLower(strings.TrimSpace(input.Language))
	if len(input.Code) > maxCodeBytes {
		return RunResult{}, fmt.Errorf("code exceeds maximum size of %d bytes", maxCodeBytes)
	}
	if len(input.Stdin) > maxStdinBytes {
		return RunResult{}, fmt.Errorf("stdin exceeds maximum size of %d bytes", maxStdinBytes)
	}

	dir, err := os.MkdirTemp("", "practice-run-*")
	if err != nil {
		return RunResult{}, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	spec, err := buildLanguageSpec(lang, input.Code)
	if err != nil {
		return RunResult{}, err
	}

	sourcePath := filepath.Join(dir, spec.filename)
	if err := os.WriteFile(sourcePath, []byte(spec.source), 0o600); err != nil {
		return RunResult{}, fmt.Errorf("write source file: %w", err)
	}

	runCtx, cancel := context.WithTimeout(ctx, runTimeout)
	defer cancel()

	cmd := exec.CommandContext(runCtx, spec.command, spec.args...)
	cmd.Dir = dir
	if spec.workDir != "" {
		cmd.Dir = spec.workDir
	}
	if input.Stdin != "" {
		cmd.Stdin = strings.NewReader(input.Stdin)
	}

	var stdoutBuf, stderrBuf limitedBuffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	start := time.Now()
	runErr := cmd.Run()
	durationMs := time.Since(start).Milliseconds()

	result := RunResult{
		Stdout:     stdoutBuf.String(),
		Stderr:     stderrBuf.String(),
		DurationMs: durationMs,
	}

	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		result.TimedOut = true
		result.ExitCode = -1
		if result.Stderr == "" {
			result.Stderr = fmt.Sprintf("execution timed out after %s", runTimeout)
		}
		return result, nil
	}

	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
			if result.Stderr == "" {
				result.Stderr = runErr.Error()
			}
		}
		return result, nil
	}

	result.ExitCode = 0
	return result, nil
}

type runSpec struct {
	filename string
	source   string
	command  string
	args     []string
	workDir  string
}

func buildLanguageSpec(lang, code string) (runSpec, error) {
	switch lang {
	case "javascript":
		return runSpec{
			filename: "main.js",
			source:   code,
			command:  "node",
			args:     []string{"main.js"},
		}, nil
	case "typescript":
		cmd, args, err := typescriptCommand()
		if err != nil {
			return runSpec{}, err
		}
		return runSpec{
			filename: "main.ts",
			source:   code,
			command:  cmd,
			args:     args,
		}, nil
	case "python":
		cmd, err := pythonCommand()
		if err != nil {
			return runSpec{}, err
		}
		return runSpec{
			filename: "main.py",
			source:   code,
			command:  cmd,
			args:     []string{"main.py"},
		}, nil
	case "go":
		return runSpec{
			filename: "main.go",
			source:   prepareGoSource(code),
			command:  "go",
			args:     []string{"run", "main.go"},
		}, nil
	case "ruby":
		return runSpec{
			filename: "main.rb",
			source:   code,
			command:  "ruby",
			args:     []string{"main.rb"},
		}, nil
	default:
		return runSpec{}, fmt.Errorf("unsupported language: %s", lang)
	}
}

func typescriptCommand() (string, []string, error) {
	if path, err := exec.LookPath("tsx"); err == nil {
		return path, []string{"main.ts"}, nil
	}
	if path, err := exec.LookPath("npx"); err == nil {
		return path, []string{"--yes", "tsx", "main.ts"}, nil
	}
	return "", nil, errors.New("typescript runner not found (install tsx or ensure npx is available)")
}

func pythonCommand() (string, error) {
	for _, name := range []string{"python3", "python"} {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}
	return "", errors.New("python interpreter not found (install python3)")
}

func prepareGoSource(code string) string {
	if strings.Contains(code, "package main") {
		return code
	}

	hasMain := mainFuncPattern.MatchString(code)

	var b strings.Builder
	b.WriteString("package main\n\n")
	if !hasMain {
		b.WriteString("import \"fmt\"\n\n")
	}
	b.WriteString(code)
	if !hasMain {
		b.WriteString("\n\nfunc main() {\n\t// Example: fmt.Println(yourFunction(args))\n}\n")
	}
	return b.String()
}

type limitedBuffer struct {
	buf bytes.Buffer
}

func (l *limitedBuffer) Write(p []byte) (int, error) {
	remaining := maxOutputBytes - l.buf.Len()
	if remaining <= 0 {
		return len(p), nil
	}
	if len(p) > remaining {
		_, _ = l.buf.Write(p[:remaining])
		_, _ = l.buf.WriteString("\n… output truncated …\n")
		return len(p), nil
	}
	return l.buf.Write(p)
}

func (l *limitedBuffer) String() string {
	return l.buf.String()
}
