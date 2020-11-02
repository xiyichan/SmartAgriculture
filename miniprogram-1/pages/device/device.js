// pages/device/device.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    // 数据源
    pageSize: "2",
    pageIndex: "1",
    device: "",
    //token:"Bearer "+wx.getStorageSync('token'),
  },
 
  waterSwitch(e){


    console.log(e);
  },
 
  /**
   * 生命周期函数--监听页面加载
   */
  device: function device() {
    var that = this;
   var token = "Bearer " + wx.getStorageSync('token');
    wx.request({
      url: 'http://47.100.108.193:8080/api/device/pi/user/list',
      data: {
        PageSize: that.data.pageSize,
        PageIndex: that.data.pageIndex,
      },

      method: "post",
      header: {
        "content-type": "application/x-www-form-urlencoded",
        "Authorization": token
      },
      success: function success(res) {
        var info = res.data;
        console.log(info);
        that.setData({
          device:info.data
        })
       
      }
    })
  },

  onLoad: function (options) {
    this.device();
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