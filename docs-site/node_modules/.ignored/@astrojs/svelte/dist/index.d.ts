import type { Options } from '@sveltejs/vite-plugin-svelte';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import type { AstroIntegration, ContainerRenderer } from 'astro';
export declare function getContainerRenderer(): ContainerRenderer;
type OptionsCallback = (defaultOptions: Options) => Options;
export default function (options?: Options | OptionsCallback): AstroIntegration;
export { vitePreprocess };
