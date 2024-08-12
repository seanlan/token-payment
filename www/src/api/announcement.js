import request from '@/utils/request'

export function getAnnouncementList(data) {
  return request({
    url: '/api/v1/announcement/list',
    method: 'post',
    data
  })
}

export function addAnnouncement(data) {
  return request({
    url: '/api/v1/announcement/add',
    method: 'post',
    data
  })
}

export function editAnnouncement(data) {
  return request({
    url: '/api/v1/announcement/edit',
    method: 'post',
    data
  })
}
