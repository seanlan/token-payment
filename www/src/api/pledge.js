import request from '@/utils/request'

export function getPledgeTemplates(data) {
  return request({
    url: '/api/v1/pledge/templates',
    method: 'post',
    data
  })
}

export function savePledgeTemplate(data) {
  return request({
    url: '/api/v1/pledge/save_template',
    method: 'post',
    data
  })
}

export function getPledgeRecords(data) {
  return request({
    url: '/api/v1/pledge/records',
    method: 'post',
    data
  })
}
