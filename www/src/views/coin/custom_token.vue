<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">创业板币种</h4>
    </el-header>
    <el-main>
      <el-form class="searchForm" label-width="120px">
        <el-row>
          <el-col :span="8">
            <el-form-item label="币种">
              <el-input v-model="searchForm.name" placeholder="请输入币种名称" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="用户ID">
              <el-input v-model.number="searchForm.uid" placeholder="请输入用户ID" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="loadData">搜索</el-button>
          <el-button type="warning" icon="el-icon-refresh-left" @click="doReset">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="table_loading" :data="list" :row-class-name="tableRowClassName" @expand-change="tableExpand">
        <el-table-column type="expand" label="" width="100">
          <template slot-scope="props">
            <el-table
              v-loading="props.row.stats_loading"
              border
              :data="props.row.stats"
            >
              <el-table-column prop="date_stamp" label="资产名称" />
              <el-table-column prop="a_type" label="资产名称" />
              <el-table-column prop="holders" label="持币用户" />
              <el-table-column prop="new_holders" label="新增持币用户" />
              <el-table-column prop="order_new" label="新增挂单" />
              <el-table-column prop="trans_order" label="成单数" />
              <el-table-column prop="trans_value" label="交易量" />
              <el-table-column prop="trans_amount" label="成单额(USDT)" />
              <el-table-column prop="recharge_order" label="充值订单" />
              <el-table-column prop="recharge_amount" label="充值数量" />
              <el-table-column prop="withdraw_order" label="提币订单" />
              <el-table-column prop="withdraw_amount" label="提币数量" />
            </el-table>
          </template>

        </el-table-column>
        <el-table-column prop="assets.id" label="ID" width="100" />
        <el-table-column prop="user.id" label="User ID" width="100" />
        <el-table-column label="图标" width="100">
          <template slot-scope="scope">
            <div class="tablecell">
              <CustomImage
                style="width: 50px; height: 50px"
                :src="scope.row.assets.icon_url ? scope.row.assets.icon_url : ''"
                fit="contain"
                class="shoplogo"
                :preview-src-list="[scope.row.assets.icon_url ? scope.row.assets.icon_url : '']"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="user.display_name" label="昵称" />
        <el-table-column prop="telegram.id" label="TG ID" />
        <el-table-column prop="telegram.username" label="TG账号" />
        <el-table-column prop="assets.whatsapp" label="Whatsapp" />
        <el-table-column prop="assets.name" label="币种" />
        <el-table-column prop="assets.a_type" label="币种代码" />
        <el-table-column
          prop="introduce"
          label="简介"
          width="300"
        >
          <template slot-scope="scope">
            {{ scope.row.assets.introduce.length > 20 ? scope.row.assets.introduce.substring(0, 20) + '...' : scope.row.assets.introduce }} <br>
            <a v-if="scope.row.assets.introduce.length > 20" class="primary-span" @click="showRemakrs(scope.row.assets.introduce)">查看更多</a>
          </template>
        </el-table-column>
        <el-table-column prop="assets.white_paper" label="白皮书" />
        <el-table-column prop="assets.site" label="网址" />
        <el-table-column prop="status" width="100" label="操作" fixed="right">
          <template slot-scope="scope">
            <el-button type="danger" size="mini" @click="topUp(scope.row)"> 置顶 </el-button>
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
import { statsCustomAssets } from '@/api/stats'
import { getCustomTokenList, topUpCustomToken } from '@/api/coin_application'
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
        page: this.page,
        size: this.pagesize
      })
      const res = await getCustomTokenList(params)
      this.totalnum = res.total
      this.list = res.list.map(item => {
        item.telegram = JSON.parse(item.user_auth?.auth_ext || '{}')
        item.stats = []
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
      this.searchForm = {}
      this.loadData()
    },
    async topUp(row) {
      await topUpCustomToken({
        id: row.assets.id
      })
      this.loadData()
    },
    showRemakrs(remarks) {
      this.$alert('<div style="white-space: pre-line;"> ' + remarks + '</div>', '简介', {
        confirmButtonText: '确定',
        closeOnClickModal: true,
        closeOnPressEscape: true,
        dangerouslyUseHTMLString: true
      })
    },
    async tableExpand(row, expand) {
      if (expand) {
        row.assets_loading = true
        const res = await statsCustomAssets({
          a_type: row.assets.a_type,
          page: 1,
          size: 20
        })
        // eslint-disable-next-line require-atomic-updates
        row.stats = res.list
        // eslint-disable-next-line require-atomic-updates
        row.total = res.total
        // eslint-disable-next-line require-atomic-updates
        row.stats_loading = false
      }
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
  .primary-span {
    color: #409EFF;
  }
  .text-danger {
    color: #F56C6C;
  }
</style>
