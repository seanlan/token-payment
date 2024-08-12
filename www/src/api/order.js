import request from '@/utils/request'

export function getOrderList(data) {
  return request({
    url: '/api/v1/order/list',
    method: 'post',
    data
  })
}

export function getOrderLogs(data) {
  return request({
    url: '/api/v1/order/logs',
    method: 'post',
    data
  })
}
