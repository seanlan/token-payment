<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">质押配置</h4>
    </el-header>
    <el-main>
      <el-row type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="showEditModal({})">添加配置</el-button>
      </el-row>
      <el-table v-loading="table_loading" :data="list">
        <el-table-column prop="id" label="ID" width="100" />
        <el-table-column prop="in_a_type" label="质押币种" width="180" />
        <el-table-column prop="out_a_type" label="奖励币种" />
        <el-table-column prop="out_per_second" label="每秒奖励" />
        <el-table-column prop="ranking" label="排序值" />
        <el-table-column prop="status" width="100" label="无效过滤" fixed="right">
          <template slot-scope="scope">
            <el-button size="mini" type="primary" @click="showEditModal(scope.row)">编辑</el-button>
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
    <el-dialog
      :title="!!editForm.id ? '修改配置': '添加配置'"
      :visible.sync="editDialogVisible"
      width="50%"
    >
      <el-form ref="editForm" :model="editForm" :rules="editFormRules" label-width="140px">
        <el-form-item label="质押币种" prop="in_a_type">
          <el-select
            v-model="editForm.in_a_type"
            placeholder="请选择质押币种"
          >
            <el-option
              v-for="item in assetsOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="奖励币种" prop="out_a_type">
          <el-select
            v-model="editForm.out_a_type"
            placeholder="请选择质押币种"
          >
            <el-option
              v-for="item in assetsOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="每秒奖励" prop="out_per_second">
          <el-input
            ref="weight"
            v-model="editForm.out_per_second"
            tabindex="1"
            placeholder="每秒奖励 (实际奖励=每秒奖励*质押数量*质押时间)"
          />
        </el-form-item>
        <el-form-item label="排序值" prop="ranking">
          <el-input
            ref="weight"
            v-model.number="editForm.ranking"
            tabindex="1"
            placeholder="排序值 (值越大越靠前)"
          />
        </el-form-item>

      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="editDialogVisible = false">取 消</el-button>
        <el-button v-if="!!editForm.id" type="primary" @click="doSubmit">编 辑</el-button>
        <el-button v-else type="primary" @click="doSubmit">添 加</el-button>
      </div>
    </el-dialog>
  </el-container>
</template>
<script>
import _ from 'lodash'
import {
  makeAjaxParamData
} from '@/utils/util'
import { Message } from 'element-ui'
import { allAssets } from '@/api/coin_application'
import { getPledgeTemplates, savePledgeTemplate } from '@/api/pledge'
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
      editForm: {},
      editFormRules: {
        in_a_type: [{ required: true, trigger: 'blur', message: '质押币种不能为空' }],
        out_a_type: [{ required: true, trigger: 'blur', message: '奖励币种不能为空' }],
        out_per_second: [{ required: true, trigger: 'blur', message: '每秒奖励不能为空' }]
      },
      assetsOptions: [],
      editDialogVisible: false
    }
  },
  computed: {},
  async mounted() {
    this.loadData()
    this.loadAllAssets()
  },
  methods: {
    async loadAllAssets() {
      const res = await allAssets({ is_system: true })
      console.log('allAssets', res)
      const assetsOptions = res.list.map(item => {
        return {
          value: item.a_type,
          label: item.name
        }
      })
      this.assetsOptions = assetsOptions
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
      const res = await getPledgeTemplates(params)
      console.log(res)
      this.totalnum = res.total
      this.list = res.list
      this.table_loading = false
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
    showEditModal(row) {
      this.editDialogVisible = true
      this.editForm = {
        ...(row || {}),
        out_a_type: row.out_a_type || 'diamond'
      }
    },
    doSubmit() {
      this.$refs.editForm.validate(valid => {
        if (valid) {
          const configInfo = {
            ...(_.cloneDeep(this.editForm)),
            out_per_second: parseFloat(this.editForm.out_per_second) || 0
          }

          savePledgeTemplate(configInfo).then(() => {
            Message({
              message: !this.editForm.id ? '添加成功' : '修改成功',
              type: 'success'
            })
            this.editDialogVisible = false
            this.loadData()
          }).catch((e) => {
            console.log(e)
          })
        } else {
          console.log('error submit!!')
          return false
        }
      })
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
