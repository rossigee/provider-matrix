# Contributing to Provider Matrix

Thank you for your interest in contributing to `provider-matrix`! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Documentation](#documentation)

## Code of Conduct

This project adheres to the [Crossplane Code of Conduct](https://github.com/crossplane/crossplane/blob/master/CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/provider-matrix.git
   cd provider-matrix
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/crossplane-contrib/provider-matrix.git
   ```

## Development Setup

### Prerequisites

- **Go 1.21+**: Install from [golang.org](https://golang.org/dl/)
- **Docker**: For building container images
- **kubectl**: For Kubernetes interaction
- **make**: For running build tasks
- **Access to a Matrix homeserver**: For testing (can use matrix.org for basic testing)

### Local Environment

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Install development tools**:
   ```bash
   make setup-tools
   ```

3. **Verify setup**:
   ```bash
   make test
   make lint
   ```

### Running Locally

To run the provider locally against a Kubernetes cluster:

```bash
# Build the provider
make build

# Run locally (requires valid kubeconfig)
make run
```

## Making Changes

### Branch Management

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Keep your branch updated**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

### Types of Contributions

#### Bug Fixes
- Include a clear description of the bug
- Add a test case that reproduces the bug
- Ensure the fix doesn't break existing functionality

#### New Features
- Discuss the feature in an issue first
- Follow existing patterns and conventions
- Include comprehensive tests
- Update documentation

#### Documentation
- Fix typos, clarify unclear sections
- Add examples and use cases
- Ensure code comments are clear and helpful

### Commit Guidelines

- Use clear, descriptive commit messages
- Follow conventional commit format when possible:
  ```
  type(scope): description
  
  Longer description if needed
  
  Fixes #123
  ```
- Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`
- Keep commits focused and atomic

## Testing

### Unit Tests

```bash
# Run all unit tests
make test

# Run tests with coverage
make test-coverage

# Run tests for specific package
go test ./internal/clients/...
```

### Integration Tests

```bash
# Run integration tests (requires test environment)
make test-integration
```

### Manual Testing

1. **Set up test environment**:
   ```bash
   # Create test namespace
   kubectl create namespace provider-matrix-test
   
   # Create test credentials secret
   kubectl create secret generic matrix-creds \
     --from-literal=credentials="your_access_token" \
     -n provider-matrix-test
   ```

2. **Deploy provider locally**:
   ```bash
   make build docker-build
   # Load image into cluster or push to registry
   kubectl apply -f examples/provider/config.yaml
   ```

3. **Test resources**:
   ```bash
   kubectl apply -f examples/provider/providerconfig.yaml
   kubectl apply -f examples/user/user.yaml
   ```

### Test Guidelines

- Write tests for all new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Mock external dependencies appropriately
- Test both success and error conditions

## Submitting Changes

### Pull Request Process

1. **Ensure your branch is up to date**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run all checks locally**:
   ```bash
   make ci
   ```

3. **Push your branch**:
   ```bash
   git push origin feature/your-feature-name
   ```

4. **Create a Pull Request** on GitHub with:
   - Clear title and description
   - Reference to related issues
   - Screenshots/examples if applicable
   - Checklist of changes made

### Pull Request Requirements

- [ ] All CI checks pass
- [ ] Tests added/updated as needed
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)
- [ ] Code follows project conventions
- [ ] Commit messages are clear

### Review Process

- All PRs require at least one approval from a maintainer
- Address all review feedback
- Keep discussions constructive and professional
- Be patient - reviews may take time

## Code Style

### Go Code

- Follow standard Go formatting (`gofmt`, `goimports`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Follow Go naming conventions
- Use Go modules for dependencies

### Kubernetes Resources

- Follow Kubernetes API conventions
- Use appropriate validation tags
- Include comprehensive examples
- Document all fields clearly

### Error Handling

- Use pkg/errors for error wrapping
- Provide meaningful error messages
- Handle errors at appropriate levels
- Log errors with appropriate context

## Documentation

### Code Documentation

- Document all exported functions and types
- Use clear, concise language
- Include examples in godoc comments
- Keep documentation up to date with code changes

### User Documentation

- Update README.md for user-facing changes
- Add examples for new features
- Update API documentation
- Include troubleshooting information

### API Documentation

All CRDs should include:
- Field descriptions
- Validation constraints  
- Usage examples
- Default values

## Matrix-Specific Guidelines

### Matrix API Usage

- Use mautrix-go library for Matrix operations
- Handle both standard and admin API endpoints
- Implement proper error handling for Matrix errors
- Respect rate limiting and server capabilities

### Resource Design

- Follow Matrix specification naming
- Support common Matrix patterns
- Handle Matrix ID validation properly
- Consider federation implications

### Security Considerations

- Never log access tokens or sensitive data
- Validate all user inputs
- Use appropriate Matrix permissions
- Handle admin operations securely

## Getting Help

### Community

- **Slack**: [Crossplane Community Slack](https://slack.crossplane.io/)
- **GitHub Discussions**: For questions and general discussion
- **GitHub Issues**: For bugs and feature requests

### Matrix Resources

- **Matrix Specification**: https://spec.matrix.org/
- **mautrix-go Documentation**: https://pkg.go.dev/maunium.net/go/mautrix
- **Matrix Developer Documentation**: https://matrix.org/docs/develop/

### Development Resources

- **Crossplane Developer Guide**: https://docs.crossplane.io/contribute/
- **Controller Runtime**: https://book.kubebuilder.io/
- **Go Documentation**: https://golang.org/doc/

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes for significant contributions
- GitHub contributor statistics

Thank you for contributing to provider-matrix! 🚀