# HELM SDK — Registry Publishing Setup

Step-by-step instructions for configuring trusted publishing across all SDK registries.

> [!IMPORTANT]
> All publish jobs use GitHub Environments with required approvals. Create the environments before configuring registry connections.

---

## Prerequisites

1. **GitHub Environments** — create protected environments in GitHub repo settings:
   - `pypi-publish` — required reviewers, deployment branch `main`
   - `npm-publish` — required reviewers, deployment branch `main`
   - `crates-publish` — required reviewers, deployment branch `main`
   - `maven-publish` — required reviewers, deployment branch `main`

---

## PyPI — Trusted Publishing (OIDC)

No password or token stored in GitHub. GitHub Actions authenticates directly with PyPI via OIDC.

### Setup Steps

1. Go to https://pypi.org → Your Projects → `helm-sdk` → Settings → Publishing
2. Click "Add a new pending publisher" (first time) or "Add publisher"
3. Fill in:
   - **Owner**: `Mindburn-Labs`
   - **Repository**: `helm`
   - **Workflow name**: `release.yml`
   - **Environment name**: `pypi-publish`
4. Save

### Workflow Configuration

The release workflow uses `pypa/gh-action-pypi-publish@release/v1` with `id-token: write` permission. No secrets needed — OIDC handles authentication.

---

## npm — Trusted Publishing (OIDC + Provenance)

npm supports provenance attestations via OIDC.

### Setup Steps

1. Go to https://www.npmjs.com → Package `@mindburn/helm-sdk` → Settings → Publishing access
2. Under "Configure trusted publishing":
   - **Repository**: `Mindburn-Labs/helm`
   - **Workflow**: `release.yml`
   - **Environment**: `npm-publish`
3. Save

### Workflow Configuration

The release workflow runs `npm publish --provenance --access public`. The `--provenance` flag generates OIDC-signed provenance. Requires `id-token: write` and `packages: write` permissions.

---

## crates.io — OIDC Trusted Publishing

### Setup Steps

1. Go to https://crates.io → `helm-sdk` → Settings → Trusted Publishing
2. Add a new configuration:
   - **Owner**: `Mindburn-Labs`
   - **Repository**: `helm`
   - **Workflow**: `release.yml`
   - **Environment**: `crates-publish`
3. Save

### Workflow Configuration

The release workflow runs `cargo publish` with OIDC token exchange. Requires `id-token: write` permission.

---

## Maven Central — Central Portal Token Flow

Maven Central does not support OIDC. Token-based authentication is required.

### Setup Steps

1. Go to https://central.sonatype.com → Log in → User Settings
2. Generate a publishing token (username + password pair)
3. In GitHub repo → Settings → Environments → `maven-publish`:
   - Add secret `MAVEN_USERNAME` (token username)
   - Add secret `MAVEN_PASSWORD` (token password)
4. Ensure your `pom.xml` has:
   ```xml
   <distributionManagement>
     <repository>
       <id>central</id>
       <url>https://central.sonatype.com/api/v1/publisher/deployments/upload</url>
     </repository>
   </distributionManagement>
   ```

### Additional POM Requirements for Central

```xml
<licenses>
  <license>
    <name>BSL-1.1</name>
    <url>https://github.com/Mindburn-Labs/helm/blob/main/LICENSE</url>
  </license>
</licenses>
<developers>
  <developer>
    <name>Mindburn Labs</name>
    <email>oss@mindburn.org</email>
    <organization>Mindburn Labs</organization>
  </developer>
</developers>
<scm>
  <connection>scm:git:git://github.com/Mindburn-Labs/helm.git</connection>
  <developerConnection>scm:git:ssh://github.com/Mindburn-Labs/helm.git</developerConnection>
  <url>https://github.com/Mindburn-Labs/helm</url>
</scm>
```

### Workflow Configuration

The release workflow runs `mvn deploy -DskipTests` with `MAVEN_USERNAME` and `MAVEN_PASSWORD` injected from the environment secrets.

---

## Go Module — Tag-Based Publishing

Go modules are published by tagging a commit. No registry upload step needed.

### Setup Steps

1. Tag the release: `git tag sdk/go/v0.1.0`
2. Push: `git push origin sdk/go/v0.1.0`
3. The Go module proxy (proxy.golang.org) automatically picks up the tag

No secrets or environment setup required.

---

## GitHub Environments — Protection Rules

For each environment above, configure:

| Environment | Required Reviewers | Deployment Branch | Wait Timer |
|-------------|-------------------|-------------------|------------|
| `pypi-publish` | 1+ maintainer | `main` only | 0 min |
| `npm-publish` | 1+ maintainer | `main` only | 0 min |
| `crates-publish` | 1+ maintainer | `main` only | 0 min |
| `maven-publish` | 1+ maintainer | `main` only | 0 min |

This prevents accidental publishes from feature branches or unauthorized actors.
