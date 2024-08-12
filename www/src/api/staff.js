import request from '@/utils/request'

export function getAllPermissions() {
  return request({
    url: '/api/v1/permission/all',
    method: 'post'
  })
}

export function getAllGroups() {
  return request({
    url: '/api/v1/group/all',
    method: 'post'
  })
}

export function addGroups(data) {
  return request({
    url: '/api/v1/group/add',
    method: 'post',
    data
  })
}

export function editGroups(data) {
  return request({
    url: '/api/v1/group/edit',
    method: 'post',
    data
  })
}

export function getStaffList(data) {
  return request({
    url: '/api/v1/staff/list',
    method: 'post',
    data
  })
}

export function addStaff(data) {
  return request({
    url: '/api/v1/staff/add',
    method: 'post',
    data
  })
}

export function editStaff(data) {
  return request({
    url: '/api/v1/staff/edit',
    method: 'post',
    data
  })
}

export function staffInfo(data) {
  return request({
    url: '/api/v1/staff/info',
    method: 'post',
    data
  })
}

export function staffFrozen(data) {
  return request({
    url: '/api/v1/staff/ban',
    method: 'post',
    data
  })
}

export function staffRePassword(data) {
  return request({
    url: '/api/v1/staff/reset_password',
    method: 'post',
    data
  })
}

export function staffLogs(data) {
  return request({
    url: '/api/v1/staff/logs',
    method: 'post',
    data
  })
}
