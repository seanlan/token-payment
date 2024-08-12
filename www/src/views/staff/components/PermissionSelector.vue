<template>
  <el-checkbox-group v-model="checkedlist">
    <el-form>
      <el-form-item v-for="t in permissionTypes" :key="t.type">
        <el-form-item :label="t.type" />
        <div style="margin: 15px 0;" />
        <el-row align="middle" type="flex" style="flex-wrap: wrap">
          <el-col v-for="p in t.permissions" :key="p.id" :sm="{ span:6 }" :xs="{ span: 24}">
            <div class="permission-item">
              <el-checkbox :label="p.id" border :disabled="disabledlist.includes(p.id)">{{ p.name }}</el-checkbox>
            </div>
          </el-col>
        </el-row>
      </el-form-item>
    </el-form>
  </el-checkbox-group>
</template>

<script>
// import _ from 'lodash'
import { getAllPermissions } from '@/api/staff'

export default {
  props: {
    checked: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      permissionTypes: [
        {
          type: '员工管理',
          permissions: [
            {
              code: 'user_view',
              id: 1,
              name: '查看员工',
              type: '员工管理'
            },
            {
              code: 'user_add',
              id: 2,
              name: '添加员工',
              type: '员工管理'
            },
            {
              code: 'user_edit',
              id: 3,
              name: '编辑员工',
              type: '员工管理'
            },
            {
              code: 'user_delete',
              id: 4,
              name: '删除员工',
              type: '员工管理'
            }
          ]
        }
      ],
      permissions: [// 所有的权限
        // { id: 1, name: '查看员工', code: 'user_view' },
        // { id: 2, name: '添加员工', code: 'user_add' },
        // { id: 3, name: '编辑员工', code: 'user_edit' },
        // { id: 4, name: '删除员工', code: 'user_delete' },
      ],
      disabledlist: [], // 禁用的权限
      checkedlist: [] // 选中的权限
    }
  },
  async mounted() {
    this.checkedlist = this.checked
    const resp = await getAllPermissions()
    this.permissionTypes = resp.data.list
  },
  methods: {
    changeDisabled(value) { // 修改禁止修改的权限
      // 改变禁止修改的权限
      const disableSet = new Set(this.disabledlist)
      const checkedSet = new Set(this.checkedlist)
      const newDisableSet = new Set(value)
      let newCheckedSet = new Set([...checkedSet].filter(x => !disableSet.has(x)))
      newCheckedSet = new Set([...newDisableSet, ...newCheckedSet])
      this.checkedlist = [...newCheckedSet]
      this.disabledlist = [...newDisableSet]
    },
    changeChecked(value) { // 修改选中的权限
      // 改变选中的权限
      const newCheckedSet = new Set([...this.checkedlist, ...value])
      this.checkedlist = [...newCheckedSet]
    },
    cleareChecked() { // 清空选中的权限
      // 清空选中的权限
      this.checkedlist = []
    },
    getNewChecked() { // 获取新的选中的权限
      const disableSet = new Set(this.disabledlist)
      const checkedSet = new Set(this.checkedlist)
      const newCheckedSet = new Set([...checkedSet].filter(x => !disableSet.has(x)))
      return [...newCheckedSet]
    },
    getChecked() { // 获取选中的权限
      return this.checkedlist
    }
  }
}
</script>
<style lang="scss" scoped>
.permission-item {
  margin-bottom: 10px;
}
.permission-group {
  margin-bottom: 20px;
  font-weight: 100 !important;
}
</style>
