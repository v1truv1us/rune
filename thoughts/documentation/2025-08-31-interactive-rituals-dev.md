---
title: Interactive Rituals - Developer Notes
audience: developer
version: Phase 3 Complete
date: 2025-08-31
feature: Interactive Rituals with Tmux Integration
---

# Interactive Rituals - Developer Notes

This document provides technical implementation details for Interactive Rituals, architectural decisions, and guidance for extending the system.

## Architecture Overview

### Component Hierarchy

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Command Layer                       │
│  cmd/rune/main.go → internal/commands/start.go             │
└─────────────────────┬───────────────────────────────────────┘
                     │
┌─────────────────────▼───────────────────────────────────────┐
│                  Ritual Engine Core                        │
│  internal/rituals/engine.go                                │
│  • Command dispatch logic                                  │
│  • Interactive vs standard execution routing               │
│  • Environment filtering and security                      │
└─────────────┬───────────────────────────┬───────────────────┘
             │                           │
    ┌────────▼────────┐         ┌────────▼────────┐
    │  PTY Execution  │         │ Tmux Integration │
    │  pty.go         │         │  tmux.go        │
    │  • Direct TTY   │         │  • Session mgmt │
    │  • Raw mode     │         │  • Templates    │
    │  • I/O handling │         │  • Multi-pane   │
    └─────────────────┘         └─────────┬───────┘
                                         │
                              ┌──────────▼──────────┐
                              │   Tmux Client API   │
                              │ internal/tmux/      │
                              │ • gotmux wrapper    │
                              │ • Session lifecycle │
                              │ • Template rendering│
                              └─────────────────────┘
```

### Data Flow

```
1. User: rune start
      ↓
2. CLI: Parse command, load config
      ↓
3. Engine: ExecuteStartRituals(project)
      ↓
4. Engine: For each command → executeCommand()
      ↓
5. Decision: Interactive flag set?
   ├─ No:  executeStandardCommand()
   └─ Yes: executeInteractiveCommand()
             ↓
6. Decision: Tmux config present?
   ├─ Template: executeTmuxCommand() → CreateFromTemplate()
   ├─ Session:  executeTmuxCommand() → CreateSession()
   └─ Neither:  executePTYCommand()
             ↓
7. Result: Interactive session attached or PTY active
```

## Key Design Decisions

### 1. Graceful Degradation Strategy

**Decision**: Interactive commands fall back through multiple execution modes:
```
Tmux Template → Tmux Session → PTY → Standard Execution
```

**Rationale**:
- Maximizes compatibility across environments
- Maintains functionality when advanced features unavailable
- Provides consistent user experience

**Implementation**:
```go
func (e *Engine) executeInteractiveCommand(cmd config.Command, scope string) error {
    if cmd.TmuxTemplate != "" {
        return e.executeTmuxTemplate(cmd, scope)
    }
    if cmd.TmuxSession != "" {
        return e.executeTmuxSession(cmd, scope)  
    }
    return e.executePTYCommand(cmd)  // Universal fallback
}
```

### 2. Security Model Preservation

**Decision**: All interactive commands maintain Rune's security model.

**Implementation**:
- Environment filtering applied to all execution modes
- PTY and tmux inherit filtered environment
- No privilege escalation in interactive sessions

```go
// Environment filtering applied consistently
execCmd.Env = filterEnvironment(os.Environ())
```

### 3. Template Variable System

**Decision**: Use Go template-like syntax with limited scope.

**Syntax**: `{{.Variable}}`
**Variables**: Currently only `{{.Project}}`
**Rationale**: 
- Familiar to Go developers
- Easy to parse and replace
- Extensible for future variables

**Implementation**:
```go
func (c *Client) replaceVariables(command string, variables map[string]string) string {
    result := command
    for key, value := range variables {
        result = strings.ReplaceAll(result, "{{."+key+"}}", value)
    }
    return result
}
```

### 4. Session Lifecycle Management

**Decision**: Sessions are created per-ritual but persist beyond ritual execution.

**Rationale**:
- Supports long-running development sessions
- Allows detach/reattach workflow
- Maintains state between ritual invocations

**Implementation**:
- Sessions created with descriptive names
- Automatic attachment after creation
- Manual cleanup required (by design)

### 5. Error Handling Philosophy

**Decision**: Interactive command failures should be graceful and informative.

**Principles**:
1. Always provide fallback execution mode
2. Log errors with context for debugging
3. Continue ritual execution when commands are optional
4. Provide actionable error messages

## Component Details

### PTY Execution (`internal/rituals/pty.go`)

**Purpose**: Direct terminal allocation for interactive commands without tmux.

**Key Features**:
- Raw terminal mode handling
- Signal forwarding (SIGWINCH for resize)
- Bidirectional I/O coordination
- Proper cleanup on exit

**Critical Code Paths**:
```go
// Terminal setup
ptmx, err := pty.Start(execCmd)
oldState, err := term.MakeRaw(fd)

