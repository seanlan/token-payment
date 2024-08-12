<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">用户管理</h4>
    </el-header>
    <el-main>
      <el-form class="searchForm" label-width="120px">
        <el-row>
          <el-col :span="8">
            <el-form-item label="用户昵称">
              <el-input v-model="searchForm.display_name" placeholder="请输入昵称" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="注册时间">
              <el-date-picker
                v-model="searchForm.register_time"
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
        <el-row>
          <el-col :span="16">
            <el-form-item label="用户ID">
              <el-input v-model="searchForm.ids" type="textarea" placeholder="请输入用户ID多个用,分割" rows="4" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="doSearch">搜索</el-button>
          <el-button type="warning" icon="el-icon-refresh-left" @click="doReset">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="table_loading" :data="list" @expand-change="tableExpand">
        <el-table-column type="expand">
          <template slot-scope="props">
            <el-header>
              <h4 class="links">{{ props.row.user.display_name }} 的资产</h4>
            </el-header>
            <el-table
              v-loading="props.row.assets_loading"
              border
              :data="props.row.assets"
            >
              <el-table-column
                prop="a_type"
                label="资产名称"
              />
              <el-table-column
                prop="amount"
                label="资产数量"
              />
              <el-table-column
                prop="is_lock"
                label="是否冻结"
              >
                <template slot-scope="scope">
                  <el-tag v-if="scope.row.is_lock" type="danger">是</el-tag>
                  <el-tag v-else type="success">否</el-tag>
                </template></el-table-column>
            </el-table>
          </template>
        </el-table-column>
        <el-table-column prop="user.id" label="用户ID" />
        <el-table-column label="头像">
          <template slot-scope="scope">
            <div class="tablecell">
              <CustomImage
                style="width: 50px; height: 50px"
                :src="scope.row.user.photo_url ? scope.row.user.photo_url : 'https://goplay-cdn.kanlian.la/goplay_avatar.png'"
                fit="contain"
                class="shoplogo"
                :preview-src-list="[scope.row.user.photo_url ? scope.row.user.photo_url : 'https://goplay-cdn.kanlian.la/goplay_avatar.png']"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="user.display_name" label="昵称" />
        <el-table-column prop="telegram.id" label="TG ID" />
        <el-table-column prop="telegram.username" label="TG Account" />
        <el-table-column prop="create_at" label="创建时间">
          <template slot-scope="scope">
            <span>{{ scope.row.user.create_at | utcTimeFormat }}</span>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        layout="sizes, prev, pager, next, jumper, ->, total"
        :current-page="page"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pagesize"
        :total="totalnum"
        @size-change="handleSizeChange"
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
import { userList, userAssets } from '@/api/user'
export default {
  data() {
    return {
      list: [],
      pagesize: 20,
      totalnum: 0,
      page: 1,
      table_loading: false,
      searchForm: {
        ids: '',
        display_name: '',
        register_time: [],
        register_start: '',
        register_end: ''
      },
      formRules: {},
      formFilter: false,
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
      }
    }
  },
  computed: {},
  async mounted() {
    this.channels = config.CHANNELS
    this.loadData()
  },
  methods: {
    handleSizeChange(val) {
      this.pagesize = val
      this.page = 1
      this.loadData()
    },
    tableRowClassName({ row, rowIndex }) {
      if (row.is_warning > 0) {
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
      let ids = []
      if (this.searchForm.ids) {
        ids = this.searchForm.ids.split(',').map(item => {
          return parseInt(item)
        })
      }
      var params = makeAjaxParamData(this.searchForm, {
        page: this.page,
        size: this.pagesize,
        order_by: this.order_by,
        order: this.order,
        user_ids: ids
      })
      const res = await userList(params)
      this.totalnum = res.total
      this.list = res.list.map(item => {
        item.assets = []
        item.assets_loading = false
        item.telegram = JSON.parse(item.auth.auth_ext)
        return item
      })
      console.log('this.list', this.list)
      this.table_loading = false
    },
    async tableExpand(row, expand) {
      console.log(row, expand)
      if (expand) {
        row.assets_loading = true
        const res = await userAssets({
          uid: row.user.id
        })
        // eslint-disable-next-line require-atomic-updates
        row.assets = res.list
        // eslint-disable-next-line require-atomic-updates
        row.assets_loading = false
      }
    },
    doSearch() {
      console.log(this.searchForm)
      if (this.searchForm.register_time) {
        this.searchForm.register_start = this.searchForm.register_time[0]
        this.searchForm.register_end = this.searchForm.register_time[1]
      }
      if (this.searchForm.online_time) {
        this.searchForm.online_start = this.searchForm.online_time[0]
        this.searchForm.online_end = this.searchForm.online_time[1]
      }
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
  span.danger {
    color: #F56C6C;
  }
  span.warning {
    color: #E6A23C;
  }
  span.success {
    color: #67C23A;
  }
  .el-button {
    margin-bottom: 5px;
  }
</style>
