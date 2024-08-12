import { Message } from 'element-ui'

export function makeItemLikeper(likes) {
  const { GoodCount, BadCount } = likes
  let likeper
  let unlikeper
  if (GoodCount === 0 && BadCount === 0) {
    likeper = unlikeper = 0
  } else {
    likeper = parseFloat((GoodCount / (GoodCount + BadCount) * 100).toFixed(1))
    unlikeper = parseFloat((100 - likeper).toFixed(1))
  }
  likes.likeper = likeper
  likes.unlikeper = unlikeper
  return likes
}

/**
 * 检测文件大小及格式（图片）
 * @param {object} file 文件对象
 */
export function checkImages(file, size = 2 * 1000 * 1000) {
  // console.log(file)
  if (file.type.indexOf('image') < 0) {
    console.log('只能上传图片文件')
    Message({
      message: '只能上传图片文件',
      type: 'error'
    })
    return false
  }
  if (file.size > size) {
    console.log('图片大小超过限制')
    Message({
      message: '图片大小超过限制',
      type: 'error'
    })
    return false
  }
  return true
}

/**
 * 检测文件大小及格式（图片及视频）
 * @param {object} file 文件对象
 */
export function checkGifVideo(file, size = 10 * 1000 * 1000) {
  // console.log(file)
  if (file.type.indexOf('image') < 0 && file.type.indexOf('video') < 0) {
    console.log('只能上传图片或视频文件')
    Message({
      message: '只能上传图片或视频文件',
      type: 'error'
    })
    return false
  }
  if (file.size > size) {
    console.log('图片或视频大小超过限制')
    Message({
      message: '图片或视频大小超过限制',
      type: 'error'
    })
    return false
  }
  return true
}

/**
 * 价格转换，默认不保留小数点
 * @param {number} nums 价格（分）
 */
export function getPriceFromNum(nums, len = 0) {
  return parseFloat(nums / 100).toFixed(len)
}

/**
 * 处理权限列表数据
 * @param {array} datas 权限数据
 */
export function makeAllPermissions(datas) {
  for (const i in datas) {
    datas[i].disabled = false
    datas[i].checked = false
  }
  return datas
}

/**
 * 获得所有选择的权限
 * @param {arrya} pers 所有权限数据
 */
export function getCheckedPermissions(pers) {
  const tmp = []
  pers.filter((item) => {
    return item.checked && !item.disabled
  }).forEach((p, i) => {
    tmp.push(p.id)
  })
  return tmp
}

export function makeAjaxParamData(param, objs) {
  // console.log(objs)
  for (const i in objs) {
    if (objs[i] !== '') {
      param[i] = objs[i]
    }
  }
  return param
}

export function getObjFromArray(datas, key) {
  const tmp = {}
  for (const i in datas) {
    tmp[datas[i][key]] = datas[i]
  }
  return tmp
}

/**
 * 获取文件后缀名
 * @param {string} str 文件名
 */
export function getFileType(str) {
  return getFileExt(str)[0].replace('.', '')
}

/**
 * 获取文件后缀
 * @param {string} str 文件名
 */
export function getFileExt(str) {
  return /\.[^\.]+$/.exec(str)
}

export function is(objname, item) {
  return objname === Object.prototype.toString.call(item).slice(8, -1)
}

export function validateIp4(ip) {
  return /^((\d|[1-9]\d|1\d\d|2([0-4]\d|5[0-5]))(\.|$)){4}$/.test(ip)
}

export function validateCidr(cidr) {
  return /^((\d|[1-9]\d|1\d\d|2([0-4]\d|5[0-5]))\.){3}((\d|[1-9]\d|1\d\d|2([0-4]\d|5[0-5]))\/){1}(\d|[1-3]\d)$/.test(cidr)
}

export const ip4ToInt = ip => ip.split('.').reduce((sum, part) => (sum << 8) + parseInt(part, 10), 0) >>> 0

export const isIp4InCidr = ip => cidr => {
  if (validateIp4(ip) && validateCidr(cidr)) {
    const [range, bits] = cidr.split('/')
    const mask = ~(2 ** (32 - bits) - 1)
    return (ip4ToInt(ip) & mask) === (ip4ToInt(range) & mask)
  } else {
    return false
  }
}

/**
 * 判断ip是否在cidr 列表内
 * @param {*} ip
 * @param {*} cidr
 * @returns
 */
export const isIp4InCidrList = (ip, cidrList) => {
  for (const cidr of cidrList) {
    if (isIp4InCidr(ip)(cidr)) {
      return true
    }
  }
  return false
}

