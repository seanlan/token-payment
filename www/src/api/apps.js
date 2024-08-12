import request from '@/utils/request'

export function getAppList(data) {
  return request({
    url: '/api/v1/app/list',
    method: 'post',
    data
  })
}
