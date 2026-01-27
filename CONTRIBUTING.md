# Contributing to OpenTelemetry Plugin for Genkit Go

Thank you for your interest in contributing to the OpenTelemetry Plugin for Genkit Go! We welcome contributions from the community.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

This project and everyone participating in it is governed by our commitment to creating a welcoming and inclusive environment. By participating, you are expected to uphold this standard.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your feature or bug fix
4. Make your changes
5. Test your changes
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.24.1 or later
- Git

### Setup Instructions

1. Clone the repository:
   ```bash
   git clone https://github.com/genkit-ai/opentelemetry-go-plugin.git
   cd genkit-opentelemetry-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

4. Run examples to verify setup:
   ```bash
   cd examples/basic
   go run main.go
   ```

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Go version and OS

### Suggesting Enhancements

Enhancement suggestions are welcome! Please:

- Use a clear and descriptive title
- Provide a detailed description of the suggested enhancement
- Explain why this enhancement would be useful
- Include code examples if applicable

### Contributing Code

1. **Choose an Issue**: Look for issues labeled `good first issue` or `help wanted`
2. **Create a Branch**: Create a feature branch from `main`
3. **Make Changes**: Implement your changes following our coding standards
4. **Add Tests**: Include tests for new functionality
5. **Update Documentation**: Update README.md and code comments as needed
6. **Test**: Ensure all tests pass and new functionality works correctly

## Commit Message Convention

This project uses [Conventional Commits](https://conventionalcommits.org/) for automatic semantic versioning and changelog generation. All commit messages must follow this format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Commit Types

- **feat**: A new feature (triggers minor version bump)
- **fix**: A bug fix (triggers patch version bump)
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to our CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files

### Breaking Changes

Breaking changes trigger a major version bump and should be indicated with `!` after the type/scope:

```
feat(api)!: change model configuration interface

BREAKING CHANGE: Model configuration now requires explicit region specification.
```

### Examples

```bash
feat(models): add support for Claude 3.5 Sonnet
fix(streaming): resolve timeout issue with long responses
docs(readme): update installation instructions
test(examples): add integration tests for tool calling
```

### Automated Release Process

When commits are pushed to the main branch:
1. **Semantic Release** analyzes commit messages
2. **Version** is automatically bumped based on commit types
3. **Changelog** is generated and updated
4. **GitHub Release** is created with release notes
5. **GoReleaser** builds and publishes artifacts

## Pull Request Process

1. **Update Documentation**: Ensure the README.md and other documentation are updated
2. **Add Tests**: Include tests that cover your changes
3. **Follow Commit Standards**: Use conventional commit format for all commits
4. **Create Pull Request**: Submit a pull request with a clear title and description
5. **Address Feedback**: Respond to review comments promptly

### Pull Request Template

When creating a pull request, please include:

- **Description**: What changes does this PR introduce?
- **Type of Change**: Bug fix, new feature, breaking change, documentation update
- **Testing**: How has this been tested?
- **Checklist**: Confirm you've followed the contribution guidelines

## Coding Standards

### Go Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Include comments for exported functions and types
- Follow Go naming conventions
- Keep functions focused and reasonably sized

### Code Organization

- Group related functionality together
- Use clear package structure
- Separate concerns appropriately
- Follow existing patterns in the codebase

### Error Handling

- Use Go's standard error handling patterns
- Provide meaningful error messages
- Handle errors at appropriate levels
- Don't ignore errors unless explicitly justified

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

### Test Requirements

- All new code should include appropriate tests
- Tests should cover both success and error cases
- Use table-driven tests where appropriate

### Test Structure

- Unit tests: Test individual functions and methods
- Example tests: Ensure examples continue to work

## Documentation

### Code Documentation

- All exported functions and types must have comments
- Comments should explain what the code does, not how
- Include examples in comments where helpful
- Use standard Go documentation conventions

### README Updates

- Keep the README.md up to date with new features
- Include examples for new functionality
- Update supported models list when adding new models
- Maintain accurate installation and usage instructions

## Release Process

Releases are handled by maintainers. The process includes:

1. Version bumping following semantic versioning
2. Updating CHANGELOG.md
3. Creating GitHub releases with release notes
4. Ensuring all tests pass before release

## Questions?

If you have questions about contributing, please:

1. Check existing issues and discussions
2. Create a new issue with the `question` label
3. Reach out to maintainers

Thank you for contributing to OpenTelemetry Plugin for Genkit Go!
