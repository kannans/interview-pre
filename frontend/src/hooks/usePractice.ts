"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { questionsApi, attemptsApi, progressApi, isApiError } from "@/lib/api";
import { queryKeys } from "@/lib/query/keys";
import type { AttemptInput, QuestionListFilter } from "@/lib/types/api";

export function useQuestions(filter?: QuestionListFilter) {
  return useQuery({
    queryKey: queryKeys.questions(filter),
    queryFn: () => questionsApi.list(filter),
  });
}

export function useQuestion(id: string) {
  return useQuery({
    queryKey: queryKeys.question(id),
    queryFn: () => questionsApi.get(id),
    enabled: !!id,
  });
}

export function useAttempts(questionId?: string) {
  return useQuery({
    queryKey: queryKeys.attempts(questionId),
    queryFn: () => attemptsApi.list(questionId),
    retry: (_, error) => !isApiError(error) || !error.isNotFound,
  });
}

export function useSubmitAttempt() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: AttemptInput) => attemptsApi.submit(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["attempts"] });
    },
  });
}

export function useEvaluateAttempt() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (attemptId: string) => attemptsApi.evaluate(attemptId),
    onSuccess: async (_, attemptId) => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.evaluation(attemptId),
      });

      try {
        await progressApi.createSnapshot();
      } catch (error) {
        if (!isApiError(error) || !error.isConflict) {
          // Evaluation succeeded; snapshot refresh is best-effort.
        }
      }

      queryClient.invalidateQueries({ queryKey: queryKeys.progressLatest });
      queryClient.invalidateQueries({ queryKey: ["progress", "snapshots"] });
    },
  });
}

export function useEvaluation(attemptId: string | null) {
  return useQuery({
    queryKey: queryKeys.evaluation(attemptId ?? ""),
    queryFn: () => attemptsApi.getEvaluation(attemptId!),
    enabled: !!attemptId,
    retry: (_, error) => !isApiError(error) || !error.isNotFound,
  });
}
