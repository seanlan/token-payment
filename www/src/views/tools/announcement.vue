<template>
  <el-container>
    <el-header class="pagetab">
      <h4 class="links">活动公告</h4>
    </el-header>
    <el-main>
      <el-row type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="showAddModal">添加公告</el-button>
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
          width="160"
        />
        <el-table-column
          prop="title"
          label="标题"
          width="160"
        />
        <el-table-column
          prop="content"
          label="内容"
        />
        <el-table-column
          prop="url"
          label="跳转链接"
        />
        <el-table-column
          prop="is_hide"
          label="是否隐藏"
        >
          <template slot-scope="scope">
            <el-switch
              v-model="scope.row.is_hide"
              active-color="#F56C6C"
              :active-value="true"
              :inactive-value="false"
              @change="editAnnouncement(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column
          prop="open_in_tg"
          label="是否在TG中打开"
        >
          <template slot-scope="scope">
            <el-switch
              v-model="scope.row.open_in_tg"
              :active-value="true"
              :inactive-value="false"
              @change="editAnnouncement(scope.row)"
            />
          </template>
        </el-table-column>
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
          prop="end_at"
          label="结束时间（北京时间）"
          width="160"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.end_at | utcTimeFormat }}</span>
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
        <el-form-item label="标题" prop="title">
          <el-input
            ref="title"
            v-model="editForm.title"
            tabindex="1"
            placeholder="请输入标题"
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
        <el-form-item label="跳转地址" prop="url">
          <el-input
            ref="url"
            v-model="editForm.url"
            tabindex="2"
            placeholder="请输入跳转地址"
          />
        </el-form-item>
        <el-form-item label="活动时间(北京时间)" prop="time_slot">
          <el-date-picker
            v-model="editForm.time_slot"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="timestamp"
          />
        </el-form-item>
        <el-form-item label="是否隐藏" prop="is_hide">
          <el-switch
            v-model="editForm.is_hide"
            active-color="#F56C6C"
            :active-value="true"
            :inactive-value="false"
          />
        </el-form-item>
        <el-form-item label="是否在TG中打开" prop="open_in_tg">
          <el-switch
            v-model="editForm.open_in_tg"
            :active-value="true"
            :inactive-value="false"
          />
        </el-form-item>
        <el-form-item label="排序值" prop="weight">
          <el-input
            ref="weight"
            v-model.number="editForm.weight"
            tabindex="1"
            placeholder="排序值"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="addDialogVisible = false">取 消</el-button>
        <el-button v-if="modaltype === 'edit'" type="primary" @click="doSubmit">编 辑</el-button>
        <el-button v-else type="primary" @click="doSubmit">添 加</el-button>
      </div>
    </el-dialog>
  </el-container>
</template>

<script>
import _ from 'lodash'
import { Message } from 'element-ui'
import { getAnnouncementList, addAnnouncement, editAnnouncement } from '@/api/announcement'

export default {
  components: {
  },
  data() {
    return {
      table_loading: false,
      addDialogVisible: false,
      modaltype: 'add',
      editForm: {
        time_slot: []
      },
      list: [],
      pagesize: 20,
      totalnum: 0,
      page: 1,
      editFormRules: {
        title: [{ required: true, trigger: 'blur', message: '标题不能为空' }],
        content: [{ required: true, trigger: 'blur', message: '内容不能为空' }],
        time_slot: [{ required: true, trigger: 'blur', message: '活动时间不能为空' }]
      }
    }
  },
  computed: {
  },
  async mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      this.table_loading = true
      var params = {
        page: this.page,
        size: this.pagesize
      }
      const res = await getAnnouncementList(params)
      this.list = res.list || []
      this.table_loading = false
    },
    doSubmit() {
      this.$refs.editForm.validate(valid => {
        if (valid) {
          const configInfo = _.cloneDeep(this.editForm)
          console.log(configInfo)
          configInfo.start_at = configInfo.time_slot[0] / 1000
          configInfo.end_at = configInfo.time_slot[1] / 1000
          const method = this.modaltype === 'add' ? addAnnouncement : editAnnouncement
          method(configInfo).then(() => {
            Message({
              message: this.modaltype === 'add' ? '添加成功' : '修改成功',
              type: 'success'
            })
            this.addDialogVisible = false
            this.loadData()
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
        title: '',
        content: '',
        url: '',
        time_slot: []
      }
    },
    showEditModal(row) {
      this.modaltype = 'edit'
      this.addDialogVisible = true
      this.editForm = _.cloneDeep(row)
      this.editForm.time_slot = [new Date(row.start_at).getTime(), new Date(row.end_at).getTime()]
    },
    editAnnouncement(row) {
      const params = _.cloneDeep(row)
      params.start_at = new Date(row.start_at).getTime() / 1000
      params.end_at = new Date(row.end_at) / 1000
      editAnnouncement(params).then(() => {
        Message({
          message: '修改成功',
          type: 'success'
        })
        this.loadData()
      }).catch((e) => {
        console.log(e)
      })
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
