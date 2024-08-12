import request from '@/utils/request'

export function getMessagePushList(data) {
  return request({
    url: '/api/v1/message/list',
    method: 'post',
    data
  })
}

export function addMessagePush(data) {
  return request({
    url: '/api/v1/message/add',
    method: 'post',
    data
  })
}

export function editMessagePush(data) {
  return request({
    url: '/api/v1/message/edit',
    method: 'post',
    data
  })
}
