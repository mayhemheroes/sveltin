import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import adapter from '@sveltejs/adapter-static';
import preprocess from 'svelte-preprocess';
import { mdsvex } from 'mdsvex';
import mdsvexConfig from './mdsvex.config.js';

const sveltinJsonFile = fileURLToPath(new URL('sveltin.json', import.meta.url));
const sveltinJson = readFileSync(sveltinJsonFile, 'utf8');
const { sveltekit } = JSON.parse(sveltinJson);

/** @type {import('@sveltejs/kit').Config} */
const config = {
	extensions: ['.svelte', ...mdsvexConfig.extensions],
	preprocess: [
		mdsvex(mdsvexConfig),
		preprocess({
			preserve: ['ld+json'],
			scss: {
				prependData: '@use "src/_variables.scss" as *;',
			},
		}),
	],
	kit: {
		adapter: adapter({
			pages: sveltekit.adapter.pages,
			assets: sveltekit.adapter.assets,
			fallback: sveltekit.adapter.fallback,
			precompress: sveltekit.adapter.precompress || false,
			strict: sveltekit.adapter.strict || true,
		}),
		prerender: {
			crawl: true,
			entries: ['*'],
		},
	},
};

export default config;
