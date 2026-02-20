import { defineConfig } from 'tsup';

export default defineConfig({
  entry: ['src/index.ts'],
  format: ['cjs', 'esm'],
  dts: true,
  splitting: false,
  sourcemap: true,
  clean: true,
  external: [],
  noExternal: [],
  treeshake: true,
  minify: false,
  metafile: false,
  platform: 'node',
  target: 'node20',
  outDir: 'dist',
  outExtension({ format }) {
    return {
      js: format === 'cjs' ? '.cjs' : '.js',
      dts: '.d.ts'
    };
  }
});
