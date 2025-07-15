import { defineConfig } from 'astro/config';
import svelte from '@astrojs/svelte';
import tailwind from '@astrojs/tailwind';
import mdx from '@astrojs/mdx';

export default defineConfig({
  site: 'https://runecli.dev',
  integrations: [
    svelte(),
    tailwind(),
    mdx()
  ],
  markdown: {
    shikiConfig: {
      theme: 'github-dark-dimmed',
      langs: ['bash', 'yaml', 'go', 'javascript', 'typescript']
    }
  }
});