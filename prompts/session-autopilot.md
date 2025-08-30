# Autopilot Session Initialization

## Quick Start
Use this prompt to initialize an autopilot session for rapid development:

```
I'm starting Session [N]: [Session Title] in autopilot mode on the autopilot branch.

Please:
1. Start in plan mode
2. Review the session requirements from PROJECT.md
3. Create a comprehensive execution plan
4. Request any clarifications needed
5. Once I approve, I'll place you in auto-accept mode to implement

Focus on functionality over performance optimization. Use agent peer reviews at key implementation points.
```

## Process Overview

### Phase 1: Planning (Plan Mode)
- Review session requirements from PROJECT.md
- Examine existing codebase state
- Identify dependencies and integration points
- Submit detailed execution plan with clear milestones

### Phase 2: Implementation (Auto-Accept Mode)
- Build functionality incrementally
- Use TodoWrite to track progress
- Engage specialized agents for peer review:
  - **go-engineer**: Go code patterns and idioms
  - **color-science-specialist**: Color algorithms and theory
  - **docs-consistency-checker**: Documentation alignment
- Test each component as it's built

### Phase 3: Validation
- Run all execution tests with `go run tests/test-*/main.go`
- Ensure 100% test success (functionality, not performance)
- Run `go vet ./...` for type checking
- Verify no regressions introduced

### Phase 4: Finalization
- Update PROJECT.md session log
- Update test README files with outputs
- Document any architectural decisions
- Prepare for next session

## Key Principles

1. **Functionality First**: Get it working, optimize later
2. **Incremental Testing**: Test as you go, not all at end
3. **Agent Validation**: Use peer review for quality assurance
4. **Clear Documentation**: Update docs as part of implementation

## Validation Checklist

- [ ] All execution tests pass
- [ ] Code follows Go idioms (go-engineer validated)
- [ ] Color science accurate (specialist validated)
- [ ] Documentation updated
- [ ] No regressions introduced
- [ ] SESSION.md log updated

## Agent Review Protocol

```
Implementation Step
    ↓
Agent Review (if applicable)
    ↓
Test Execution
    ↓
Fix Issues (if any)
    ↓
Next Step
```

## Session Success Criteria

A session is complete when:
1. All planned functionality is implemented
2. All tests pass with clear output
3. Documentation is updated
4. No compilation errors (`go vet ./...` clean)
5. Session log in PROJECT.md is updated

## Example Session Flow

```markdown
User: I'm starting Session 4: Color Synthesis & Palette Generation in autopilot mode on the autopilot branch.