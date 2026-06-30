"use client";

import { useMemo } from "react";
import { useTodayTask } from "./usePlan";
import { useQuestions } from "./usePractice";
import {
  pickRecommendedQuestion,
  practiceHrefForTask,
} from "@/lib/utils/recommended-question";

/** Today's plan task plus a matching practice question from the question bank. */
export function useRecommendedPractice() {
  const today = useTodayTask();
  const questions = useQuestions();

  const task = today.data?.task;
  const allQuestions = questions.data?.questions ?? [];

  const recommended = useMemo(() => {
    if (!task || allQuestions.length === 0) return null;
    return pickRecommendedQuestion(allQuestions, task);
  }, [task, allQuestions]);

  const practiceHref = useMemo(
    () => practiceHrefForTask(allQuestions, task),
    [allQuestions, task],
  );

  return {
    today,
    questions,
    task,
    recommended,
    practiceHref,
    isLoading: today.isLoading || questions.isLoading,
  };
}
