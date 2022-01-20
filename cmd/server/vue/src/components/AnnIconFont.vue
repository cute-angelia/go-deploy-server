<template>
  <svg class="icon" aria-hidden="true">
    <use :xlink:href="'#' + icon"></use>
  </svg>
</template>

<script>
export default {
  props: {
    icon: {
      default: "ann-icon-zhaofuwu",
      type: String,
    },
  },
  data() {
    return {};
  },
  mounted() {
    let s = document.getElementById("iconfont");
    // if (s) {
    //     s.remove()
    // }
    if (!s) {
      s = document.createElement("script");
      s.id = "iconfont";
      s.type = "text/javascript";
      s.src =
        "https://best-ecology.oss-cn-hangzhou.aliyuncs.com/tool/28/iconfont.js?a=3";
      document.body.appendChild(s);
    }
  },
};
</script>

<style scoped>
.icon {
  width: 1em;
  height: 1em;
  vertical-align: -0.15em;
  fill: currentColor;
  overflow: hidden;
}
</style>