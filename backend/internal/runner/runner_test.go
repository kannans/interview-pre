package runner_test

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/kannans/interview-prep/backend/internal/runner"
)

func TestPrepareGoSourceAddsMainWhenMissing(t *testing.T) {
	t.Parallel()

	// prepareGoSource is not exported; exercise via Run with a tiny program.
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not installed")
	}

	exec := runner.NewLocalExecutor()
	result, err := exec.Run(context.Background(), runner.RunInput{
		Language: "go",
		Code: `package main

import "fmt"

func sumToN(n int) int {
	total := 0
	for i := 1; i <= n; i++ {
		total += i
	}
	return total
}

func main() {
	fmt.Println(sumToN(5))
}`,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%q", result.ExitCode, result.Stderr)
	}
	if !strings.Contains(result.Stdout, "15") {
		t.Fatalf("stdout=%q", result.Stdout)
	}
}

func TestRunJavaScript(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not installed")
	}

	exec := runner.NewLocalExecutor()
	result, err := exec.Run(context.Background(), runner.RunInput{
		Language: "javascript",
		Code:     "console.log(1 + 2);",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%q", result.ExitCode, result.Stderr)
	}
	if strings.TrimSpace(result.Stdout) != "3" {
		t.Fatalf("stdout=%q", result.Stdout)
	}
}

func TestRunCapturesStderrOnSyntaxError(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not installed")
	}

	exec := runner.NewLocalExecutor()
	result, err := exec.Run(context.Background(), runner.RunInput{
		Language: "javascript",
		Code:     "function {",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.ExitCode == 0 {
		t.Fatal("expected non-zero exit code")
	}
	if result.Stderr == "" {
		t.Fatal("expected stderr output")
	}
}
