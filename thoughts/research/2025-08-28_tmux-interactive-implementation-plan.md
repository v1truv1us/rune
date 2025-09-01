---
date: 2025-08-28T19:30:00-08:00
researcher: Claude Code (Sonnet 4)
git_commit: b1900184
branch: main
repository: rune
topic: "Tmux Interactive Implementation Plan - Deep Research for Ritual System Enhancement"
tags: [research, implementation, tmux, interactive, rituals, terminal, automation, golang]
status: complete
last_updated: 2025-08-28
last_updated_by: Claude Code (Sonnet 4)
---

## Ticket Synopsis

Conduct deep research into tmux interactive capabilities and develop a thorough implementation plan for enhancing Rune's ritual system with interactive session management. Research tmux scripting, Go libraries, terminal handling patterns, and similar CLI automation tools to inform architectural decisions.

## Executive Summary

**Research Conclusion**: Tmux provides excellent interactive capabilities that can significantly enhance Rune's ritual system. The research reveals clear implementation paths using existing Go libraries, established patterns from similar tools, and architectural approaches that maintain Rune's security-first design while adding powerful interactive session management.

**Key Finding**: A phased implementation approach leveraging `github.com/creack/pty` for TTY handling, `github.com/GianlucaP106/gotmux` for tmux control, and Rune's existing auto-tmux patterns can deliver sophisticated interactive ritual capabilities without compromising the current architecture.

## Comprehensive Research Findings

### 1. Tmux Session Management and Scripting Capabilities

**Advanced Session Orchestration**
- **Multi-window/pane layouts**: Programmatic creation of complex development environments
- **send-keys automation**: Interactive command sequences with real-time process control
- **Session templates**: YAML-based configuration similar to tmuxinator/tmuxp patterns
- **State persistence**: Session restoration with tmux-resurrect integration
- **Control mode**: Protocol-level automation through stdin commands

**Production-Ready Patterns**:
```bash
# Automated session with complex layout
SESSION="rune-dev-project"
tmux new-session -d -s $SESSION
tmux split-window -h -t $SESSION:0
tmux split-window -v -t $SESSION:0.1
tmux send-keys -t $SESSION:0.0 'vim src/main.go' C-m
tmux send-keys -t $SESSION:0.1 'go run main.go' C-m  
tmux send-keys -t $SESSION:0.2 'git log --oneline' C-m
```

**Key Benefits for Rune**:
- Session isolation per project
- Persistent workspace across reboots  
- Multi-service orchestration
- Interactive process management
- Development environment consistency

### 2. Go Tmux Libraries Analysis

**Primary Library: gotmux (`github.com/GianlucaP106/gotmux`)**
- **Status**: Actively maintained, updated 2024
- **Features**: Session/window/pane management, command execution, attachment control
- **API Design**: Clean Go interfaces with error handling
- **Integration**: Works through tmux CLI with structured control

**Alternative Libraries**:
- `github.com/jubnzv/go-tmux` - Session/window/pane management
- `github.com/wricardo/gomux` - Session creation wrapper
- `github.com/philipgraf/libtmux-go` - Go port of Python libtmux

**Recommended Approach**: Use `gotmux` as primary library with fallback to direct tmux CLI commands for maximum compatibility.

### 3. Current Ritual Engine Extension Points

**Architecture Analysis** (`internal/rituals/engine.go`):
- **Current Structure**: Simple command execution with 30-second timeout
- **Extension Points**: `executeCommand()` method can be enhanced with command type detection
- **Security Model**: Environment filtering maintains security with interactive commands
- **Background Support**: Existing `cmd.Background` flag provides foundation for session management

**Critical Limitations**:
```go
// Line 160: Blocks interactive sessions
output, err := execCmd.CombinedOutput()

// Line 134: Prevents long-running interactive processes  
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
```

**Configuration Enhancement Required**:
```go
// Current Command struct (config.go:60-65)
type Command struct {
    Name       string `yaml:"name"`
    Command    string `yaml:"command"`  
    Optional   bool   `yaml:"optional"`
    Background bool   `yaml:"background"`
    // New fields needed:
    // Interactive bool   `yaml:"interactive"`
    // TmuxSession string `yaml:"tmux_session,omitempty"`
}
```

### 4. Interactive Command Patterns from Similar Tools

