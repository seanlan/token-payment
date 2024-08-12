<!-- eslint-disable vue/no-parsing-error -->
<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">权限列表</h4>
    </el-header>
    <el-main>
      <PermissionSelector ref="pselector" class="permission-selector" />
    </el-main>
  </el-container>
</template>

<script>
import PermissionSelector from './components/PermissionSelector'
export default {
  components: {
    PermissionSelector
  },
  data() {
    return {
    }
  },
  computed: {
  },
  async mounted() {
  },
  methods: {
    submitForm() {
      this.$refs.pselector.changeDisabled([3, 2])
    },
    getChecked() {
      console.log(this.$refs.pselector.getNewChecked())
    }
  }
}

</script>

<style lang="scss" scoped>
.permission-selector {
  margin-left: 100px;
}
</style>
