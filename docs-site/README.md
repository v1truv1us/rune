# Rune CLI Documentation Site

This is the documentation site for Rune CLI, built with Astro + Svelte + TailwindCSS.

## Development

### Prerequisites

- Node.js 20+
- pnpm

### Local Development

```bash
# Install dependencies
pnpm install

# Start development server
pnpm dev

# Build for production
pnpm build

# Preview production build
pnpm preview
```

### Docker Development

```bash
# Start development server with Docker
docker-compose --profile dev up docs-dev

# Build and run production container
docker-compose up docs
```

## Deployment

The site is configured to deploy automatically via GitHub Actions when changes are pushed to the main branch.

### Manual Deployment

```bash
# Build the site
pnpm build

# Deploy using your preferred method
# Examples:
# - Vercel: vercel deploy --prod
# - Netlify: netlify deploy --prod
# - Custom server: rsync -av dist/ user@server:/var/www/runecli.dev/
```

## Site Structure

```
src/
├── layouts/
│   └── BaseLayout.astro     # Main layout template
├── pages/
│   ├── index.astro          # Homepage
│   ├── getting-started.md   # Getting started guide
│   ├── docs/
│   │   ├── index.md         # Documentation index
│   │   ├── commands.md      # Commands reference
│   │   └── configuration.md # Configuration guide
│   └── examples/
│       └── index.md         # Examples and workflows
├── components/              # Reusable components
└── content/                 # Content collections
```

## Content Guidelines

### Writing Documentation

- Use clear, concise language
- Include code examples for all features
- Follow the existing structure and style
- Test all code examples before publishing

### Code Examples

- Use proper syntax highlighting
- Include both simple and advanced examples
- Show expected output when helpful
- Use realistic project names and paths

### Navigation

- Update navigation in BaseLayout.astro when adding new pages
- Ensure all pages are accessible from the main navigation
- Use descriptive page titles and meta descriptions

## Deployment Options

### Vercel (Recommended)

1. Connect your GitHub repository to Vercel
2. Set build command: `cd docs-site && pnpm build`
3. Set output directory: `docs-site/dist`
4. Configure custom domain: runecli.dev

### Netlify

1. Connect repository to Netlify
2. Set build command: `cd docs-site && pnpm build`
3. Set publish directory: `docs-site/dist`
4. Configure custom domain

### Self-Hosted

Use the provided Dockerfile and nginx configuration:

```bash
# Build and deploy
docker build -t runecli-docs .
docker run -d -p 80:80 runecli-docs
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes in the `docs-site/` directory
4. Test locally with `pnpm dev`
5. Submit a pull request

## License

MIT License - see the main project LICENSE file.
EOF < /dev/null