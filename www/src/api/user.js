import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/api/v1/admin/login',
    method: 'post',
    data
  })
}

export function logout() {
  return request({
    url: '/api/v1/admin/logout',
    method: 'post'
  })
}

export function rePassword(data) {
  return request({
    url: '/api/v1/admin/change_password',
    method: 'post',
    data
  })
}

export function getInfo(token) {
  return request({
    url: '/api/v1/admin/info',
    method: 'post',
    data: { token }
  })
}
