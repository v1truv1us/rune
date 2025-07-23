import { fileURLToPath } from "node:url";
import { svelte, vitePreprocess } from "@sveltejs/vite-plugin-svelte";
import { VERSION } from "svelte/compiler";
const isSvelte5 = Number.parseInt(VERSION.split(".").at(0)) >= 5;
function getRenderer() {
  return {
    name: "@astrojs/svelte",
    clientEntrypoint: isSvelte5 ? "@astrojs/svelte/client-v5.js" : "@astrojs/svelte/client.js",
    serverEntrypoint: isSvelte5 ? "@astrojs/svelte/server-v5.js" : "@astrojs/svelte/server.js"
  };
}
function getContainerRenderer() {
  return {
    name: "@astrojs/svelte",
    serverEntrypoint: isSvelte5 ? "@astrojs/svelte/server-v5.js" : "@astrojs/svelte/server.js"
  };
}
async function svelteConfigHasPreprocess(root) {
  const svelteConfigFiles = ["./svelte.config.js", "./svelte.config.cjs", "./svelte.config.mjs"];
  for (const file of svelteConfigFiles) {
    const filePath = fileURLToPath(new URL(file, root));
    try {
      const config = (await import(
        /* @vite-ignore */
        filePath
      )).default;
      return !!config.preprocess;
    } catch {
    }
  }
}
async function getViteConfiguration({
  options,
  isDev,
  root
}) {
  const defaultOptions = {
    emitCss: true,
    compilerOptions: { dev: isDev }
  };
  if (!isSvelte5) {
    defaultOptions.compilerOptions.hydratable = true;
  }
  if (!isDev) {
    defaultOptions.hot = false;
  }
  let resolvedOptions;
  if (!options) {
    resolvedOptions = defaultOptions;
  } else if (typeof options === "function") {
    resolvedOptions = options(defaultOptions);
  } else {
    resolvedOptions = {
      ...options,
      ...defaultOptions,
      compilerOptions: {
        ...options.compilerOptions,
        // Always use dev and hydratable from defaults
        ...defaultOptions.compilerOptions
      }
    };
  }
  if (!resolvedOptions.preprocess && !await svelteConfigHasPreprocess(root)) {
    resolvedOptions.preprocess = vitePreprocess();
  }
  return {
    optimizeDeps: {
      include: [isSvelte5 ? "@astrojs/svelte/client-v5.js" : "@astrojs/svelte/client.js"],
      exclude: [isSvelte5 ? "@astrojs/svelte/server-v5.js" : "@astrojs/svelte/server.js"]
    },
    plugins: [svelte(resolvedOptions)]
  };
}
function src_default(options) {
  return {
    name: "@astrojs/svelte",
    hooks: {
      // Anything that gets returned here is merged into Astro Config
      "astro:config:setup": async ({ command, updateConfig, addRenderer, config }) => {
        addRenderer(getRenderer());
        updateConfig({
          vite: await getViteConfiguration({
            options,
            isDev: command === "dev",
            root: config.root
          })
        });
      }
    }
  };
}
export {
  src_default as default,
  getContainerRenderer,
  vitePreprocess
};
