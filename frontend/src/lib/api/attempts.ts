import { apiRequest } from "./client";
import type {
  Attempt,
  AttemptInput,
  Evaluation,
  ListAttemptsResponse,
} from "@/lib/types/api";

export const attemptsApi = {
  submit: (input: AttemptInput) =>
    apiRequest<Attempt>("/attempts", { method: "POST", body: input }),
  list: (questionId?: string) => {
    const qs = questionId ? `?question_id=${encodeURIComponent(questionId)}` : "";
    return apiRequest<ListAttemptsResponse>(`/attempts${qs}`);
  },
  get: (id: string) => apiRequest<Attempt>(`/attempts/${id}`),
  evaluate: (id: string) =>
    apiRequest<Evaluation>(`/attempts/${id}/evaluate`, { method: "POST" }),
  getEvaluation: (id: string) =>
    apiRequest<Evaluation>(`/attempts/${id}/evaluation`),
};
