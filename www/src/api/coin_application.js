import request from '@/utils/request'

export function coinApplicationList(data) {
  return request({
    url: '/api/v1/coin/list',
    method: 'post',
    data
  })
}

export function examineCoinApplication(data) {
  return request({
    url: '/api/v1/coin/examine',
    method: 'post',
    data
  })
}

export function allAssets(data) {
  return request({
    url: '/api/v1/coin/all',
    method: 'post',
    data
  })
}

export function getCustomTokenList(data) {
  return request({
    url: '/api/v1/coin/custom_token_list',
    method: 'post',
    data
  })
}

export function topUpCustomToken(data) {
  return request({
    url: '/api/v1/coin/top_up_custom_token',
    method: 'post',
    data
  })
}
