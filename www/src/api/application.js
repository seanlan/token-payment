import request from '@/utils/request'

export function getApplicationList(data) {
  return request({
    url: `/api/v1/application/list`,
    method: 'post',
    data
  })
}
