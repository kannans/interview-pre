"use client";

import Link from "next/link";
import { TopBar } from "@/components/layout/TopBar";
import { Card } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { ErrorAlert } from "@/components/ui/ErrorAlert";
import { EmptyState } from "@/components/ui/EmptyState";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { useRecommendedPractice } from "@/hooks/useRecommendedPractice";
import { formatCategory, formatLevel } from "@/lib/utils/format";
import { isApiError } from "@/lib/api";

export default function PracticePage() {
  const { today, questions, task, recommended, practiceHref, isLoading } =
    useRecommendedPractice();

  if (isLoading) {
    return <LoadingSpinner label="Loading practice…" />;
  }

  const { data, isError, error } = questions;
  const todayError =
    today.isError && isApiError(today.error) && today.error.isNotFound
      ? undefined
      : today.isError
        ? today.error.message
        : undefined;

  const sortedQuestions = data?.questions
    ? [...data.questions].sort((a, b) => {
        if (recommended?.id === a.id) return -1;
        if (recommended?.id === b.id) return 1;
        return 0;
      })
    : [];

  return (
    <>
      <TopBar
        title="Daily practice"
        subtitle="Work through programming and DSA questions aligned with your recovery plan."
      />

      {todayError && (
        <ErrorAlert className="mb-4" message={todayError} />
      )}

      {task && (
        <Card title="Today's recovery task" className="mb-4">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <div className="flex flex-wrap gap-2">
                <Badge tone="info">{formatCategory(task.category)}</Badge>
                <Badge>{formatLevel(task.level)}</Badge>
              </div>
              <p className="mt-2 font-medium text-zinc-900 dark:text-zinc-50">
                {task.title}
              </p>
              {recommended && (
                <p className="mt-2 text-sm text-zinc-500">
                  Recommended: {recommended.title}
                </p>
              )}
            </div>
            <Button size="sm" href={practiceHref}>
              Start recommended
            </Button>
          </div>
        </Card>
      )}

      {isError && (
        <ErrorAlert message={error?.message ?? "Failed to load questions"} />
      )}

      {!isError && (!data?.questions || data.questions.length === 0) && (
        <EmptyState
          title="No questions available"
          description="Ensure the backend is running and the question bank is seeded."
        />
      )}

      <ul className="grid gap-3 sm:grid-cols-2">
        {sortedQuestions.map((q) => {
          const isRecommended = recommended?.id === q.id;
          return (
            <li key={q.id}>
              <Link href={`/practice/${q.id}`}>
                <Card
                  className={
                    isRecommended
                      ? "border-teal-300 ring-1 ring-teal-200 transition-shadow hover:shadow-md dark:border-teal-800 dark:ring-teal-900/50"
                      : "transition-shadow hover:shadow-md"
                  }
                >
                  <div className="flex flex-wrap items-center gap-2">
                    {isRecommended && (
                      <Badge tone="success">Recommended</Badge>
                    )}
                    <Badge tone="info">{formatCategory(q.type)}</Badge>
                    <Badge>{formatCategory(q.category)}</Badge>
                    <Badge tone="default">{formatLevel(q.level)}</Badge>
                  </div>
                  <h3 className="mt-3 font-medium text-zinc-900 dark:text-zinc-50">
                    {q.title}
                  </h3>
                  <p className="mt-2 line-clamp-2 text-sm text-zinc-500">
                    {q.content}
                  </p>
                </Card>
              </Link>
            </li>
          );
        })}
      </ul>
    </>
  );
}
