import request from '@/utils/request'

export function getQiniuToken() {
  return request({
    url: '/api/v1/utils/qiniu_token',
    method: 'post'
  })
}
