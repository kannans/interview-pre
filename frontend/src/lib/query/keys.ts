export const queryKeys = {
  profile: ["profile"] as const,
  diagnosis: ["diagnosis"] as const,
  plan: ["plan"] as const,
  todayTask: ["plan", "today"] as const,
  questions: (filter?: { type?: string; category?: string; level?: string }) =>
    ["questions", filter ?? {}] as const,
  question: (id: string) => ["questions", id] as const,
  attempts: (questionId?: string) => ["attempts", questionId ?? "all"] as const,
  attempt: (id: string) => ["attempts", "detail", id] as const,
  evaluation: (attemptId: string) => ["evaluation", attemptId] as const,
  progressLatest: ["progress", "latest"] as const,
  progressSnapshots: (params?: Record<string, unknown>) =>
    ["progress", "snapshots", params ?? {}] as const,
};
