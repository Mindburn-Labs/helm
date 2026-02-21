// ─── HELM CLI v3 Types ───────────────────────────────────────────────────────
// Stable JSON schema for all CLI inputs and outputs.

/** Conformance level shortcut. */
export type ConformanceLevel = "L1" | "L2";

/** CLI invocation options parsed from argv. */
export interface CLIOptions {
    bundle?: string;
    level?: ConformanceLevel;
    ci: boolean;
    json: boolean;
    depth: number;
    report?: string;
    noCache: boolean;
    cacheDir?: string;
    help: boolean;
    version: boolean;
}

// ─── Bundle Schema ───────────────────────────────────────────────────────────

/** §3.2 Index entry in 00_INDEX.json. */
export interface IndexEntry {
    path: string;
    sha256: string;
    size_bytes: number;
    schema_version?: string;
    content_type: string;
}

/** §3.1 Index manifest (00_INDEX.json). */
export interface IndexManifest {
    run_id: string;
    profile: string;
    created_at: string;
    topo_order_rule: string;
    entries: IndexEntry[];
}

/** §6.1 Gate result from 01_SCORE.json. */
export interface GateResult {
    gate_id: string;
    pass: boolean;
    reasons: string[];
    evidence_paths: string[];
    metrics: {
        duration_ms: number;
        counts?: Record<string, number>;
    };
    details?: Record<string, unknown>;
}

/** Conformance report from 01_SCORE.json. */
export interface ConformanceReport {
    run_id: string;
    profile: string;
    timestamp: string;
    pass: boolean;
    gate_results: GateResult[];
    duration: number | string;
    metadata?: Record<string, unknown>;
}

// ─── Attestation Schema ──────────────────────────────────────────────────────

/** v3 release attestation — signed with Ed25519. */
export interface Attestation {
    format: "helm-attestation-v3";
    release_tag: string;
    asset_name: string;
    asset_sha256: string;
    manifest_root_hash: string;
    merkle_root: string;
    created_at: string;
    profiles_version?: string;
}

// ─── Verification Result Types ───────────────────────────────────────────────

export interface StructureCheck {
    pass: boolean;
    dirCount: number;
    hasIndex: boolean;
    hasScore: boolean;
    missingDirs: string[];
    extraEntries: string[];
}

export interface HashChainCheck {
    pass: boolean;
    totalEntries: number;
    verifiedEntries: number;
    failedEntries: Array<{ path: string; expected: string; actual: string }>;
}

export interface SignatureCheck {
    pass: boolean;
    signerID?: string;
    signedAt?: string;
    reason?: string;
}

export interface GateCheck {
    pass: boolean;
    level: ConformanceLevel;
    totalGates: number;
    passedGates: number;
    failedGates: GateResult[];
    gateResults: GateResult[];
}

export interface AttestationCheck {
    pass: boolean;
    verified: boolean;
    reason?: string;
    attestation?: Attestation;
}

/** Complete verification result — stable JSON output schema. */
export interface VerificationResult {
    tool: "@mindburn/helm";
    artifact: string;
    verdict: "PASS" | "FAIL";
    profile: string;
    timing_ms: number;
    roots: {
        manifest_root_hash: string;
        merkle_root: string;
    };
    structure: StructureCheck;
    hash_chain: HashChainCheck;
    signature: SignatureCheck;
    gates: GateCheck;
    attestation: AttestationCheck;
    evidence: string;
}
