---
date: 2025-08-28T18:15:00-08:00
researcher: Claude Code (Sonnet 4)
git_commit: b1900184
branch: main
repository: rune
topic: "Tmux Terminal Integration for Interactive Scripts in Ritual System"
tags: [research, codebase, tmux, rituals, interactive, terminal, automation]
status: complete
last_updated: 2025-08-28
last_updated_by: Claude Code (Sonnet 4)
---

## Ticket Synopsis

Investigate how tmux terminals work with interactive scripts when rituals run interactive commands in the Rune project. Research existing tmux integration patterns and terminal handling for ritual command execution.

## Summary

**Key Finding**: Rune has a sophisticated **auto-tmux functionality** built into shell scripts that automatically creates/attaches tmux sessions for interactive script execution, but the **ritual engine itself does not have native tmux support** for interactive commands. The ritual system uses basic `exec.Command` with 30-second timeouts, making it unsuitable for long-running interactive processes.

**Two Distinct Systems**:
1. **Script-level tmux integration**: Automatic tmux session management for shell scripts
2. **Ritual engine limitations**: Basic command execution without terminal/tmux awareness

## Detailed Findings

### Auto-Tmux Functionality in Shell Scripts

- **Implementation**: All major scripts have auto-tmux functionality at lines 3-23
- **Pattern**: Detects if running inside tmux (`$TMUX`), creates/attaches sessions automatically
- **Session Management**: Creates named sessions (e.g., `test-telemetry-build`) with proper isolation
- **Fallback Behavior**: Falls back to regular execution if tmux unavailable, requires interactive TTY
- **Skip Mechanism**: `SKIP_AUTO_TMUX` environment variable bypasses tmux integration

**Auto-tmux Code Pattern** (`scripts/test-telemetry-build.sh:3-23`):
```bash
SESSION_NAME="test-telemetry-build"
if [ -z "$TMUX" ] && [ -z "$SKIP_AUTO_TMUX" ]; then
  if command -v tmux >/dev/null 2>&1; then
    if tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
      exec tmux attach-session -t "$SESSION_NAME"
    else
      exec tmux new-session -s "$SESSION_NAME" "$0" "$@"
    fi
  fi
fi
```

### Ritual Engine Command Execution

- **Basic Implementation**: Uses `exec.CommandContext` with 30-second timeout (`internal/rituals/engine.go:144`)
- **No TTY Support**: No terminal allocation for interactive commands
- **Output Handling**: `CombinedOutput()` captures stdout/stderr, not suitable for interactive input
- **Background Mode**: Supports background processes but no session management
- **Environment Filtering**: Removes sensitive environment variables for security

**Limitations for Interactive Commands**:
- 30-second timeout makes long-running interactive processes impossible
- No stdin/terminal allocation prevents user input
- `CombinedOutput()` method blocks interactive session flow

### Tmux Usage in Configuration Examples

**Session Capture Example** (`docs-site/src/pages/examples/index.md:392`):
```yaml
- name: "Save session"
  command: "tmux capture-session -p > ~/logs/session-$(date +%Y%m%d).log"
  optional: true
```

**Demonstrates**: Tmux commands work in rituals when used as non-interactive utilities (capture, kill, status), but not for interactive session management.

### Terminal/TTY Handling Patterns

- **TTY Detection**: Scripts check `[ -t 0 ]` for interactive terminal detection
- **Error Messages**: Proper fallback messages when terminal unavailable
- **Cross-platform**: Works on macOS, Linux, Windows with consistent behavior
- **Test Coverage**: Comprehensive test suite in `scripts/test-auto-tmux.sh` and `scripts/test-auto-tmux-simple.sh`

**TTY Detection Pattern**:
```bash
if [ ! -t 0 ]; then
  echo "This script must be run in an interactive terminal (tmux not available)."
  exit 1
fi
```

## Code References

- `scripts/test-telemetry-build.sh:3-23` - Auto-tmux implementation pattern
- `internal/rituals/engine.go:144` - Ritual command execution with timeout
- `internal/rituals/engine.go:160` - `CombinedOutput()` usage preventing interactivity
- `scripts/test-auto-tmux.sh:75-104` - Tmux session creation test function
- `docs-site/src/pages/examples/index.md:392` - Tmux usage in ritual configuration

## Architecture Insights

**Two-Layer Design**: Rune separates script-level terminal management (tmux) from application-level command execution (rituals), creating a clean but limited architecture.

**Security-First Approach**: Environment filtering (`internal/rituals/engine.go:14-63`) prevents secret leakage to subprocesses, prioritizing security over flexibility.

**Test-Driven Development**: Comprehensive test coverage for tmux functionality demonstrates production-ready terminal handling at the script level.

**Platform Consistency**: Auto-tmux functionality works consistently across platforms with appropriate fallbacks and error handling.

## Implementation Gap Analysis

### Current Limitations

1. **No Interactive Ritual Support**: Ritual engine cannot handle interactive commands (editors, prompts, long-running processes)
2. **Fixed Timeout**: 30-second limit prevents meaningful interactive sessions  
3. **No TTY Allocation**: Cannot pass terminal control to child processes
4. **Session Isolation**: No mechanism for ritual commands to share tmux sessions

### Potential Enhancement Paths

**Option 1: Enhanced Ritual Engine**
- Add `interactive: true` flag to command configuration
- Implement TTY allocation and pass-through for interactive commands
- Remove timeout restrictions for interactive commands
- Add tmux session management to ritual engine

**Option 2: Script Integration Pattern**
- Generate tmux-enabled scripts from ritual configurations
- Use existing auto-tmux patterns for ritual command execution
- Maintain separation between script-level and application-level concerns

## Related Research

This research complements the previous Context Switching Assistance PRD analysis by identifying terminal interaction limitations that may affect mode switching and interactive planning features.

## Open Questions

1. **Interactive Mode Design**: How should the ritual engine handle interactive vs. non-interactive commands?
2. **Session Management**: Should rituals create persistent tmux sessions for project-specific work environments?
3. **Background Process Integration**: How to integrate long-running interactive processes with Rune's session tracking?
4. **User Experience**: What's the optimal way to switch between Rune CLI and interactive ritual sessions?

## Conclusion

While Rune has excellent tmux integration at the shell script level with auto-session management, the ritual engine lacks support for interactive commands. This creates a clear architectural boundary but limits the ability to automate interactive development workflows through rituals. Any enhancement for interactive ritual support would require fundamental changes to the command execution architecture.