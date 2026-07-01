"use client";

import { FormEvent } from "react";
import { Button } from "@/components/ui/Button";

type PracticeSubmitBarProps = {
  isConcept: boolean;
  loading: boolean;
  disabled?: boolean;
  onSubmit: (e: FormEvent) => void;
};

export function PracticeSubmitBar({
  isConcept,
  loading,
  disabled,
  onSubmit,
}: PracticeSubmitBarProps) {
  return (
    <div className="sticky bottom-0 z-10 -mx-4 mt-4 border-t border-zinc-200 bg-white/95 px-4 py-4 backdrop-blur supports-[backdrop-filter]:bg-white/80 dark:border-zinc-800 dark:bg-zinc-900/95 dark:supports-[backdrop-filter]:bg-zinc-900/80 lg:-mx-8 lg:px-8">
      <form
        onSubmit={onSubmit}
        className="mx-auto flex max-w-6xl flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
      >
        <div>
          <p className="text-sm font-medium text-zinc-800 dark:text-zinc-200">
            Ready for AI evaluation?
          </p>
          <p className="mt-0.5 text-sm text-zinc-500 dark:text-zinc-400">
            {isConcept
              ? "Submit when your written answer is complete."
              : "Run and debug first — submit only when you want OpenAI scoring."}
          </p>
        </div>
        <Button
          type="submit"
          className="shrink-0 self-start sm:self-center"
          loading={loading}
          disabled={disabled}
        >
          Submit & evaluate
        </Button>
      </form>
    </div>
  );
}
