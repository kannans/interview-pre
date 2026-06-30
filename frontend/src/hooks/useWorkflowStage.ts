"use client";

import { useMemo } from "react";
import { useProfile } from "./useProfile";
import { useDiagnosis } from "./useDiagnosis";
import { usePlan } from "./usePlan";
import { isApiError } from "@/lib/api";

export type WorkflowStage =
  | "profile"
  | "diagnosis"
  | "plan"
  | "practice"
  | "progress";

export function useWorkflowStage() {
  const profile = useProfile();
  const diagnosis = useDiagnosis();
  const plan = usePlan();

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
    return "practice";
  }, [profile.isError, profile.error, diagnosis.isError, diagnosis.error, plan.isError, plan.error]);

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
      default:
        return {
          href: "/practice",
          label: "Start today's practice",
          description: "Work through your next task and get AI feedback.",
        };
    }
  }, [stage]);

  return { stage, nextAction, isLoading: profile.isLoading };
}
