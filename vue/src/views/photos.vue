<template>
  <div>
    <div>
      <el-input
        placeholder="kindid"
        v-model="search.kindid"
        size="mini"
        style="margin: 5px 5px 5px 0px; width: 80px"
      />
      <el-input
        placeholder="numid"
        v-model="search.numid"
        size="mini"
        style="margin: 5px 5px; width: 80px"
      />

      <el-select
        style="margin: 5px 5px"
        v-model="searchStatus"
        size="mini"
        placeholder="Select"
      >
        <el-option
          v-for="item in statusOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        >
        </el-option>
      </el-select>
      <el-button size="mini" @click="getData(1, 'search')" type="success" round
        >搜索</el-button
      >
    </div>
    <div style="width: 100%">
      <div class="imgbox" v-for="(item, index) in imgList" :key="index">
        <el-image
          :src="item.url_userinfo_face_url"
          style="width: 100px; height: 100px"
        >
          <template #placeholder>
            <div class="image-slot">Loading<span class="dot">...</span></div>
          </template>
        </el-image>
        <div class="covericon">
          <ann-icon-font
            title="缩小图片"
            v-if="searchStatus == '4'"
            @click="doUpdate(item, index, 3)"
            style="float: right; margin: 2px"
            icon="iconsmall-copy"
          ></ann-icon-font>
          <ann-icon-font
            title="检测通过"
            @click="doUpdate(item, index, 0)"
            style="float: right; margin: 2px"
            icon="iconroundcheck-copy"
          ></ann-icon-font>

          <el-popconfirm
            title="确认还原头像?"
            @confirm="doUpdate(item, index, 2)"
          >
            <template #reference>
              <ann-icon-font
                title="重置误杀"
                style="float: right; margin: 2px"
                icon="iconReset-Outlined-copy"
              ></ann-icon-font>
            </template>
          </el-popconfirm>

          <el-popconfirm title="确认删除?" @confirm="doUpdate(item, index, 1)">
            <template #reference>
              <ann-icon-font
                title="删除图片"
                style="float: right; margin: 2px"
                icon="iconic-delete-copy"
              ></ann-icon-font>
            </template>
          </el-popconfirm>

          <ann-icon-font
            title="查看原图"
            @click="showOrigin(item)"
            style="float: right; margin: 2px"
            icon="iconcover"
          ></ann-icon-font>

          <span
            :title="item.url_userinfo_face_url"
            style="
              font-size: 12px;
              width: 97px;
              overflow: hidden;
              height: 20px;
              white-space: nowrap;
              display: inline-block;
              color: #ccc;
            "
            >{{ item.url_userinfo_face_url }}</span
          >
        </div>
      </div>
    </div>

    <br />

    <el-pagination
      style="display: inline-block; margin: 5px 0px"
      :current-page="pagination.page"
      :page-size="pagination.perpage"
      layout="total,prev, pager, next"
      :total="pagination.total"
      prev-text="上一页"
      next-text="下一页"
      @currentPage="currentChange"
      @current-change="currentChange"
    >
    </el-pagination>

    <el-dialog v-model="dialogVisible" title="查看原图" width="80%">
      <span>
        id:{{ currentItem.id }} kindid:{{ currentItem.kindid }} numid:{{
          currentItem.numid
        }}</span
      >

      <div style="box-sizing: border-box">
        <div
          style="
            padding: 30 px 0;
            text-align: center;
            border-right: solid 1 px var(--el-border-color-base);
            display: inline-block;
            width: 200px;
            box-sizing: border-box;
            vertical-align: top;
            margin-right: 10px;
          "
        >
          <span
            style="
              display: block;
              color: var(--el-text-color-secondary);
              font-size: 14px;
              margin-bottom: 20 px;
            "
            >入库检测图片</span
          >
          <el-image :src="currentItem.url" style="width: 200px; height: 200px">
          </el-image>
          <span
            style="
              display: block;
              color: var(--el-text-color-secondary);
              font-size: 14px;
              margin-bottom: 20 px;
            "
            >{{ currentItem.url }}</span
          >
        </div>
        <div
          style="
            padding: 30 px 0;
            text-align: center;
            border-right: solid 1 px var(--el-border-color-base);
            display: inline-block;
            width: 20%;
            box-sizing: border-box;
            vertical-align: top;
          "
        >
          <span
            style="
              display: block;
              color: var(--el-text-color-secondary);
              font-size: 14px;
              margin-bottom: 20 px;
            "
            >用户当前使用图片</span
          >
          <el-image
            :src="currentItem.url_userinfo_face_url"
            style="width: 200px; height: 200px"
          >
          </el-image>
          <span
            style="
              display: block;
              color: var(--el-text-color-secondary);
              font-size: 14px;
              margin-bottom: 20 px;
            "
            >{{ currentItem.url_userinfo_face_url }}</span
          >
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="dialogVisible = false"
            >关闭</el-button
          >
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
// import { inject } from "vue";
// var ipost = inject("$post");

