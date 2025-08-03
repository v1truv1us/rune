# AGENTS.md - Global Development Guidelines

## Specialized Agents Directory

Rune includes specialized agents organized by department to assist with various aspects of software development:

### Product Strategy (`agents/product-strategy/`)
- **Product Strategist** - Analyzes features and provides strategic insights
- **Growth Engineer** - Identifies user hooks and builds viral loops
- **User Researcher** - Analyzes user flows and identifies drop-off points
- **Revenue Optimizer** - Identifies monetization opportunities
- **Market Analyst** - Compares features to competitors and provides market intelligence

### Development (`agents/development/`)
- **System Architect** - Transforms codebases and designs scalable architectures
- **API Builder** - Creates developer-friendly APIs with proper documentation
- **Database Expert** - Optimizes database queries and designs efficient data models
- **Integration Master** - Connects services and builds seamless integrations
- **Mobile Optimizer** - Enhances mobile experience and platform optimization
- **Performance Engineer** - Improves app speed, conducts load testing, and resolves performance bottlenecks
- **Accessibility Pro** - Ensures WCAG compliance and accessibility standards

### Design & UX (`agents/design-ux/`)
- **UX Optimizer** - Simplifies user flows, enhances user experience, and optimizes conversion rates
- **UI Polisher** - Adds premium design elements and visual enhancements
- **Content Writer** - Improves app messaging and creates compelling copy
- **Design System Builder** - Creates consistent component libraries and design systems

### Quality & Testing (`agents/quality-testing/`)
- **Test Generator** - Creates comprehensive tests and ensures code quality
- **Security Scanner** - Identifies vulnerabilities and implements security best practices
- **Code Reviewer** - Provides engineering-level code feedback, refactoring, and quality improvement

### Operations (`agents/operations/`)
- **Deployment Wizard** - Sets up CI/CD pipelines and deployment automation
- **Infrastructure Builder** - Designs scalable cloud architecture and infrastructure as code
- **Monitoring Expert** - Implements operational monitoring, alerting, and observability infrastructure
- **Release Manager** - Handles smooth deployments and release coordination
- **Cost Optimizer** - Reduces cloud expenses and optimizes resource utilization

### Business & Analytics (`agents/business-analytics/`)
- **Analytics Engineer** - Tracks user behavior and implements analytics solutions
- **Email Automator** - Creates effective email flows and marketing automation
- **Support Builder** - Reduces support tickets and implements self-service tools
- **Compliance Expert** - Handles regulatory requirements and compliance standards
- **SEO Master** - Improves search engine optimization and organic visibility
- **Community Features** - Adds user engagement tools and community systems

### AI & Innovation (`agents/ai-innovation/`)
- **AI Integration Expert** - Adds AI features and machine learning capabilities
- **Automation Builder** - Creates workflow automations and process optimization
- **Innovation Lab** - Experiments with cutting-edge tech and emerging technologies

## Build Commands
- **TypeScript/Node**: `bun build`, `bun run typecheck`, `bun test`, `bun run lint`
- **Go**: `go build`, `go test`, `go vet`, `golangci-lint run`, `gofmt`
- **Python**: `pytest`, `ruff check`, `mypy`, `black --check`
- **Rust**: `cargo build`, `cargo test`, `cargo clippy`, `cargo fmt`
- **Single test**: `bun test -- <test-name>` or `go test <test-name>`

## Code Style Guidelines

### Test-Driven Development (TDD) - REQUIRED
- Write tests BEFORE implementing features
- Follow AAA pattern: Arrange, Act, Assert
- Use descriptive test names and group related tests
- Test both success and error paths, including edge cases
- Target 80-90% code coverage for critical paths
- Mock external dependencies appropriately

### Security-First Development
- Validate ALL user inputs before processing
- Use parameterized queries, never string concatenation
- Store secrets in environment variables, never in code
- Implement proper authentication/authorization at boundaries
- Use secure defaults and principle of least privilege
- Never expose stack traces or sensitive info in production
- Regular dependency security audits

### Code Simplicity
- Readability first - use clear, descriptive names
- Keep functions under 30 lines, avoid deep nesting (max 3 levels)
- Use early returns to reduce complexity
- Follow Single Responsibility Principle
- Start with simplest solution, add complexity only when justified
- DRY (Don't Repeat Yourself) and YAGNI (You Aren't Gonna Need It)

### Language-Specific Conventions
- **TypeScript**: Use strict mode, prefer `const`, explicit types for public APIs
- **Go**: Use gofmt, follow effective Go guidelines, handle errors explicitly
- **Python**: Follow PEP 8, use type hints, prefer f-strings
- **Rust**: snake_case variables, PascalCase types, explicit error handling with `Result`
- **Imports**: Group std/external imports first, then local imports
- **Formatting**: Use project's formatter (Prettier/rustfmt/gofmt) - no manual formatting

### Error Handling
- Use specific error types per module with descriptive names
- Never expose stack traces or sensitive info in production errors
- Log detailed errors server-side, return generic messages to users
- Handle ALL error cases explicitly
- Implement proper error boundaries and recovery mechanisms

## Package Management - CRITICAL
- Use ONE package manager consistently per project
- **JavaScript/Node**: Bun (preferred), pnpm, npm, or yarn
- **Go**: go mod with go.mod and go.sum
- **Python**: pip with requirements.txt, poetry, or conda
- **Rust**: cargo with Cargo.toml and Cargo.lock
- Maintain lock files in version control
- Regular dependency audits for vulnerabilities
- Use exact versions for production dependencies

## File Organization & Project Structure

