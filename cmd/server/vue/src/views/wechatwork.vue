<template>
  <div>
    <el-card class="box-card" v-for="(item, index) in tableList" :key="index">
      <template #header>
        <div class="card-header">
          <span>{{ item.name }}({{ item.node.length }}台)</span>
          <!-- <el-button class="button" type="text">Operation button</el-button> -->
          <ann-icon-font
            v-if="item.type == 'svn'"
            style="margin: 2px"
            icon="icon-wugui"
            title="svn"
          ></ann-icon-font>
          <ann-icon-font
            v-else
            style="margin: 2px"
            icon="icongithub-copy"
            title="git"
          ></ann-icon-font>
        </div>
      </template>
      <div v-for="(inode, nindex) in item.node" :key="nindex" class="text item">
        <ann-icon-font
          v-if="inode.online"
          :title="inode.alias"
          style="margin: 2px; color: #0e932e"
          icon="iconfuwuqi-green-copy"
        ></ann-icon-font>
        <ann-icon-font
          v-else
          :title="inode.alias"
          style="margin: 2px"
          icon="iconfuwuqi-red-copy"
        ></ann-icon-font>
        {{ inode.alias }}
      </div>
      <div class="bottom" style="margin: 10px 0px">
        <time class="time">历史</time>
        <el-select
          v-model="item.reversion"
          style="margin: 0px 5px"
          class="m-2"
          size="mini"
        >
          <el-option
            v-for="item in item.logs"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          >
          </el-option>
        </el-select>
        <el-button
          @click="getLog(item.groupid, index)"
          type="text"
          class="button"
          >刷新</el-button
        >
        <el-popconfirm
          v-if="item.reversion"
          @confirm="rollback(item)"
          :title="'回滚代码到[' + item.reversion + ']' + '?'"
        >
          <template #reference>
            <el-button size="mini" type="danger">回滚</el-button>
          </template>
        </el-popconfirm>
      </div>
      <div class="bottom" style="margin: 10px 0px">
        <el-row>
          <el-popconfirm
            @confirm="update(item)"
            :title="'确定更新代码到[' + item.name + ']?'"
          >
            <template #reference>
              <el-button type="primary">更新</el-button>
            </template>
          </el-popconfirm>
        </el-row>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ipost } from "@/utils/api";
import AnnIconFont from "@/components/AnnIconFont.vue";
import { h } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";

export default {
  components: { AnnIconFont },
  data() {
    return {
      tableList: [],
    };
  },
  mounted() {
    this.getData(1);
    // console.log("get token:", this.$store.getters.getToken);
  },
  methods: {
    getData(page) {
      var that = this;
      ipost("List", {
        page: page,
      }).then((data) => {
        if (data.code == 0) {
          that.tableList = data.data;
        }
      });
    },
    getLog(groupId, index) {
      var that = this;
      var logs = [];
      ipost("Showlog", {
        groupid: groupId,
      }).then((data) => {
        if (data.code == 0) {
          var dataz = data.data;
          for (let i = 0; i < dataz.length; i++) {
            const element = dataz[i];
            var label =
              element["Reversion"] +
              " | " +
              element["Author"] +
              " | " +
              element["Time"] +
              " | " +
              element["Content"];
            var value = element["Reversion"];

            logs.push({
              label: label,
              value: value,
            });
          }
          var tempv = that.tableList[index];
          tempv.logs = logs;
          tempv.reversion = logs.length > 0 ? logs[0]["value"] : "";
          that.tableList[index] = tempv;
        }
      });
    },
    update(item) {
      ipost("Deploy", {
        groupid: item.groupid,
      }).then((data) => {
        if (data.code == 0) {
          ElMessageBox({
            title: "通知",
            message: h("div", null, [
              h(
                "font",
                {
                  style: "color: teal",
                  display: "inherit",
                },
                "运行时间：" + data.data.time_cos + "s"
              ),
              h(
                "span",
                {
                  style: "white-space: pre-wrap",
                },
                "\n\n" + data.data.msg
              ),
            ]),
            confirmButtonText: "OK",
          });
        }
      });
    },
    rollback(item) {
      if (!item.reversion) {
        ElMessage({
          message: "未能回滚版本，请点击刷新获取历史版本",
          type: "error",
          duration: 5 * 1000,
        });
        return;
      }
      ipost("Rollback", {
        groupid: item.groupid,
        reversion: item.reversion,
      }).then((data) => {
        if (data.code == 0) {
          ElMessageBox({
            title: "通知",
            message: h("div", null, [
              h(
                "font",
                {
                  style: "color: teal",
                  display: "inherit",
                },
                "运行时间：" + data.data.time_cos + "s"
              ),
              h(
                "span",
                {
                  style: "white-space: pre-wrap",
                },
                "\n\n" + data.data.msg
              ),
            ]),
            confirmButtonText: "OK",
          });
        }
      });
    },
  },
};
</script>


<style scoped>
.box-card {
  width: 390px;
  display: inline-block;
  margin: 10px 5px 0px 0px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.item {
  display: inline-block;
  margin-right: 5px;
}
</style>