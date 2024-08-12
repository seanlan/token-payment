import request from '@/utils/request'

export function statsList(data) {
  return request({
    url: `/api/v1/stats/date`,
    method: 'post',
    data
  })
}

export function statsToday(data) {
  return request({
    url: `/api/v1/stats/today`,
    method: 'post',
    data
  })
}

export function statsAssets(data) {
  return request({
    url: `/api/v1/stats/assets`,
    method: 'post',
    data
  })
}

export function statsCustomAssets(data) {
  return request({
    url: `/api/v1/stats/custom`,
    method: 'post',
    data
  })
}
