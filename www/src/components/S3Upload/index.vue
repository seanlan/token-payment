<template>
  <div class="flex flex-y-center">
    <gallery v-if="showPreview" :images="images" :index="index" @close="()=>{index = null;showPreview=false}" />
    <el-upload
      ref="upload"
      drag
      action=""
      :accept="accept"
      :class="{ hide: hide }"
      list-type="picture-card"
      :http-request="uploadFile"
      :before-upload="beforeUpload"
      :file-list="filesList"
      :multiple="!isSingle"
      :on-change="handleChange"
      :limit="maxlimit"
    >
      <i slot="default" class="el-icon-plus" />
      <div slot="file" slot-scope="{ file }">
        <img v-if="isimg" class="el-upload-list__item-thumbnail" :src="file.url" alt="">
        <img v-else class="el-upload-list__item-thumbnail" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAAAAXNSR0IArs4c6QAAAW5JREFUaEPtWUESgjAMDC9TX6a+TH2ZujPgMDUpJI0BJL1woXQ3SZdt09HGR7dx/PT3BI5EdCYiPJcYdyJ6vBe+SIvXMoBJAL+GcXoHEWS+Ro3AbcHIl0ABHiRUBJ5rCP0IAxvsWgY4AojE9cfEhn1XLtNMAMDFzeRMitt/myKALGAPjkcSiCyhzABqr1ShzIBCqcJKqNUrDf+V0iaEEfDwS5xNCCPg4Ze4vRVGAGWNLBwU9T1+FeA5lxlKwIi9Oi0JaP8D1hMbe0jpzyEhXshDhbjTVlgJeRBYVEY9CKxCRq1KhOinjErR240btUooApcyqrlW4dTCw8wtKqPcD0erSCmju1chbcnMeT/MC80BY3knCWjPA5Yo1+ZkBrgMiF0S7/D3FwVle6v5dho4h+aG5GE8uEhNRRcCHgCt31AT8PA5VrDlPFOTz8PneBEQb8anOvUgITXdvMDVviNd/n7mTBGIANm0RhJoCp/D5BffNZExoYt3gwAAAABJRU5ErkJggg==" alt="">
        <span class="el-upload-list__item-actions">
          <span v-if="!disabled" class="el-upload-list__item-delete" @click="handleRemove(file)">
            <i class="el-icon-delete" />
          </span>
          <span v-if="preview" class="el-upload-list__item-delete" @click="handlePreview(file)">
            <i class="el-icon-search" />
          </span>
        </span>
      </div>
    </el-upload>
    <b> * 拖拽文件或点击；支持文件格式：{{ accept }} * </b>
  </div>
</template>
<script>
// import * as imageConversion from 'image-conversion'
// import * as qiniu from 'qiniu-js'
import { uuid } from 'vue-uuid'
import { cdnHost, R2Endpoint, R2AccessKeyId, R2SecretAccessKey } from '@/settings'
import VueGallery from 'vue-gallery'
import S3 from 'aws-sdk/clients/s3.js'

