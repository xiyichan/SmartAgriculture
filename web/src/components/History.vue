<template>
  <el-container style="display:block">
    <el-header style="margin-top: 20px;">
      <span>
        <el-button type="primary" size="small" @click="goBack()">返回</el-button>
      </span>
      <span style="float: right;">
        <el-date-picker v-model="value1" type="date" placeholder="选择查看日期" @change="getTableList()" size="small"></el-date-picker>
      </span>
    </el-header>
    <el-main>
      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="Temperature" label="温度"></el-table-column>
        <el-table-column prop="Humidity" label="湿度"></el-table-column>
        <el-table-column prop="SoilMoisture" label="土壤湿度"></el-table-column>
        <el-table-column prop="LightIntensity" label="光照强度"></el-table-column>
        <el-table-column prop="Time" label="时间"></el-table-column>
      </el-table>
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="page"
        :page-sizes="[5, 10, 15, 20]"
        :page-size="limit"
        layout="sizes, prev, pager, next"
        :total="total"
        style="float: right;margin-top: 20px;"
      ></el-pagination>
    </el-main>
  </el-container>
</template>
<script>
export default {
  data() {
    return {
      value1: "",
      limit: 10,
      page: 1,
      tableData: [],
      total: 0
    }
  },
  computed: {
    id() {
      return this.$route.params.task_id
    }
  },
  mounted() {
    this.getTableList()
  },
  methods: {
    getTableList() {
      this.$axios
        .piHistory({
          IotId: this.id,
          Day: new Date(this.value1).Format("yyyy-MM-dd"),
          Limit: this.limit,
          Page: this.page
        })
        .then(res => {
          this.tableData = res.data.map(row => ({
            ...row,
            Time: new Date(row.Time).Format("yyyy-MM-dd hh:mm:ss")
          }))
          this.total = res.count
        })
    },
    handleSizeChange(val) {
      this.limit = val
      this.getTableList()
    },
    handleCurrentChange(val) {
      this.page = val
      this.getTableList()
    },
    goBack() {
      this.$router.push({ path: "/Pi" })
    }
  }
}
</script>