**Task (Taskfile) Patterns**:
- **Interactive prompts**: `--yes` flag for automation-friendly execution
- **Watch mode**: Live file change detection with re-execution
- **Dynamic configuration**: Runtime task generation and execution
- **Shell integration**: Native sh/bash compatibility across platforms

**Tmuxinator/tmuxp Session Management**:
- **YAML configuration**: Declarative session definitions
- **Layout templates**: Reusable window/pane configurations
- **Startup scripts**: Pre-session initialization commands
- **Project isolation**: Session-per-project patterns

**Key Architectural Insights**:
- **Declarative configuration** preferred over imperative scripting
- **Template-based approach** for session reusability
- **Graceful fallback** when interactive features unavailable
- **Background/foreground coordination** for complex workflows

### 5. Terminal/TTY Handling in Go Applications

**Core Library: `github.com/creack/pty`** (Updated Oct 2024)
- **Pseudoterminal creation**: Full TTY allocation for interactive processes
- **Standard integration**: Works with `exec.Cmd` for seamless command execution
- **Cross-platform support**: Linux, macOS, Windows compatibility

**Implementation Pattern**:
```go
import "github.com/creack/pty"

// Start command with PTY
cmd := exec.Command("bash", "-c", ritualCommand)
ptmx, err := pty.Start(cmd)
if err != nil {
    return err
}
defer ptmx.Close()

// Enable raw mode for terminal interaction  
oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
if err != nil {
    return err
}
defer term.Restore(int(os.Stdin.Fd()), oldState)

// Coordinate I/O between terminal and process
go io.Copy(ptmx, os.Stdin)
io.Copy(os.Stdout, ptmx)
```

**Security Considerations**:
- **PTY allocation** requires careful permission handling
- **Environment filtering** remains critical for security
- **Session isolation** prevents cross-project contamination

### 6. Background Process Management Patterns

**Go Context Patterns for Process Coordination**:
- **context.WithCancel**: Hierarchical process termination
- **context.WithTimeout**: Selective timeout control for interactive vs automated commands
- **WaitGroup coordination**: Multi-process synchronization
- **Channel communication**: Process status and error propagation

**Modern 2024 Patterns**:
```go
// Context-aware background process management
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

cmd := exec.CommandContext(ctx, command, args...)
if interactive {
    // Use PTY for interactive commands (no timeout)
    return executeWithPTY(cmd)
} else {
    // Use existing timeout logic for automation
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    return executeWithTimeout(cmd, ctx)
}
```

## Detailed Implementation Plan

### Phase 1: Enhanced Command Types and Configuration

**Configuration Extensions**:
```yaml
# Enhanced command configuration
rituals:
  start:
    global:
      - name: "Setup Development Environment"
        command: "tmux-session-template"
        interactive: true
        tmux_session: "dev-{{.Project}}"
        tmux_template: "development-layout"
    
    templates:
      development-layout:
        session_name: "dev-{{.Project}}"
        windows:
          - name: "editor"
            layout: "main-horizontal"
            panes:
              - command: "vim ."
              - command: "git status"
          - name: "server"  
            panes:
              - command: "npm run dev"
              - command: "tail -f logs/app.log"
```

**Go Implementation**:
```go
// Enhanced Command struct
type Command struct {
    Name         string       `yaml:"name" mapstructure:"name"`
    Command      string       `yaml:"command" mapstructure:"command"`
    Optional     bool         `yaml:"optional" mapstructure:"optional"`
    Background   bool         `yaml:"background" mapstructure:"background"`
    Interactive  bool         `yaml:"interactive" mapstructure:"interactive"`
    TmuxSession  string       `yaml:"tmux_session,omitempty" mapstructure:"tmux_session"`
    TmuxTemplate string       `yaml:"tmux_template,omitempty" mapstructure:"tmux_template"`
}

// Template configuration
type TmuxTemplate struct {
    SessionName string       `yaml:"session_name" mapstructure:"session_name"`
    Windows     []TmuxWindow `yaml:"windows" mapstructure:"windows"`
}

type TmuxWindow struct {
    Name     string   `yaml:"name" mapstructure:"name"`
    Layout   string   `yaml:"layout,omitempty" mapstructure:"layout"`
    Panes    []string `yaml:"panes" mapstructure:"panes"`
}
```

### Phase 2: Interactive Execution Engine

