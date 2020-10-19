<template>
  <div v-loading="loading">
    <!-- 搜索框+注册 -->
    <div style="margin-top: 20px;width:100%;">
      <el-button
        type="primary"
        icon="el-icon-connection"
        style="width: 110px;height: 34px;padding: 0px;background-color: rgb(48, 208, 172);border-color: rgb(48, 208, 172);line-height: 34px;position: relative;left:80%;}"
        @click="newWeather"
      >新增气象</el-button>
    </div>
    <!-- table -->
    <div style="width: 95%;margin-left: 2.5%;margin-top: 30px;border-top: 1px solid #EBEEF5;">
      <el-table
        :data="list"
        :header-cell-style="{'text-align':'center'}"
        :cell-style="{'text-align':'center'}"
        style="width: 100%"
        id="userTable"
      >
      <el-table-column prop="DeviceName" label="设备名字"></el-table-column>
        <el-table-column prop="DeviceSecret" label="设备密码"></el-table-column>
        <el-table-column prop="NickName" label="昵称"></el-table-column>
        <el-table-column prop="IotId" label="IotId"></el-table-column>
        <el-table-column prop="Status" label="状态"></el-table-column>
        <el-table-column prop="UserId" label="所属用户ID"></el-table-column>
        <el-table-column label="操作" prop="Status">
          <template slot-scope="scope">
            <el-button type="text" size="small" style="font-size: 15px;" @click="del(scope)">删除气象</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        layout="prev, pager, next"
        :total="count"
        :page-size="6"
        :current-page="1"
        @current-change="pageChange"
      ></el-pagination>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      requestData: {
        PageIndex: 1,
        PageSize: 6,
      },
      list: [],
      count: 0,
      currentPage: 1,
      dialogFormVisible: false,
      form: {
        Title: "",
        Detail: "",
      },
      formLabelWidth: "120px",
    };
  },
  created() {
    this.loading = true;
    this.$axios.piList(this.requestData).then((res) => {
      if (res.code == 200) {
        for (let i = 0; i < res.data.length; i++) {
          for (let key in res.data[i]) {
            if (res.data[i][key] == "") {
              //res.data[i][key] = "无";
            }
          }
        }
        this.list = res.data;
        this.count = res.count;
      }
      this.loading = false;
    });
  },
  methods: {
    updateList() {
      this.loading = true;
      this.$axios
        .piList(this.requestData)
        .then((res) => {
          if (res.code == 200) {
            for (let i = 0; i < res.data.length; i++) {
              for (let key in res.data[i]) {
                if (res.data[i][key] == "") {
                 // res.data[i][key] = "无";
                }
              }
            }
            this.list = res.data;
            this.count = res.count;
          }
          this.loading = false;
        });
    },
    pageChange(e) {
      this.requestData.PageIndex = e;
      this.$axios.piList(this.requestData).then((res) => {
        if (res.code == 200) {
          for (let i = 0; i < res.data.length; i++) {
            for (let key in res.data[i]) {
              if (res.data[i][key] == "") {
                //res.data[i][key] = "无";
              }
            }
          }
          this.list = res.data;
          this.count = res.count;
          this.currentPage = e;
        }
      });
    },
    del(e) {
      this.$confirm("此操作将永久删除气象, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        this.loading = true;
        console.log(e);
        var d = {
          IotId:e.row.IotId
        };
        this.$axios
          .piDelete(d)
          .then((res) => {
            if (res.code == 200) {
              this.$message({
                type: "success",
                message: "删除成功!",
              });
              this.updateList();
              this.loading = false;
            }
          })
          .catch((error) => {
            if (error.response) {
              console.log(error.response.data);
              this.$alert(error.response.data.msg);
              this.loading = false;
            } else if (error.request) {
              console.log(error.request);
            } else {
              console.log("Error", error.message);
            }
            console.log(error.config);
          });
      });
    },
    newWeather(){
        this.$axios.piNew().then(res=>{
            if(res.code==200){
                this.$message({
                type: "success",
                message: "新增气象成功",
              });
            this.updateList()
            }
        })
    }
  },
};
</script>

<style scoped>
.el-dropdown-link {
  cursor: pointer;
  color: #409eff;
}
.el-icon-arrow-down {
  font-size: 12px;
}
</style>