### Required Files
- **README.md**: Setup instructions, usage examples, contribution guidelines
- **TODO.md**: Task tracking with priority levels
- **CHANGELOG.md**: Version history and changes
- **LICENSE**: Appropriate license file
- **.editorconfig**: Consistent coding styles across editors

### Directory Structure
- Follow language/framework conventions
- Separate concerns: `/src`, `/tests`, `/docs`, `/config`
- Group related functionality together
- Use meaningful directory and file names
- Keep modules focused and dependencies clear

## Development Tools Setup

### Required Linting & Formatting
- **JavaScript/TypeScript**: ESLint + Prettier + eslint-plugin-jsx-a11y
- **Go**: golangci-lint + gofmt
- **Python**: ruff (linting) + black (formatting) + mypy (type checking)
- **Rust**: clippy (linting) + rustfmt (formatting)
- Configure pre-commit hooks for automatic formatting

### IDE Configuration
- EditorConfig for consistent styles
- Language server setup (LSP)
- Debugger configuration
- Extension recommendations file

## Accessibility Standards - MANDATORY

### CLI Applications
- Clear, descriptive help text and error messages
- Screen reader friendly terminal output
- Semantic formatting that degrades gracefully
- Keyboard navigation for interactive elements
- Support for high contrast themes

### Web Applications (WCAG 2.2 AA)
- Minimum 4.5:1 contrast ratio for text
- Keyboard navigation for all interactive elements
- Semantic HTML with proper ARIA labels
- Alt text for images and media
- Minimum 44px touch targets
- Screen reader compatibility testing

### Mobile Applications
- Platform accessibility guidelines (VoiceOver, TalkBack)
- Dynamic text sizing support
- Minimum touch target sizes (44pt iOS, 48dp Android)
- Meaningful content descriptions

## UI/Styling Standards

### Web Development
- **Use REM for sizing**: Ensures consistent scaling (base: 16px)
- **Use HSLA for colors**: Better control over hue, saturation, lightness, alpha
- **CSS Variables**: Store design tokens in :root
- **Responsive Design**: Mobile-first approach with relative units
- **Semantic class names**: Follow BEM or similar methodology

### Example CSS Variables
```css
:root {
  --spacing-xs: 0.25rem; /* 4px */
  --spacing-sm: 0.5rem;  /* 8px */
  --spacing-md: 1rem;    /* 16px */
  --spacing-lg: 1.5rem;  /* 24px */
  --spacing-xl: 2rem;    /* 32px */
  
  --primary-color: hsla(210, 100%, 50%, 1);
  --text-color: hsla(0, 0%, 20%, 1);
  --background-color: hsla(0, 0%, 100%, 1);
}
```

## Performance Guidelines
- Profile before optimizing
- Implement proper caching strategies
- Minimize bundle sizes and lazy load resources
- Use appropriate data structures
- Monitor Core Web Vitals for web apps
- Load testing for critical paths

## Documentation Requirements
- Document public APIs with examples
- Keep README current with setup instructions
- Document architectural decisions
- Include security considerations
- Maintain changelog for version tracking
- Comment complex algorithms only

## Common Commands by Language

### TypeScript/Node.js (Bun)
```bash
bun install          # Install dependencies
bun run build        # Build project
bun test             # Run tests
bun run lint         # Run linter
bun run typecheck    # Type checking
bun run dev          # Development server
bun run format       # Format code
```

### Go
```bash
go mod tidy          # Clean dependencies
go build             # Build project
go test ./...        # Run tests
golangci-lint run    # Linting
gofmt -w .          # Formatting
go vet              # Static analysis
```

### Python
```bash
pip install -r requirements.txt  # Install dependencies
pytest                           # Run tests
ruff check                       # Linting
black .                          # Formatting
mypy .                          # Type checking
```

### Rust
```bash
cargo build          # Build project
cargo test           # Run tests
cargo clippy         # Linting
cargo fmt            # Formatting
cargo run            # Run project
```

## Review Checklist
- [ ] Tests written first and passing
- [ ] Security boundaries validated
- [ ] Code is readable and simple
- [ ] Error handling is comprehensive
- [ ] No secrets in code
- [ ] Dependencies are justified
- [ ] Accessibility requirements met
- [ ] Performance considerations addressed
- [ ] Documentation updated

## Git Workflow
- Clear, descriptive commit messages
- Keep changes focused and atomic
- Update tests alongside code changes
- Run linting and tests before commits
- Use semantic versioning for releases
- Document breaking changes clearly

## Emergency Procedures
- **Security Incident**: Immediately revoke compromised credentials, assess impact
- **Production Bug**: Implement hotfix, deploy, then address root cause
- **Dependency Vulnerability**: Update immediately, test thoroughly
- **Performance Degradation**: Identify bottleneck, implement fix, monitor

## Project-Specific Guidelines

### Rune CLI Development
- **Primary Language**: Go 1.21+
- **CLI Framework**: Cobra for command structure
- **Configuration**: Viper for YAML parsing
- **Testing**: Go testing + Testify
- **Performance**: <200ms startup time, <50MB memory
- **Security**: Command sandboxing, OS keychain for credentials
- **Cross-Platform**: macOS, Linux, Windows support

### Frontend Components (if applicable)
- **Runtime**: Bun for package management and execution
- **Framework**: To be determined based on requirements
- **Styling**: CSS-in-JS or CSS modules
- **Testing**: Bun's built-in test runner
- **Build**: Bun's bundler for production builds

Remember: Security and accessibility are non-negotiable. Code quality and testing are essential for maintainability. Always prioritize user safety and experience.