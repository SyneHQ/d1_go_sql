# Contributing to D1 Go SQL Driver

Thank you for your interest in contributing to the D1 Go SQL Driver! We welcome contributions from the community.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/d1_go_sql.git
   cd d1_go_sql
   ```
3. **Install dependencies**:
   ```bash
   make deps
   ```

## Development Workflow

### Setting Up Your Environment

1. Ensure you have Go 1.21 or later installed
2. Set up your Cloudflare credentials for testing:
   ```bash
   export CLOUDFLARE_ACCOUNT_ID="your-account-id"
   export CLOUDFLARE_API_TOKEN="your-api-token"
   export D1_DATABASE_ID="your-database-id"
   ```

### Making Changes

1. **Create a new branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards

3. **Format your code**:
   ```bash
   make fmt
   ```

4. **Run tests**:
   ```bash
   make test
   ```

5. **Run linters**:
   ```bash
   make lint
   ```

6. **Commit your changes** with clear, descriptive commit messages:
   ```bash
   git commit -m "Add feature: description of your changes"
   ```

### Commit Message Guidelines

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

Examples:
```
Add support for named parameters
Fix connection leak in transaction rollback
Update documentation for DSN format
```

## Coding Standards

### Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` for formatting (run `make fmt`)
- Keep functions focused and single-purpose
- Write clear, self-documenting code with meaningful variable names

### Documentation

- Document all exported types, functions, and constants
- Include examples in documentation where appropriate
- Keep README.md up to date with new features
- Add inline comments for complex logic

### Testing

- Write unit tests for all new functionality
- Maintain or improve code coverage
- Test edge cases and error conditions
- Use table-driven tests where appropriate

Example test structure:
```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Error Handling

- Return errors rather than panicking
- Wrap errors with context using `fmt.Errorf` with `%w`
- Define custom error types for specific error conditions
- Use clear, actionable error messages

### Performance

- Avoid unnecessary allocations
- Use efficient algorithms and data structures
- Profile code for performance-critical sections
- Document any performance trade-offs

## Pull Request Process

1. **Update documentation** for any changed functionality
2. **Add tests** for new features or bug fixes
3. **Ensure all tests pass**: `make test`
4. **Ensure code is formatted**: `make fmt`
5. **Ensure no linter errors**: `make lint`
6. **Update CHANGELOG.md** if applicable
7. **Submit your pull request** with a clear description

### Pull Request Description Template

```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
Describe the tests you ran and how to reproduce them

## Checklist
- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
```

## Code Review Process

1. At least one maintainer will review your pull request
2. Address any feedback or requested changes
3. Once approved, a maintainer will merge your pull request

## Reporting Bugs

When reporting bugs, please include:

1. **Description**: Clear description of the bug
2. **Steps to reproduce**: Detailed steps to reproduce the issue
3. **Expected behavior**: What you expected to happen
4. **Actual behavior**: What actually happened
5. **Environment**: Go version, OS, and relevant configuration
6. **Code samples**: Minimal code to reproduce the issue
7. **Error messages**: Full error messages and stack traces

## Feature Requests

When requesting features, please include:

1. **Use case**: Describe the problem you're trying to solve
2. **Proposed solution**: Your suggested approach
3. **Alternatives**: Other approaches you've considered
4. **Additional context**: Any other relevant information

## Questions and Support

- **GitHub Issues**: For bug reports and feature requests
- **GitHub Discussions**: For questions and general discussion

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be recognized in the project's README.md and CHANGELOG.md.

Thank you for contributing to D1 Go SQL Driver! 🎉

