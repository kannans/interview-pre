# Phase 12 — 30-Day Personal Validation Playbook

**Purpose:** Decide honestly whether Interview Preparation Coach helps *you* recover interview performance before building Phases 13–16.

**Primary user:** ~15 YOE engineer — strong at work, weak under interview pressure.

**First success:** Personal improvement and a sustainable daily habit — not revenue, not feature count.

**Motto:** Turn interview rust into measurable recovery before your next real round.

---

## How to use this doc

| File | Use |
|------|-----|
| [VALIDATION.md](./VALIDATION.md) | This playbook — read once, reference weekly |
| [validation/baseline-metrics.md](./validation/baseline-metrics.md) | Day 0 snapshot — fill once |
| [validation/session-logs.md](./validation/session-logs.md) | Append after every practice session |
| [validation/friction-log.md](./validation/friction-log.md) | Log blockers and annoyances as they happen |
| [validation/weekly-reviews.md](./validation/weekly-reviews.md) | Sunday reviews + Day 30 decision |
| [REVIEW.md](./REVIEW.md) | Day 0 manual QA checklist |

**Rule:** Validation failure is valuable data. Do not rationalize a GO if the numbers and honest answers say otherwise.

---

## 1. Pre-validation setup (Day 0)

Complete before your first counted practice day.

### 1.1 Environment ready

Choose one stack and stick to it for 30 days (switching mid-validation confounds friction data).

```bash
# Option A — local dev
./runserver.sh

# Option B — Docker (recommended for daily dogfood)
cp .env.example .env   # set OPENAI_API_KEY
docker compose up --build -d
curl http://localhost:8080/health
```

Open http://localhost:3000. Confirm API health and OpenAI key work.

### 1.2 Run full manual QA

