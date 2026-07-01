import type { PracticeLanguage } from "@/lib/practice/languages";

export type PracticeDraft = {
  code: string;
  stdin: string;
  updatedAt: number;
};

export function practiceDraftKey(
  questionId: string,
  language: PracticeLanguage,
): string {
  return `${questionId}:${language}`;
}

export function isMeaningfulDraft(
  draft: PracticeDraft | undefined,
): draft is PracticeDraft {
  if (!draft) return false;
  return draft.code.trim().length > 0 || draft.stdin.trim().length > 0;
}