import { ipost } from "@/utils/api";
import AnnIconFont from "@/components/AnnIconFont.vue";
import { ElMessage } from "element-plus";

// 0 => 未处理
// 1 => 已处理完成
// 2 => 异常：上传头像失败情况
// 3 => 状态1的补充，给用户恢复了默认头像，但是未删除云端违规图片， 脚本 delete 删除改状态图片
// 4 => 图片尺寸大于 1m 情况， 需要脚本处理小头像 ?x-oss-process=image/resize,w_400/format,jpg/quality,q_90

export default {
  components: { AnnIconFont },
  data() {
    return {
      currentPage: 1,
      imgList: [], // {"url_current" "url_old"}
      pagination: {
        page: 1,
        total: 0,
        perpage: 36,
      },
      search: {
        kindid: "",
        numid: "",
      },
      searchStatus: "3",
      statusOptions: [
        { value: "-1", label: "全部头像" },
        { value: "3", label: "3违规图片" },
        { value: "4", label: "4大尺寸图片" },
        { value: "1", label: "1已处理完成" },
        { value: "0", label: "0未处理图片" },
        { value: "2", label: "2异常：上传头像失败情况" },
      ],
      dialogVisible: false,
      currentItem: null,
    };
  },
  mounted() {
    // console.log("  mounted()", this.$route.query.auth);
    this.getData(this.pagination.page);
    console.log("get token:", this.$store.getters.getToken);
  },
  methods: {
    // 分页
    currentChange: function (val) {
      console.log("currentChange", val);
      this.getData(val, "currentChange");
    },
    showOrigin(item) {
      console.log(item);
      this.currentItem = item;
      this.dialogVisible = true;
    },
    getData(page, from = "") {
      var that = this;
      console.log("page:", page, "status:", that.searchStatus, "from", from);
      this.pagination.page = page;
      that.imgList = [];
      ipost("FaceLists", {
        page: this.pagination.page,
        status: that.searchStatus,
        perpage: that.pagination.perpage,
        auth: that.$route.query.auth,
        kindid: that.search.kindid,
        numid: that.search.numid,
      }).then((data) => {
        if (data.code == 0) {
          // 分页数据
          that.pagination.total = data.data.count;
          for (let i = 0; i < data.data.lists.length; i++) {
            const element = data.data.lists[i];
            that.imgList.push(element);
          }
        }
      });
    },
    doUpdate(item, index, flag) {
      console.log(item, index, flag);
      var that = this;
      ipost("FaceUpdate", {
        flag: flag,
        id: item.id,
        current_face_url: item.url_userinfo_face_url,
        auth: that.$route.query.auth,
      }).then((data) => {
        if (data.code == 0) {
          console.log(data);
          that.imgList.splice(index, 1);

          ElMessage({
            message: data.message || "Error",
            type: "success",
            duration: 3 * 1000,
          });
        }
      });
    },
  },
};
</script>

<style scoped>
.div-inline {
  display: inline;
}
.imgbox {
  width: 100px;
  height: 100px;
  margin: 1px;
  display: inline-block;
}
.covericon {
  height: 20px;
  width: 100px;
  position: relative;
  top: 0px;
  z-index: 100;
}
</style>