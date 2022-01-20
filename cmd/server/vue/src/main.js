import ElementPlus from 'element-plus'
import 'element-plus/theme-chalk/index.css'
import {
  createApp
} from 'vue'
import App from './App.vue'

import {
  createStore
} from 'vuex'
import router from './router/router'


const app = createApp(App);


// 创建一个新的 store 实例
const store = createStore({
  state() {
    return {
      token: "",
      userInfo: {
        uid: 0,
        nickname: "",
        mobile: "",
        head: "",
      },
    }
  },
  mutations: {
    setToken(state, token) {
      localStorage['token'] = token
      this.state.token = token;
    },
    setUserInfo(state, info) {
      this.state.userInfo = info;
    }
  },
  getters: {
    getToken: (state) => {
      return state.token
    },
    getUserInfo: (state) => {
      return state.userInfo
    }
  }
})
// 将 store 实例作为插件安装
app.use(store)


import {
  iget,
  ipost,
  upload
} from "@/utils/api";
app.provide('$post', ipost);
app.provide('$get', iget);
app.provide('$upload', upload);

app.use(router)
app.use(ElementPlus)

app.mount('#app')