import { apiRequest } from "./client";
import type {
  ListProgressSnapshotsResponse,
  ProgressSnapshot,
} from "@/lib/types/api";

export type SnapshotListParams = {
  limit?: number;
  start_date?: string;
  end_date?: string;
};

function toQuery(params?: SnapshotListParams): string {
  if (!params) return "";
  const search = new URLSearchParams();
  if (params.limit) search.set("limit", String(params.limit));
  if (params.start_date) search.set("start_date", params.start_date);
  if (params.end_date) search.set("end_date", params.end_date);
  const qs = search.toString();
  return qs ? `?${qs}` : "";
}

export const progressApi = {
  createSnapshot: () =>
    apiRequest<ProgressSnapshot>("/progress/snapshot", { method: "POST" }),
  getLatest: () => apiRequest<ProgressSnapshot>("/progress/latest"),
  listSnapshots: (params?: SnapshotListParams) =>
    apiRequest<ListProgressSnapshotsResponse>(
      `/progress/snapshots${toQuery(params)}`,
    ),
};
