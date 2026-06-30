"use client";

import Link from "next/link";
import { TopBar } from "@/components/layout/TopBar";
import { Card } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { ErrorAlert } from "@/components/ui/ErrorAlert";
import { EmptyState } from "@/components/ui/EmptyState";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { useQuestions } from "@/hooks/usePractice";
import { formatCategory, formatLevel } from "@/lib/utils/format";

export default function PracticePage() {
  const { data, isLoading, isError, error } = useQuestions();

  if (isLoading) {
    return <LoadingSpinner label="Loading questions…" />;
  }

  return (
    <>
      <TopBar
        title="Daily practice"
        subtitle="Work through programming and DSA questions aligned with your recovery plan."
      />

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
        {data?.questions.map((q) => (
          <li key={q.id}>
            <Link href={`/practice/${q.id}`}>
              <Card className="transition-shadow hover:shadow-md">
                <div className="flex flex-wrap gap-2">
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
        ))}
      </ul>
    </>
  );
}
