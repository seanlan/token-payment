<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">成交记录</h4>
    </el-header>
    <el-main>
      <el-form class="searchForm" label-width="120px">
        <el-row>
          <el-col :span="8">
            <el-form-item label="用户ID">
              <el-input v-model.number="searchForm.uid" placeholder="请输入用户ID" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="币种">
              <el-select v-model="searchForm.a_type" placeholder="币种">
                <el-option
                  v-for="item in assetsOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="成交时间">
              <el-date-picker
                v-model="searchForm.create_time"
                type="daterange"
                range-separator="至"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                placeholder="选择时间范围"
                value-format="yyyy-MM-dd"
                :picker-options="pickerOptions"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="doSearch">搜索</el-button>
          <el-button type="warning" icon="el-icon-refresh-left" @click="doReset">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="table_loading" :data="list" :row-class-name="tableRowClassName" @sort-change="sortChange">
        <el-table-column prop="id" label="ID" width="100" />
        <el-table-column prop="order_id" label="挂单ID" width="100" />
        <el-table-column prop="order_id" label="挂单用户[ID]" width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.order_user.display_name }}[{{ scope.row.order_user.id }}]</span>
          </template>
        </el-table-column>
        <el-table-column label="挂单数量/挂单额(USDT)">
          <template slot-scope="scope">
            <span>{{ scope.row.order.amount_max | floatNormalDisplay }}</span>/
            <span>{{ (scope.row.order.amount_max * scope.row.order.price) | floatNormalDisplay }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="user_id" label="User ID" width="100" />
        <el-table-column prop="user.display_name" label="昵称" width="180" />
        <el-table-column column-key="a_type" prop="a_type" label="币种" width="100" />
        <el-table-column prop="price" label="成交价格(USDT)">
          <template slot-scope="scope">
            <span>{{ scope.row.price | floatNormalDisplay }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="成交数量" />
        <el-table-column prop="amount" label="成交额(USDT)">
          <template slot-scope="scope">
            <span>{{ (scope.row.price * scope.row.amount) | floatNormalDisplay }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="create_at" label="成交时间(北京时间)" width="180">
          <template slot-scope="scope">
            <span>{{ scope.row.create_at | utcTimeFormat }}</span>
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
import { getOrderLogs } from '@/api/order'
import { allAssets } from '@/api/coin_application'
export default {
  data() {
    return {
      list: [],
      pagesize: 20,
      totalnum: 0,
      page: 1,
      table_loading: false,
      searchForm: {
        a_type: '',
        create_time: [],
        create_start: '',
        create_end: ''
      },
      formRules: {},
      formFilter: false,
      formRealPrice: 0,
      order_by: '`user`.`id`',
      order: 'desc',
      assetsOptions: [],
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
    this.loadAllAssets()
  },
  methods: {
    async loadAllAssets() {
      const res = await allAssets()
      const assetsOptions = res.list.map(item => {
        return {
          value: item.a_type,
          label: item.name
        }
      })
      this.assetsOptions = [{
        value: '',
        label: '全部'
      }, ...assetsOptions]
    },
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
        page: this.page,
        size: this.pagesize,
        from: this.searchForm.create_start,
        to: this.searchForm.create_end
      })
      const res = await getOrderLogs(params)
      this.totalnum = res.total
      this.list = res.list
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
      if (this.searchForm.create_time) {
        this.searchForm.create_start = this.searchForm.create_time[0]
        this.searchForm.create_end = this.searchForm.create_time[1]
      }
      console.log(this.searchForm)
      this.page = 1
      this.loadData()
    },
    doReset() {
      this.searchForm = {
        nickname: '',
        channel: '',
        create_time: [],
        online_time: [],
        create_start: '',
        create_end: '',
        online_start: '',
        online_end: ''
      }
      this.loadData()
    },
    async frozenChange(row) {
      await this.$store.dispatch('user/userFrozen', {
        uid: row.user_id,
        is_frozen: row.is_frozen
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
