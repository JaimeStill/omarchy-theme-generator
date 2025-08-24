# Development Session Finalization

## Session Closeout Objectives
This prompt ensures proper documentation, artifact updates, and preparation for future sessions following the Intelligent Development methodology.

## Finalization Responsibilities

### Claude Code Tasks
- **Documentation Comments**: Add comprehensive comments to all new/revised infrastructure
- **Repository Artifacts**: Update PROJECT.md, technical specifications, and cross-references
- **Progress Log Update**: Update PROJECT.md with session progress entry
- **Context Optimization**: Update CLAUDE.md with key decisions, remove outdated details
- **Reference Validation**: Ensure all internal links and file paths are accurate

### User Tasks
- **Final Validation**: Run complete test suite one final time
- **Commit Preparation**: Stage changes and prepare commit message
- **Next Session Input**: Provide direction for subsequent session focus

## Finalization Process

### 1. Documentation Review
**Infrastructure Comments**
- [ ] All new public functions have godoc comments
- [ ] All new types have purpose documentation  
- [ ] Complex algorithms include implementation notes
- [ ] Test files have clear purpose statements

**Code Quality Check**
- [ ] Run `go vet ./...` for type validation
- [ ] Run `go fmt ./...` for consistent formatting
- [ ] All execution tests run successfully
- [ ] No compilation errors or warnings

### 2. Repository Artifact Updates

**PROJECT.md Updates**
- [ ] Mark completed tasks with ✅ 
- [ ] Add session entry to Progress Log following template:
  ```markdown
  ### Session N: [Date]
  **Completed:**
  - ✅ Task description - reference to code/test
  
  **Insights:**
  - Key learning or discovery
  
  **Decision:**
  - Architectural choice made - link to docs/decisions/
  
  **Next:**
  - What to tackle in next session
  ```
- [ ] Update metrics tracking table if applicable

**Technical Specification Updates**
- [ ] Add new algorithms to specification if implemented
- [ ] Update performance targets based on test results
- [ ] Document any architectural changes

**CLAUDE.md Context Updates**
- [ ] Add essential new insights to Implementation Status
- [ ] Update Key Technical Decisions with any changes
- [ ] Remove outdated details to maintain context efficiency
- [ ] Ensure Next Session Focus is accurate

### 3. Cross-Reference Validation

**Link Accuracy Check**
- [ ] All `See pkg/[package]/[file].go` references are correct
- [ ] All documentation links point to existing files
- [ ] All absolute paths in documentation are accurate
- [ ] Import statements match directory structure

**Terminology Consistency**
- [ ] Binary name `omarchy-theme-gen` used consistently
- [ ] Import path `github.com/JaimeStill/omarchy-theme-generator` correct
- [ ] Go version references accurate (Go 1.25)
- [ ] Package names match directory structure

### 4. PROJECT.md Progress Log Update

Update PROJECT.md Progress Log with session entry following the established template format:

```markdown
### Session N: [Date]
**Completed:**
- ✅ Task description - reference to code/test

**Insights:**
- Key learning or discovery

**Decision:**
- Architectural choice made - link to docs/decisions/

**Next:**
- What to tackle in next session
```

### 5. Context Optimization

**CLAUDE.md Maintenance**
- Remove implementation details now captured in code
- Preserve architectural decisions and their rationale
- Keep essential performance targets and constraints
- Maintain concise command reference
- Update Implementation Status accurately

**Information Density Check**
- Target: CLAUDE.md under 2000 lines
- Method: Reference rather than repeat
- Focus: Essential outcomes and forward momentum
- Remove: Outdated details and resolved issues

### 6. Next Session Preparation

**Handoff Documentation**
- [ ] Clear objectives identified for next session
- [ ] Dependencies and prerequisites documented
- [ ] Any required research or preparation noted
- [ ] Implementation approach considerations captured

**State Validation**  
- [ ] Repository in clean state for next developer
- [ ] All tests passing and documented
- [ ] Context optimized for efficient onboarding
- [ ] Progress accurately reflected in all documentation

## Success Criteria

Session considered properly finalized when:

- **Documentation Complete**: All new code properly commented
- **Artifacts Updated**: PROJECT.md and technical docs reflect current state  
- **References Valid**: All links and paths accurate
- **Context Optimized**: CLAUDE.md concise but complete
- **Forward Momentum**: Next session has clear objectives
- **Quality Assured**: All tests pass, code formatted and validated

## Final Checklist

### User Actions Required
- [ ] Run final test validation: `go vet ./... && go fmt ./...`
- [ ] Execute all relevant tests: `go run tests/test-*/main.go`
- [ ] Review updated documentation for accuracy
- [ ] Provide next session direction or priorities

### Claude Actions Required  
- [ ] Complete all documentation updates
- [ ] Validate all cross-references and links
- [ ] Update PROJECT.md Progress Log with session entry
- [ ] Optimize CLAUDE.md context for next session
- [ ] Confirm all artifacts properly updated

## References

- **Development Methodology**: `/home/jaime/code/omarchy-theme-generator/docs/development-methodology.md`
- **Testing Strategy**: `/home/jaime/code/omarchy-theme-generator/docs/testing-strategy.md`
- **Progress Tracking**: `/home/jaime/code/omarchy-theme-generator/PROJECT.md`
- **Context Management**: `/home/jaime/code/omarchy-theme-generator/CLAUDE.md`

The session is complete when all finalization tasks are verified and next session objectives are clearly established.