Work through every row in [REVIEW.md](./REVIEW.md#manual-test-checklist). Mark pass/fail in that doc or copy the table into your notes.

**Minimum bar to start Day 1:**

- [ ] Profile → diagnosis → plan → practice → evaluate → progress completes without manual API calls
- [ ] Today's task links to a recommended question
- [ ] Evaluation scores and feedback render
- [ ] Progress snapshot updates after evaluate
- [ ] Plan exhaustion (404 on `/plan/today`) routes to progress/adjust — not a dead end

**If QA fails:** Log P0 items in [friction-log.md](./validation/friction-log.md). Fix only what blocks the loop. Do not start Phase 13 features.

### 1.3 Record baseline metrics

Fill [baseline-metrics.md](./validation/baseline-metrics.md):

| Source | What to capture |
|--------|-----------------|
| `/diagnosis` report | `overall_readiness`, all five `dimension_scores` (0–10), `priority_focus` |
| `/progress` latest snapshot | `readiness_index`, `weak_patterns` (may be empty on Day 0) |
| Self | **Self-rated confidence (1–10):** "If I had a technical screen tomorrow, how ready do I feel?" |

If you have not practiced yet, run one session on Day 0 so you have a post-evaluate `readiness_index`. Otherwise record diagnosis only and note "readiness_index pending first session."

### 1.4 Set validation calendar

- **Day 0:** Today — QA + baseline
- **Days 1–30:** Daily practice (see §2)
- **Sundays:** Weekly review (see §3)
- **Day 30:** Decision gate (see §7)

### 1.5 Friction log ready

[friction-log.md](./validation/friction-log.md) is your running backlog. Log issues the moment they cost you time or trust — not at end of week when memory fades.

---

## 2. Daily usage routine (Mon–Sun)

### 2.1 Standard day (30–60 min) — target 5–6 days/week

**Morning block (best before work distractions):**

1. **Open dashboard** (`/`) — note workflow stage card and today's task (~30 sec)
2. **Start today's task** — click through to recommended practice (`/practice/[id]`)
3. **Time yourself:** record "open app → cursor in editor" (target **< 2 min**)
4. **Practice one question** — treat it like a real round:
   - Think out loud (even if typing)
   - State approach before diving into code
   - Call out edge cases explicitly
5. **Submit → evaluate** — read all four feedback dimensions
6. **Log session** in [session-logs.md](./validation/session-logs.md) (5 min)
7. **Optional:** skim `/progress` if snapshot refreshed

**Do not** during validation: build new features, tweak prompts, or expand question bank unless P0 loop breakage.

### 2.2 When `GET /plan/today` returns 404 (plan exhausted)

This is expected behavior, not a bug. The dashboard workflow stage moves to **progress** or **adjust**.

**Do this:**

1. Go to **Progress** (`/progress`) — read latest snapshot, weak patterns, adjustment recommendation
2. Decide:
   - **Regenerate plan** (`POST /plan` via UI on `/plan`) if you want a fresh task list aligned to current gaps
   - **Free practice** — pick any question from `/practice` if you still want reps before regenerating
3. Log in session notes: "plan exhausted — regenerated / free practice only"
4. If exhaustion happened before Day 14, note in friction log — plan may be too short for your pace

**Do not** treat 404 as "app is broken." Treat it as a natural loop transition.

### 2.3 Minimum viable day (15 min) — busy days

Goal: **never zero.** Streak matters more than duration for habit validation.

1. Open dashboard → today's task or last recommended question
2. One question — shorter time box (10 min work, 5 min log)
3. Still log: confidence, one wrong thing, one win
4. Count as a practice day ✓

**Busy-week floor:** ≥ 4 practice days/week even if some are 15-min MVPs.

### 2.4 What not to count as practice

- Reading diagnosis or progress without submitting an attempt
- Fixing deploy issues (log to friction, don't count as practice day)
- Re-reading old evaluations without a new attempt

### 2.5 Daily checklist (copy mentally or pin)

```
[ ] Opened app
[ ] Completed ≥1 evaluated attempt OR logged why not
[ ] Session log appended
[ ] Friction logged if anything blocked >2 min
```

---

## 3. Weekly review checklist (4 weeks)

Complete each review on **Sunday** in [weekly-reviews.md](./validation/weekly-reviews.md).

### Week 1 — Loop reliability and friction

**Question:** Can I complete the loop daily without fighting the app?

| Check | Pass criteria |
|-------|---------------|
| Practice days | ≥ 4/7 |
| Loop completion rate | ≥ 80% of attempts reach evaluate |
| Open → editor | Median < 2 min |
| P0 friction items | 0 unresolved |
| Plan task completion | At least 3 tasks done |
| 404 handling | You successfully recovered from exhaustion once (or noted N/A) |

**Review actions:**

- Triage [friction-log.md](./validation/friction-log.md) — fix P0 only
- Note top 3 friction themes (UX, AI latency, content, deploy)
- **Week 1 fail signal:** Loop breaks more than 2×/week or you skip practice because the app feels broken

---

### Week 2 — AI feedback quality and plan relevance

**Question:** Does the coach give feedback I'd act on, and does the plan match my gaps?

| Check | Pass criteria |
|-------|---------------|
| Practice days | ≥ 4/7 |
| Feedback usefulness | ≥ 60% sessions logged "would use in real interview" = yes or maybe |
| Specificity | Feedback cites *your* answer, not generic tips |
| Plan relevance | Tasks map to diagnosis `priority_focus` (subjective 1–5, target ≥ 3) |
| Behavior change | At least 1 session where you explicitly applied prior feedback |

**Review actions:**

- Pick best and worst evaluation of the week — paste one line from each into weekly review
- Rate average AI feedback quality 1–5
- Log weak feedback examples in friction log (area: AI feedback)
- **Week 2 fail signal:** Feedback feels like ChatGPT boilerplate; plan tasks feel random

---

### Week 3 — Progress measurability

**Question:** Can I *see* improvement, or is this just activity?

| Check | Pass criteria |
|-------|---------------|
| Practice days | ≥ 4/7 |
| readiness_index | Up ≥ 3 pts from Week 1 end OR flat with explanation (harder questions) |
| Evaluation dimensions | At least 2 of 4 dimensions trending up vs Week 1 average |
| Snapshot history | ≥ 3 snapshots; trend interpretable |
| weak_patterns | Match your self-known gaps |
| Confidence | Self-rated confidence up ≥ 1 pt vs baseline |

**Review actions:**

- Fill Week 1 vs Week 3 table in [weekly-reviews.md](./validation/weekly-reviews.md)
- Answer honestly: "Am I improving?" — require evidence (scores, confidence, replay of old attempt)
- **Week 3 fail signal:** No metric moved; you can't tell if you're better than Day 0

---

### Week 4 — Real interview readiness

**Question:** Does training transfer under pressure?

| Check | Pass criteria |
|-------|---------------|
| Practice days | ≥ 4/7 |
| Pressure exposure | ≥ 1 timed mock (peer, self, or real screen) |
| Transfer | App feedback themes (edge cases, communication, approach) surfaced in mock |
| Confidence | Self-rated ≥ baseline + 1 OR honest explanation why not |
| Recommendation | "Would recommend to senior engineer" = yes or maybe |

**Mock options (pick one):**

- **Peer mock** — 45 min with a friend; use a seeded question you haven't done
- **Self mock** — 35 min timer, camera on, no AI until after; then submit recording as attempt
- **Real round** — if scheduled, debrief in weekly review (best data)

**Review actions:**

- Complete Day 30 qualitative summary in weekly-reviews.md
- Run decision gate (§7)
- **Week 4 fail signal:** Mock felt like Day 0; app didn't prepare you for pressure-specific failure

---

## 4. Session log template

Append to [session-logs.md](./validation/session-logs.md) after every session.

```markdown
### YYYY-MM-DD — [question title or id]

- Category: programming_basics | dsa | system_design | behavioral | other
- Time taken: ___ min
- Open app → editor: ___ min
- Scores: correctness __ | approach __ | edge_cases __ | communication __
- Self-rated confidence before submit: /10
- One thing that went wrong:
- One win:
- Would use this feedback in a real interview? yes / no / maybe
- Notes:
```

**Confidence scale:**

| Score | Meaning |
|-------|---------|
| 1–3 | Guessing; would freeze in a real round |
| 4–6 | Shaky; could muddle through |
| 7–8 | Solid under practice conditions |
| 9–10 | Would volunteer this answer in a staff loop |

Track confidence **before** submit — not after seeing scores (calibration data for future Phase 13).

---

## 5. Success metrics (quantitative)

Track weekly; final comparison is **Week 1 end vs Week 4 end**.

| Metric | Target | How to measure |
|--------|--------|----------------|
| Practice days/week | **≥ 4** | Count days with ≥1 evaluated attempt in session logs |
| Avg evaluation scores by dimension | **Upward trend** | Mean of correctness, approach, edge_cases, communication per week |
| readiness_index delta | **≥ +5** (Week 1 → Week 4) | `/progress` latest snapshot vs Week 1 Sunday reading |
| Plan tasks completed | **≥ 80%** of generated tasks | Plan page or dashboard task count |
| Time open app → start practice | **< 2 min median** | Session log field |
| Loop completion rate | **≥ 90%** | Attempts submitted / sessions started |
| Sessions/week | **≥ 4** | Session log count |

### Score aggregation cheat sheet

**Week N averages** (from session logs):

```
avg_correctness_N = mean of all correctness scores that week
(same for approach, edge_cases, communication)
```

**readiness_index:** Use snapshot at end of Week 1 Sunday and Week 4 Sunday. If you practice daily, same-day upsert means latest snapshot reflects that day.

### Quantitative GO bar (all should be true at Day 30)

- ≥ 16 practice days total (avg 4/week)
- readiness_index Week 4 ≥ Week 1 + 5 **OR** ≥ 2 evaluation dimensions up ≥ 0.5 pts with stable readiness
- Median open → editor < 2 min
- ≥ 70% plan tasks completed across the validation period

---

## 6. Success metrics (qualitative)

Answer on Day 30 with evidence — not hope.

| Question | Pass | Evidence required |
|----------|------|-------------------|
| **Am I improving?** | yes | Score trends, confidence delta, or replay comparison |
| **Does the coach tell me what to do next?** | yes | Dashboard task + plan drove practice without you deciding what to study |
| **Would I recommend this to another senior engineer?** | yes or maybe | You'd send them the link with honest caveats |

**Qualitative GO bar:** First two = **yes**. Third = **yes** or **maybe** (not **no**).

If "Am I improving?" is **no** with honest evidence, outcome is **PIVOT** or **STOP** regardless of habit.

---

## 7. Decision gates at Day 30

Be strict. The point is to avoid building Phase 13+ on wishful thinking.

### GO → Proceed to Phase 13

**All true:**

- Quantitative GO bar met (§5)
- Qualitative GO bar met (§6)
- Habit sustainable — you *want* to open the app tomorrow
- Friction log has no unresolved P0 items

**Action:** Update [NEXT.md](./NEXT.md) to Phase 13. Carry friction P1/P2 items into Phase 13 scope only if they block measurability.

---

### PIVOT → Fix core loop / AI before Phase 13

**Pattern:**

- Loop works; you practiced ≥ 4 days/week
- Feedback or plan quality failed Week 2 bar
- OR readiness/scores flat despite consistent effort (Week 3 fail)
- You still believe the *concept* helps

**Action:**

1. Document top 3 failures in weekly-reviews.md
2. Fix prompts, evaluation rubric, plan generation, or UX — **not** failure patterns / velocity / pressure mode yet
3. Re-run a **14-day mini-validation** before Phase 13
4. Log "idea" items in friction log for later — do not build them now

---

### STOP → Pause feature build

**Pattern (any):**

- < 3 practice days/week sustained — habit didn't form
- "Would recommend" = **no** and "Am I improving?" = **no**
- You avoid the app or practice without it instead
- P0 friction never resolved; loop still feels broken

**Action:**

1. Write 1-page post-mortem in weekly-reviews.md: why it didn't work (product, life context, wrong user, wrong problem)
2. Do **not** start Phase 13 until STOP reasons are addressed or scope is redefined
3. STOP is success — you saved weeks of building the wrong thing

---

### Decision matrix (quick reference)

| Habit (≥4 d/wk) | Improving? | Feedback useful? | Outcome |
|-----------------|------------|------------------|---------|
| yes | yes | yes | **GO** |
| yes | no | yes | **PIVOT** (content/plan/scoring) |
| yes | no | no | **PIVOT** (AI/prompts) |
| yes | yes | no | **PIVOT** (feedback quality) |
| no | * | * | **STOP** (habit first) |

---

## 8. Friction and improvement backlog template

Live doc: [validation/friction-log.md](./validation/friction-log.md)

| Date | Area | Severity | keep/fix/cut/idea | Notes |
|------|------|----------|-------------------|-------|
| 2026-07-01 | practice | P1 | fix | Draft lost on refresh |
| 2026-07-05 | AI feedback | P2 | fix | Edge case notes too generic |
| 2026-07-10 | progress | P3 | keep | Chart would help — **idea** for Phase 13 |

**Severity:**

- **P0** — Blocks loop (can't evaluate, 502 loop, data loss)
- **P1** — Workaround exists but painful; kills daily habit
- **P2** — Annoyance; fix after validation if GO
- **P3** — Polish; log as **idea**, don't fix during validation

**Disposition:**

- **keep** — Accept for MVP; not worth fixing now
- **fix** — Must fix before GO or during PIVOT
- **cut** — Remove/simplify scope
- **idea** — Phase 13+ candidate; **do not implement during validation**

---

## 9. File structure summary

```
docs/progress/
├── VALIDATION.md              ← this playbook
├── REVIEW.md                  ← Day 0 manual QA
├── STATUS.md                  ← phase checklist
├── NEXT.md                    ← post-validation next phase
└── validation/
    ├── baseline-metrics.md    ← Day 0 snapshot (fill once)
    ├── session-logs.md        ← append per session
    ├── friction-log.md        ← running backlog
    └── weekly-reviews.md      ← 4 weekly + Day 30 decision
```

---

## Appendix A — Metrics reference (app fields)

**Diagnosis (`diagnostic_report`):**

- `overall_readiness` (0–100)
- `dimension_scores`: programming_basics, dsa, system_design_concepts, communication, behavioral (0–10 each)

**Evaluation (per attempt):**

- `score_correctness`, `score_approach`, `score_edge_cases`, `score_communication`

**Progress snapshot (`stats`):**

- `readiness_index` (0–100) — computed from evaluations, plan completion, attempt volume
- `weak_patterns` — lowest average dimensions
- `adjustment_recommendation` — shown when plan exhausted

**Not in MVP (track manually in session logs):**

- `self_rated_confidence` — no API field yet; Phase 13 may add it

---

## Appendix B — Ideas that may emerge (do not build during validation)

If friction log disposition = **idea**, examples:

- Recovery velocity chart
- Failure pattern memory across attempts
- Confidence slider on submit
- Pressure/timed mode
- Mock interview backend
- More questions in weak categories

Log only. Phase 13+ scope waits for GO gate.

---

## Appendix C — Quick start (Day 0 today)

1. `docker compose up -d` or `./runserver.sh`
2. Run [REVIEW.md](./REVIEW.md) checklist
3. Fill [baseline-metrics.md](./validation/baseline-metrics.md)
4. One practice session + session log
5. Calendar Day 1–30
6. Tomorrow: §2 standard day

**After Day 30:** Decision in [weekly-reviews.md](./validation/weekly-reviews.md) → update [NEXT.md](./NEXT.md).
