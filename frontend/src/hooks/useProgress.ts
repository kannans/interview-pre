"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { progressApi, isApiError } from "@/lib/api";
import { queryKeys } from "@/lib/query/keys";

export function useLatestProgress() {
  return useQuery({
    queryKey: queryKeys.progressLatest,
    queryFn: () => progressApi.getLatest(),
    retry: (_, error) => !isApiError(error) || !error.isNotFound,
  });
}

export function useProgressSnapshots(limit = 14) {
  return useQuery({
    queryKey: queryKeys.progressSnapshots({ limit }),
    queryFn: () => progressApi.listSnapshots({ limit }),
  });
}

export function useCreateSnapshot() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => progressApi.createSnapshot(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.progressLatest });
      queryClient.invalidateQueries({ queryKey: queryKeys.progressSnapshots({}) });
    },
  });
}
