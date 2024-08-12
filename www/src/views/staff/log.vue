<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">员工日志</h4>
    </el-header>
    <el-main>
      <el-form class="searchForm" label-width="120px">
        <el-row>
          <el-col>
            <el-form-item label="员工">
              <el-select v-model="searchForm.staff_id" filterable placeholder="请选择">
                <el-option
                  v-for="staff in staffs"
                  :key="staff.id"
                  :label="staff.username"
                  :value="staff.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="8">
            <el-form-item label="记录时间">
              <el-date-picker
                v-model="searchForm.date_time"
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
      <el-divider />
      <el-table
        ref="multipleTable"
        v-loading="table_loading"
        :data="list"
        tooltip-effect="dark"
        style="width: 100%"
      >
        <el-table-column
          prop="id"
          label="ID"
          width="100"
        />
        <el-table-column
          prop="staff.account"
          label="账号"
          width="200"
        >
          <template slot-scope="scope">
            <span> {{ scope.row.staff.account }} </span>
            <el-tag v-if="scope.row.staff.is_super===1" type="danger">超级用户</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="staff.username"
          label="姓名"
          width="200"
        />
        <el-table-column
          prop="remark"
          label="记录"
        />
        <el-table-column
          prop="create_at"
          label="记录时间"
          width="200"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.create_at | toDate }}</span>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        hide-on-single-page
        layout="prev, pager, next, jumper, ->, total"
        :current-page="current_page"
        :page-size="page_size"
        :total="total"
        @current-change="currentPageChange"
      />
    </el-main>
  </el-container>
</template>

<script>
import { staffLogs, getStaffList } from '@/api/staff'
import { DatePickerOptions } from '@/utils/config'
export default {
  components: {
  },
  data() {
    return {
      list: [],
      current_page: 1,
      page_size: 20,
      total: 0,
      table_loading: false,
      searchForm: {
        date_time: [],
        date_start: '',
        date_end: ''
      },
      staffs: [],
      pickerOptions: DatePickerOptions
    }
  },
  async mounted() {
    this.loadStaffs()
    this.loadData()
  },
  methods: {
    currentPageChange(page) {
      this.current_page = page
      this.loadData()
    },
    async loadStaffs() {
      const res = await getStaffList()
      this.staffs = res.data.list
    },
    async loadData() {
      this.table_loading = true
      const res = await staffLogs({
        staff_id: this.searchForm.staff_id,
        date_start: this.searchForm.date_start,
        date_end: this.searchForm.date_end,
        page: this.current_page,
        size: this.page_size
      })
      console.log(res)
      this.total = res.data.total
      this.list = res.data.list
      this.table_loading = false
    },
    doReset() {
      this.searchForm = {
        date_time: []
      }
      this.loadData()
    },
    doSearch() {
      console.log(this.searchForm)
      if (this.searchForm.date_time) {
        this.searchForm.date_start = this.searchForm.date_time[0]
        this.searchForm.date_end = this.searchForm.date_time[1]
      }
      this.current_page = 1
      this.loadData()
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
