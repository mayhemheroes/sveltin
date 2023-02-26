import { mdsvex } from 'mdsvex';
import mdsvexConfig from './mdsvex.config.js';

import preprocess from 'svelte-preprocess';
import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	extensions: ['.svelte', ...mdsvexConfig.extensions],
	preprocess: [
		mdsvex(mdsvexConfig),
		preprocess({
			preserve: ['ld+json'],
		}),
	],
	kit: {
		adapter: adapter({
			pages: 'build',
			assets: 'build',
			fallback: '200.html',
			precompress: true,
		}),
		prerender: {
			crawl: true,
			entries: ['*'],
		},
	},
};

export default config;
