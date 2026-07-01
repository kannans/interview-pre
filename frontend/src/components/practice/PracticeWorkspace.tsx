"use client";

import { useEffect, useState } from "react";
import clsx from "clsx";
import { CodeEditor } from "@/components/practice/CodeEditor";
import { LanguageSelect } from "@/components/practice/LanguageSelect";
import { RunOutput } from "@/components/practice/RunOutput";
import { Button } from "@/components/ui/Button";
import { useVerticalResize } from "@/hooks/useVerticalResize";
import type { PracticeLanguage } from "@/lib/practice/languages";
import type { RunCodeResult } from "@/lib/types/api";

type PracticeWorkspaceProps = {
  language: PracticeLanguage;
  onLanguageChange: (language: PracticeLanguage) => void;
  editorLanguage: string;
  answer: string;
  onAnswerChange: (value: string) => void;
  stdin: string;
  onStdinChange: (value: string) => void;
  runResult: RunCodeResult | null;
  runPending: boolean;
  onRun: () => void;
  readOnly?: boolean;
  draftReady?: boolean;
};

export function PracticeWorkspace({
  language,
  onLanguageChange,
  editorLanguage,
  answer,
  onAnswerChange,
  stdin,
  onStdinChange,
  runResult,
  runPending,
  onRun,
  readOnly = false,
  draftReady = false,
}: PracticeWorkspaceProps) {
  const [stdinOpen, setStdinOpen] = useState(stdin.length > 0);
  const [terminalOpen, setTerminalOpen] = useState(true);
  const {
    height: editorHeight,
    onPointerDown,
    onPointerMove,
    onPointerUp,
  } = useVerticalResize({ disabled: readOnly });

  useEffect(() => {
    if (stdin.length > 0) setStdinOpen(true);
  }, [stdin]);

  useEffect(() => {
    if (runPending || runResult) setTerminalOpen(true);
  }, [runPending, runResult]);

  return (
    <div className="flex min-h-[36rem] flex-col overflow-hidden rounded-xl border border-zinc-200 bg-white shadow-sm dark:border-zinc-800 dark:bg-zinc-900">
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-zinc-200 px-4 py-3 dark:border-zinc-800">
        <div className="flex flex-wrap items-center gap-3">
          <LanguageSelect
            value={language}
            onChange={onLanguageChange}
            disabled={readOnly}
          />
          {draftReady && !readOnly && (
            <span className="text-xs text-zinc-500">Draft saved locally</span>
          )}
        </div>
        <Button
          type="button"
          variant="primary"
          size="sm"
          onClick={onRun}
          loading={runPending}
          disabled={readOnly || runPending}
        >
          {!runPending && <PlayIcon />}
          {runPending ? "Running…" : "Run"}
        </Button>
      </div>

      <div className="shrink-0 p-4 pb-0">
        <CodeEditor
          value={answer}
          onChange={onAnswerChange}
          language={language}
          editorLanguage={editorLanguage}
          readOnly={readOnly}
          height={`${editorHeight}px`}
        />
      </div>

      <div
        role="separator"
        aria-orientation="horizontal"
        aria-valuenow={editorHeight}
        tabIndex={0}
        onPointerDown={onPointerDown}
        onPointerMove={onPointerMove}
        onPointerUp={onPointerUp}
        onPointerCancel={onPointerUp}
        className={clsx(
          "group flex h-2 shrink-0 cursor-row-resize items-center justify-center",
          "border-y border-zinc-200 bg-zinc-100 hover:bg-zinc-200",
          "dark:border-zinc-800 dark:bg-zinc-900 dark:hover:bg-zinc-800",
          readOnly && "pointer-events-none opacity-50",
        )}
      >
        <span className="h-1 w-10 rounded-full bg-zinc-300 group-hover:bg-zinc-400 dark:bg-zinc-600 dark:group-hover:bg-zinc-500" />
      </div>

      <div className="border-t border-zinc-200 px-4 py-2 dark:border-zinc-800">
        <CollapsibleHeader
          open={stdinOpen}
          onToggle={() => setStdinOpen((open) => !open)}
          disabled={readOnly}
          label="Stdin (optional)"
        />
        {stdinOpen && (
          <textarea
            value={stdin}
            onChange={(e) => onStdinChange(e.target.value)}
            disabled={readOnly}
            rows={2}
            placeholder="Input passed to your program on stdin"
            className="mt-2 w-full rounded-lg border border-zinc-300 bg-white px-3 py-2 font-mono text-sm text-zinc-900 disabled:opacity-60 dark:border-zinc-700 dark:bg-zinc-950 dark:text-zinc-100"
          />
        )}
      </div>

      <div
        className={clsx(
          "flex flex-1 flex-col border-t border-zinc-200 dark:border-zinc-800",
          "bg-zinc-50 dark:bg-zinc-950/80",
          terminalOpen && "min-h-[8rem]",
        )}
      >
        <div className="px-4 py-2">
          <CollapsibleHeader
            open={terminalOpen}
            onToggle={() => setTerminalOpen((open) => !open)}
            disabled={readOnly}
            label="Terminal"
            uppercase
          />
        </div>
        {terminalOpen && (
          <div className="flex-1 px-4 pb-4">
            <RunOutput result={runResult} loading={runPending} compact />
          </div>
        )}
      </div>
    </div>
  );
}

function CollapsibleHeader({
  open,
  onToggle,
  disabled,
  label,
  uppercase = false,
}: {
  open: boolean;
  onToggle: () => void;
  disabled?: boolean;
  label: string;
  uppercase?: boolean;
}) {
  return (
    <button
      type="button"
      onClick={onToggle}
      disabled={disabled}
      aria-expanded={open}
      className="flex w-full items-center gap-2 py-1 text-left text-xs font-medium text-zinc-500 hover:text-zinc-700 disabled:cursor-not-allowed dark:hover:text-zinc-300"
    >
      <span
        aria-hidden
        className="inline-flex h-4 w-4 shrink-0 items-center justify-center rounded text-sm leading-none"
      >
        {open ? "−" : "+"}
      </span>
      <span className={clsx(uppercase && "uppercase tracking-wide")}>{label}</span>
    </button>
  );
}

function PlayIcon() {
  return (
    <svg
      viewBox="0 0 16 16"
      fill="currentColor"
      className="h-3.5 w-3.5"
      aria-hidden
    >
      <path d="M4 2.5v11l9-5.5-9-5.5z" />
    </svg>
  );
}
