import request from '@/utils/request'

export function getChainList(data) {
  return request({
    url: `/api/v1/chain/list`,
    method: 'post',
    data
  })
}
