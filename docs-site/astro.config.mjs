import { defineConfig } from "astro/config";
import svelte from "@astrojs/svelte";
import mdx from "@astrojs/mdx";
import sitemap from "@astrojs/sitemap";

export default defineConfig({
  site: "https://runecli.dev",
  integrations: [
    svelte(), 
    mdx(),
    sitemap({
      changefreq: 'weekly',
      priority: 0.7,
      lastmod: new Date(),
      entryLimit: 10000,
    })
  ],
  output: 'static',
  markdown: {
    shikiConfig: {
      theme: "github-dark-dimmed",
      langs: ["bash", "yaml", "go", "javascript", "typescript"],
    },
  },
});
