## 1. Create CI workflow file

- [x] 1.1 Create `.github/workflows/ci.yml` with push/PR triggers on `main`, Go 1.25 setup, and steps: `go vet`, `go test -race -coverprofile`, `go tool cover -func`
- [x] 1.2 Push branch and verify the workflow triggers on the PR (check Actions tab on GitHub)
- [x] 1.3 Ensure all steps pass green (vet, test, coverage)
