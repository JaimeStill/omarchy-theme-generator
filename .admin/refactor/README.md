# Development Adjustments

This directory represents a substantial adjustment to the development workflow, project management approach, and code structure. The artifacts it contains were output as the result of an extensive review and planning phase with Claude Opus 4.1 on https://claude.ai.

The ultimate execution order is to:

1. Update [PROJECT.md](../../PROJECT.md), [README.md](../../README.md), and the root [docs/](../) artifacts to align our context with the propsed adjustments.
2. Tune [CLAUDE.md](../../CLAUDE.md), [.claude/agents](../../.claude/agents) (if necessary), and [prompts/](../../prompts/) to ensure our future development efforts align with the proposed adjustments and are as smooth as possible.
3. Incorporate the immediate changes laid out in [PROJECT.md Restructure](./project-restructure-and-formats-design.md).
4. Optimize the extractor package and embedded strategies features to be effectively layered into separate packages.
5. Extract the settings and configuration features into their own package (settings for system operations values, configuration for theme user configuration).
6. Making the requisite adjustments to the image processing and metadata structures proposed in [Extraction Architecture Decisions](./extraction-architecture-decisions.md).

## Description of Artifacts

1. [Extraction Architecture Decisions](./extraction-architecture-decisions.md): This artifact captures the shift from frequency-based to purpose-driven color extraction, where colors are categorized into their intended role (backgrounds, foregrounds, accents) rather than just their occurrence count. It defines 7 core architectural decisions include mode-aware role assignment, settings-driven configuration, and edge case handling for minimal-color images. The document provides detailed implementation examples for the new `ThemeColorMap` structure and role-based organization system. It also includes updated session tasks for PROJECT.md and the documentation sections to add to the technical specification.

2. [PROJECT.md Restructure & Formats Package Design](./project-restructure-and-formats-design.md): This artifact restructure PROJECT.md from session-based planning to a component-based architecture organized by dependency layers. It introduces the simplified `pkg/formats` package that replaces the over-engineered `pkg/color`, keeping only essential functions like `RGBToHSL()`, `ContrastRatio()`, and `ToHex()`. The document defines 5 architectural layers (Foundation -> Analysis -> Extraction -> Generation -> Application) with clear dependencies between packages. It also provides the complete implementation for the new formats package using standard Go `color.RGBA` types instead of custom types.

3. [Implementation Roadmap - Architecture Refactoring](./implementation-roadmap.md): This artifact provides a comprehensive 5-phase implementation plan for refactoring the entire codebase to align with the new architecture. Phase 1 covers foundation refactoring (pkg/color -> pkg/formats), Phase 2 handles extractor decomposition, Phase 3 implements purpose-driven extraction, Phase 4 updates documentation, and Phase 5 handles testing. Each phase includes specific code changes, migration examples, and success criteria with clear implementation sessions. The roadmap emphasizes removing 90% of unused code, leveraging standard library types, and creating clear separation between system settings and user preferences.

4. [Documentation Infrastructure Cleanup](./documentation-cleanup.md): This artifact defines Phase 0, which must be executed before any code changes to ensure development sessions have accurate references. It provides complete rewrites for PROJECT.md, README.md, technical documentation, CLAUDE.md, and development prompts that align with the new architecture. The documentation emphasizes the layered architecture, purpose-driven extraction, and clear separation between settings (HOW the tool operates) and configuration (WHAT the user wants). All documentation is updated to reflect standard library usage, role-based color organization, and the removal of speculative features, ensuring future development sessions proceed smoothly without confusion from outdated references.
