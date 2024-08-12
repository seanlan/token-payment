<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">{{ staff_id ? '修改信息':'添加员工' }}</h4>
    </el-header>
    <el-main>
      <el-form ref="editForm" :model="editForm" :rules="formRules" label-width="150px">
        <div class="formtitle">基本信息</div>
        <el-row>
          <el-col :xs="{span:24}" :sm="{span: 8}">
            <el-form-item label="账号" prop="account">
              <el-input v-model="editForm.account" placeholder="请输入账号" name="account" :disabled="staff_id>0" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row v-if="staff_id==0">
          <el-col :xs="{span:24}" :sm="{span: 8}">
            <el-form-item label="密码" prop="password">
              <el-input v-model="editForm.password" placeholder="请输入密码" type="password" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :xs="{span:24}" :sm="{span: 8}">
            <el-form-item label="姓名" prop="username">
              <el-input v-model="editForm.username" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :xs="{span:24}" :sm="{span: 8}">
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="editForm.phone" placeholder="请输入手机号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :xs="{span:24}" :sm="{span: 8}">
            <el-form-item label="Email" prop="email">
              <el-input v-model="editForm.email" placeholder="请输入Email" />
            </el-form-item>
          </el-col>
        </el-row>
        <div class="formtitle">员工组</div>
        <el-row>
          <el-col :xs="{span:24}" :sm="{span: 8}">
            <el-form-item label="员工角色" prop="member_base">
              <el-select v-model="editForm.group" placeholder="请选择员工角色" @change="groupChanged">
                <el-option
                  v-for="g in groups"
                  :key="g.id"
                  :label="g.name"
                  :value="g.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <div class="formtitle">权限分配</div>
        <el-form-item>
          <PermissionSelector ref="pselector" />
        </el-form-item>
        <el-form-item class="form-footer">
          <el-button type="primary" @click="submitForm">{{ staff_id ? '保存信息':'立即创建' }}</el-button>
        </el-form-item>
      </el-form>
    </el-main>
  </el-container>
</template>

<script>
import PermissionSelector from './components/PermissionSelector'
import { getAllGroups, addStaff, editStaff, staffInfo } from '@/api/staff'
export default {
  components: {
    PermissionSelector
  },
  data() {
    return {
      staff_id: 0,
      editForm: {
        account: '',
        username: '',
        phone: '',
        email: '',
        password: '',
        group: 0,
        permissions: []
      },
      groups: [
        {
          id: 0,
          name: '无',
          permissions: []
        },
        {
          id: 1,
          name: '普通员工',
          permissions: [1]
        },
        {
          id: 2,
          name: '管理员',
          permissions: [1, 2, 3, 4]
        }
      ],
      formRules: {
        account: [
          { required: true, message: '请输入登录账户', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入登录密码', trigger: 'blur' }
        ],
        username: [
          { required: true, message: '请输入姓名', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
  },
  async mounted() {
    this.staff_id = this.$route.params.id ? parseInt(this.$route.params.id) : 0
    this.loadGroups()
    if (this.staff_id > 0) {
      await this.loadData()
    }
  },
  methods: {
    async loadData() {
      const { data } = await staffInfo({ staff_id: this.staff_id })
      this.editForm = {
        account: data.staff.account,
        username: data.staff.username,
        phone: data.staff.phone,
        email: data.staff.email,
        password: '',
        group: data.staff.group
      }
      this.$refs.pselector.changeChecked(data.staff.permissions)
      this.$refs.pselector.changeDisabled(data.staff.group_permissions)
    },
    async loadGroups() {
      const res = await getAllGroups()
      console.log(res)
      this.groups = [
        {
          id: 0,
          name: '无',
          permissions: []
        }, ...res.data.list]
    },
    groupChanged() {
      let permissions = []
      this.groups.forEach(g => {
        if (g.id === this.editForm.group) {
          permissions = g.permissions
        }
      })
      this.$refs.pselector.changeDisabled(permissions)
    },
    async submitForm() {
      this.$refs.editForm.validate(async valid => {
        if (valid) {
          this.table_loading = true
          this.editForm.permissions = this.$refs.pselector.getChecked()
          if (this.staff_id === 0) {
            await addStaff(this.editForm)
          } else {
            await editStaff({ ...this.editForm, staff_id: this.staff_id })
          }
          this.$message({
            type: 'success',
            message: '保存成功'
          })
          this.$router.push({ name: 'staff-list' })
        }
      })
    }
  }
}

</script>

<style lang="scss" scoped>
.formtitle {
  margin: 20px;
  font-size: 14px;
  font-weight: bold;
  color: #333;
}
.form-footer {
  margin-top: 50px;
}
</style>
