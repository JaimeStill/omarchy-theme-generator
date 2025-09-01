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

### Testing Requirements
All functionality must be validated with standard Go tests following these principles:

- **Standard Go Testing**: Use `*_test.go` files in tests/ directory
- **Layered Testing**: Each package tested in isolation with clear dependencies
- **Transparent Output**: Clear test failures showing expected vs actual results
- **Real-World Validation**: Integration tests with actual image samples
- **Immediate Feedback**: Run with `go test ./tests -v`

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
- ✅ Settings-driven configuration
- 🔄 Architecture refactoring in progress
- ⏳ Purpose-driven extraction pending
- ⏳ Color scheme generation pending

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

- **Architecture Overview**: `../docs/architecture.md`
- **Development Process**: `../docs/development-methodology.md`
- **Testing Strategy**: `../docs/testing-strategy.md`  
- **Terminology Reference**: `../docs/glossary.md`
- **Project Context**: `../CLAUDE.md`
- **Progress Tracking**: `../PROJECT.md`

Begin session by reviewing current objectives and establishing clear implementation goals.