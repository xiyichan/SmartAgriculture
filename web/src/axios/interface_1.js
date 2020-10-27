import { request } from './index'

export default {
    //管理员登录
    adminLogin: (data) => {
        return request('/api/admin/login/password', 'post', data)
    },
    adminRegister:(data)=>{
        return request('/api/admin/register', 'post', data)
    },
    verify: (data) => {
        return request('/api/email/captcha', 'post', data)
    },
    //管理员修改密码
    adminUpdatePassword: (data) => {
        return request('/api/admin/update/password', 'post', data)
    },
    adminUpdateName:(data)=>{
        return request('/api/admin/update/name', 'post', data)
    },
    adminList:(data)=>{
        return request('/api/admin/list', 'post', data)
    },
    adminInfo:(data)=>{
        return request('api/admin/info','get',data)
    },
    //用户列表
    userList: (data) => {
        return request('/api/user/list', 'post', data)
    },
    //删除用户
    userDelete: (data) => {
        return request('/api/user/delete', 'post', data)
    },


    // 气象站
    piList: (data) => {
        return request('/api/device/pi/list', 'post', data)
    },
    //删除气象站
    piDelete: (data) => {
        return request('/api/device/pi/delete', 'post', data)
    },
    piNew: () => {
        return request('/api/device/pi/register', 'post', '')
    },

}