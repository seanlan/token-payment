<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">上币申请</h4>
    </el-header>
    <el-main>
      <el-form class="searchForm" label-width="120px">
        <el-row>
          <el-col :span="8">
            <el-form-item label="过滤无效申请">
              <el-switch
                v-model="searchForm.filter"
                @change="loadData"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <el-table v-loading="table_loading" :data="list" :row-class-name="tableRowClassName" @sort-change="sortChange">
        <el-table-column prop="user.id" label="User ID" width="100" />
        <el-table-column prop="user.display_name" label="昵称" width="180" />
        <el-table-column prop="telegram.id" label="TG ID" />
        <el-table-column prop="telegram.username" label="TG账号" />
        <el-table-column prop="coin.telegram" label="Telegram" width="180" />
        <el-table-column prop="coin.name" label="币种" width="100" />
        <el-table-column prop="coin.chain" label="币种网络" />
        <el-table-column prop="coin.contract" label="合约地址" />
        <el-table-column prop="coin.website" label="网址" />
        <el-table-column prop="status" width="100" label="无效过滤" fixed="right">
          <template slot-scope="scope">
            <el-switch
              v-model="scope.row.status"
              active-color="#F56C6C"
              :active-value="2"
              :inactive-value="1"
              @change="examineChange(scope.row)"
            />
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        layout="prev, pager, next, jumper, ->, total"
        :current-page="page"
        :page-size="pagesize"
        :total="totalnum"
        @current-change="currentChange"
      />
    </el-main>
  </el-container>
</template>
<script>
import config from '@/utils/config'
import {
  makeAjaxParamData
} from '@/utils/util'
import { coinApplicationList, examineCoinApplication } from '@/api/coin_application'
export default {
  data() {
    return {
      list: [],
      pagesize: 20,
      totalnum: 0,
      page: 1,
      table_loading: false,
      searchForm: {
        filter: false
      },
      formRules: {},
      formFilter: false,
      formRealPrice: 0,
      order_by: '`user`.`id`',
      order: 'desc',
      pickerOptions: {
        shortcuts: [{
          text: '今天',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime())
            picker.$emit('pick', [start, end])
          }
        },
        {
          text: '昨天',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 1)
            end.setTime(end.getTime() - 3600 * 1000 * 24 * 1)
            picker.$emit('pick', [start, end])
          }
        },
        {
          text: '最近一周',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近一个月',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
            picker.$emit('pick', [start, end])
          }
        }
        ]
      },
      channels: []
    }
  },
  computed: {},
  async mounted() {
    this.channels = config.CHANNELS
    this.loadData()
  },
  methods: {
    tableRowClassName({ row, rowIndex }) {
      console.log(row)
      console.log(rowIndex)
      if (row.channel === 'googleplay') {
        return 'ignore-row'
      }
      if (row.same_ip >= config.WARNING_SAME_IP || row.same_device_no >= config.WARNING_SAME_DEVICE) {
        console.log('warning')
        return 'warning-row'
      }
      return ''
    },
    currentChange(e) {
      this.page = e
      this.loadData()
    },
    async loadData() {
      this.table_loading = true
      var params = makeAjaxParamData(this.searchForm, {
        filter: this.searchForm.filter,
        page: this.page,
        size: this.pagesize,
        from: this.searchForm.recharge_start,
        to: this.searchForm.recharge_end
      })
      const res = await coinApplicationList(params)
      this.totalnum = res.total
      this.list = res.list.map(item => {
        item.telegram = JSON.parse(item.user_auth?.auth_ext || '{}')
        return item
      })
      this.table_loading = false
    },
    sortChange(col) {
      if (col.column.columnKey) {
        this.order_by = col.column.columnKey
      } else {
        this.order_by = ''
      }
      switch (col.order) {
        case 'ascending':
          this.order = 'asc'
          break
        case 'descending':
          this.order = 'desc'
          break
        default:
          this.order = 'desc'
          break
      }
      this.page = 1
      this.loadData()
    },
    gotoUserDetail(user_id) {
      // this.$router.push({
      //   name: 'user-detail',
      //   params: {
      //     id: user_id
      //   }
      // })
      window.open(this.$router.resolve({
        name: 'user-detail',
        params: {
          id: user_id
        }
      }).href, '_blank')
    },
    doSearch() {
      if (this.searchForm.recharge_time) {
        this.searchForm.recharge_start = this.searchForm.recharge_time[0]
        this.searchForm.recharge_end = this.searchForm.recharge_time[1]
      }
      console.log(this.searchForm)
      this.page = 1
      this.loadData()
    },
    doReset() {
      this.searchForm = {
        nickname: '',
        channel: '',
        register_time: [],
        online_time: [],
        register_start: '',
        register_end: '',
        online_start: '',
        online_end: ''
      }
      this.loadData()
    },
    async examineChange(row) {
      await examineCoinApplication({
        id: row.id,
        status: row.status
      })
      this.loadData()
    }
  }
}

</script>
<style>
  .el-table .warning-row {
    background: #e47470;
  }
  .el-table .ignore-row {
    background: #e7c650;
  }
</style>
