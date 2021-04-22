// pages/device/device.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    // 数据源
    pageSize: "5",
    pageIndex: "1",
    scanCodeMsg: "",

    device: "",
    array:[],
    index:0,
    count:0,
    //token:"Bearer "+wx.getStorageSync('token'),
  },

  waterSwitch(e) {
    var id = e.currentTarget.id;
    // var that=this;
    this.data.device[id].WaterSwitch = e.detail.value;
    var token = "Bearer " + wx.getStorageSync('token');
    wx.request({
      url: 'http://47.100.108.193:8080/api/device/pi/set/property',
      data: {
        IotId: this.data.device[id].IotId,
        WaterSwitch: this.data.device[id].WaterSwitch,
        FanSwitch: this.data.device[id].FanSwitch,
        lightSwitch: this.data.device[id].LightSwitch,
        autoSwitch:this.data.device[id].AutoSwitch,
      },
      method: "post",
      header: {
        "content-type": "application/x-www-form-urlencoded",
        "Authorization": token
      },
      success: function success(res) {
        var info = res.data;
        console.log(info);

      }
    })
    //console.log( that.data.device[id]);
    console.log(e);
  },
  lightSwitch(e) {
    var id = e.currentTarget.id;
    var that = this;
    this.data.device[id].LightSwitch = e.detail.value;
    var token = "Bearer " + wx.getStorageSync('token');
    wx.request({
      url: 'http://47.100.108.193:8080/api/device/pi/set/property',
      data: {
        IotId: this.data.device[id].IotId,
        WaterSwitch: this.data.device[id].WaterSwitch,
        FanSwitch: this.data.device[id].FanSwitch,
        lightSwitch: this.data.device[id].LightSwitch,
        autoSwitch:this.data.device[id].AutoSwitch,
      },
      method: "post",
      header: {
        "content-type": "application/x-www-form-urlencoded",
        "Authorization": token
      },
      success: function success(res) {
        var info = res.data;
        console.log(info);

      }
    })
    console.log(that.data.device[id]);
    console.log(e);
  },
  fanSwitch(e) {
    var id = e.currentTarget.id;
    var that = this;
    this.data.device[id].FanSwitch = e.detail.value;
    var token = "Bearer " + wx.getStorageSync('token');
    wx.request({
      url: 'http://47.100.108.193:8080/api/device/pi/set/property',
      data: {
        IotId: this.data.device[id].IotId,
        WaterSwitch: this.data.device[id].WaterSwitch,
        FanSwitch: this.data.device[id].FanSwitch,
        lightSwitch: this.data.device[id].LightSwitch,
        autoSwitch:this.data.device[id].AutoSwitch,
      },
      method: "post",
      header: {
        "content-type": "application/x-www-form-urlencoded",
        "Authorization": token
      },
      success: function success(res) {
        var info = res.data;
        console.log(info);

      }
    })
    console.log(that.data.device[id]);
    console.log(e);
  },
  autoSwitch(e) {
    var id = e.currentTarget.id;
    var that = this;
    this.data.device[id].FanSwitch = e.detail.value;
    var token = "Bearer " + wx.getStorageSync('token');
    wx.request({
      url: 'http://47.100.108.193:8080/api/device/pi/set/property',
      data: {
        IotId: this.data.device[id].IotId,
        WaterSwitch: this.data.device[id].WaterSwitch,
        FanSwitch: this.data.device[id].FanSwitch,
        lightSwitch: this.data.device[id].LightSwitch,
        autoSwitch:this.data.device[id].AutoSwitch,
      },
      method: "post",
      header: {
        "content-type": "application/x-www-form-urlencoded",
        "Authorization": token
      },
      success: function success(res) {
        var info = res.data;
        console.log(info);

      }
    })
    console.log(that.data.device[id]);
    console.log(e);
  },
  addDevice: function addDevice() {
    var that = this;
    var token = "Bearer " + wx.getStorageSync('token');
    wx.scanCode({ //扫描API
      success(res) { //扫描成功
        console.log(res) //输出回调信息
        that.setData({
          scanCodeMsg: res.result
        });
        wx.request({
          url: 'http://47.100.108.193:8080/api/device/pi/bind/user',
          data: {
            DeviceName: res.result,
          },
          method: "post",
          header: {
            "content-type": "application/x-www-form-urlencoded",
            "Authorization": token
          },
          success: function success(asd) {
            var info = asd.data;
            console.log(info);
          }
        })
        //this.device();
        wx.showToast({
          title: '成功',
          duration: 1000
        })

      }

    })
    that.deviceList();
  },
  /**
   * 生命周期函数--监听页面加载
   */
  deviceList: function deviceList() {
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
        let a = [];
       // console.log(info.count,info.count/that.data.pageSize);
        for (var i = 0; i < info.count/that.data.pageSize; i++) {
          a.push(i+1);
         }
         console.log(a)
        that.setData({
          device: info.data,
          count:info.count,
          array:a
        })

      }
    })
  },
  bindPickerChange:function(e){
    var that =this;
    console.log(e.detail.value);
    this.setData({
      index:e.detail.value,
      pageIndex:e.detail.value+1
    })
    this.deviceList()
  },
  item:function(e){
    var that=this
    console.log(e.currentTarget.value)
  },

  historyData:function(e){
    var that=this
    console.log(that.data.device)
    console.log(e.currentTarget.dataset.alphaBeta)
    wx.navigateTo({
      url: '/pages/deviceHistoryData/deviceHistoryData?IotId='+e.currentTarget.dataset.alphaBeta,
    })
  },
  onLoad: function (options) {
    var that = this;
    this.deviceList();
    setInterval(function () {
      that.deviceList();
    }, 10000);
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
    this.device();
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