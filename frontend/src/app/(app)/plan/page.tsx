"use client";

import { TopBar } from "@/components/layout/TopBar";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import { ErrorAlert } from "@/components/ui/ErrorAlert";
import { EmptyState } from "@/components/ui/EmptyState";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import {
  useCompleteTask,
  useGeneratePlan,
  usePlan,
} from "@/hooks/usePlan";
import { formatCategory, formatLevel } from "@/lib/utils/format";
import { isApiError } from "@/lib/api";

export default function PlanPage() {
  const plan = usePlan();
  const generate = useGeneratePlan();
  const completeTask = useCompleteTask();

  if (plan.isLoading) {
    return <LoadingSpinner label="Loading recovery plan…" />;
  }

  const hasPlan = plan.data && plan.data.tasks.length > 0;
  const completedCount =
    plan.data?.tasks.filter((t) => t.is_completed).length ?? 0;
  const totalCount = plan.data?.tasks.length ?? 0;

  return (
    <>
      <TopBar
        title="Recovery plan"
        subtitle="Structured daily tasks generated from your baseline diagnosis."
        actions={
          <Button
            onClick={() => generate.mutate()}
            loading={generate.isPending}
            variant={hasPlan ? "secondary" : "primary"}
          >
            {hasPlan ? "Regenerate plan" : "Generate plan"}
          </Button>
        }
      />

      {plan.isError && isApiError(plan.error) && plan.error.isNotFound && (
        <EmptyState
          title="No recovery plan yet"
          description="Complete your profile and baseline diagnosis, then generate an AI recovery plan."
          action={
            <Button onClick={() => generate.mutate()} loading={generate.isPending}>
              Generate recovery plan
            </Button>
          }
        />
      )}

      {plan.isError && !(isApiError(plan.error) && plan.error.isNotFound) && (
        <ErrorAlert message={plan.error?.message ?? "Failed to load plan"} />
      )}

      {hasPlan && plan.data && (
        <>
          <Card className="mb-4" title="Plan overview">
            <p className="text-sm text-zinc-600 dark:text-zinc-300">
              {completedCount} of {totalCount} tasks completed
            </p>
            <div className="mt-2 h-2 overflow-hidden rounded-full bg-zinc-200 dark:bg-zinc-800">
              <div
                className="h-full rounded-full bg-teal-600 transition-all"
                style={{
                  width: `${totalCount ? (completedCount / totalCount) * 100 : 0}%`,
                }}
              />
            </div>
          </Card>

          <ul className="space-y-3">
            {plan.data.tasks.map((task) => (
              <li key={task.id}>
                <Card>
                  <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                    <div>
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="text-xs font-medium text-zinc-500">
                          #{task.sequence_order}
                        </span>
                        <Badge tone={task.is_completed ? "success" : "default"}>
                          {task.is_completed ? "Done" : formatCategory(task.category)}
                        </Badge>
                        <Badge tone="info">{formatLevel(task.level)}</Badge>
                      </div>
                      <p className="mt-2 font-medium text-zinc-900 dark:text-zinc-50">
                        {task.title}
                      </p>
                    </div>
                    {!task.is_completed && (
                      <Button
                        size="sm"
                        variant="secondary"
                        loading={
                          completeTask.isPending &&
                          completeTask.variables === task.id
                        }
                        onClick={() => completeTask.mutate(task.id)}
                      >
                        Mark complete
                      </Button>
                    )}
                  </div>
                </Card>
              </li>
            ))}
          </ul>
        </>
      )}
    </>
  );
}
