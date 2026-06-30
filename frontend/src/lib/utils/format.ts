export function formatCategory(value: string): string {
  return value.replace(/_/g, " ");
}

export function formatLevel(level: string): string {
  return level.replace("level_", "Level ");
}

export function formatPercent(value: number): string {
  return `${Math.round(value)}%`;
}

export function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  if (mins === 0) return `${secs}s`;
  return `${mins}m ${secs}s`;
}
