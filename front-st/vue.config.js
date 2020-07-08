
const path = require("path");

module.exports = {
    outputDir: path.resolve(__dirname, "../public"),
    productionSourceMap: false,
    css: {
        loaderOptions: {
            sass: {
                prependData: `@import "@/app.scss";`
            }
        }
    },
    chainWebpack: config => {
        config.module
            .rule('i18n')
            .resourceQuery(/blockType=i18n/)
            .type('javascript/auto')
            .use('i18n')
            .loader('@intlify/vue-i18n-loader')
    },
    configureWebpack: config => {
        config.watchOptions = {
            aggregateTimeout: 700,
            poll: 2000,
            ignored: ['node_modules']
        },
        config.optimization = {
            minimize: false
        }
    }
}
