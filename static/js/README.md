### Frontend Architecture

Tapestry's front-end is built in JavaScript (Babel/ES6) using React and LESS.

Module loading, source map generation, polyfills, source compilation/transpilation is managed via WebPack.

To perform a one-off compile the JSX and LESS files to `app.js` (and the sourcemap), simply run webpack:

```sh
$ webpack
```

If you're hacking on the frontend, though (CSS or JSX), you'll want to run WebPack in "watch" mode, so it will continuously re-compile the site as soon as the source files are changed.

```sh
$ webpack -w --progress
```
