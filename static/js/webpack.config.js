var path = require('path');
var webpack = require('webpack');

module.exports = {
  entry: './index.jsx',
  devtool: 'cheap-module-source-map',
  output: {
    filename: 'app.js',
    sourceMapFilename: 'index.js.map',
  },
  response: './',
  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
        query: {
          presets: ['es2015', 'react'],
          plugins: ['syntax-decorators'],
        }
      },
      {
        test: /\.less$/,
        loader: 'style!css!less',
      },
    ]
  },
};
