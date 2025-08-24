# Development Session Initiation

## Session Objectives
This prompt establishes the development methodology and role distribution for implementing features in the Omarchy Theme Generator project.

## Role Distribution

### Claude Code Responsibilities
- **Implementation Guides**: Provide step-by-step implementation instructions
- **Review Gates**: Analyze user implementations before proceeding
- **Test Generation**: Create execution tests following transparent test patterns
- **Documentation**: Update infrastructure docs and cross-references
- **Technical Precision**: Use correct domain terminology and reference existing code

### User Responsibilities  
- **Core Implementation**: Write the actual feature code
- **Project Direction**: Make architectural decisions and feature choices
- **Completion Confirmation**: Explicitly confirm when implementations are done
- **Code Execution**: Run tests and validation commands

## Development Process

### Step-by-Step Implementation
1. **Guide Phase**: Claude provides detailed implementation guide for current step
2. **Implementation Phase**: User implements the guided functionality
3. **Confirmation Phase**: User explicitly states "Implementation complete" 
4. **Review Phase**: Claude analyzes implementation before next step
5. **Test Phase**: Claude creates/updates execution tests
6. **Validation Phase**: User runs tests and reports results

### Execution Test Requirements
All functionality must be validated with execution tests following these principles:

- **Transparent Output**: Show initial state, transformations, expected vs actual results
- **Minimal Structure**: No frameworks, direct `go run` execution
- **Educational Value**: Tests serve as living documentation
- **Immediate Feedback**: Run with `go run tests/test-[concept]/main.go`

Example transparent test output:
```
AA compliance testing:
  Testing: RGB(119,119,119) on RGB(255,255,255) background
  Calculated contrast: 4.48:1
  Required for AA: 4.5:1
  Result: FAIL ✗ (4.48 < 4.5, difference: 0.02)
```

## Implementation Guidelines

### Technical Precision
- Reference existing implementations: "See `pkg/color/space.go`"
- Use correct domain terminology consistently
- Validate against established standards (CSS Color Module, WCAG)
- Maintain Go idioms and performance considerations

### Context Management
- Link to relevant methodology documents
- Reference technical specifications rather than repeat
- Update CLAUDE.md with essential outcomes only
- Preserve forward momentum from previous sessions

### Quality Gates
Each step must meet criteria before proceeding:
- **Correctness**: Implementation matches specifications
- **Performance**: Meets established targets (4K < 2s, Memory < 100MB)
- **Testability**: Execution test provides clear validation
- **Documentation**: Code includes appropriate godoc comments

## Session Flow

### Initiation Checklist
- [ ] Review current PROJECT.md status
- [ ] Identify session objectives from roadmap
- [ ] Confirm development environment: `go vet ./...` 
- [ ] Set explanatory output mode: `/output-style explanatory`

### Per-Step Process
1. Claude: Provide implementation guide with specific requirements
2. User: Implement guided functionality  
3. User: Confirm "Implementation complete"
4. Claude: Review implementation quality and correctness
5. Claude: Create/update execution test
6. User: Run test and report results
7. Claude: Proceed to next step or address issues

### Communication Protocol
- User uses explicit confirmation: "Implementation complete"
- Claude waits for confirmation before proceeding
- All file paths in responses must be absolute
- Reference existing code and documentation extensively

## Success Criteria

Session considered successful when:
- All planned functionality implemented and tested
- Execution tests provide transparent validation  
- Code demonstrates technical precision and understanding
- Documentation updated with essential insights
- Forward momentum preserved for next session

## References

- **Development Process**: `/home/jaime/code/omarchy-theme-generator/docs/development-methodology.md`
- **Testing Strategy**: `/home/jaime/code/omarchy-theme-generator/docs/testing-strategy.md`  
- **Technical Specification**: `/home/jaime/code/omarchy-theme-generator/docs/technical-specification.md`
- **Project Context**: `/home/jaime/code/omarchy-theme-generator/CLAUDE.md`
- **Progress Tracking**: `/home/jaime/code/omarchy-theme-generator/PROJECT.md`

Begin session by reviewing current objectives and establishing clear implementation goals.