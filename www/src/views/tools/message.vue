<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">消息推送</h4>
    </el-header>
    <el-main>
      <el-row type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="showAddModal">添加消息</el-button>
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
          prop="image"
          label="图片"
        />
        <el-table-column
          prop="content"
          label="内容"
        >
          <template slot-scope="scope">
            <div style="white-space: pre-line;">
              <span>{{ scope.row.content }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          prop="ext"
          label="附加选项"
        />
        <el-table-column
          prop="start_at"
          label="开始时间（北京时间）"
          width="160"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.start_at | utcTimeFormat }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="state"
          label="状态"
          width="100"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.state | messagePushState }}</span>
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          width="160"
        >
          <template slot-scope="scope">
            <el-button type="success" size="mini" @click="showEditModal(scope.row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-main>
    <el-dialog
      :title="modaltype=='edit'? '修改公告': '添加公告'"
      :visible.sync="addDialogVisible"
      width="50%"
    >
      <el-form ref="editForm" :model="editForm" :rules="editFormRules" label-width="140px">
        <el-form-item label="图片" prop="image">
          <el-input
            ref="image"
            v-model="editForm.image"
            tabindex="1"
            placeholder="请输入图片地址"
          />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input
            ref="content"
            v-model="editForm.content"
            tabindex="2"
            placeholder="请输入内容"
            type="textarea"
            :rows="4"
          />
        </el-form-item>

        <!-- <el-form-item label="附加选项" prop="ext">
          <el-input
            ref="ext"
            v-model="editForm.ext"
            tabindex="2"
            placeholder="请输入附加选项"
            type="textarea"
            :rows="4"
          />
        </el-form-item> -->

        <el-form-item label="附加选项" prop="ext">
          <el-button type="primary" icon="el-icon-plus" size="small" @click="addExt">添加</el-button>
        </el-form-item>

        <el-form-item
          v-for="(ext, index) in editForm.extList"
          :key="index"
        >
          <el-row>
            <el-col :span="10">
              <el-input
                v-model="ext.text"
                placeholder="请输入文本"
                @change="changeExt"
              />
            </el-col>
            <el-col :span="10">
              <el-input
                v-model="ext.url"
                placeholder="请输入链接"
                @change="changeExt"
              />
            </el-col>
            <el-col :span="4">
              <el-button @click.prevent="removeExt(index)">删除</el-button>
            </el-col>
          </el-row>
        </el-form-item>

        <el-form-item label="开始时间" prop="start">
          <el-date-picker
            v-model="editForm.start"
            type="datetime"
            placeholder="选择日期时间"
            value-format="timestamp"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="addDialogVisible = false">取 消</el-button>
        <el-button v-if="modaltype === 'edit'" type="primary" @click="doAddGroup">编 辑</el-button>
        <el-button v-else type="primary" @click="doAddGroup">添 加</el-button>
      </div>
    </el-dialog>
  </el-container>
</template>

<script>
import _ from 'lodash'
import { Message } from 'element-ui'
import { getMessagePushList, addMessagePush, editMessagePush } from '@/api/message'
// import Needaddshop from '@/components/Needaddshop'

export default {
  components: {
  },
  data() {
    return {
      table_loading: false,
      addDialogVisible: false,
      modaltype: 'add',
      editForm: {
        ext: '[]',
        start: new Date().getTime()
      },
      list: [],
      pagesize: 20,
      totalnum: 0,
      page: 1,
      editFormRules: {
        content: [{ required: true, trigger: 'blur', message: 'VALUE不能为空' }]
      }
    }
  },
  computed: {
  },
  async mounted() {
    this.loadConfigList()
  },
  methods: {
    async loadConfigList() {
      this.table_loading = true
      var params = {
        page: this.page,
        size: this.pagesize
      }
      const res = await getMessagePushList(params)
      this.list = res.list || []
      this.table_loading = false
    },
    doAddGroup() {
      this.$refs.editForm.validate(valid => {
        console.log(this.editForm)
        if (valid) {
          const configInfo = _.cloneDeep(this.editForm)
          // 过滤掉extList中的空项
          configInfo.extList = configInfo.extList.filter(item => item.text && item.url)
          configInfo.ext = JSON.stringify(configInfo.extList)
          const handler = this.modaltype === 'add' ? addMessagePush : editMessagePush
          handler(configInfo).then(() => {
            Message({
              message: this.modaltype === 'add' ? '添加成功' : '修改成功',
              type: 'success'
            })
            this.addDialogVisible = false
            this.loadConfigList()
          }).catch((e) => {
            console.log(e)
          })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    showAddModal() {
      this.modaltype = 'add'
      this.addDialogVisible = true
      this.editForm = {
        ext: '[]',
        content: '',
        start: new Date().getTime(),
        extList: []
      }
    },
    showEditModal(row) {
      this.modaltype = 'edit'
      this.addDialogVisible = true
      this.editForm = {
        ...row,
        start: new Date(row.start_at).getTime(),
        extList: JSON.parse(row.ext)
      }
    },
    addExt() {
      const extList = _.cloneDeep(this.editForm.extList) || []
      extList.push({
        text: '',
        url: ''
      })
      this.editForm.extList = extList
      this.editForm.ext = JSON.stringify(this.editForm.extList)
    },
    removeExt(index) {
      const extList = _.cloneDeep(this.editForm.extList) || []
      extList.splice(index, 1)
      this.editForm.extList = extList
      this.editForm.ext = JSON.stringify(this.editForm.extList)
    },
    changeExt() {
      this.editForm.ext = JSON.stringify(this.editForm.extList)
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
