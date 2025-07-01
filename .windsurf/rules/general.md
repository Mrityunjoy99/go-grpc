---
trigger: always_on
---

# General Guidelines for go-grpc Project

## 1. Code Quality
- Follow Go best practices and idioms.
- Ensure all code is formatted using `gofmt` before committing.
- Write clear, concise, and self-explanatory code.
- Avoid large functions; break logic into smaller, reusable units.
- Use meaningful variable, function, and type names.

## 2. gRPC and Protobuf
- Define all service contracts in `.proto` files under the `proto/` directory.
- Regenerate Go code from `.proto` files using the provided scripts before committing changes.
- Keep backward compatibility in mind when updating proto definitions.

## 3. Documentation
- Document all public functions, types, and packages using Go doc comments.
- Keep the `README.md` up to date with build, run, and usage instructions.
- Add comments to complex or non-obvious code sections.

## 4. Testing
- Write unit tests for all major logic and components.
- Use table-driven tests where appropriate.
- Run all tests and ensure they pass before submitting a PR.

## 5. Collaboration
- Use feature branches for new work; do not commit directly to `main`.
- Submit pull requests with clear descriptions of changes.
- Review and address feedback promptly.
- Resolve merge conflicts as soon as possible.

## 6. Dependency Management
- Use Go modules (`go.mod`, `go.sum`) to manage dependencies.
- Do not commit unnecessary dependencies or vendor files.

## 7. Security and Confidentiality
- Do not commit secrets, credentials, or sensitive information.
- Follow security best practices for Go and gRPC.

## 8. Scripts and Automation
- Use scripts in the `scripts/` directory for code generation or repetitive tasks.
- Keep scripts well-documented and cross-platform where possible.