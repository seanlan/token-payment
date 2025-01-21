// import parseTime, formatTime and set to filter
export { parseTime, formatTime } from '@/utils'
import moment from 'moment'

/**
 * Show plural label if time is plural number
 * @param {number} time
 * @param {string} label
 * @return {string}
 */
function pluralize(time, label) {
  if (time === 1) {
    return time + label
  }
  return time + label + 's'
}

/**
 * @param {number} time
 */
export function timeAgo(time) {
  const between = Date.now() / 1000 - Number(time)
  if (between < 3600) {
    return pluralize(~~(between / 60), ' minute')
  } else if (between < 86400) {
    return pluralize(~~(between / 3600), ' hour')
  } else {
    return pluralize(~~(between / 86400), ' day')
  }
}

/**
 * Number formatting
 * like 10000 => 10k
 * @param {number} num
 * @param {number} digits
 */
export function numberFormatter(num, digits) {
  const si = [
    { value: 1E18, symbol: 'E' },
    { value: 1E15, symbol: 'P' },
    { value: 1E12, symbol: 'T' },
    { value: 1E9, symbol: 'G' },
    { value: 1E6, symbol: 'M' },
    { value: 1E3, symbol: 'k' }
  ]
  for (let i = 0; i < si.length; i++) {
    if (num >= si[i].value) {
      return (num / si[i].value).toFixed(digits).replace(/\.0+$|(\.[0-9]*[1-9])0+$/, '$1') + si[i].symbol
    }
  }
  return num.toString()
}

/**
 * 10000 => "10,000"
 * @param {number} num
 */
export function toThousandFilter(num) {
  return (+num || 0).toString().replace(/^-?\d+/g, m => m.replace(/(?=(?!\b)(\d{3})+$)/g, ','))
}

/**
 * Upper case first char
 * @param {String} string
 */
export function uppercaseFirst(string) {
  return string.charAt(0).toUpperCase() + string.slice(1)
}

/**
 * @param {number} time (seconds)
 * @returns {string}
 */
export function toDate(value) {
  return moment(value * 1000).format('YYYY/MM/DD HH:mm:ss')
}

// "2024-05-10T21:51:48.946124-04:00" -> "2024-05-10 21:51:48"
export function utcTimeFormat(value) {
  return moment((new Date(value)).getTime()).format('YYYY-MM-DD HH:mm:ss')
}

export function floatNormalDisplay(num) {
  // 如果数值很小，可以用 toFixed 保留足够的小数位数
  if (num < 1e-6) {
    return num.toFixed(10)
  }
  return num.toString()
}

export function percent(value) {
  return (value * 100).toFixed(2)
}

// messagePushState
export function messagePushState(value) {
  const stateMap = {
    0: '未发送',
    1: '发送中',
    2: '发送完成'
  }
  return stateMap[value]
}

export function toGwei(value) {
  return value / 1e9
}
