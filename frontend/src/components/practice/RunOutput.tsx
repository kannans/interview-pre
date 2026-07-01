"use client";

import clsx from "clsx";
import type { RunCodeResult } from "@/lib/types/api";

type RunOutputProps = {
  result: RunCodeResult | null;
  loading?: boolean;
  className?: string;
};

export function RunOutput({ result, loading, className }: RunOutputProps) {
  if (loading) {
    return (
      <div
        className={clsx(
          "rounded-lg border border-zinc-300 bg-zinc-950 p-3 font-mono text-sm text-zinc-300 dark:border-zinc-700",
          className,
        )}
      >
        Running…
      </div>
    );
  }

  if (!result) {
    return (
      <div
        className={clsx(
          "rounded-lg border border-dashed border-zinc-300 p-3 text-sm text-zinc-500 dark:border-zinc-700",
          className,
        )}
      >
        Add print/console.log calls to test your code, then click Run.
      </div>
    );
  }

  const hasStdout = result.stdout.trim().length > 0;
  const hasStderr = result.stderr.trim().length > 0;

  return (
    <div className={clsx("space-y-2", className)}>
      <div className="flex flex-wrap items-center gap-2 text-xs text-zinc-500">
        <span>Exit code: {result.exit_code}</span>
        <span>·</span>
        <span>{result.duration_ms} ms</span>
        {result.timed_out && (
          <>
            <span>·</span>
            <span className="text-amber-600 dark:text-amber-400">Timed out</span>
          </>
        )}
      </div>

      {hasStdout && (
        <OutputBlock label="stdout" tone="stdout">
          {result.stdout}
        </OutputBlock>
      )}

      {hasStderr && (
        <OutputBlock label="stderr" tone="stderr">
          {result.stderr}
        </OutputBlock>
      )}

      {!hasStdout && !hasStderr && (
        <OutputBlock label="output" tone="stdout">
          {result.exit_code === 0
            ? "Program finished with no output."
            : "Program finished with no output. Check stderr above if shown."}
        </OutputBlock>
      )}
    </div>
  );
}

function OutputBlock({
  label,
  tone,
  children,
}: {
  label: string;
  tone: "stdout" | "stderr";
  children: string;
}) {
  return (
    <div>
      <p className="mb-1 text-xs font-medium uppercase tracking-wide text-zinc-500">
        {label}
      </p>
      <pre
        className={clsx(
          "max-h-48 overflow-auto rounded-lg border p-3 font-mono text-sm whitespace-pre-wrap",
          tone === "stderr"
            ? "border-red-300/40 bg-red-950/40 text-red-200"
            : "border-zinc-700 bg-zinc-950 text-emerald-100",
        )}
      >
        {children}
      </pre>
    </div>
  );
}
