# Phased Modernization and Security Review Plan for Trident

This document outlines a phased plan to modernize the Trident codebase, transition it to modern Go standards, establish a robust testing harness with high code coverage, and conduct a thorough security review.

---

## Phase 1: Build System & Dependency Modernization (Immediate)

**Objective**: Transition the codebase fully to modern Go Modules, eliminating legacy `GOPATH` requirements and ensuring the project can be built and tested using standard Go toolchain commands.

### 1.1. Go Module Transition & Clean Up
*   **Issue**: The current build system relies on a legacy `GOPATH` hack (`ext/_gopath`) and clones dependencies manually via Git in the `Makefile`.
*   **Tasks**:
    *   Permanently adopt Go Modules (`go.mod`).
    *   Update the `go.mod` file to use a modern Go version (e.g., `go 1.21` or higher) to leverage modern language features and security updates.
    *   Configure `replace` directives in `go.mod` to point to local checkouts of internal dependencies (`trident.li/pitchfork`, `trident.li/keyval`, `trident.li/go`) during transition, or migrate them to standard GitHub import paths (e.g., `github.com/tridentli/...`) if the vanity domain `trident.li` remains unreliable.
    *   Ensure all local dependency repositories (`pitchfork`, `go`, `keyval`) have valid, tidied `go.mod` files.
    *   Clean up the `Makefile` to remove `GOPATH` exports and replace legacy commands with standard `go build`, `go test`, and `go vet` commands.

### 1.2. Code Quality & Static Analysis (Style Fixes)
*   **Issue**: `go vet` currently flags multiple warnings regarding unkeyed struct literals in menu definitions (e.g., `PfMEntry` and `PfUIMentry`).
*   **Tasks**:
    *   Refactor all struct initializations to use keyed fields. This prevents silent breakage if the underlying struct definitions in `pitchfork` or other dependencies change in the future.
    *   Example Refactoring:
        ```diff
        - {"reset", user_pw_reset, 1, 2, []string{"username", "nominator"}, pf.PERM_USER, "Send a recovery password..."},
        + {Cmd: "reset", Fun: user_pw_reset, Args_min: 1, Args_max: 2, Args: []string{"username", "nominator"}, Perms: pf.PERM_USER, Desc: "Send a recovery password..."},
        ```
    *   Integrate `go vet` and `go fmt` checks as blocking steps in the local development loop.

---

## Phase 2: Test Harness & Unit Testing Strategy

**Objective**: Introduce a comprehensive test suite targeting near 100% code coverage for core business logic in `src/lib/`, and establish an automated test execution environment.

### 2.1. Establish the Test Harness
*   **Database Dependency**: Trident relies heavily on PostgreSQL. Running unit tests requires a strategy to handle database operations without polluting production or requiring a complex local DB setup for every developer.
*   **Proposed Solutions**:
    1.  **Unit Testing with Mocking**: Use `github.com/DATA-DOG/go-sqlmock` to mock database responses for unit tests in `src/lib/` (e.g., testing `user.go`, `vouch.go` logic without a real DB). This is fast and suitable for pure unit tests.
    2.  **Integration Testing with Testcontainers**: Use `github.com/testcontainers/testcontainers-go` to automatically spin up a temporary PostgreSQL Docker container during test execution. This allows running integration tests against a real database with the actual schema applied, ensuring query compatibility.
    *   **Decision**: We will implement **both**: `go-sqlmock` for rapid unit tests of business logic, and `testcontainers` for database-centric integration tests.

### 2.2. Unit Testing Targets (`src/lib/`)
*   **User Management (`src/lib/user.go`)**:
    *   Test `IsNominator` and `BestNominator` with mocked DB states.
    *   Test `user_pw_send` (split password recovery): Verify that the recovery token is correctly generated, split into two unique halves, hashed/stored correctly, and that the correct emails are dispatched to the user and nominator.
    *   Test `user_merge` transactional integrity.
*   **Vouching Registry (`src/lib/vouch.go`)**:
    *   Test vouch validation rules: Ensure a user cannot vouch for themselves, cannot vouch without meeting group attestations, and that positive/negative vouches transition member states correctly.
*   **Group Attestations (`src/lib/group_attestation.go`)**:
    *   Test parsing and validation of custom group attestations.

### 2.3. UI Testing (`src/ui/`)
*   *   Use `net/http/httptest` to mock HTTP requests and responses.
    *   Verify that UI handlers correctly parse form inputs, enforce session permissions, and call the underlying library functions.
    *   Verify that HTML outputs are correctly sanitized using `Blue Monday`.

