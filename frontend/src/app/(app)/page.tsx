"use client";

import { TopBar } from "@/components/layout/TopBar";
import {
  NextActionCard,
  ProgressSummary,
  TodayTaskCard,
  WeakAreasPanel,
} from "@/components/dashboard/DashboardCards";
import { useCompleteTask, useTodayTask } from "@/hooks/usePlan";
import { useLatestProgress } from "@/hooks/useProgress";
import { useDiagnosis } from "@/hooks/useDiagnosis";
import { useWorkflowStage } from "@/hooks/useWorkflowStage";
import { parseDiagnosticReport } from "@/lib/utils/diagnostic-report";
import { isApiError } from "@/lib/api";

export default function HomePage() {
  const today = useTodayTask();
  const progress = useLatestProgress();
  const diagnosis = useDiagnosis();
  const { nextAction } = useWorkflowStage();
  const completeTask = useCompleteTask();

  const report = diagnosis.data
    ? parseDiagnosticReport(diagnosis.data.diagnostic_report)
    : null;

  const todayError =
    today.isError && isApiError(today.error) && today.error.isNotFound
      ? "All tasks completed or no plan yet."
      : today.isError
        ? today.error.message
        : undefined;

  return (
    <>
      <TopBar
        title="Today"
        subtitle="Your recovery command center — task, progress, and next step at a glance."
      />

      <div className="grid gap-4 lg:grid-cols-2">
        <div className="space-y-4 lg:col-span-2">
          <TodayTaskCard
            task={today.data?.task}
            loading={today.isLoading}
            error={todayError}
            onComplete={
              today.data?.task
                ? () => completeTask.mutate(today.data!.task.id)
                : undefined
            }
            completing={completeTask.isPending}
          />
        </div>

        <ProgressSummary
          snapshot={progress.data}
          loading={progress.isLoading}
        />

        <WeakAreasPanel
          weakAreas={
            report?.weak_areas ??
            progress.data?.stats.weak_patterns ??
            progress.data?.stats.focus_areas
          }
          loading={diagnosis.isLoading && progress.isLoading}
        />

        <div className="lg:col-span-2">
          <NextActionCard {...nextAction} />
        </div>
      </div>
    </>
  );
}
