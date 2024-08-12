<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">员工角色</h4>
    </el-header>
    <el-main>
      <el-row v-if="checkPermission(['group_add'])" type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="showDialog('add')">添加角色</el-button>
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
          prop="name"
          label="角色名称"
        />
        <el-table-column
          prop="name"
          label="操作"
        >
          <template slot-scope="scope">
            <el-button v-if="checkPermission(['group_edit'])" type="primary" size="mini" @click="showDialog('edit', scope.row)">编辑</el-button>
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
      :title="modaltype=='edit'? '编辑角色': '添加角色'"
      :visible.sync="editDialogVisible"
      width="50%"
    >
      <el-form ref="editForm" :model="editForm" label-width="140px" :rules="editFormRules">
        <el-form-item label="角色名" prop="name">
          <el-input v-model="editForm.name" class="simple_input" />
        </el-form-item>
        <el-form-item label="权限分配" prop="permissions">
          <PermissionSelector ref="pselector" class="permission-selector" />
        </el-form-item>
        <el-form-item class="form-footer">
          <el-button type="primary" @click="submitEditForm"> 确定 </el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </el-container>
</template>

<script>
import { mapGetters } from 'vuex'
import { getAllGroups, addGroups, editGroups } from '@/api/staff'
import { hasPermission } from '@/utils/permission'
import PermissionSelector from './components/PermissionSelector'

export default {
  components: {
    PermissionSelector
  },
  data() {
    return {
      list: [],
      current_page: 1,
      page_size: 20,
      total: 0,
      table_loading: false,
      modaltype: 'add',
      editDialogVisible: false,
      editForm: {
        group_id: '',
        name: ''
      },
      editFormRules: {
        name: [
          { required: true, message: '请输入角色名', trigger: 'blur' }
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
      const res = await getAllGroups({
        page: this.current_page,
        size: this.page_size
      })
      console.log(res)
      this.total = res.data.total
      this.list = res.data.list
      this.table_loading = false
    },
    showDialog(type, row) {
      this.modaltype = type
      this.editDialogVisible = true
      if (type === 'add') {
        this.editForm = {
          name: ''
        }
        setTimeout(() => {
          this.$refs.pselector.cleareChecked()
        }, 50)
      } else {
        this.editForm = {
          group_id: row.id,
          name: row.name
        }
        setTimeout(() => {
          this.$refs.pselector.changeChecked(row.permissions)
        }, 50)
      }
    },
    submitEditForm() {
      this.$refs.editForm.validate(async valid => {
        if (valid) {
          this.editDialogVisible = false
          const param = {
            group_id: this.editForm.group_id,
            name: this.editForm.name,
            permissions: this.$refs.pselector.getChecked()
          }
          if (this.modaltype === 'add') {
            await addGroups(param)
          } else {
            await editGroups(param)
          }
          this.$message({
            message: '提交成功',
            type: 'success'
          })
          this.loadData()
        }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.permission-selector {
  margin-top: 20px;
  // margin-left: 20px;
}
.simple_input {
  width: 40%;
}
</style>