### 2.4. CI/CD Integration (GitHub Actions)
*   Define a GitHub Actions workflow (`.github/workflows/test.yml`) that:
    1.  Runs on every push and Pull Request.
    2.  Sets up a modern Go environment.
    3.  Starts a PostgreSQL service container (or relies on Testcontainers).
    4.  Runs `go vet ./...` and `go fmt` checks.
    5.  Runs all tests with race detection enabled: `go test -race -coverprofile=coverage.out ./...`.
    6.  Reports code coverage metrics.

---

## Phase 3: Code Refactoring & Idiomatic Go Adoption

**Objective**: Clean up legacy code patterns, improve safety, and adopt modern Go idioms.

### 3.1. Language Feature Upgrades (Go 1.21+)
*   **Structured Logging**: Replace custom verbose logging/debugging with the standard `log/slog` package for structured, levels-based logging.
*   **Slices & Maps**: Replace manual slice manipulation loops with functions from the standard `slices` and `maps` packages.
*   **Context Propagation**: Ensure `context.Context` is passed through all database and network boundary calls (e.g., propagating context from HTTP requests down to `pgx` database queries) to support cancellation and timeouts.

### 3.2. Refactoring Package Names and Cycle Resolution
*   Resolve the duplicate `crypt` package naming (in `osutil-crypt` and its `common` subdirectory) by renaming the internal packages to be distinct (e.g., `common` to `cryptcommon` or merging them) to prevent import shadowing and confusion.
*   Ensure explicit error wrapping using `fmt.Errorf("...: %w", err)` to preserve error chains for debugging.

---

## Phase 4: Deep Security Review Plan & Findings Mitigation

**Objective**: Conduct a rigorous security assessment of Trident's high-trust model and establish a clear protocol for addressing vulnerabilities.

### 4.1. Static & Dynamic Security Scanning
*   **SAST (Static Application Security Testing)**:
    *   Integrate `github.com/securego/gosec/v2` into the CI/CD pipeline to scan for Go-specific security issues (e.g., hardcoded credentials, weak cryptography, unsafe SQL query construction).
*   **Vulnerability Scanning**:
    *   Use `govulncheck` to identify known vulnerabilities in imported dependencies.

### 4.2. Targeted Manual Security Audits
Given Trident's threat model (high-trust threat intelligence exchange), we must manually audit these critical areas:

1.  **Split Password Recovery Flow**:
    *   **Entropy**: Verify that the random generator used for `user_portion` and `nom_portion` (`pw.GenPass`) uses a cryptographically secure source (`crypto/rand`) and has sufficient entropy.
    *   **Storage**: Ensure the combined hash stored in the DB uses a modern secure hashing algorithm (e.g., bcrypt, Argon2id) and is not vulnerable to timing attacks during comparison.
    *   **Transit**: Verify that the transit channels (emails) are encrypted (enforced STARTTLS in Postfix configuration).
2.  **SQL Injection Prevention**:
    *   Review all raw SQL strings in `src/lib/` and `pitchfork`. Ensure *all* queries use parameterized placeholders (`$1`, `$2`) and that no string concatenation of user input occurs.
3.  **Cross-Site Scripting (XSS) in Wiki & UI**:
    *   Review the rendering pipeline of the Wiki. Ensure that `Blue Monday` sanitization is applied *immediately before* rendering to the user and that no unescaped templates allow bypasses.
    *   Audit the dual Markdown engine setup to ensure no discrepancies allow "mutation XSS" (where preview behaves differently from final rendered output in a exploitable way).
4.  **Session & Authentication (JWT)**:
    *   Review JWT signing key generation and storage. Ensure keys are loaded securely and rotated.
    *   Verify JWT expiration times and revocation mechanisms (if any).
    *   Ensure secure cookie flags (`Secure`, `HttpOnly`, `SameSite`) are enforced in production (disabling them only via explicitly flagged development arguments like `--insecurecookies`).

### 4.3. Handling Security Findings
*   **Triage Process**: Findings from scans and audits will be logged in a secure, restricted tracker and triaged by severity (Critical, High, Medium, Low).
*   **Remediation SLA**:
    *   *Critical/High*: Fix within 14 days.
    *   *Medium*: Fix within 30 days.
    *   *Low*: Fix within 90 days.
    *   **Coordinated Disclosure**: Since Trident is a community project, establish a secure security contact (currently `project@trident.li`) and a policy for coordinated disclosure, allowing existing operators of Trident instances to update before details are made public.
