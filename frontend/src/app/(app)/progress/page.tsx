"use client";

import { TopBar } from "@/components/layout/TopBar";
import { Card } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { ErrorAlert } from "@/components/ui/ErrorAlert";
import { EmptyState } from "@/components/ui/EmptyState";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import {
  useCreateSnapshot,
  useLatestProgress,
  useProgressSnapshots,
} from "@/hooks/useProgress";
import { formatCategory, formatDate } from "@/lib/utils/format";
import { isApiError } from "@/lib/api";

export default function ProgressPage() {
  const latest = useLatestProgress();
  const history = useProgressSnapshots(14);
  const createSnapshot = useCreateSnapshot();

  if (latest.isLoading) {
    return <LoadingSpinner label="Loading progress…" />;
  }

  const snapshot = latest.data;
  const isDegraded = snapshot?.stats.status === "degraded_ai_summary";

  return (
    <>
      <TopBar
        title="Progress dashboard"
        subtitle="Track readiness, weak patterns, and weekly AI summaries."
        actions={
          <Button
            onClick={() => createSnapshot.mutate()}
            loading={createSnapshot.isPending}
            variant="secondary"
          >
            Refresh snapshot
          </Button>
        }
      />

      {latest.isError &&
        isApiError(latest.error) &&
        !latest.error.isNotFound && (
          <ErrorAlert
            className="mb-4"
            message={latest.error.message}
            onRetry={() => latest.refetch()}
          />
        )}

      {latest.isError && isApiError(latest.error) && latest.error.isNotFound && (
        <EmptyState
          title="No progress data yet"
          description="Complete practice attempts and evaluations, then create your first daily snapshot."
          action={
            <Button
              onClick={() => createSnapshot.mutate()}
              loading={createSnapshot.isPending}
            >
              Create first snapshot
            </Button>
          }
        />
      )}

      {snapshot && (
        <>
          {isDegraded && (
            <Card className="mb-4 border-amber-200 bg-amber-50 dark:border-amber-900/50 dark:bg-amber-950/20">
              <p className="text-sm text-amber-800 dark:text-amber-300">
                Weekly AI summary unavailable — stats saved without narrative.
                Retry tomorrow or when the AI service recovers.
              </p>
            </Card>
          )}

          <div className="grid gap-4 lg:grid-cols-3">
            <Card title="Readiness index">
              <p className="text-4xl font-bold text-teal-600">
                {snapshot.stats.readiness_index}%
              </p>
              <p className="mt-2 text-sm text-zinc-500">
                Based on evaluations, plan completion, and practice volume.
              </p>
            </Card>

            <Card title="Average scores">
              <dl className="space-y-2 text-sm">
                <ScoreRow label="Correctness" value={snapshot.stats.average_scores.correctness} />
                <ScoreRow label="Approach" value={snapshot.stats.average_scores.approach} />
                <ScoreRow label="Edge cases" value={snapshot.stats.average_scores.edge_cases} />
                <ScoreRow label="Communication" value={snapshot.stats.average_scores.communication} />
              </dl>
            </Card>

            <Card title="Activity">
              <dl className="space-y-2 text-sm">
                <StatRow label="Attempts" value={snapshot.stats.attempts_count} />
                <StatRow label="Evaluations" value={snapshot.stats.evaluations_count} />
                <StatRow
                  label="Plan tasks"
                  value={`${snapshot.stats.plan_tasks_completed}/${snapshot.stats.plan_tasks_total}`}
                />
              </dl>
            </Card>
          </div>

          {snapshot.stats.weekly_summary && (
            <Card title="Weekly summary" className="mt-4">
              <p className="text-sm text-zinc-700 dark:text-zinc-300">
                {snapshot.stats.weekly_summary}
              </p>
            </Card>
          )}

          <div className="mt-4 grid gap-4 lg:grid-cols-2">
            <Card title="Weak patterns">
              <div className="flex flex-wrap gap-2">
                {snapshot.stats.weak_patterns.length === 0 ? (
                  <p className="text-sm text-zinc-500">No patterns yet.</p>
                ) : (
                  snapshot.stats.weak_patterns.map((p) => (
                    <Badge key={p} tone="warning">
                      {formatCategory(p)}
                    </Badge>
                  ))
                )}
              </div>
            </Card>

            <Card title="Focus areas">
              {snapshot.stats.focus_areas?.length ? (
                <ul className="list-inside list-disc text-sm text-zinc-600 dark:text-zinc-300">
                  {snapshot.stats.focus_areas.map((area) => (
                    <li key={area}>{area}</li>
                  ))}
                </ul>
              ) : (
                <p className="text-sm text-zinc-500">No focus areas in this snapshot.</p>
              )}
            </Card>
          </div>

          {snapshot.stats.adjustment_recommendation && (
            <Card title="Adjustment recommendation" className="mt-4">
              <p className="text-sm text-zinc-700 dark:text-zinc-300">
                {snapshot.stats.adjustment_recommendation}
              </p>
            </Card>
          )}
        </>
      )}

      {history.data && history.data.snapshots.length > 0 && (
        <Card title="Snapshot history" className="mt-4">
          <ul className="divide-y divide-zinc-200 dark:divide-zinc-800">
            {history.data.snapshots.map((s) => (
              <li
                key={s.id}
                className="flex items-center justify-between py-3 text-sm"
              >
                <span>{formatDate(s.created_at)}</span>
                <span className="font-medium">
                  Readiness {s.stats.readiness_index}%
                </span>
              </li>
            ))}
          </ul>
        </Card>
      )}
    </>
  );
}

function ScoreRow({ label, value }: { label: string; value: number }) {
  return (
    <div className="flex justify-between">
      <dt className="text-zinc-500">{label}</dt>
      <dd className="font-medium">{value.toFixed(1)}/10</dd>
    </div>
  );
}

function StatRow({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="flex justify-between">
      <dt className="text-zinc-500">{label}</dt>
      <dd className="font-medium">{value}</dd>
    </div>
  );
}
