
const path = require('path');

const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin

const ThreeExamples = require('import-three-examples')
function resolve(dir) {
    return path.join(__dirname,'.', dir);
}

module.exports = {
    baseUrl : '/',
    publicPath : '/',
    productionSourceMap: true,
    chainWebpack: config => {
        config.resolve.alias
        .set('@', resolve('src'))
        .set('@assets',resolve('src/assets'))
        .set('styles',resolve('src/styles'))
        .set('@utils',resolve('src/utils'))
        .set('@views',resolve('src/views'))
    },
    devServer:{
        // 设置代理
        // before: require('./src/mock'),
        host: '0.0.0.0',
        port: 8080,
        hot: true,
        https: false,
        open: false,
        disableHostCheck: true,
        proxy: {
            "/config": {
                target: "http://127.0.0.1:8080", // 访问数据的计算机域名192.168.9.155:3001
                // target: "http://122.51.143.73:8080", // 访问数据的计算机域名192.168.9.155:3001
                ws: true, // 是否启用websockets
                changOrigin: true, //开启代理
                //将api替换为空
                // pathRewrite:{
                //     '^/api':''
                // },
            },
        }
        
    },
    configureWebpack: config => {
        if (process.env.NODE_ENV === 'production') {
            return {
                plugins: [new BundleAnalyzerPlugin()]
            }
        }
    },
     // 第三方插件配置
     pluginOptions: {
        // ...
        ...ThreeExamples
    },
    pages: {
        login: new PageReset('login', 'pixiu控制台管理系统'),
    }
}
/**
 * 页面构造器
 * @param {String} name 页面名称
 * @param {String} title 页面title
 */
 function PageReset (name, title) {
    // page 的入口
    this.entry = `src/entry/${name}.js`
    // 模板来源
    this.template = 'public/index.html'
    // 在 dist/index.html 的输出
    this.filename = `${name}.html`
    // 当使用 title 选项时，
    // template 中的 title 标签需要是 <title><%= htmlWebpackPlugin.options.title %></title>
    this.title = title
    // 在这个页面中包含的块，默认情况下会包含
    // 提取出来的通用 chunk 和 vendor chunk。
    this.chunks = ['chunk-vendors', 'chunk-common', name]
  }
  