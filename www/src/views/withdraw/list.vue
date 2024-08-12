<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">提现审核</h4>
    </el-header>
    <el-main>
      <el-form v-model="searchForm" label-width="100px">
        <el-row>
          <el-col :span="10">
            <el-form-item label="用户ID">
              <el-input v-model.number="searchForm.uid" placeholder="请输入用户ID" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="10">
            <el-form-item label="提现地址">
              <el-input v-model="searchForm.address" placeholder="请输入提现地址" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="提现币种">
          <el-select v-model="searchForm.atypes" multiple placeholder="请选择">
            <el-option
              v-for="atype in assetsOptions"
              :key="atype.value"
              :label="atype.label"
              :value="atype.value"
            />
          </el-select>
        </el-form-item>
        <el-row>
          <el-col :span="10">
            <el-form-item label="提现数额">
              <el-row>
                <el-col :span="8">
                  <el-input v-model.number="searchForm.amount_gt" placeholder="最小金额 >" />
                </el-col>
                <el-col :span="2" style="text-align:center;">
                  {{ '  至' }}
                </el-col>
                <el-col :span="8">
                  <el-input v-model.number="searchForm.amount_lt" placeholder="最大金额 <" />
                </el-col>
              </el-row>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="提现状态">
          <el-select v-model="searchForm.status" placeholder="请选择">
            <el-option
              v-for="tag in status_option"
              :key="tag.value"
              :label="tag.label"
              :value="tag.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="提现时间">
          <el-date-picker
            v-model="searchForm.withdraw_time"
            type="daterange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            placeholder="选择时间范围"
            value-format="yyyy-MM-dd"
            :picker-options="pickerOptions"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="doSearch">搜索</el-button>
          <el-button type="warning" icon="el-icon-refresh-left" @click="doReset">重置</el-button>
        </el-form-item>
      </el-form>
      <el-row type="flex" justify="end">
        <el-button type="primary" icon="el-icon-download" @click="exportLog">导出表格</el-button>
      </el-row>
      <el-divider />
      <el-divider />
      <el-table
        ref="multipleTable"
        v-loading="table_loading"
        max-height="800"
        :data="list"
      >
        <el-table-column
          type="selection"
          width="55"
          fixed="left"
        />
        <el-table-column
          prop="id"
          label="ID"
          width="100"
        />
        <el-table-column
          prop="user_id"
          label="用户ID"
          width="100"
        >
          <template slot-scope="scope">
            <el-button type="text" @click="gotoUserDetail(scope.row.user_id)">
              {{ scope.row.user_id }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column
          prop="user.display_name"
          label="昵称"
          width="120"
        />
        <el-table-column
          prop="symbol"
          label="币种"
        />
        <el-table-column
          prop="chain"
          label="提币网络"
        />
        <el-table-column
          prop="amount"
          label="数量"
        />
        <el-table-column
          prop="address"
          label="地址"
        />
        <el-table-column
          prop="status"
          label="状态"
          width="120"
        >
          <template slot-scope="scope">
            <el-tag
              v-if="scope.row.handle_status === 0"
              type="primary"
            >
              待审核
            </el-tag>
            <el-tag
              v-if="scope.row.handle_status === 1"
              type="warning"
            >
              处理中
            </el-tag>
            <el-tag
              v-if="scope.row.handle_status === 2"
              type="warning"
            >
              待确认
            </el-tag>
            <el-tag
              v-if="scope.row.handle_status === 3"
              type="success"
            >
              提现成功
            </el-tag>
            <el-tag
              v-if="scope.row.handle_status === 4"
              type="danger"
            >
              提现失败
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="create_at ｜ times"
          label="申请时间"
          width="180"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.create_at | utcTimeFormat }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="outserial"
          label="订单号"
          width="100"
        />
        <el-table-column
          prop="tx_id"
          label="TxHash"
          width="300"
        />
        <el-table-column
          prop="reason"
          label="原因"
          width="150"
        />
      </el-table>
      <el-pagination
        background
        layout="sizes, prev, pager, next, jumper, ->, total"
        :current-page="page"
        :page-sizes="[20, 50, 100]"
        :page-size="pagesize"
        :total="totalnum"
        @size-change="handleSizeChange"
        @current-change="currentPageChange"
      />
    </el-main>
  </el-container>
</template>
<script>
import { getWithdrawList, downloadWithdrawLog } from '@/api/withdraw'
import { allAssets } from '@/api/coin_application'
export default {
  data() {
    return {
      selected: [],
      list: [],
      pagesize: 20,
      totalnum: 0,
      page: 1,
      table_loading: false,
      searchForm: {
        uid: '',
        status: -1,
        withdraw_time: [],
        withdraw_start: '',
        withdraw_end: '',
        atypes: [],
        from: '',
        to: '',
        amount_gt: null,
        amount_lt: null
      },
      formRules: {
      },
      status_option: [
        {
          label: '全部',
          value: -1
        },
        {
          label: '待审核',
          value: 0
        },
        {
          label: '处理中',
          value: 1
        },
        {
          label: '待确认',
          value: 2
        },
        {
          label: '提现成功',
          value: 3
        },
        {
          label: '提现失败',
          value: 4
        }
      ],
      pickerOptions: {
        shortcuts: [
          {
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
      }
    }
  },
  computed: {
  },
  async mounted() {
    this.loadData()
    this.loadAllAssets()
  },
  methods: {
    exportLog() {
      downloadWithdrawLog({
        status: this.searchForm.status,
        uid: this.searchForm.uid || 0,
        address: this.searchForm.address,
        from: this.searchForm.from,
        to: this.searchForm.to,
        atypes: this.searchForm.atypes,
        amount_gt: this.searchForm.amount_gt,
        amount_lt: this.searchForm.amount_lt
      })
    },
    async loadAllAssets() {
      const res = await allAssets({ is_system: true })
      this.assetsOptions = res.list.map(item => {
        return {
          value: item.a_type,
          label: item.name
        }
      })
    },
    handleSizeChange(val) {
      this.pagesize = val
      this.page = 1
      this.loadData()
    },
    currentPageChange(e) {
      this.page = e
      this.loadData()
    },
    async loadData() {
      console.log('this.searchForm.atypes', this.searchForm.atypes)
      this.table_loading = true
      const res = await getWithdrawList({
        page: this.page,
        size: this.pagesize,
        status: this.searchForm.status,
        uid: this.searchForm.uid || 0,
        address: this.searchForm.address,
        from: this.searchForm.from,
        to: this.searchForm.to,
        atypes: this.searchForm.atypes,
        amount_gt: this.searchForm.amount_gt,
        amount_lt: this.searchForm.amount_lt
      })
      this.totalnum = res.total
      this.list = res.list
      this.table_loading = false
    },
    doSearch() {
      if (this.searchForm.withdraw_time) {
        this.searchForm.from = this.searchForm.withdraw_time[0]
        this.searchForm.to = this.searchForm.withdraw_time[1]
      }
      this.page = 1
      this.loadData()
    },
    doReset() {
      this.searchForm = {
        status: -1,
        atypes: [],
        amount_gt: null,
        amount_lt: null
      }
      this.loadData()
    }
  }
}
</script>
<style>
  .el-table .success-row {
    background: #52d40c;
  }
 .el-table .warning-row {
    background: #e47470;
  }
 .el-table .ignore-row {
    background: #e7c650;
  }
  .el-table__body tr.hover-row.current-row>td, .el-table__body tr.hover-row.el-table__row--striped.current-row>td, .el-table__body tr.hover-row.el-table__row--striped>td, .el-table__body tr.hover-row>td{
    background-color:transparent;
  }
</style>
