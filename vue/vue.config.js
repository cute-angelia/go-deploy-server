// vue.config.js
const name = "box";
const path = require("path");

function resolve(dir) {
  return path.join(__dirname, dir);
}

// 接口地址
const UserConfig = require("./config/" + process.env.NODE_ENV + ".config");
process.env.VUE_APP_BASE_API = UserConfig.BASE_API_URL;

module.exports = {
  parallel: false,
  publicPath: process.env.NODE_ENV === "production" ? "/" : "/",
  outputDir: "dist",
  configureWebpack: {
    // provide the app's title in webpack's name field, so that
    // it can be accessed in index.html to inject the correct title.
    name: name,
    resolve: {
      alias: {
        "@": resolve("src")
      }
    }
  },
}