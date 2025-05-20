# MCP IAM Pipeline

## Overview
A standardized CI/CD pipeline for deploying and managing microservices across development, staging, and production environments in the T-Mobile MCP ecosystem.

This GitLab CI/CD pipeline provides end-to-end automation for microservices with multi-environment deployment capabilities. It leverages T-Mobile's standard templates and Helm chart deployment to ensure consistent application delivery.

## Key Features
- Dynamic branch-based environments with automatic cleanup
- Intelligent workflow rules to prevent redundant pipeline runs
- Multi-stage deployment across development, staging, and production
- Automated testing with JUnit-compatible reporting
- Standardized Helm chart deployment using T-Mobile templates

## Prerequisites
- GitLab project with CI/CD enabled
- Access to T-Mobile template repositories
- Kubernetes cluster access
- Proper environment variables configured
- Go 1.19 or later installed locally
- Docker installed for local testing

## Deployment Model

### Development Environments
- Triggered by: Pushes to non-default branches
- Naming pattern: mcp-iam-{branch-name}
- URL pattern: mcp-iam-{branch-name}.azure.kube.t-mobile.com
- Lifecycle: Auto-removed after 5 days of inactivity
- Resource limits: 2 CPU, 4GB Memory

### Staging Environment
- Triggered by: Merges to the default branch
- Name: mcp-iam-stg
- URL: mcp-iam.mcp-stg.azure.kube.t-mobile.com
- Lifecycle: Persistent
- Resource limits: 4 CPU, 8GB Memory

### Production Environment
- Triggered by: Manual approval from default branch
- Name: mcp-iam-prd
- URL: mcp-iam.mcp.azure.kube.t-mobile.com
- Lifecycle: Persistent
- Resource limits: 8 CPU, 16GB Memory
- Requires approval from: Cloud Security Team

## Best Practices
- Use GitLab's environment scopes for variable management
- Create feature branches for development
- Use merge requests for code review
- Write comprehensive unit tests (minimum 80% coverage)
- Clean up unused environments regularly
- Follow Go coding standards and documentation practices

## Developer Workflow

### Local Development
1. Clone repository: `git clone <mcp_iam_repo_url>`
2. Install dependencies: `go mod download`
3. Run tests: `go test ./...`
4. Run linter: `golangci-lint run`

### Development Deployment
1. Create a feature branch: `git checkout -b feature/your-feature`
2. Ensure your code passes:
    - `golangci-lint run` for style checks
    - `go test -v ./... -cover` for unit tests (80% coverage required)
3. Push changes: `git push origin feature/your-feature`
4. Pipeline automatically deploys to dev on push: mcp-iam-{branch-name}.azure.kube.t-mobile.com
5. Dev environments auto-cleanup after 5 days

### Staging Deployment
1. Ensure all tests pass in development
2. Create Merge Request to main branch
3. Required approvals: 2 team members
4. Pipeline checks:
    - Go lint validation
    - Unit test coverage (minimum 80%)
    - Security scanning
    - dev endpoint
5. After merge, automatic deployment to staging
6. Verify at: mcp-iam.mcp-stg.azure.kube.t-mobile.com

### Quality Gates
Pipeline enforces:
- Go linting using golangci-lint
- Unit test coverage >75%
- Security vulnerability scanning
- Dependency vulnerability checks
- Container image scanning

## Troubleshooting
Common issues and solutions:
- Image pull failures: Check registry credentials
- Deployment failures: Verify resource limits
- Environment configuration: Validate environment variables
- Pipeline failures: Check logs in GitLab CI/CD interface

## Contact
For pipeline issues contact:
- MCP Cloud Security Team