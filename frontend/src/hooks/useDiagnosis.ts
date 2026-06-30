"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { diagnosisApi, planApi, isApiError } from "@/lib/api";
import { queryKeys } from "@/lib/query/keys";
import type { DiagnosisInput } from "@/lib/types/api";

export function useDiagnosis() {
  return useQuery({
    queryKey: queryKeys.diagnosis,
    queryFn: () => diagnosisApi.get(),
    retry: (_, error) => !isApiError(error) || !error.isNotFound,
  });
}

export function useSubmitDiagnosis() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: DiagnosisInput) => diagnosisApi.submit(input),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.diagnosis, data);
    },
  });
}

/** Submit baseline diagnosis, then generate or replace the recovery plan. */
export function useSubmitDiagnosisAndPlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: DiagnosisInput) => {
      const diagnosis = await diagnosisApi.submit(input);
      const plan = await planApi.generate();
      return { diagnosis, plan };
    },
    onSuccess: ({ diagnosis, plan }) => {
      queryClient.setQueryData(queryKeys.diagnosis, diagnosis);
      queryClient.setQueryData(queryKeys.plan, plan);
      queryClient.invalidateQueries({ queryKey: queryKeys.todayTask });
    },
  });
}
