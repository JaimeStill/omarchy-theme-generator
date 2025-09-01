# Development Methodology - Intelligent Development Principles

> "To build intelligently, you must speak precisely. To speak precisely, you must understand deeply."

## Core Principles

### 1. Precise Technical Language
- Use correct domain terminology, not approximations
- Every abstraction must be grounded in concrete understanding
- If you cannot explain HOW something works technically, you don't understand it

### 2. Reference, Don't Repeat
- Your repository is a single source of truth
- Link to existing code: `See pkg/formats/color.go`
- Reference architecture: `docs/architecture.md`
- Point to test results and validation evidence

### 3. User-Driven Development
- All code modifications require explicit user direction
- Claude Code operates in Explanatory mode for insights
- No changes without user consent

### 4. Standard Go Tests as Truth
- Use standard Go test files (`*_test.go`)
- Unit tests for each package in isolation
- Integration tests for complete workflows
- Validate immediately with `go test`

### 5. Context Optimization
- Each Claude instance is like a Mr. Meeseeks - specific purpose
- Preserve forward momentum and finalized decisions
- Remove outdated details, keep essential outcomes
- Link rather than duplicate

## Development Process

### Component-Based Development
- **Foundation Layer**: Core utilities with no dependencies
- **Analysis Layer**: Image understanding and profile detection
- **Processing Layer**: Extraction algorithms and orchestration
- **Generation Layer**: Theme creation and output formatting
- **Application Layer**: CLI interface and user interaction

### Layer Development Rules
1. **Dependencies flow downward only** - higher layers depend on lower layers
2. **Standard library first** - prefer Go standard types over custom implementations
3. **Settings-driven** - no hardcoded values in business logic
4. **Purpose-driven** - organize by intent (role-based colors, not frequency)
5. **Extensible design** - interfaces for pluggable components

### Implementation Workflow
1. Identify target layer and dependencies
2. Implement functionality with user direction
3. Create standard Go tests (`*_test.go`)
4. Validate with `go test` and `go vet`
5. Update documentation and architecture references

### Knowledge Requirements

Before starting:
1. **Assess Foundation**: What do you genuinely understand?
2. **Map Required Knowledge**: What algorithms and structures are needed?
3. **Bridge Gaps**: Study fundamentals, not just libraries
4. **Demonstrate Mastery**: Use precise technical language

## Success Metrics

### Technical Precision
- Specification reviewable by domain expert
- Correct terminology throughout
- Algorithms match literature

### Process Efficiency
- New AI instance productive in < 5 minutes
- CLAUDE.md under 2000 lines
- Execution tests validate concepts immediately

### Knowledge Transfer
- Code demonstrates understanding
- Decisions have technical rationale
- Future developers can extend work

## Tools & Configuration

### Claude Code Setup
```bash
# Set output style for educational insights
/output-style explanatory
```

### Validation
```bash
# Validate code
go vet ./...

# Run standard Go tests
go test ./...

# Run specific package tests
go test ./pkg/formats

# Run tests with verbose output
go test -v ./tests
```

## References

- Architecture overview: `docs/architecture.md`
- Testing approach: `docs/testing-strategy.md`
- Progress tracking: `PROJECT.md`
- Terminology reference: `docs/glossary.md`