// I/O coordination
go func() { io.Copy(ptmx, os.Stdin) }()
go func() { io.Copy(os.Stdout, ptmx) }()

// Cleanup
defer term.Restore(fd, oldState)
defer ptmx.Close()
```

**Platform Considerations**:
- Unix PTY semantics (works on macOS, Linux)
- Windows requires WSL2 for full functionality
- Error handling for missing PTY support

### Tmux Integration (`internal/rituals/tmux.go`)

**Purpose**: Session management and template execution for tmux-based environments.

**Key Features**:
- Template variable expansion
- Session existence checking
- Command routing to sessions
- Template validation

**Template Processing Flow**:
```go
1. Load template from config
2. Expand variables in session name
3. Create session if not exists
4. For each window:
   - Create or rename window
   - Set layout if specified
   - Create panes and execute commands
5. Attach to session
```

**Session Management**:
```go
// Session lifecycle
if !c.SessionExists(sessionName) {
    c.CreateSession(sessionName)
}
return c.AttachSession(sessionName)
```

### Tmux Client (`internal/tmux/client.go`)

**Purpose**: Wrapper around gotmux library providing Rune-specific functionality.

**Abstractions**:
- Simplified session management
- Template-based session creation
- Variable expansion
- Error normalization

**gotmux Integration**:
- Uses `github.com/GianlucaP106/gotmux` for tmux control
- Wraps gotmux APIs with Rune-specific logic
- Handles version compatibility issues

**Template Rendering**:
```go
func (c *Client) CreateFromTemplate(template *config.TmuxTemplate, variables map[string]string) error {
    // 1. Create session
    session := tmux.NewSession(&gotmux.SessionOptions{Name: sessionName})
    
    // 2. Process windows
    for _, window := range template.Windows {
        // Create/configure window
        // Create panes with commands
        // Apply layout
    }
    
    // 3. Attach
    return session.Attach()
}
```

## Extension Points

### Adding New Variables

**Current**: Only `{{.Project}}` supported
**Extension**: Add new variables to expansion system

**Implementation**:
1. Add variable to `variables` map in execution context
2. Variables are automatically available in templates
3. Update documentation and examples

```go
// In executeTmuxCommand()
variables := map[string]string{
    "Project": scope,
    "User":    os.Getenv("USER"),        // New variable
    "Date":    time.Now().Format("2006-01-02"), // New variable
}
```

### Adding New Tmux Layouts

**Current**: Supports standard tmux layouts
**Extension**: Add custom layout handling

**Implementation**:
1. Extend layout string validation in config
2. Add layout mapping in tmux client
3. Update documentation

```go
// In client.go CreateFromTemplate()
var layout gotmux.WindowLayout
switch window.Layout {
case "custom-layout":
    layout = gotmux.WindowLayoutCustom  // New layout
    // Custom layout logic
}
```

### Adding New Execution Modes

**Current**: Standard, PTY, Tmux Template, Tmux Session
**Extension**: Add new interactive execution backends

**Pattern**:
1. Create new execution function in `rituals/`
2. Add detection logic in `executeInteractiveCommand()`
3. Add configuration fields to `config.Command`
4. Update validation and documentation

```go
// New execution mode
func (e *Engine) executeDockerCommand(cmd config.Command, scope string) error {
    // Implementation
}

// Integration in dispatch
if cmd.DockerContainer != "" {
    return e.executeDockerCommand(cmd, scope)
}
```

### Session Persistence Enhancement

**Current**: Sessions persist but no state tracking
**Extension**: Add session state management

**Implementation Areas**:
1. Session metadata storage
2. State restoration on rune restart  
3. Automatic cleanup policies
4. Session sharing between projects

**Code Structure**:
```go
type SessionManager struct {
    stateStore *SessionStateStore
    client     *tmux.Client
}

