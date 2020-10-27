<template>
  <div v-loading="loading">



    <!-- table -->

    <div style="width: 95%;margin-left: 2.5%;margin-top: 30px;border-top: 1px solid #EBEEF5;">
      <el-table
          :data="list"
          :header-cell-style="{'text-align':'center'}"
          :cell-style="{'text-align':'center'}"
          style="width: 100%"
          id="adminTable"
      >
        <el-table-column prop="Avatar" label="头像"></el-table-column>
        <el-table-column prop="Name" label="用户名"></el-table-column>
        <el-table-column prop="Email" label="邮箱"></el-table-column>
        <el-table-column prop="Phone" label="手机"></el-table-column>

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
      dialogVisible:false,
      Email:"",
      Password:"",
      Name:"",
      Captcha:"",
      rePassword:"",
    };
  },
  created() {
    this.loading = true;
    this.$axios.adminList(this.requestData).then((res) => {
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
  methods: {
    newAdmin(){

    },
    updateList() {
      this.loading = true;
      this.$axios
          .adminList(this.requestData)
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
              this.currentPage = e;
            }
            this.loading = false;
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
    },
    pageChange(e) {
      this.requestData.PageIndex = e;
      this.$axios.adminList(this.requestData).then((res) => {
        if (res.code == 200) {
          for (let i = 0; i < res.data.length; i++) {
            for (let key in res.data[i]) {
              if (res.data[i][key] == "") {
              //  res.data[i][key] = "无";
              }
            }
          }
          this.list = res.data;
          this.count = res.count;
          this.currentPage = e;
        }
      });
    },
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