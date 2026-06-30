"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { diagnosisApi, isApiError } from "@/lib/api";
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
