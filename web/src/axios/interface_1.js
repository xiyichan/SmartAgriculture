import { request } from './index'

export default {
    //管理员登录
    adminLogin: (data) => {
        return request('/api/admin/login/password', 'post', data)
    },
    //管理员修改密码
    adminUpdatePassword: (data) => {
        return request('/api/admin/update/password', 'post', data)
    },
    //用户列表
    userList: (data) => {
        return request('/api/admin/user/list', 'post', data)
    },
    //删除用户
    userDelete: (data) => {
        return request('/api/admin/userdelete', 'post', data)
    },
    //专家列表
    expertList: (data) => {
        return request('/api/admin/expertlist', 'get', data)
    },
    //删除专家
    expertDelete: (data) => {
        return request('/api/admin/expertdelete', 'post', data)
    },

    //公告
    announceList: (data) => {
        return request('/api/announcement/list', 'post', data)
    },
    // 提交
    announceNew: (data) => {
        return request('/api/announcement/new', 'post', data)
    },
    // 删除
    announceDelete: (data) => {
        return request('/api/announcement/delete', 'post', data)
    },

    // 气象站
    weatherList: (data) => {
        return request('/api/admin/pi/list', 'post', data)
    },
    //删除气象站
    weatherDelete: (data) => {
        return request('/api/weatherStation/delete', 'post', data)
    },
    weatherNew: () => {
        return request('/api/weatherStation/new', 'post', '')
    },
    //新闻
    newNew: (data) => {
        return request('/api/news/new', 'post', data)
    },
    newList: (data) => {
        return request('/api/news/list', 'get', data)
    },
    newDelete: (data) => {
        return request('/api/news/delete', 'post', data)
    }
}