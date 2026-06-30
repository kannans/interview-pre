"use client";

import { useMemo } from "react";
import { useProfile } from "./useProfile";
import { useDiagnosis } from "./useDiagnosis";
import { usePlan, useTodayTask } from "./usePlan";
import { useLatestProgress } from "./useProgress";
import { isApiError } from "@/lib/api";

export type WorkflowStage =
  | "profile"
  | "diagnosis"
  | "plan"
  | "practice"
  | "progress"
  | "adjust";

export function useWorkflowStage() {
  const profile = useProfile();
  const diagnosis = useDiagnosis();
  const plan = usePlan();
  const today = useTodayTask();
  const progress = useLatestProgress();

  const stage = useMemo((): WorkflowStage => {
    if (profile.isError && isApiError(profile.error) && profile.error.isNotFound) {
      return "profile";
    }
    if (
      diagnosis.isError &&
      isApiError(diagnosis.error) &&
      diagnosis.error.isNotFound
    ) {
      return "diagnosis";
    }
    if (plan.isError && isApiError(plan.error) && plan.error.isNotFound) {
      return "plan";
    }
    if (
      today.isError &&
      isApiError(today.error) &&
      today.error.isNotFound
    ) {
      if (progress.data?.stats.adjustment_recommendation) {
        return "adjust";
      }
      return "progress";
    }
    return "practice";
  }, [
    profile.isError,
    profile.error,
    diagnosis.isError,
    diagnosis.error,
    plan.isError,
    plan.error,
    today.isError,
    today.error,
    progress.data?.stats.adjustment_recommendation,
  ]);

  const nextAction = useMemo(() => {
    switch (stage) {
      case "profile":
        return {
          href: "/onboarding",
          label: "Set up profile",
          description: "Tell the coach about your experience and interview goals.",
        };
      case "diagnosis":
        return {
          href: "/diagnosis",
          label: "Run baseline diagnosis",
          description: "Submit sample answers so we can map your weak areas.",
        };
      case "plan":
        return {
          href: "/plan",
          label: "Generate recovery plan",
          description: "Turn your diagnosis into a structured daily plan.",
        };
      case "progress":
        return {
          href: "/progress",
          label: "Review progress",
          description:
            "All plan tasks are done for now — check your readiness dashboard and snapshot history.",
        };
      case "adjust":
        return {
          href: "/progress",
          label: "Review plan adjustment",
          description:
            progress.data?.stats.adjustment_recommendation ??
            "See AI recommendations to tune your recovery plan.",
        };
      default:
        return {
          href: "/practice",
          label: "Start today's practice",
          description: "Work through your next task and get AI feedback.",
        };
    }
  }, [stage, progress.data?.stats.adjustment_recommendation]);

  return {
    stage,
    nextAction,
    isLoading:
      profile.isLoading ||
      diagnosis.isLoading ||
      plan.isLoading ||
      today.isLoading ||
      progress.isLoading,
  };
}
