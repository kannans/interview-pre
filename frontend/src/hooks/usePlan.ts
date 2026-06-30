"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { planApi, isApiError } from "@/lib/api";
import { queryKeys } from "@/lib/query/keys";

export function usePlan() {
  return useQuery({
    queryKey: queryKeys.plan,
    queryFn: () => planApi.get(),
    retry: (_, error) => !isApiError(error) || !error.isNotFound,
  });
}

export function useTodayTask() {
  return useQuery({
    queryKey: queryKeys.todayTask,
    queryFn: () => planApi.getToday(),
    retry: (_, error) => !isApiError(error) || !error.isNotFound,
  });
}

export function useGeneratePlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => planApi.generate(),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.plan, data);
      queryClient.invalidateQueries({ queryKey: queryKeys.todayTask });
    },
  });
}

export function useCompleteTask() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (taskId: string) => planApi.completeTask(taskId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.plan });
      queryClient.invalidateQueries({ queryKey: queryKeys.todayTask });
    },
  });
}
