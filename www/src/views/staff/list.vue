<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">员工列表</h4>
    </el-header>
    <el-main>
      <el-row v-if="checkPermission(['user_add'])" type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="gotoStaffAdd(0)">添加员工</el-button>
      </el-row>
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
        />
        <el-table-column
          prop="account"
          label="账号"
          width="200"
        >
          <template slot-scope="scope">
            <span> {{ scope.row.account }} </span>
            <el-tag v-if="scope.row.is_super===1" type="danger">超级用户</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="username"
          label="姓名"
        />
        <el-table-column
          prop="phone"
          label="手机号"
        />
        <el-table-column
          prop="email"
          label="Email"
        />
        <el-table-column
          prop="group_name"
          label="角色"
        />
        <el-table-column prop="is_frozen" width="100" label="封禁">
          <template slot-scope="scope">
            <el-switch
              v-model="scope.row.is_frozen"
              active-color="#F56C6C"
              inactive-color="#409EFF"
              :active-value="1"
              :inactive-value="0"
              :disabled="scope.row.is_super== 1 || !checkPermission(['user_edit'])"
              @change="frozenChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column
          prop="create_at"
          label="创建时间"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.create_at | toDate }}</span>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="200">
          <template slot-scope="scope">
            <el-button v-if="checkPermission(['user_edit'])" type="primary" size="mini" :disabled="!!scope.row.is_super" @click="gotoStaffAdd(scope.row.id)">编辑</el-button>
            <el-button v-if="checkPermission(['user_edit'])" type="warning" size="mini" :disabled="!!scope.row.is_super" @click="gotoRepassword(scope.row)">修改密码</el-button>
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
    <el-dialog
      title="修改密码"
      :visible.sync="repasswordDialogVisible"
      width="40%"
    >
      <el-form ref="repasswordForm" :model="repasswordForm" label-width="140px" :rules="repasswordFormRules">
        <el-form-item label="账号" prop="account">
          <el-input v-model="repasswordForm.account" disabled />
        </el-form-item>
        <el-form-item label="姓名" prop="username">
          <el-input v-model="repasswordForm.username" disabled />
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input v-model="repasswordForm.password" type="password" placeholder="请输入新密码" @keyup.enter.native="rePassword" />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm_password">
          <el-input v-model="repasswordForm.confirm_password" type="password" placeholder="请输入确认密码" @keyup.enter.native="rePassword" />
        </el-form-item>
        <el-form-item class="form-footer">
          <el-button type="primary" @click.native.prevent="rePassword">修改密码</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </el-container>
</template>

<script>
import { mapGetters } from 'vuex'
import { getStaffList, staffFrozen, staffRePassword } from '@/api/staff'
import { hasPermission } from '@/utils/permission'

export default {
  components: {
  },
  data() {
    var validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.repasswordForm.password) {
        callback(new Error('两次输入密码不一致!'))
      } else {
        callback()
      }
    }
    return {
      list: [],
      current_page: 1,
      page_size: 20,
      total: 0,
      table_loading: false,
      repasswordDialogVisible: false,
      repasswordForm: {
        staff_id: 0,
        password: '',
        confirm_password: ''
      },
      repasswordFormRules: {
        password: [
          { required: true, message: '请输入登录密码', trigger: 'blur' }
        ],
        confirm_password: [
          { validator: validatePass2, required: true, trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    ...mapGetters([
      'roles'
    ])
  },
  async mounted() {
    this.loadData()
    console.log(this.roles)
  },
  methods: {
    checkPermission(args) {
      return hasPermission(args)
    },
    currentPageChange(page) {
      this.current_page = page
      this.loadData()
    },
    async loadData() {
      this.table_loading = true
      const res = await getStaffList({
        page: this.current_page,
        size: this.page_size
      })
      console.log(res)
      this.total = res.data.total
      this.list = res.data.list
      this.table_loading = false
    },
    gotoStaffAdd(id) {
      this.$router.push({ name: 'staff-edit', params: { id }})
    },
    gotoRepassword(row) {
      this.repasswordDialogVisible = true
      this.repasswordForm = {
        staff_id: row.id,
        account: row.account,
        username: row.username,
        password: '',
        confirm_password: ''
      }
      console.log(this.repasswordForm)
    },
    rePassword() {
      this.$refs.repasswordForm.validate(async valid => {
        if (valid) {
          await staffRePassword({
            staff_id: this.repasswordForm.staff_id,
            password: this.repasswordForm.password
          })
          this.repasswordDialogVisible = false
          this.$message({
            message: '修改密码成功',
            type: 'success'
          })
          this.loadData()
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    async frozenChange(row) {
      await staffFrozen({
        staff_id: row.id,
        is_frozen: row.is_frozen
      })
      this.loadData()
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
