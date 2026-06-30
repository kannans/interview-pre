"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { profileApi, isApiError } from "@/lib/api";
import { queryKeys } from "@/lib/query/keys";
import type { ProfileInput } from "@/lib/types/api";

export function useProfile() {
  return useQuery({
    queryKey: queryKeys.profile,
    queryFn: () => profileApi.get(),
    retry: (_, error) => !isApiError(error) || !error.isNotFound,
  });
}

export function useUpsertProfile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: ProfileInput) => profileApi.upsert(input),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.profile, data);
    },
  });
}
