const { defineConfig } = require("@vue/cli-service");
module.exports = defineConfig({
  publicPath: '{{.BasePath}}static/',
  transpileDependencies: [
    'vuetify',
    '@prometheus-io/codemirror-promql',
    'sql-formatter'
  ],
  chainWebpack: (config) => {
    // Use vue-svg-loader for .svg files
    config.module.rules.delete('svg');// Clear existing loaders

    config.module
      .rule('svg')
      .test(/\.svg$/)
      .use('vue-svg-loader')
      .loader('vue-svg-loader')
      .end();
  },
});
