<template>
  <el-menu
    :default-active="activeIndex"
    class="el-menu-demo"
    mode="horizontal"
    @select="handleSelect"
  >
    <!-- <el-menu-item index="1">
      <router-link style="text-decoration: none" to="/">{{
        title
      }}</router-link></el-menu-item
    > -->
    <el-menu-item index="1">
      <router-link style="text-decoration: none" to="/">主页</router-link>
    </el-menu-item>
  </el-menu>
  <!-- 路由匹配到的组件将渲染在这里 -->
  <router-view></router-view>
</template>

<script>
import { ipost } from "@/utils/api";
import { ElMessage } from "element-plus";

export default {
  name: "App",
  components: {},
  data() {
    return {
      title: "发布系统",
      activeIndex: 1,
      code: "",
    };
  },
  watch: {
    // 如果路由有变化，会再次执行该方法
    $route: "getCode",
  },
  mounted() {
    // 本地开发
    if (this.isLocal()) {
      this.getCode();
    }
  },
  methods: {
    isLocal() {
      return window.location.href.indexOf("localhost") >= 0;
    },
    handleSelect(key, keyPath) {
      console.log(key, keyPath);
    },
    getCode() {
      var code = this.$route.query.code;
      if (code == undefined) {
        // 本地模拟登录
        if (this.isLocal()) {
          console.log(" isLocal get token:", this.$store.getters.getToken);

          if (
            this.$store.getters.getToken &&
            this.$store.getters.getToken.length > 0
          ) {
            return;
          }

          let userInfo = {
            head: "",
            mobile: "",
            nickname: "游客_493782",
            openid: "oe0-i4qJWvZkRcqhrwV_ziNxM5Dc",
            sex: 0,
            token:
              "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBpZCI6IjIwMjAxMjAxMTQyNjE4NTAxMTMwIiwiY2lkIjoiMSIsImV4cCI6MTk3MjQ1Nzk5MiwiZnJvbSI6Im9yaWdpbiIsIm5pY2tuYW1lIjoi5ri45a6iXzQ5Mzc4MiIsInRva2VuIjoiZGRiMTI0NTQwOTY3NjJlY2NjODg2NDUxZTc3ZDg1MmYiLCJ1aWQiOiI2MTgwMCJ9.a_-L8H3IVCoMAPyi7XYQl8bHt20IpU4EPRRocjN4B2w",
            uid: 61800,
          };

          console.log("本地开发登录");

          this.$store.commit("setUserInfo", userInfo);
          this.$store.commit("setToken", userInfo.token);

          // 跳转
          setTimeout(() => {
            this.$router.push({ name: "yqgame", query: {} });
          }, 200);
        }
      } else {
        this.code = code;
        this.getLoginToken();
      }
    },
    getLoginToken() {
      var that = this;
      ipost(
        "https://api-game-common.yqgame.online/api-auth/wxWork/getToken?debug=1",
        {
          code: this.code,
        }
      )
        .then((data) => {
          console.log(data);
          that.$store.commit("setToken", data.data.token);
          that.$store.commit("setUserInfo", data.data);

          // 跳转
          setTimeout(() => {
            this.$router.push({ name: "yqgame", query: {} });
          }, 200);
        })
        .catch((error) => {
          // 跳转授权页面
          this.$router.push({ name: "home", query: {} });
          console.log(error);

          ElMessage({
            message: "请在企业微信打开",
            type: "success",
            duration: 3 * 1000,
          });
        });
    },
  },
};
</script>

<style>
</style>
