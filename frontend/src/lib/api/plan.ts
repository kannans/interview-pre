import { apiRequest } from "./client";
import type {
  PlanTask,
  PlanWithTasks,
  TodayTaskResponse,
} from "@/lib/types/api";

export const planApi = {
  get: () => apiRequest<PlanWithTasks>("/plan"),
  generate: () => apiRequest<PlanWithTasks>("/plan", { method: "POST" }),
  getToday: () => apiRequest<TodayTaskResponse>("/plan/today"),
  completeTask: (taskId: string) =>
    apiRequest<PlanTask>(`/plan/tasks/${taskId}/complete`, { method: "PATCH" }),
};
