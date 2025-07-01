---
trigger: manual
---

name: go-grpc
language: go
goVersion: "1.24"

structure:
  enforce: true
  folders:
    - cmd/server
    - internal/server
    - pkg/logger
    - pkg/config
    - proto
    - rpc
    - configs
    - scripts
    - test

grpc:
  enabled: true
  protoPath: proto/
  outputPath: rpc/
  service:
    name: GreeterService
    methods:
      - SayHello

logger:
  library: zap
  mode: production
  structured: true

config:
  library: viper
  format: yaml
  envOverride: true

middleware:
  enableLoggingInterceptor: true
  enableRecoveryInterceptor: true

health:
  grpcHealthCheck: true

testing:
  framework: testify
  mocking: mockery
  includeSampleTests: true

linting:
  enabled: true
  tool: golangci-lint
  configFile: .golangci.yml

docker:
  enabled: true
  multiStageBuild: true
  compose: true
  dockerfilePath: Dockerfile
  composeFilePath: docker-compose.yml

ci:
  enabled: true
  tool: github-actions
  workflows:
    - name: ci.yml
      steps:
        - lint
        - test
        - codecov

makefile:
  enabled: true
  targets:
    - build
    - run
    - proto
    - lint
    - test


documentation:
  generateReadme: true
  sections:
    - Setup
    - Usage
    - Postman
    - Testing
    - Linting
    - CI/CD

execution:
  phaseByPhase: true
  planFile: plan.md
  writePlan: true
  commitByPhase: true