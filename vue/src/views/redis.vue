<template>
  <div>
    <div>
      <el-input
        placeholder="numids"
        v-model="search.numids"
        size="mini"
        style="margin: 5px 5px; width: 120px"
        @keyup.enter="getData(1, 'search')"
      />
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
            title="设置头像"
            @click="settingHead(item, index)"
            style="float: right; margin: 2px"
            icon="iconReset-Outlined-copy"
          ></ann-icon-font>

          <el-popconfirm title="确认删除?" @confirm="doUpdate(item, index, 4)">
            <template #reference>
              <ann-icon-font
                title="删除图片"
                style="float: right; margin: 2px"
                icon="iconic-delete-copy"
              ></ann-icon-font>
            </template>
          </el-popconfirm>

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

    <el-dialog v-model="dialog.show" title="设置头像" width="80%">
      <el-form ref="formRef" :model="dialog.form" label-width="120px">
        <el-form-item label="设置新头像地址">
          <el-input v-model="dialog.form.newUrl"></el-input>
        </el-form-item>
        <el-form-item label="设置新头像">
          <el-image
            :src="dialog.form.newUrl"
            style="width: 100px; height: 100px"
          />
        </el-form-item>
        <el-form-item label="当前头像地址">
          <el-input disabled v-model="dialog.form.currentUrl"></el-input>
        </el-form-item>
        <el-form-item label="当前头像">
          <el-image
            :src="dialog.form.currentUrl"
            style="width: 100px; height: 100px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit(dialog.form.index)"
            >保存</el-button
          >
          <el-button @click="dialog.show = false">取消</el-button>
        </el-form-item>
      </el-form>
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
        numids: "",
      },
      currentItem: null,
      dialog: {
        show: false,
        dform: {
          index: 0,
          currentUrl: "",
          newUrl: "",
        },
        form: {
          index: 0,
          currentUrl: "",
          newUrl: "",
        },
      },
    };
  },
  mounted() {
    this.getData(this.pagination.page);
  },
  methods: {
    // 分页
    currentChange: function (val) {
      this.getData(val, "currentChange");
    },
    showOrigin(item) {
      this.currentItem = item;
      this.dialogVisible = true;
    },
    settingHead(item, index) {
      console.log(item, index);
      this.dialog.show = true;
      this.dialog.form.currentUrl = item.url_userinfo_face_url;
      this.dialog.form.newUrl = "";
      this.dialog.form.index = index;
    },
    getData(page) {
      var that = this;
      this.pagination.page = page;
      that.imgList = [];
      ipost("FaceSearchLists", {
        page: this.pagination.page,
        perpage: that.pagination.perpage,
        numids: that.search.numids,
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
        id: 0,
        current_face_url: item.url_userinfo_face_url,
      }).then((data) => {
        if (data.code == 0) {
          console.log(data);
          // that.imgList.splice(index, 1);
          that.imgList[index].url_userinfo_face_url =
            data.data.url_userinfo_face_url;

          ElMessage({
            message: data.message || "Error",
            type: "success",
            duration: 10 * 1000,
          });
        }
      });
    },
    onSubmit(index) {
      var that = this;
      ipost("FaceUpdate", {
        flag: 5,
        id: 0,
        current_face_url: this.dialog.form.currentUrl,
        new_face_url: this.dialog.form.newUrl,
      }).then((data) => {
        if (data.code == 0) {
          console.log(data);
          // that.imgList.splice(index, 1);
          that.imgList[index].url_userinfo_face_url =
            data.data.url_userinfo_face_url;

          ElMessage({
            message: data.message || "Error",
            type: "success",
            duration: 10 * 1000,
          });

          this.dialog.show = false;
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
  top: -29px;
  z-index: 100;
}
</style>