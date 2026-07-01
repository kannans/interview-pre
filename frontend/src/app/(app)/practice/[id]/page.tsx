"use client";

import { FormEvent, useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { TopBar } from "@/components/layout/TopBar";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import { ErrorAlert } from "@/components/ui/ErrorAlert";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { CodeEditor } from "@/components/practice/CodeEditor";
import { PracticeSubmitBar } from "@/components/practice/PracticeSubmitBar";
import { PracticeWorkspace } from "@/components/practice/PracticeWorkspace";
import {
  useEvaluateAttempt,
  useAttempts,
  useQuestion,
  useRunCode,
  useSubmitAttempt,
} from "@/hooks/usePractice";
import { useCompleteTask, useTodayTask } from "@/hooks/usePlan";
import { usePracticeDraftsStore, usePracticeDraftsHydrated } from "@/lib/store/practice-drafts";
import { useSettingsStore } from "@/lib/store/settings";
import { formatCategory, formatLevel } from "@/lib/utils/format";
import { isApiError } from "@/lib/api";
import { attemptsApi } from "@/lib/api/attempts";
import {
  getStarterCode,
  monacoLanguageId,
  type PracticeLanguage,
} from "@/lib/practice/languages";
import type { Attempt, RunCodeResult } from "@/lib/types/api";

function latestAttemptForLanguage(
  attempts: Attempt[] | undefined,
  language: PracticeLanguage,
): Attempt | undefined {
  if (!attempts?.length) return undefined;
  return attempts
    .filter((a) => a.language === language)
    .sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    )[0];
}

