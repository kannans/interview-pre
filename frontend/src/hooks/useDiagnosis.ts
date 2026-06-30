"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { diagnosisApi, planApi, isApiError } from "@/lib/api";
import { queryKeys } from "@/lib/query/keys";
import { toastError, toastSuccess } from "@/lib/toast";
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
      toastSuccess("Diagnosis saved");
    },
    onError: (error) => toastError(error, "Diagnosis failed"),
  });
}

/** Submit baseline diagnosis, then generate or replace the recovery plan. */
export function useSubmitDiagnosisAndPlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: DiagnosisInput) => {
      const diagnosis = await diagnosisApi.submit(input);
      try {
        const plan = await planApi.generate();
        return { diagnosis, plan };
      } catch (planErr) {
        queryClient.setQueryData(queryKeys.diagnosis, diagnosis);
        if (isApiError(planErr)) {
          throw new Error(
            `Diagnosis saved, but plan generation failed: ${planErr.message}`,
          );
        }
        throw planErr;
      }
    },
    onSuccess: ({ diagnosis, plan }) => {
      queryClient.setQueryData(queryKeys.diagnosis, diagnosis);
      queryClient.setQueryData(queryKeys.plan, plan);
      queryClient.invalidateQueries({ queryKey: queryKeys.todayTask });
      toastSuccess("Diagnosis complete — recovery plan generated");
    },
    onError: (error) => toastError(error, "Diagnosis failed"),
  });
}
