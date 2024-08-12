import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/api/v1/login/login',
    method: 'post',
    data
  })
}

export function logout() {
  return request({
    url: '/api/v1/login/logout',
    method: 'post'
  })
}

export function rePassword(data) {
  return request({
    url: '/api/v1/login/change_password',
    method: 'post',
    data
  })
}

export function getInfo(token) {
  return request({
    url: '/api/v1/login/info',
    method: 'post',
    data: { token }
  })
}

export function userList(data) {
  return request({
    url: '/api/v1/user/list',
    method: 'post',
    data
  })
}

export function userAssets(data) {
  return request({
    url: '/api/v1/user/assets',
    method: 'post',
    data
  })
}
