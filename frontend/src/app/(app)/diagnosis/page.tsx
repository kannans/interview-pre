"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { TopBar } from "@/components/layout/TopBar";
import { Card } from "@/components/ui/Card";
import { Textarea } from "@/components/ui/Textarea";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { useDiagnosis, useSubmitDiagnosisAndPlan } from "@/hooks/useDiagnosis";
import { parseDiagnosticReport } from "@/lib/utils/diagnostic-report";
import { formatCategory, formatPercent } from "@/lib/utils/format";

const CODING_PLACEHOLDER =
  "Two Sum: iterate once, store each number in a hash map with its index. For each value, check if (target - value) is in the map — O(n) time, O(n) space. I'd call out empty array and duplicate-index edge cases before coding.";

const CONCEPTUAL_PLACEHOLDER =
  "URL shortener at scale: stateless API behind a load balancer, Postgres for slug→URL mappings, Redis cache for hot links, base62 IDs. Trade-offs: cache consistency vs read latency, shard by slug prefix, rate limit creates at the edge.";

const BEHAVIORAL_PLACEHOLDER =
  "Situation: checkout p99 hit 5s during a flash sale. Task: lead on-call fix. Action: traced to a missing index + cold cache, added index, warmed cache, rolled back bad deploy. Result: p99 under 200ms; added load test gate in CI.";

export default function DiagnosisPage() {
  const router = useRouter();
  const existing = useDiagnosis();
  const submit = useSubmitDiagnosisAndPlan();

  const [codingAnswer, setCodingAnswer] = useState("");
  const [conceptualAnswer, setConceptualAnswer] = useState("");
  const [behavioralAnswer, setBehavioralAnswer] = useState("");

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    submit.mutate(
      {
        coding_answer: codingAnswer,
        conceptual_answer: conceptualAnswer,
        behavioral_answer: behavioralAnswer,
      },
      {
        onSuccess: () => router.push("/"),
      },
    );
  }

  const report = existing.data
    ? parseDiagnosticReport(existing.data.diagnostic_report)
    : null;

  if (existing.isLoading) {
    return <LoadingSpinner label="Loading diagnosis…" />;
  }

  return (
    <>
      <TopBar
        title="Baseline diagnosis"
        subtitle="Submit sample answers across coding, concepts, and behavioral — AI maps your starting point."
        actions={
          report && (
            <Badge tone="info">
              Readiness {formatPercent(report.overall_readiness)}
            </Badge>
          )
        }
      />

      {report && (
        <Card title="Latest diagnosis" className="mb-4">
          <p className="text-sm text-zinc-600 dark:text-zinc-300">
            {report.summary}
          </p>
          {report.weak_areas.length > 0 && (
            <div className="mt-3 flex flex-wrap gap-2">
              {report.weak_areas.map((area) => (
                <Badge key={area} tone="warning">
                  {formatCategory(area)}
                </Badge>
              ))}
            </div>
          )}
          <p className="mt-3 text-xs text-zinc-500">
            Recommended start: {formatCategory(report.recommended_start_level)}
          </p>
        </Card>
      )}

      <Card
        title={report ? "Re-run diagnosis" : "Submit baseline answers"}
        description="Be honest — this sets your recovery plan focus."
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <Textarea
            label="Coding answer"
            hint="Describe your approach to a medium DSA problem (e.g. two-sum with hash map)."
            placeholder={CODING_PLACEHOLDER}
            value={codingAnswer}
            onChange={(e) => setCodingAnswer(e.target.value)}
            required
          />
          <Textarea
            label="Conceptual answer"
            hint="Explain a system design or architecture concept briefly."
            placeholder={CONCEPTUAL_PLACEHOLDER}
            value={conceptualAnswer}
            onChange={(e) => setConceptualAnswer(e.target.value)}
            required
          />
          <Textarea
            label="Behavioral answer"
            hint="Use STAR format for a leadership or incident story."
            placeholder={BEHAVIORAL_PLACEHOLDER}
            value={behavioralAnswer}
            onChange={(e) => setBehavioralAnswer(e.target.value)}
            required
          />

          <Button type="submit" loading={submit.isPending}>
            {report ? "Update diagnosis" : "Run AI diagnosis"}
          </Button>
        </form>
      </Card>
    </>
  );
}