func (sm *SessionManager) SaveSessionState(sessionName, template, project string) error
func (sm *SessionManager) RestoreSession(sessionName string) error
func (sm *SessionManager) CleanupOrphanedSessions() error
```

## Testing Strategy

### Unit Tests

**Coverage Areas**:
- Template variable expansion
- Session lifecycle management  
- Error handling and fallbacks
- Configuration validation

**Test Utilities**:
```go
// Test helpers in internal/tmux/
func SetupMockTmux() *MockTmuxClient
func CreateTestTemplate() config.TmuxTemplate
func ValidateSessionExists(name string) bool
```

### Integration Tests

**Test Scenarios**:
- Full ritual execution with interactive commands
- Multi-mode fallback behavior
- Template processing with real tmux
- Cross-platform compatibility

**Test Organization**:
```
internal/rituals/integration_test.go
internal/tmux/integration_test.go
```

### Manual Testing

**Testing Checklist**:
- [ ] PTY commands work in terminal
- [ ] Tmux sessions create correctly
- [ ] Template layouts apply properly
- [ ] Variable expansion works
- [ ] Fallback modes activate appropriately
- [ ] Error messages are helpful
- [ ] Session cleanup works

## Performance Considerations

### Session Creation Overhead

**Measurement**: Session creation takes 0.5-2s depending on complexity
**Optimization**: 
- Cache tmux client connections
- Minimize template processing
- Parallel pane setup where possible

### Memory Usage

**Current**: ~5MB per active session (mostly tmux overhead)
**Monitoring**: Track session count and memory usage
**Limits**: Consider session count limits for resource management

### I/O Performance

**PTY**: Direct I/O coordination, minimal overhead
**Tmux**: Additional layer but acceptable latency (<100ms)
**Optimization**: Buffering strategies for high-throughput scenarios

## Security Considerations

### Environment Isolation

**Implementation**: `filterEnvironment()` removes sensitive variables
**Coverage**: AWS keys, tokens, passwords, secrets
**Testing**: Verify no secrets leak to interactive sessions

**Code**:
```go
func filterEnvironment(env []string) []string {
    sensitiveSubstrings := []string{
        "AWS_", "SECRET", "TOKEN", "KEY", "PASSWORD", // ...
    }
    // Filter logic
}
```

### Process Security

**Concerns**: Interactive sessions may run with different privileges
**Mitigation**: 
- No privilege escalation in ritual execution
- Same user context as parent rune process
- Proper signal handling for cleanup

### Session Isolation

**Current**: Sessions are user-scoped (tmux default)
**Security**: Sessions visible to same user only
**Future**: Consider project-based session isolation

## Maintenance Guidelines

### Dependency Management

**gotmux**: Pin to stable version, test updates thoroughly
**creack/pty**: Stable library, updates rare
**golang.org/x/term**: Part of Go extended libs, follows Go release cycle

### Error Monitoring

**Key Metrics**:
- Interactive command failure rates
- Tmux availability across environments
- PTY allocation failures
- Session creation latency

### Backwards Compatibility

**Commitment**: Existing ritual configurations continue working
**Testing**: Comprehensive backwards compatibility test suite
**Migration**: Provide migration guides for major changes

### Documentation Maintenance

**Areas Requiring Updates**:
- New tmux layout support
- Additional template variables  
- Platform-specific installation guides
- Troubleshooting for new error cases

## Future Enhancements

### Planned Features (Phase 4+)

1. **Session State Persistence**
   - Save/restore session configuration
   - Automatic session recovery
   - Session sharing between users

2. **Enhanced Template System**
   - Conditional template logic
   - Template inheritance
   - Dynamic pane creation

3. **Integration Improvements**
   - IDE integration hooks
   - Container environment support
   - Cloud development environment support

### Research Areas

1. **Alternative Terminal Backends**
   - Screen support as tmux alternative
   - Windows Terminal integration
   - Browser-based terminal support

2. **Performance Optimization**
   - Session pooling
   - Template compilation caching
   - Async session initialization

3. **Security Enhancements**
   - Session encryption
   - Fine-grained permission control
   - Audit logging for interactive sessions

This developer documentation provides comprehensive technical coverage for maintaining, extending, and optimizing the Interactive Rituals system.