export default {
  name: 'ImgUpload',
  components: {
    'gallery': VueGallery
  },
  props: {
    checkFile: {
      type: Function,
      default: () => { (file) => { return true } }
    },
    disabled: {
      type: Boolean,
      default: false
    },
    preview: {
      type: Boolean,
      default: true
    },
    isimg: {
      type: Boolean,
      default: true
    },
    accept: {
      type: String,
      default: '.jpg,.jpeg,.png,.gif'
    },
    // 张数限制
    maxlimit: {
      type: Number,
      default: () => 1
    },
    // 固定参数名修改
    paramsName: {
      type: String,
      default: () => 'image'
    },
    // 已上传文件list
    fileList: {
      type: Array,
      default: () => []
    },
    // 是否需要隐藏上传icon
    isHide: {
      type: Boolean,
      default: () => false
    },
    // 是否是单图上传
    isSingle: {
      type: Boolean,
      default: false
    }
    // qnToken: {
    //   type: String,
    //   default: null
    // },
    // 七牛JavaScript SDK API: qiniu.upload(file: blob, key: string, token: string, putExtra: object, config: object) 里的 config
    // 具体参数查看 https://developer.qiniu.com/kodo/sdk/1283/javascript#3
    // qnConfig: {
    //   type: Object,
    //   default() {
    //     return {
    //       useCdnDomain: true,
    //       disableStatisticsReport: false,
    //       retryCount: 6,
    //       region: qiniu.region.z2
    //     }
    //   }
    // }
  },
  data() {
    return {
      images: [],
      index: null,
      showPreview: false,
      hide: this.isHide,
      uploadParams: this.otherParams,
      filesList: this.fileList,
      defaultParams: this.paramsName,
      loadingObj: null
    }
  },
  computed: {
  },
  watch: {
    otherParams: {
      handler(newVal) {
        this.uploadParams = newVal
      },
      deep: true,
      immediate: true
    },
    limit: {
      handler(newVal) {
        console.log('newVal', newVal)
        this.maxlimit = newVal
      }
    },
    paramsName: {
      handler(newVal) {
        this.defaultParams = newVal
      }
    },
    fileList: {
      handler(newVal) {
        this.filesList = newVal
        this.hide = this.fileList.length >= this.maxlimit
      },
      deep: true,
      immediate: true
    },
    isHide: {
      handler(newVal) {
        this.hide = newVal
      }
    }
  },
  methods: {
    beforeUpload(file) {
      return this.checkFile(file)
    },
    // 图片上传改变
    handleChange(file) {
      this.hide = this.fileList.length >= this.maxlimit
    },
    // 图片移除
    handleRemove(file) {
      const arr = this.$refs.upload.uploadFiles
      const index = arr.indexOf(file)
      this.filesList.splice(index, 1)
      let num = 0
      this.$refs.upload.uploadFiles.map((item) => {
        if (item.uid === file.uid) {
          this.$refs.upload.uploadFiles.splice(num, 1)
        }
        num++
      })
      this.hide = this.fileList.length >= this.maxlimit
      this.$emit('handleRemove', this.filesList, this.isSingle)
    },
    // perview
    handlePreview(file) {
      this.showPreview = true
      setTimeout(() => {
        this.images = [file.url]
        this.index = 0
      }, 100)
    },
    // 图片上传
    handleSuccess(res, file, filesList) {
      console.log('res', res)
      this.loadingObj && this.loadingObj.close()
      file.url = cdnHost + res.key
      if (this.isSingle) {
        this.filesList.push(file)
        this.$emit('handleSuccess', this.filesList, 'single')
        if (this.$refs.upload.uploadFiles.length === this.maxlimit) {
          this.hide = true
        }
      } else {
        this.filesList.push(file)
        this.$emit('handleSuccess', this.filesList, 'group')
        if (this.$refs.upload.uploadFiles.length === this.maxlimit) {
          this.hide = true
        }
      }
    },
    async uploadFile(option) {
      this.loadingObj = this.$loading({
        lock: true,
        text: '文件上传中，请稍后...',
        background: 'rgba(0, 0, 0, 0.8)'
      })
      const s3 = new S3({
        endpoint: R2Endpoint,
        accessKeyId: R2AccessKeyId,
        secretAccessKey: R2SecretAccessKey,
        signatureVersion: 'v4'
      })
      await s3.upload({
        Bucket: 'www',
        Key: this.changeFileName(option.file.name),
        Body: option.file
      }, (err, data) => {
        if (err) {
          this.$message.error('上传失败')
          this.loadingObj.close()
          return
        }
        console.log('data', data)
        this.handleSuccess({ key: data.Key }, option.file)
      })
    },
    // 修改原文件名，给文件名添加一个时间戳
    changeFileName(filename) {
      // return filename.replace(/.[a-zA-Z0-9]+$/, (match) => {
      //   return `-${Date.now()}${match}`
      // })
      return uuid.v1() + filename.match(/.[a-zA-Z0-9]+$/)[0]
    }
  }
}
</script>
<style>
.hide .el-upload--picture-card {
  display: none;
}
.el-upload-dragger {
    width: 145px;
    height: 145px;
    border: 0;
}
</style>