export default function PracticeQuestionPage() {
  const params = useParams<{ id: string }>();
  const questionId = params.id;

  const question = useQuestion(questionId);
  const attempts = useAttempts(questionId);
  const today = useTodayTask();
  const completeTask = useCompleteTask();
  const submitAttempt = useSubmitAttempt();
  const evaluateAttempt = useEvaluateAttempt();
  const runCode = useRunCode();
  const defaultMinutes = useSettingsStore((s) => s.defaultPracticeMinutes);
  const defaultLanguage = useSettingsStore((s) => s.defaultPracticeLanguage);
  const getDraft = usePracticeDraftsStore((s) => s.getDraft);
  const setDraft = usePracticeDraftsStore((s) => s.setDraft);
  const clearDraft = usePracticeDraftsStore((s) => s.clearDraft);
  const draftsHydrated = usePracticeDraftsHydrated();

  const [language, setLanguage] = useState<PracticeLanguage>(defaultLanguage);
  const [answer, setAnswer] = useState("");
  const [stdin, setStdin] = useState("");
  const [runResult, setRunResult] = useState<RunCodeResult | null>(null);
  const [loadedKey, setLoadedKey] = useState<string | null>(null);
  const [elapsed, setElapsed] = useState(0);
  const [evaluation, setEvaluation] = useState<{
    score_correctness: number;
    score_approach: number;
    score_edge_cases: number;
    score_communication: number;
    feedback_correctness: string;
    feedback_approach: string;
    feedback_edge_cases: string;
    feedback_communication: string;
  } | null>(null);

  const loadKey = `${questionId}:${language}`;
  const draftReady = loadedKey === loadKey && !evaluation;
  const contentReady = loadedKey === loadKey;

  useEffect(() => {
    setLanguage(defaultLanguage);
  }, [defaultLanguage]);

  useEffect(() => {
    setLoadedKey(null);
    setRunResult(null);
    setAnswer("");
    setStdin("");
  }, [questionId]);

  useEffect(() => {
    if (!question.data || !draftsHydrated || loadedKey === loadKey) return;

    const draft = getDraft(questionId, language);
    if (draft) {
      setAnswer(draft.code);
      setStdin(draft.stdin);
      setRunResult(null);
      setLoadedKey(loadKey);
      return;
    }

    if (attempts.isLoading) return;

    const previous = latestAttemptForLanguage(attempts.data?.attempts, language);
    setAnswer(
      previous?.submitted_content ??
        getStarterCode(question.data, language),
    );
    setStdin("");
    setRunResult(null);
    setLoadedKey(loadKey);
  }, [
    question.data,
    language,
    questionId,
    attempts.data,
    attempts.isLoading,
    loadedKey,
    loadKey,
    getDraft,
    draftsHydrated,
  ]);

  useEffect(() => {
    if (!question.data || evaluation || loadedKey !== loadKey) return;
    setDraft(questionId, language, { code: answer, stdin });
  }, [
    answer,
    stdin,
    questionId,
    language,
    question.data,
    evaluation,
    loadedKey,
    loadKey,
    setDraft,
  ]);

  useEffect(() => {
    const timer = window.setInterval(() => setElapsed((s) => s + 1), 1000);
    return () => window.clearInterval(timer);
  }, []);

  function handleRun() {
    setRunResult(null);
    runCode.mutate(
      { language, code: answer, stdin: stdin || undefined },
      { onSuccess: setRunResult },
    );
  }

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    const timeTaken = Math.max(elapsed, defaultMinutes * 60);

    submitAttempt.mutate(
      {
        question_id: questionId,
        submitted_content: answer,
        time_taken_seconds: timeTaken,
        language,
      },
      {
        onSuccess: (attempt) => {
          clearDraft(questionId, language);
          evaluateAttempt.mutate(attempt.id, {
            onSuccess: (ev) => {
              setEvaluation(ev);
              const task = today.data?.task;
              if (task && !task.is_completed) {
                completeTask.mutate(task.id);
              }
            },
            onError: async (error) => {
              if (isApiError(error) && error.isConflict) {
                try {
                  const ev = await attemptsApi.getEvaluation(attempt.id);
                  setEvaluation(ev);
                } catch {
                  // keep existing error UI
                }
              }
            },
          });
        },
      },
    );
  }

  if (question.isLoading || !draftsHydrated) {
    return <LoadingSpinner label="Loading question…" />;
  }

  if (question.isError) {
    return (
      <ErrorAlert
        message={
          isApiError(question.error)
            ? question.error.message
            : "Question not found"
        }
      />
    );
  }

  const q = question.data!;
  const isConcept = q.type === "concept";
  const editorLanguage = isConcept ? "markdown" : monacoLanguageId(language);
  const readOnly = !!evaluation;
  const submitLoading = submitAttempt.isPending || evaluateAttempt.isPending;

  return (
    <>
      <TopBar
        title={q.title}
        subtitle="Run locally to debug, then submit for OpenAI evaluation."
        actions={
          <div className="flex items-center gap-2">
            <Badge tone="info">{formatCategory(q.type)}</Badge>
            <Badge>{formatLevel(q.level)}</Badge>
            <span className="text-sm tabular-nums text-zinc-500">
              {Math.floor(elapsed / 60)}:{String(elapsed % 60).padStart(2, "0")}
            </span>
          </div>
        }
      />

      <div className="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.25fr)] lg:items-start">
        <Card
          title="Problem"
          className="lg:sticky lg:top-4"
        >
          <p className="whitespace-pre-wrap text-sm text-zinc-700 dark:text-zinc-300">
            {q.content}
          </p>
          {q.constraints && (
            <p className="mt-3 text-sm text-zinc-500">
              <span className="font-medium">Constraints:</span> {q.constraints}
            </p>
          )}
        </Card>

        {isConcept ? (
          contentReady ? (
            <Card title="Your answer">
              <div className="space-y-4">
                <CodeEditor
                  value={answer}
                  onChange={setAnswer}
                  language={language}
                  editorLanguage={editorLanguage}
                  readOnly={readOnly}
                  minHeight="20rem"
                />
                {draftReady && !readOnly && (
                  <p className="text-xs text-zinc-500">Draft saved locally</p>
                )}
              </div>
            </Card>
          ) : (
            <Card title="Your answer">
              <LoadingSpinner label="Loading workspace…" />
            </Card>
          )
        ) : contentReady ? (
          <PracticeWorkspace
            language={language}
            onLanguageChange={setLanguage}
            editorLanguage={editorLanguage}
            answer={answer}
            onAnswerChange={setAnswer}
            stdin={stdin}
            onStdinChange={setStdin}
            runResult={runResult}
            runPending={runCode.isPending}
            onRun={handleRun}
            readOnly={readOnly}
            draftReady={draftReady}
          />
        ) : (
          <div className="flex min-h-[36rem] items-center justify-center rounded-xl border border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-900">
            <LoadingSpinner label="Loading workspace…" />
          </div>
        )}
      </div>

      {!evaluation && (
        <PracticeSubmitBar
          isConcept={isConcept}
          loading={submitLoading}
          disabled={readOnly}
          onSubmit={handleSubmit}
        />
      )}

      {evaluation && (
        <div className="mt-4 flex justify-end">
          <Button variant="secondary" href="/">
            Back to dashboard
          </Button>
        </div>
      )}

      {evaluation && (
        <Card title="AI feedback" className="mt-4">
          <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
            <Score label="Correctness" value={evaluation.score_correctness} />
            <Score label="Approach" value={evaluation.score_approach} />
            <Score label="Edge cases" value={evaluation.score_edge_cases} />
            <Score label="Communication" value={evaluation.score_communication} />
          </div>
          <dl className="mt-4 space-y-3 text-sm">
            <FeedbackItem label="Correctness" text={evaluation.feedback_correctness} />
            <FeedbackItem label="Approach" text={evaluation.feedback_approach} />
            <FeedbackItem label="Edge cases" text={evaluation.feedback_edge_cases} />
            <FeedbackItem label="Communication" text={evaluation.feedback_communication} />
          </dl>
        </Card>
      )}
    </>
  );
}

function Score({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-lg bg-zinc-100 p-3 dark:bg-zinc-800">
      <p className="text-xs text-zinc-500">{label}</p>
      <p className="text-xl font-semibold">{value}/10</p>
    </div>
  );
}

function FeedbackItem({ label, text }: { label: string; text: string }) {
  return (
    <div>
      <dt className="font-medium text-zinc-700 dark:text-zinc-300">{label}</dt>
      <dd className="mt-1 text-zinc-600 dark:text-zinc-400">{text}</dd>
    </div>
  );
}
