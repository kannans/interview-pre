"use client";

import { Card } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import type { PlanTask, ProgressSnapshot } from "@/lib/types/api";
import { formatCategory, formatLevel } from "@/lib/utils/format";

type TodayTaskCardProps = {
  task?: PlanTask;
  loading?: boolean;
  error?: string;
  onComplete?: () => void;
  completing?: boolean;
  practiceHref?: string;
};

export function TodayTaskCard({
  task,
  loading,
  error,
  onComplete,
  completing,
  practiceHref = "/practice",
}: TodayTaskCardProps) {
  if (loading) {
    return (
      <Card title="Today's task">
        <div className="h-20 animate-pulse rounded-lg bg-zinc-200 dark:bg-zinc-800" />
      </Card>
    );
  }

  if (error) {
    return (
      <Card title="Today's task">
        <p className="text-sm text-zinc-500">{error}</p>
        <Button variant="secondary" size="sm" className="mt-3" href="/plan">
          View recovery plan
        </Button>
      </Card>
    );
  }

  if (!task) {
    return (
      <Card title="Today's task">
        <p className="text-sm text-zinc-500">
          No active task yet. Complete onboarding and generate your plan.
        </p>
        <Button size="sm" className="mt-3" href="/plan">
          Set up plan
        </Button>
      </Card>
    );
  }

  return (
    <Card
      title="Today's task"
      description="Your next recovery step"
      action={
        <Badge tone="info">
          {formatCategory(task.category)} · {formatLevel(task.level)}
        </Badge>
      }
    >
      <p className="text-lg font-medium text-zinc-900 dark:text-zinc-50">
        {task.title}
      </p>
      <div className="mt-4 flex flex-wrap gap-2">
        <Button size="sm" href={practiceHref}>
          Start practice
        </Button>
        {onComplete && (
          <Button
            variant="secondary"
            size="sm"
            loading={completing}
            onClick={onComplete}
          >
            Mark complete
          </Button>
        )}
      </div>
    </Card>
  );
}

type ProgressSummaryProps = {
  snapshot?: ProgressSnapshot | null;
  loading?: boolean;
};

export function ProgressSummary({ snapshot, loading }: ProgressSummaryProps) {
  if (loading) {
    return (
      <Card title="Progress">
        <div className="h-16 animate-pulse rounded-lg bg-zinc-200 dark:bg-zinc-800" />
      </Card>
    );
  }

  if (!snapshot) {
    return (
      <Card title="Progress">
        <p className="text-sm text-zinc-500">
          No progress snapshot yet. Practice and evaluate attempts to build your
          dashboard.
        </p>
        <Button variant="secondary" size="sm" className="mt-3" href="/progress">
          Open progress
        </Button>
      </Card>
    );
  }

  const { stats } = snapshot;
  const planPct =
    stats.plan_tasks_total > 0
      ? Math.round(
          (stats.plan_tasks_completed / stats.plan_tasks_total) * 100,
        )
      : 0;

  return (
    <Card title="Progress" description={`Updated ${snapshot.snapshot_date}`}>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        <Stat label="Readiness" value={`${stats.readiness_index}%`} />
        <Stat label="Attempts" value={String(stats.attempts_count)} />
        <Stat label="Evaluations" value={String(stats.evaluations_count)} />
        <Stat label="Plan done" value={`${planPct}%`} />
      </div>
    </Card>
  );
}

function Stat({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <p className="text-xs uppercase tracking-wide text-zinc-500">{label}</p>
      <p className="mt-1 text-xl font-semibold text-zinc-900 dark:text-zinc-50">
        {value}
      </p>
    </div>
  );
}

type WeakAreasPanelProps = {
  weakAreas?: string[];
  loading?: boolean;
};

export function WeakAreasPanel({ weakAreas, loading }: WeakAreasPanelProps) {
  if (loading) {
    return (
      <Card title="Weak areas">
        <div className="h-12 animate-pulse rounded-lg bg-zinc-200 dark:bg-zinc-800" />
      </Card>
    );
  }

  const items = weakAreas ?? [];

  return (
    <Card title="Weak areas" description="Focus where improvement matters most">
      {items.length === 0 ? (
        <p className="text-sm text-zinc-500">
          Complete baseline diagnosis to see your weak areas.
        </p>
      ) : (
        <ul className="flex flex-wrap gap-2">
          {items.map((area) => (
            <li key={area}>
              <Badge tone="warning">{formatCategory(area)}</Badge>
            </li>
          ))}
        </ul>
      )}
      <Button variant="ghost" size="sm" className="mt-3" href="/diagnosis">
        Review diagnosis
      </Button>
    </Card>
  );
}

type NextActionCardProps = {
  href: string;
  label: string;
  description: string;
};

export function NextActionCard({ href, label, description }: NextActionCardProps) {
  return (
    <Card title="Next action">
      <p className="text-sm text-zinc-600 dark:text-zinc-300">{description}</p>
      <Button className="mt-4" href={href}>
        {label}
      </Button>
    </Card>
  );
}
