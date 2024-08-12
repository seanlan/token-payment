import request from '@/utils/request'
import axios from 'axios'
import fileDownload from 'js-file-download'
import { getToken } from '@/utils/auth'
import { Message } from 'element-ui'

export function paymentList(data) {
  return request({
    url: '/api/v1/payment/list',
    method: 'post',
    data
  })
}

export function downloadPayment(data) {
  axios
    .create({
      baseURL: process.env.VUE_APP_BASE_API,
      // withCredentials: true, // send cookies when cross-domain requests
      timeout: 30000, // request timeout
      headers: {
        'X-TOKEN': getToken()
      }
    })
    .post('/api/v1/payment/export',
      data,
      {
        responseType: 'blob'
      })
    .then((res) => {
      if (res.data.type === 'application/json') {
        // 将Blob对象转化为json对象
        const reader = new FileReader()
        reader.readAsText(res.data, 'utf-8')
        reader.onload = (e) => {
          const resData = JSON.parse(e.target.result)
          console.log('resData', resData)
          if (resData.error !== 0) {
            Message({
              message: resData.message || 'Error',
              type: 'error',
              duration: 5 * 1000
            })
          }
        }
        return
      }
      fileDownload(res.data, 'payment.xlsx')
    })
}
