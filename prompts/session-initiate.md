# Development Session Initiation

## Session Objectives
This prompt establishes the development methodology and role distribution for implementing features in the Omarchy Theme Generator project.

## Role Distribution

### AI (Claude Code) Responsibilities
- **Implementation Guides**: Provide comprehensive step-by-step implementation instructions
- **Test Generation**: Create unit tests in tests/ subdirectories
- **Documentation**: Maintain accurate documentation and cross-references
- **Code Review**: Analyze implementations for best practices and patterns
- **Technical Precision**: Use correct domain terminology and reference existing code

### User Responsibilities  
- **Core Implementation**: Develop source code based on AI guides
- **Architecture Direction**: Make design decisions and set priorities
- **Review and Refine**: Edit and approve AI outputs
- **Code Execution**: Run tests and validation commands

## Development Process

### Implementation Workflow
1. **Guide Phase**: AI provides comprehensive implementation guide
2. **Implementation Phase**: User develops code based on guide
3. **Test Creation**: AI generates unit tests for the implementation
4. **Review Phase**: AI analyzes code for quality and patterns
5. **Documentation**: AI updates relevant documentation
6. **Validation Phase**: User runs tests and confirms results

### Testing Requirements
All functionality must be validated with standard Go tests following these principles:

- **Standard Go Testing**: Use `*_test.go` files in tests/ subdirectories per package
- **Package Organization**: tests/formats/, tests/extractor/, tests/chromatic/, etc.
- **Transparent Output**: Clear test failures showing expected vs actual results
- **Real-World Validation**: Integration tests with actual image samples
- **Immediate Feedback**: Run with `go test ./tests/... -v`

Example Go test structure:
```go
func TestColorConversion(t *testing.T) {
    input := color.RGBA{255, 128, 64, 255}
    h, s, l := formats.RGBToHSL(input)
    
    expectedH, expectedS, expectedL := 0.083, 1.0, 0.627
    if h != expectedH || s != expectedS || l != expectedL {
        t.Errorf("Expected HSL(%.3f, %.3f, %.3f), got HSL(%.3f, %.3f, %.3f)", 
                 expectedH, expectedS, expectedL, h, s, l)
    }
}
```

## Architecture Context

The project uses a layered architecture with clear dependencies:

### Package Layers
1. **Foundation** (formats, settings, config) - No dependencies
2. **Analysis** (analysis) - Depends on foundation
3. **Processing** (strategies, extractor) - Depends on analysis
4. **Generation** (schemes, theme) - Depends on processing
5. **Application** (cmd) - Depends on all

### Key Principles
- Use standard library types (`color.RGBA` not custom types)
- Settings-driven (no hardcoded values)
- Purpose-driven extraction (role-based, not frequency)
- Clear separation of concerns

### Current Implementation Status
- ✅ Multi-strategy extraction (frequency/saliency)
- ✅ Foundation package structures created
- 🔄 Unit tests for all packages in development
- 🔄 Color derivation algorithms in pkg/chromatic
- ⏳ Strategy extraction from extractor pending
- ⏳ Theme generation pending

## Implementation Guidelines

### Technical Precision
- Reference existing implementations: "See `pkg/formats/color.go`"
- Use correct domain terminology from docs/glossary.md
- Follow layered architecture dependencies
- Validate against established standards (CSS Color Module, WCAG)
- Maintain Go idioms and standard library usage

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
- [ ] Verify test structure: `ls tests/`

### Per-Step Process
1. AI: Provide comprehensive implementation guide
2. User: Develop source code based on guide
3. AI: Create unit tests in appropriate tests/ subdirectory
4. AI: Review code and update documentation
5. User: Run tests with `go test ./tests/... -v`
6. User: Provide feedback on results
7. AI: Address issues or proceed to next task

### Communication Protocol
- AI provides comprehensive guides before user implementation
- User implements and provides feedback
- AI creates tests and documentation
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

- **Architecture Overview**: `../docs/architecture.md`
- **Development Process**: `../docs/development-methodology.md`
- **Testing Strategy**: `../docs/testing-strategy.md`  
- **Terminology Reference**: `../docs/glossary.md`
- **Project Context**: `../CLAUDE.md`
- **Progress Tracking**: `../PROJECT.md`

Begin session by reviewing current objectives and establishing clear implementation goals.