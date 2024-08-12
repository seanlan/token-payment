module.exports = {
  DatePickerOptions: {
    shortcuts: [
      {
        text: '今天',
        onClick(picker) {
          const end = new Date()
          const start = new Date()
          start.setTime(start.getTime())
          picker.$emit('pick', [start, end])
        }
      },
      {
        text: '昨天',
        onClick(picker) {
          const end = new Date()
          const start = new Date()
          start.setTime(start.getTime() - 3600 * 1000 * 24 * 1)
          end.setTime(end.getTime() - 3600 * 1000 * 24 * 1)
          picker.$emit('pick', [start, end])
        }
      },
      {
        text: '最近一周',
        onClick(picker) {
          const end = new Date()
          const start = new Date()
          start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
          picker.$emit('pick', [start, end])
        }
      }, {
        text: '最近一个月',
        onClick(picker) {
          const end = new Date()
          const start = new Date()
          start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
          picker.$emit('pick', [start, end])
        }
      }
    ]
  },
  orderStatus: [
    {
      value: 0,
      label: '交易中',
      type: 'primary'
    },
    {
      value: 1,
      label: '已完成',
      type: 'info'
    },
    {
      value: 2,
      label: '已取消',
      type: 'danger'
    },
    {
      value: -1,
      label: '全部',
      type: 'info'
    }
  ]
}
