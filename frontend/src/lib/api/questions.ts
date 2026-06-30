import { apiRequest } from "./client";
import type {
  ListQuestionsResponse,
  Question,
  QuestionListFilter,
} from "@/lib/types/api";

function toQuery(filter?: QuestionListFilter): string {
  if (!filter) return "";
  const params = new URLSearchParams();
  if (filter.type) params.set("type", filter.type);
  if (filter.category) params.set("category", filter.category);
  if (filter.level) params.set("level", filter.level);
  const qs = params.toString();
  return qs ? `?${qs}` : "";
}

export const questionsApi = {
  list: (filter?: QuestionListFilter) =>
    apiRequest<ListQuestionsResponse>(`/questions${toQuery(filter)}`),
  get: (id: string) => apiRequest<Question>(`/questions/${id}`),
};