**Enhanced Engine Architecture**:
```go
type Engine struct {
    config       *config.Config
    tmuxClient   *gotmux.TmuxClient
    ptySupport   bool
    activeSession map[string]*TmuxSession
}

func (e *Engine) executeCommand(cmd config.Command, scope string) error {
    if cmd.Interactive {
        return e.executeInteractiveCommand(cmd, scope)
    }
    // Existing non-interactive path
    return e.executeStandardCommand(cmd, scope)
}

func (e *Engine) executeInteractiveCommand(cmd config.Command, scope string) error {
    if cmd.TmuxTemplate != "" {
        return e.executeTemplateCommand(cmd, scope)
    } else if cmd.TmuxSession != "" {
        return e.executeTmuxCommand(cmd, scope)  
    } else {
        return e.executePTYCommand(cmd, scope)
    }
}
```

**Tmux Integration**:
```go
import "github.com/GianlucaP106/gotmux/gotmux"

func (e *Engine) executeTmuxCommand(cmd config.Command, scope string) error {
    sessionName := expandTemplate(cmd.TmuxSession, map[string]string{
        "Project": scope,
    })
    
    // Check if session exists
    sessions, err := e.tmuxClient.ListSessions()
    if err != nil {
        return err
    }
    
    var session *gotmux.Session
    for _, s := range sessions {
        if s.Name == sessionName {
            session = s
            break
        }
    }
    
    // Create session if it doesn't exist
    if session == nil {
        session, err = e.tmuxClient.NewSession(gotmux.SessionOptions{
            Name: sessionName,
        })
        if err != nil {
            return err
        }
    }
    
    // Send command to session
    return session.SendKeys(cmd.Command, true)
}
```

### Phase 3: Template System Implementation

**Template Engine**:
```go
func (e *Engine) executeTemplateCommand(cmd config.Command, scope string) error {
    template := e.config.Rituals.Templates[cmd.TmuxTemplate]
    if template == nil {
        return fmt.Errorf("template '%s' not found", cmd.TmuxTemplate)
    }
    
    sessionName := expandTemplate(template.SessionName, map[string]string{
        "Project": scope,
    })
    
    // Create session with template
    session, err := e.tmuxClient.NewSession(gotmux.SessionOptions{
        Name: sessionName,
    })
    if err != nil {
        return err
    }
    
    // Apply template configuration
    for _, window := range template.Windows {
        win, err := session.NewWindow(gotmux.WindowOptions{
            Name: window.Name,
        })
        if err != nil {
            continue
        }
        
        // Create panes and send commands
        for i, paneCmd := range window.Panes {
            if i > 0 {
                // Split window for additional panes
                win.SplitWindow(gotmux.PaneOptions{})
            }
            if paneCmd != "" {
                win.SendKeys(paneCmd, true)
            }
        }
    }
    
    return session.AttachSession()
}
```

### Phase 4: PTY Integration for Non-Tmux Interactive Commands

**Direct PTY Support**:
```go
import "github.com/creack/pty"

func (e *Engine) executePTYCommand(cmd config.Command, scope string) error {
    execCmd := exec.Command("sh", "-c", cmd.Command)
    execCmd.Env = filterEnvironment(os.Environ())
    
    // Start command with PTY
    ptmx, err := pty.Start(execCmd)
    if err != nil {
        return err
    }
    defer ptmx.Close()
    
    // Set raw mode for terminal interaction
    oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
    if err != nil {
        return err
    }
    defer term.Restore(int(os.Stdin.Fd()), oldState)
    
    // Coordinate I/O
    go func() { io.Copy(ptmx, os.Stdin) }()
    io.Copy(os.Stdout, ptmx)
    
    return execCmd.Wait()
}
```

## Integration with Existing Architecture

### Maintaining Backwards Compatibility

**Configuration Migration**:
- Existing ritual configurations continue working unchanged
- Interactive features are opt-in additions
- Graceful fallback when tmux/PTY unavailable

**Security Preservation**:
- Environment variable filtering remains active for all command types
- Session isolation prevents cross-project contamination  
- PTY allocation uses existing permission model

### Building on Auto-Tmux Pattern

**Leveraging Existing Infrastructure**:
```go
func (e *Engine) executeInteractiveCommand(cmd config.Command, scope string) error {
    // Check if already in tmux (similar to existing auto-tmux logic)
    if os.Getenv("TMUX") != "" {
        // Already in tmux, execute directly
        return e.executeTmuxCommand(cmd, scope)
    }
    
    // Not in tmux, check availability and create session
    if !e.tmuxAvailable() {
        return e.executePTYCommand(cmd, scope)
    }
    
    return e.executeTemplateCommand(cmd, scope)
}
```

