"use client";

import { TopBar } from "@/components/layout/TopBar";
import { Card } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { EmptyState } from "@/components/ui/EmptyState";

const MOCK_FLOW = [
  {
    step: 1,
    title: "Configure session",
    description: "Choose focus areas and difficulty (backend AI endpoints coming in Phase 9).",
  },
  {
    step: 2,
    title: "Answer questions",
    description: "Timed mixed coding, system design, and behavioral prompts.",
  },
  {
    step: 3,
    title: "Review grading",
    description: "AI scores strengths, gaps, and next steps.",
  },
];

export default function MockInterviewPage() {
  return (
    <>
      <TopBar
        title="Mock interview"
        subtitle="UI shell for timed interview practice — AI generation and grading stay on the backend."
        actions={<Badge tone="warning">Backend stubs only</Badge>}
      />

      <Card className="mb-4" title="Coming soon">
        <p className="text-sm text-zinc-600 dark:text-zinc-300">
          Mock interview endpoints are defined in the backend AI layer but not
          yet exposed via REST. This page previews the flow you will use once
          Phase 9 wires the full loop.
        </p>
      </Card>

      <ol className="space-y-3">
        {MOCK_FLOW.map((item) => (
          <li key={item.step}>
            <Card>
              <div className="flex gap-4">
                <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-teal-100 text-sm font-semibold text-teal-800 dark:bg-teal-900/40 dark:text-teal-300">
                  {item.step}
                </span>
                <div>
                  <h3 className="font-medium text-zinc-900 dark:text-zinc-50">
                    {item.title}
                  </h3>
                  <p className="mt-1 text-sm text-zinc-500">{item.description}</p>
                </div>
              </div>
            </Card>
          </li>
        ))}
      </ol>

      <EmptyState
        title="Mock session not available yet"
        description="Use daily practice and baseline diagnosis while mock interview APIs are being connected."
        action={
          <Button href="/practice" variant="secondary">
            Go to practice
          </Button>
        }
      />
    </>
  );
}
