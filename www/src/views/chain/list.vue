<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">链管理</h4>
    </el-header>
    <el-main>
      <el-table v-loading="table_loading" :data="list">
        <el-table-column prop="id" label="ID" width="100" />
        <el-table-column prop="chain_symbol" label="Chain" width="100" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="currency" label="币种" />
        <el-table-column prop="confirm" label="确认数" />
        <el-table-column prop="concurrent" label="并发检索" />
        <el-table-column prop="gas_price" label="Gas Price">
          <template slot-scope="scope">
            <span>{{ scope.row.gas_price | toGwei }} Gwei</span>
          </template>
        </el-table-column>
        <el-table-column prop="latest_block" label="Latest Block" width="120" />
        <el-table-column prop="chain_type" label="链类型" />
        <el-table-column prop="chain_type" label="操作" width="220">
          <template>
            <el-button type="primary" size="mini">RPC管理</el-button>
            <el-button type="primary" size="mini">Token管理</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        layout="prev, pager, next, jumper, ->, total"
        :current-page="page"
        :page-size="pagesize"
        :total="totalnum"
        @current-change="(e) => { page = e; loadData() }"
      />
    </el-main>
  </el-container>
</template>
<script>
import { getChainList } from '@/api/chain'
export default {
  data() {
    return {
      list: [],
      pagesize: 20,
      totalnum: 0,
      page: 1,
      table_loading: false
    }
  },
  computed: {},
  async mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      this.table_loading = true
      var params = {
        page: this.page,
        size: this.pagesize
      }
      const res = await getChainList(params)
      this.totalnum = res.data.total
      this.list = res.data.list
      this.table_loading = false
    }
  }
}
</script>
