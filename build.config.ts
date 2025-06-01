import { defineBuildConfig } from 'unbuild'

export default defineBuildConfig({
  entries: ['src/index.ts', 'src/dictionary.ts'],
  declaration: true /* Generates .d.ts declaration file */,
  sourcemap: false /* Generates .map sourcemap file */,
  clean: true /* Clean the dist directory before building */,
  rollup: {
    emitCJS: false /* Don't emit CommonJS files */,
    commonjs: false /* Don't use CommonJS */,
  },
})
