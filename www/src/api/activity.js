import request from '@/utils/request'

export function getActivityList(data) {
  return request({
    url: '/api/v1/activity/list',
    method: 'post',
    data
  })
}
