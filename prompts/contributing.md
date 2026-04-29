# Contributing to SentinelFlow

Welcome! As a Senior-led project, we maintain high standards for code quality and system reliability.

## 🛠 Development Workflow
1.  **Fork & Branch:** Use `feature/JIRA-ID-description`.
2.  **Atomic Commits:** Each commit should represent one logical change.
3.  **Linting:** * Go: `golangci-lint run`
    * TS: `npm run lint`

## 🧪 Testing Standards
* **Unit Tests:** Mandatory for all business logic. Use `testify/assert` in Go.
* **Integration Tests:** Required for any service interacting with Kafka or Postgres.
* **Benchmark:** If modifying the Ingestion worker pool, run `go test -bench` to ensure no performance regressions.

## 🔍 Pull Request (PR) Policy
* All PRs require **one approval** from a core maintainer.
* Documentation must be updated in the `/docs` folder if API contracts change.
* CI Pipeline (GitHub Actions) must pass 100% of tests.

> "We optimize for readability and reliability. Complex 'clever' code is discouraged in favor of maintainable, observable code."