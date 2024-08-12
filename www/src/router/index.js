import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/* Router Modules */
// import componentsRouter from './modules/components'
// import chartsRouter from './modules/charts'
// import tableRouter from './modules/table'
// import nestedRouter from './modules/nested'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    noCache: true                if set true, the page will no be cached(default is false)
    affix: true                  if set true, the tag will affix in the tags-view
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [{ // 固定路由
  path: '/redirect',
  component: Layout,
  hidden: true,
  children: [{
    path: '/redirect/:path(.*)',
    component: () => import('@/views/redirect/index')
  }]
},
{
  path: '/login',
  component: () => import('@/views/login/index'),
  hidden: true
},
{
  path: '/auth-redirect',
  component: () => import('@/views/login/auth-redirect'),
  hidden: true
},
{
  path: '/404',
  component: () => import('@/views/error-page/404'),
  hidden: true
},
{
  path: '/401',
  component: () => import('@/views/error-page/401'),
  hidden: true
},
{
  path: '/',
  component: Layout,
  redirect: '/dashboard',
  children: [{
    path: 'dashboard',
    component: () => import('@/views/dashboard/index'),
    name: 'Dashboard',
    meta: {
      title: '主页',
      icon: 'dashboard',
      affix: true
    }
  }]
}
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [
  {
    path: '/user',
    component: Layout,
    redirect: '/user/list',
    name: '用户管理',
    meta: { title: '用户管理', icon: 'el-icon-user-solid', roles: ['user_view'] },
    alwaysShow: true,
    children: [
      {
        path: 'list',
        name: 'user-list',
        component: () => import('@/views/users/list'),
        meta: { title: '用户管理', icon: 'el-icon-user-solid', roles: ['user_view'] }
      }
    ]
  },
  {
    path: '/recharge',
    component: Layout,
    redirect: '/recharge/list',
    name: '充值记录',
    meta: { title: '充值管理', icon: 'el-icon-s-goods', roles: ['recharge_view'] },
    alwaysShow: true,
    children: [
      {
        path: 'list',
        name: 'list',
        component: () => import('@/views/recharge/recharge'),
        meta: { title: '充值记录', icon: 'el-icon-s-order', roles: ['recharge_view'] }
      }
    ]
  },
  {
    path: '/order',
    component: Layout,
    redirect: '/order/list',
    name: '订单管理',
    meta: { title: '订单管理', icon: 'el-icon-s-order', roles: ['order_view', 'order_log_view'] },
    alwaysShow: true,
    children: [
      {
        path: 'list',
        name: 'list',
        component: () => import('@/views/order/order'),
        meta: { title: '交易挂单', icon: 'el-icon-s-order', roles: ['order_view'] }
      },
      {
        path: 'log',
        name: 'log',
        component: () => import('@/views/order/order-log'),
        meta: { title: '成交记录', icon: 'el-icon-tickets', roles: ['order_log_view'] }
      }
    ]
  },
  {
    path: '/withdraw',
    component: Layout,
    redirect: '/withdraw/list',
    name: '提现管理',
    meta: { title: '提现管理', icon: 'el-icon-s-finance', roles: ['withdraw_view'] },
    alwaysShow: true,
    children: [
      {
        path: 'list',
        name: 'withdraw-list',
        component: () => import('@/views/withdraw/list'),
        meta: { title: '提现记录', icon: 'el-icon-s-finance', roles: ['withdraw_view'] }
      }
    ]
  },
  {
    path: '/coin',
    component: Layout,
    redirect: '/coin/list',
    name: '上币申请',
    meta: { title: '上币申请', icon: 'example', roles: ['coin_view'] },
    alwaysShow: true,
    children: [
      {
        path: 'list',
        name: 'list',
        component: () => import('@/views/coin/coin_application'),
        meta: { title: '上币申请', icon: 'example', roles: ['coin_view'] }
      },
      {
        path: 'custom',
        name: 'custom',
        component: () => import('@/views/coin/custom_token'),
        meta: { title: '创业板', icon: 'chart', roles: ['coin_view'] }
      }
    ]
  },
  {
    path: '/pledge',
    component: Layout,
    redirect: '/pledge/templates',
    name: '代币质押',
    meta: { title: '代币质押', icon: 'lock', roles: ['pledge_view'] },
    alwaysShow: true,
    children: [
      {
        path: 'templates',
        name: 'templates',
        component: () => import('@/views/pledge/templates'),
        meta: { title: '质押配置', icon: 'lock', roles: ['pledge_edit'] }
      },
      {
        path: 'records',
        name: 'records',
        component: () => import('@/views/pledge/records'),
        meta: { title: '质押记录', icon: 'nested', roles: ['pledge_view'] }
      }
    ]
  },
  {
    path: '/tools',
    component: Layout,
    redirect: '/tools/announcement',
    name: '运营工具',
    meta: { title: '运营工具', icon: 'el-icon-star-on', roles: ['announcement_view', 'message_view'] },
    alwaysShow: true,
    children: [
      {
        path: 'announcement',
        name: 'announcement',
        component: () => import('@/views/tools/announcement'),
        meta: { title: '公告发布', icon: 'el-icon-s-comment', roles: ['announcement_view'] }
      },
      {
        path: 'message',
        name: 'message',
        component: () => import('@/views/tools/message'),
        meta: { title: '消息推送', icon: 'el-icon-s-promotion', roles: ['message_view'] }
      }
    ]
  },

  // 404 page must be placed at the end !!!
  {
    path: '*',
    redirect: '/404',
    hidden: true
  }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({
    y: 0
  }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
