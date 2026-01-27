# Contributing to Stream

Thank you for your interest in contributing to Stream! This project aims to bring Kotlin-like functional stream operations to Go, and I welcome contributions from the community.

## Ways to Contribute

- **Report bugs** - Found a bug? Open an issue
- **Suggest features** - Have an idea? Let's discuss it
- **Improve documentation** - Help others understand the library
- **Submit code** - Fix bugs or implement new features
- **Write tests** - Improve test coverage

## Getting Started

### Prerequisites

- Go 1.25.4 or later
- Git

### Setup

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/stream.git
   cd stream
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run tests to verify setup:
   ```bash
   go test ./...
   ```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Your Changes

- Write clear, idiomatic Go code
- Add tests for new functionality
- Update documentation as needed
- Run tests frequently

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Run benchmarks
go test -bench=. ./...

# Check for race conditions
go test -race ./...
```

### 4. Format and Lint

```bash
# Format code
go fmt ./...

# Run go vet
go vet ./...

# Use goimports (if installed)
goimports -w .
```

### 5. Commit Your Changes

Use clear, descriptive commit messages:

```
<type>: <short description>

<optional longer description>

Fixes #123
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

**Examples:**
```
fix: correct index out of range in Map function

The Map function was creating an empty slice and then trying to 
access indices directly, causing a panic.

Fixes #42
```

```
feat: add TakeWhile operation for sequences

Implements lazy TakeWhile that stops consuming elements once
the predicate returns false.
```

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then open a Pull Request on GitHub.

## Pull Request Guidelines

### Before Submitting

- ✅ All tests pass (`go test ./...`)
- ✅ Code is formatted (`go fmt ./...`)
- ✅ No `go vet` warnings
- ✅ New features have tests
- ✅ Bug fixes have regression tests
- ✅ Public APIs have GoDoc comments
- ✅ README updated if needed

### PR Description

Include:
- **What**: What changes does this PR make?
- **Why**: Why are these changes needed?
- **How**: How were the changes implemented?
- **Testing**: How was this tested?
- **Related Issues**: Link to related issues

### Example PR Description

```markdown
## What
Adds a TakeWhile operation for lazy sequences

## Why
Completes the set of standard functional operations and allows
users to consume sequences until a condition is met.

## How
Implemented as a lazy operation that wraps the sequence's next()
function and stops iteration when the predicate returns false.

## Testing
- Added unit tests covering various scenarios
- Verified lazy evaluation with a test counter
- Added benchmark comparing with Filter + Take

Closes #56
```

## Coding Standards

### Go Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Keep functions focused and small
- Use meaningful variable names
- Prefer table-driven tests

### GoDoc Comments

All exported types, functions, and methods must have GoDoc comments:

```go
// Map applies the mapper function to each element in the stream
// and returns a new stream with the transformed elements.
// If an error occurs during mapping, it is captured and returned
// when a terminal operation is called.
func (s *slice[T]) Map(mapper func(elem T) (T, error)) *slice[T] {
```

### Error Handling

- Return errors, don't panic (except for programmer errors)
- Use the error field pattern for stream operations
- Provide clear error messages

### Testing

- Use table-driven tests where appropriate
- Test both success and error cases
- Include edge cases (empty slices, nil, etc.)
- Use meaningful test names: `TestFilter_WithEvenNumbers`
- Use `testify/assert` for assertions

**Example:**

```go
func TestMap_TransformsElements(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		mapper   func(int) (int, error)
		expected []int
		hasError bool
	}{
		{
			name:     "doubles all elements",
			input:    []int{1, 2, 3},
			mapper:   func(n int) (int, error) { return n * 2, nil },
			expected: []int{2, 4, 6},
			hasError: false,
		},
		// more test cases...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := From(tt.input).Map(tt.mapper).ToSlice()
			
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
```

## Reporting Issues

### Bug Reports

Include:
- **Description**: Clear description of the bug
- **Steps to Reproduce**: Minimal example to reproduce
- **Expected Behavior**: What should happen
- **Actual Behavior**: What actually happens
- **Go Version**: Output of `go version`
- **OS**: Your operating system

### Feature Requests

Include:
- **Use Case**: What problem does this solve?
- **Proposed API**: How should it work?
- **Examples**: Show example usage
- **Alternatives**: What alternatives did you consider?

## Code Review Process

1. Maintainers will review your PR
2. Address feedback by pushing new commits
3. Once approved, a maintainer will merge
4. Your contribution will be in the next release!

## Questions?

Feel free to open an issue for questions or discussion. I'am here to help!

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
