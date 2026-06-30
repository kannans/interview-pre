import type { PlanTask, Question } from "@/lib/types/api";

/** Pick a practice question aligned with today's recovery task. */
export function pickRecommendedQuestion(
  questions: Question[],
  task: PlanTask,
): Question | null {
  if (questions.length === 0) return null;

  const byType = questions.filter((q) => q.type === task.category);
  const pool = byType.length > 0 ? byType : questions;

  const exactLevel = pool.find((q) => q.level === task.level);
  if (exactLevel) return exactLevel;

  const levelOrder = ["level_1", "level_2", "level_3"];
  const taskLevelIdx = levelOrder.indexOf(task.level);

  if (taskLevelIdx >= 0) {
    const sorted = [...pool].sort((a, b) => {
      const aDist = Math.abs(levelOrder.indexOf(a.level) - taskLevelIdx);
      const bDist = Math.abs(levelOrder.indexOf(b.level) - taskLevelIdx);
      return aDist - bDist;
    });
    return sorted[0] ?? null;
  }

  return pool[0] ?? null;
}

export function practiceHrefForTask(
  questions: Question[],
  task: PlanTask | undefined,
): string {
  if (!task) return "/practice";
  const recommended = pickRecommendedQuestion(questions, task);
  return recommended ? `/practice/${recommended.id}` : "/practice";
}
