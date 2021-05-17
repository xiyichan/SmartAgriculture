// pages/login/login.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    account: "2863768433@qq.com",
    password: "czx987852",
  },
  Account: function Account(e) {
    this.setData({
      account: e.detail.value
    });
  },
  Password: function Password(e) {
    this.setData({
      password: e.detail.value
    });
  },
  login: function login() {
    var that = this;
    wx.request({
      url: 'http://47.100.108.193:8080/api/user/login/password',
      data: {
        Account: that.data.account,
        Password: that.data.password
      },
      method: "post",
      header: {
        "content-type": "application/x-www-form-urlencoded"
      },
      success: function success(res) {
        //  console.log(Account);
        var info = res.data;
        console.log(info);
        if (info.code == 200) {
          wx.setStorageSync('token', info.token);
          wx.setStorageSync('userID', info.data.ID)
          wx.setStorageSync('avater', info.data.Avatar)
          wx.setStorageSync('name', info.data.Name)
          wx.setStorageSync('email', info.data.Email)
          wx.switchTab({
            url: "/pages/device/device"
          });
        } else {
          wx.showModal({
            title: "提示",
            content: info.msg,
          });
        }
      }
    })
  },
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady: function () {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide: function () {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload: function () {

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh: function () {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom: function () {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {

  }
})