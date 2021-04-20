// pages/deviceHistoryData/deviceHistoryData.js
var dateTimePicker = require('../../utils/dateTimePicker.js');
 

Page({

  /**
   * 页面的初始数据
   */
  data: {
      iotId:"",
      date: '2018-10-01',
      time: '12:00',
      dateTimeArray: null,
      dateTime: null,
      dateTimeArray1: null,
      dateTime1: null,
      startYear: 2000,
      endYear: 2050,
      limit:5,
      page:1,
      list:"",
      array:[],
      index:0,
      count:0,
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    var obj = dateTimePicker.dateTimePicker(this.data.startYear, this.data.endYear);
    var obj1 = dateTimePicker.dateTimePicker(this.data.startYear, this.data.endYear);
    // 精确到分的处理，将数组的秒去掉
    var lastArray = obj1.dateTimeArray.pop();
    var lastTime = obj1.dateTime.pop();
    var timestamp = Date.parse(new Date());
    var date = new Date(timestamp);
    //获取年份  
    var Y =date.getFullYear();
    //获取月份  
    var M = (date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1);
    //获取当日日期 
    var D = date.getDate() < 10 ? '0' + date.getDate() : date.getDate(); 
    //console.log(Y,M,D)
    this.setData({
      date:Y+'-'+M+'-'+D,
      dateTime: obj.dateTime,
      dateTimeArray: obj.dateTimeArray,
      dateTimeArray1: obj1.dateTimeArray,
      dateTime1: obj1.dateTime
    });
   console.log(options.IotId)
    this.setData({
      iotId: options.IotId
    })
    this.getList()
  },

  getList:function(){
   // console.log(date)
    var that=this;
    var token = "Bearer " + wx.getStorageSync('token');
    wx.request({
      url: 'http://47.100.108.193:8080/api/device/pi/get/history?IotId='+that.data.iotId+'&&Day='+that.data.date+'&&Page='+that.data.page+'&&Limit='+that.data.limit,
      // data: {
      //   IotId: that.data.iotId,
      //   Day: that.data.date,
      //   Page:that.data.page,
      //   Limit:that.data.limit
      // },
      header: {
        "content-type": "application/x-www-form-urlencoded",
        "Authorization": token
      },
      method:"get",
      success: function success(res) {
        var info = res.data;
        console.log(info);
        let a = [];
        for (var i = 1; i <= info.count/that.data.limit; i++) {
         a.push(i);
        }
        console.log(a)
        that.setData({
          list: info.data,
          count:info.count,
          array:a
        })

      }
    })
  },
  changeDate(e){
    this.setData({ date:e.detail.value});
    console.log(this.data.date)
    this.getList(this.data.date)
  },
  changeTime(e){
    this.setData({ time: e.detail.value });
  },
  changeDateTime(e){
    this.setData({ dateTime: e.detail.value });
  },
  changeDateTime1(e) {
    this.setData({ dateTime1: e.detail.value });
  },
  changeDateTimeColumn(e){
    var arr = this.data.dateTime, dateArr = this.data.dateTimeArray;

    arr[e.detail.column] = e.detail.value;
    dateArr[2] = dateTimePicker.getMonthDay(dateArr[0][arr[0]], dateArr[1][arr[1]]);
    
    this.setData({
      dateTimeArray: dateArr,
      dateTime: arr
    });
  },
  changeDateTimeColumn1(e) {
    var arr = this.data.dateTime1, dateArr = this.data.dateTimeArray1;

    arr[e.detail.column] = e.detail.value;
    dateArr[2] = dateTimePicker.getMonthDay(dateArr[0][arr[0]], dateArr[1][arr[1]]);

    this.setData({
      dateTimeArray1: dateArr,
      dateTime1: arr
    });
  },
  bindPickerChange:function(e){
    var that =this;
   
    console.log(e.detail.value);
    this.setData({
      index:e.detail.value,
      page:e.detail.value+1
    })
    this.getList()
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