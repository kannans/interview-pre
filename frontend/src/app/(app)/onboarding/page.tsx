"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { TopBar } from "@/components/layout/TopBar";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { Textarea } from "@/components/ui/Textarea";
import { Button } from "@/components/ui/Button";
import { ErrorAlert } from "@/components/ui/ErrorAlert";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { useProfile, useUpsertProfile } from "@/hooks/useProfile";
import { isApiError } from "@/lib/api";
import type { ProfileInput } from "@/lib/types/api";

const EMPTY_FORM: ProfileInput = {
  experience_years: 5,
  target_role: "",
  target_companies: "",
  weak_areas_self: "",
  timeline: "60 days",
};

function profileToInput(profile: {
  experience_years: number;
  target_role: string;
  target_companies: string;
  weak_areas_self: string;
  timeline: string;
}): ProfileInput {
  return {
    experience_years: profile.experience_years,
    target_role: profile.target_role,
    target_companies: profile.target_companies,
    weak_areas_self: profile.weak_areas_self,
    timeline: profile.timeline,
  };
}

export default function OnboardingPage() {
  const router = useRouter();
  const profile = useProfile();
  const upsert = useUpsertProfile();
  const [draft, setDraft] = useState<ProfileInput | null>(null);

  const values =
    draft ?? (profile.data ? profileToInput(profile.data) : EMPTY_FORM);

  function updateField<K extends keyof ProfileInput>(
    key: K,
    value: ProfileInput[K],
  ) {
    setDraft((prev) => ({ ...(prev ?? values), [key]: value }));
  }

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    upsert.mutate(values, {
      onSuccess: () => router.push("/diagnosis"),
    });
  }

  if (profile.isLoading && !profile.isError) {
    return <LoadingSpinner label="Loading profile…" />;
  }

  return (
    <>
      <TopBar
        title="Profile setup"
        subtitle="Capture your experience, targets, and timeline so the coach can personalize your recovery."
      />

      {profile.isError && !isApiError(profile.error) && (
        <ErrorAlert message={profile.error.message} className="mb-4" />
      )}

      <Card title="Preparation profile">
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Years of experience"
            type="number"
            min={0}
            value={values.experience_years}
            onChange={(e) =>
              updateField("experience_years", Number(e.target.value))
            }
            required
          />
          <Input
            label="Target role"
            placeholder="Senior Backend Engineer"
            value={values.target_role}
            onChange={(e) => updateField("target_role", e.target.value)}
            required
          />
          <Input
            label="Target companies"
            placeholder="Google, Stripe, …"
            value={values.target_companies}
            onChange={(e) => updateField("target_companies", e.target.value)}
            required
          />
          <Textarea
            label="Self-identified weak areas"
            placeholder="DSA, system design, behavioral stories…"
            value={values.weak_areas_self}
            onChange={(e) => updateField("weak_areas_self", e.target.value)}
            required
          />
          <Input
            label="Interview timeline"
            placeholder="60 days"
            value={values.timeline}
            onChange={(e) => updateField("timeline", e.target.value)}
            required
          />

          <div className="flex gap-2 pt-2">
            <Button type="submit" loading={upsert.isPending}>
              {profile.data ? "Update profile" : "Save & continue"}
            </Button>
            {profile.data && (
              <Button type="button" variant="secondary" href="/">
                Back to home
              </Button>
            )}
          </div>
        </form>
      </Card>
    </>
  );
}
