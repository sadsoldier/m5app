
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
    configureWebpack: {
        watchOptions: {
            aggregateTimeout: 500,
            poll: 1000,
            ignored: ['node_modules']
        },
        optimization: {
            minimize: false
        }
    }
}