## Benefits and Use Cases

### Enhanced Developer Productivity

**Multi-Service Development Environment**:
```yaml
rituals:
  start:
    per_project:
      web-app:
        - name: "Full Stack Development Environment"
          interactive: true
          tmux_template: "fullstack-dev"
  
  templates:
    fullstack-dev:
      session_name: "dev-{{.Project}}"
      windows:
        - name: "backend"
          panes: ["cd backend && go run main.go", "cd backend && go test ./..."]
        - name: "frontend"  
          panes: ["cd frontend && npm run dev", "cd frontend && npm run test:watch"]
        - name: "database"
          panes: ["docker-compose up postgres", "psql -d app_db"]
        - name: "monitoring"
          panes: ["htop", "tail -f logs/app.log"]
```

**Interactive Configuration Workflows**:
```yaml
rituals:
  start:
    global:
      - name: "Project Setup Wizard"
        command: "./setup-wizard.sh"
        interactive: true
        tmux_session: "setup-{{.Project}}"
```

### Advanced Automation Capabilities

**Session State Management**:
- Persistent development environments across machine restarts
- Project-specific session restoration
- Multi-developer session sharing capabilities
- Integration with existing time tracking

**Process Orchestration**:
- Interactive debugging sessions within automated workflows
- Long-running development server management
- Real-time log monitoring and analysis
- Interactive git workflows and conflict resolution

## Technical Considerations and Risks

### Implementation Complexity

**Medium Risk Areas**:
- PTY handling across different platforms (Linux, macOS, Windows)
- Tmux version compatibility and feature detection
- Session lifecycle management and cleanup
- Error handling and recovery in interactive contexts

**Mitigation Strategies**:
- Comprehensive platform testing with fallback mechanisms
- Version detection with graceful degradation
- Proper resource cleanup with defer patterns
- User-friendly error messages with troubleshooting guidance

### Performance Implications

**Resource Usage**:
- Tmux sessions consume additional memory (~1-5MB per session)
- PTY allocation requires file descriptor management
- Background process coordination increases CPU usage

**Optimization Approaches**:
- Session pooling and reuse for similar workflows
- Lazy session creation and cleanup
- Resource monitoring and limits
- Efficient I/O handling with buffered operations

### Security Considerations

**Maintained Security Model**:
- Environment variable filtering applies to all command types
- Session isolation prevents information leakage
- PTY permissions follow existing process security model
- Interactive sessions inherit Rune's security posture

**Additional Considerations**:
- Session access control and authentication
- Command injection prevention in template expansion
- Audit logging for interactive sessions
- Secure session sharing mechanisms

## Success Metrics and Validation

### Implementation Validation

**Functional Testing**:
- Cross-platform session creation and management
- Template-based environment setup
- Interactive command execution with proper TTY handling
- Background process coordination and lifecycle management

**Performance Benchmarks**:
- Session creation/destruction timing
- Memory usage with multiple active sessions
- I/O throughput for interactive operations
- Resource cleanup verification

**User Experience Validation**:
- Configuration complexity and learning curve
- Error message clarity and troubleshooting guidance
- Integration with existing Rune workflows
- Developer productivity improvements

## Conclusion and Recommendations

The research demonstrates that tmux provides robust interactive capabilities that can significantly enhance Rune's ritual system. The proposed phased implementation approach leverages proven Go libraries, established patterns from similar tools, and Rune's existing auto-tmux infrastructure to deliver sophisticated interactive session management.

**Recommended Implementation Priority**:
1. **Phase 1**: Enhanced command configuration and basic tmux integration
2. **Phase 2**: Interactive execution engine with gotmux library
3. **Phase 3**: Template system for complex session management
4. **Phase 4**: Direct PTY support for non-tmux interactive commands

**Key Success Factors**:
- Building incrementally on existing patterns maintains stability
- Comprehensive testing across platforms ensures reliability
- Clear documentation and examples drive user adoption
- Performance monitoring prevents resource degradation

The implementation will transform Rune from a workflow automator into a comprehensive interactive development environment manager while maintaining its security-first design and reliability standards.

**Next Steps**: Proceed with Phase 1 implementation, starting with configuration enhancements and basic tmux integration to validate the architectural approach before advancing to more complex interactive features.