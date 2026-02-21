// Generate the minimal golden fixture for HELM OSS.
// Run: node scripts/generate-fixture.mjs
// Output: fixtures/minimal/ with deterministic 00_INDEX.json, 01_SCORE.json, and EXPECTED.json

import { createHash } from "node:crypto";
import { writeFileSync } from "node:fs";
import { join } from "node:path";

const FIXTURE_DIR = join(import.meta.dirname, "..", "fixtures", "minimal");

// ─── Helper ──────────────────────────────────────────────────────────────────

function sha256Hex(buf) {
    return createHash("sha256").update(buf).digest("hex");
}

function sha256Raw(buf) {
    return createHash("sha256").update(buf).digest();
}

function canonicalJSON(obj) {
    return JSON.stringify(sortKeys(obj));
}

function sortKeys(value) {
    if (value === null || typeof value !== "object") return value;
    if (Array.isArray(value)) return value.map(sortKeys);
    const sorted = {};
    for (const key of Object.keys(value).sort()) {
        sorted[key] = sortKeys(value[key]);
    }
    return sorted;
}

// ─── 1. Create a minimal log file as evidence ────────────────────────────────

const logContent = `{"event":"verify","ts":"2026-02-21T00:00:00Z","status":"pass"}\n`;
const logPath = "06_LOGS/verify.jsonl";
writeFileSync(join(FIXTURE_DIR, logPath), logContent);
const logHash = sha256Hex(Buffer.from(logContent));

// ─── 2. Create 01_SCORE.json ─────────────────────────────────────────────────

const score = {
    run_id: "fixture-minimal-001",
    profile: "CORE",
    timestamp: "2026-02-21T00:00:00Z",
    pass: true,
    gate_results: [
        { gate_id: "G0", status: "pass", pass: true, reasons: [], evidence_paths: [], metrics: { duration_ms: 1 } },
        { gate_id: "G1", status: "pass", pass: true, reasons: [], evidence_paths: [], metrics: { duration_ms: 1 } },
        { gate_id: "G2", status: "pass", pass: true, reasons: [], evidence_paths: [], metrics: { duration_ms: 1 } },
        { gate_id: "G2A", status: "pass", pass: true, reasons: [], evidence_paths: [], metrics: { duration_ms: 1 } },
        { gate_id: "G3A", status: "pass", pass: true, reasons: [], evidence_paths: [], metrics: { duration_ms: 1 } },
        { gate_id: "G5", status: "pass", pass: true, reasons: [], evidence_paths: [], metrics: { duration_ms: 1 } },
        { gate_id: "G8", status: "pass", pass: true, reasons: [], evidence_paths: [], metrics: { duration_ms: 1 } },
        { gate_id: "GX_ENVELOPE", status: "pass", pass: true, reasons: [], evidence_paths: [], metrics: { duration_ms: 1 } },
        { gate_id: "GX_TENANT", status: "pass", pass: true, reasons: [], evidence_paths: [], metrics: { duration_ms: 1 } },
    ],
    duration: 9,
};

const scoreJson = canonicalJSON(score);
writeFileSync(join(FIXTURE_DIR, "01_SCORE.json"), scoreJson);
const scoreHash = sha256Hex(Buffer.from(scoreJson));

// ─── 3. Create 00_INDEX.json ─────────────────────────────────────────────────

const index = {
    format_version: "3.0.0",
    run_id: "fixture-minimal-001",
    profile: "CORE",
    created_at: "2026-02-21T00:00:00Z",
    topo_order_rule: "creation_timestamp",
    entries: [
        {
            path: "01_SCORE.json",
            sha256: scoreHash,
            size_bytes: Buffer.byteLength(scoreJson),
            kind: "helm:report",
        },
        {
            path: logPath,
            sha256: logHash,
            size_bytes: Buffer.byteLength(logContent),
            kind: "helm:log",
        },
    ],
};

const indexJson = canonicalJSON(index);
writeFileSync(join(FIXTURE_DIR, "00_INDEX.json"), indexJson);

// ─── 4. Compute expected roots ───────────────────────────────────────────────

const bundleRoot = sha256Hex(Buffer.from(indexJson));

// Merkle tree: leaves sorted by path ascending
const sortedEntries = [...index.entries].sort((a, b) => a.path.localeCompare(b.path));
const LEAF_PREFIX = Buffer.from([0x00]);
const NODE_PREFIX = Buffer.from([0x01]);

let level = sortedEntries.map(e => {
    const hashBytes = Buffer.from(e.sha256, "hex");
    return sha256Raw(Buffer.concat([LEAF_PREFIX, hashBytes]));
});

while (level.length > 1) {
    const next = [];
    for (let i = 0; i < level.length; i += 2) {
        const left = level[i];
        const right = i + 1 < level.length ? level[i + 1] : level[i];
        next.push(sha256Raw(Buffer.concat([NODE_PREFIX, left, right])));
    }
    level = next;
}

const merkleRoot = level[0].toString("hex");

// ─── 5. Write EXPECTED.json ──────────────────────────────────────────────────

const expected = {
    bundle_root: bundleRoot,
    merkle_root: merkleRoot,
    expected_verdict: "PASS",
    expected_gates: score.gate_results.map(g => g.gate_id),
    description: "Expected verification roots and outcomes for the minimal golden fixture.",
};

writeFileSync(
    join(FIXTURE_DIR, "EXPECTED.json"),
    JSON.stringify(expected, null, 2) + "\n",
);

// Remove old expected-roots.json if present
import { unlinkSync } from "node:fs";
try { unlinkSync(join(FIXTURE_DIR, "expected-roots.json")); } catch { }

console.log("Golden fixture generated:");
console.log(`  bundle_root: ${bundleRoot}`);
console.log(`  merkle_root: ${merkleRoot}`);
console.log(`  entries:     ${index.entries.length}`);
console.log(`  path:        fixtures/minimal/`);
