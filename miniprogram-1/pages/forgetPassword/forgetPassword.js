// pages/forgetPassword/forgetPassword.js
Page({

  /**
   * 页面的初始数据
   */
  data: {

  },
  Email: function Email(e) {
    this.setData({
      email: e.detail.value
    });
  },
  Password: function Password(e) {
    this.setData({
      password: e.detail.value
    });
  },
  RePassword: function RePassword(e) {
    this.setData({
      rePassword: e.detail.value
    });
  },
  Captcha: function Captcha(e) {
    this.setData({
      captcha: e.detail.value
    });
  },
  getCaptcha: function getCaptcha() {
    var that = this;
    wx.request({
      url: 'http://47.100.108.193:8080/api/email/captcha',
      data: {
        Email: that.data.email
      },
      method: "post",
      header: {
        "content-type": "application/x-www-form-urlencoded"
      },
      success: function success(res) {
        var info = res.data;
        wx.showModal({
          title: "提示",
          content: info.msg,
        })

      }
    })
  },
  updatePassword:function updatePassword(){
    var that=this;
    if (that.data.password!=that.data.password){
      wx.showModal({
        title: "提示",
        content: "密码2次输入不一样",
      });
      return
    }
    wx.request({
      url: 'http://47.100.108.193:8080/api/user/forget/password/email',
      data:{
        Email:that.data.email,
        Password:that.data.password,
        Captcha:that.data.captcha,
      },
      method: "post",
      header: {
        "content-type": "application/x-www-form-urlencoded"
      },
      success:function success(res){
        var info = res.data;
        wx.showModal({
          title: "提示",
          content: info.msg,
        })
        if (info.code==200){
          wx.navigateTo({
            url: '/pages/login/login',
          })
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