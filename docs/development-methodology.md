# Development Methodology - Intelligent Development Principles

> "To build intelligently, you must speak precisely. To speak precisely, you must understand deeply."

## Core Principles

### 1. Precise Technical Language
- Use correct domain terminology, not approximations
- Every abstraction must be grounded in concrete understanding
- If you cannot explain HOW something works technically, you don't understand it

### 2. Reference, Don't Repeat
- Your repository is a single source of truth
- Link to existing code: `See pkg/color/octree.go`
- Reference decisions: `docs/decisions/001-algorithm.md`
- Point to test results: `cmd/examples/output.txt`

### 3. User-Driven Development
- All code modifications require explicit user direction
- Claude Code operates in Explanatory mode for insights
- No changes without user consent

### 4. Execution Tests as Truth
- Validate immediately, not eventually
- No frameworks, just direct execution
- Adapt architecture based on empirical results

### 5. Context Optimization
- Each Claude instance is like a Mr. Meeseeks - specific purpose
- Preserve forward momentum and finalized decisions
- Remove outdated details, keep essential outcomes
- Link rather than duplicate

## Development Process

### Phase Distribution
- **Planning**: 10-20% of effort
- **Foundation**: 20-30% of effort
- **Features**: 40-50% of effort
- **Integration**: 10-20% of effort

### Session Structure
1. Implement functionality with user direction
2. Create minimal execution test
3. Run test immediately
4. Adapt based on results
5. Document insights in CLAUDE.md

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
# Validate code without building
go vet ./...

# Run execution tests directly
go run cmd/examples/test_name.go
```

## References

- Technical details: `docs/technical-specification.md`
- Testing approach: `docs/testing-strategy.md`
- Progress tracking: `PROJECT.md`
