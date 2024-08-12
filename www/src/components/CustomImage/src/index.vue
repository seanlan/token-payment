<script>
import { Image } from 'element-ui'
export default {
  name: 'CustomImage',
  components: {},
  extends: Image,
  data() {
    return {

    }
  },
  computed: {},

  created() {
  },
  mounted() {},

  methods: {
    clickHandler() {
      if (!this.preview) {
        return
      }
      this.prevOverflow = document.body.style.overflow
      document.body.style.overflow = 'hidden'
      this.showViewer = true
      setTimeout(function() {
        // 获取遮罩层dom
        const domImageMask = document.querySelector('.el-image-viewer__mask')
        if (!domImageMask) {
          return
        }
        domImageMask.addEventListener('click', () => {
          document.querySelector('.el-image-viewer__close').click()
        })
      }, 100)
    }
  }
}
</script>
