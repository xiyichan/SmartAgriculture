// pages/user/user.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    user: "",
    userimg: "/images/default_avatar.jpg",
    hiddenNameFlag: false
  },

  userInfo: function userInfo() {
    var that = this;
    var token = "Bearer " + wx.getStorageSync('token');
    wx.request({
      url: 'http://47.100.108.193:8080/api/user/info',
      method: "get",
      header: {
        "content-type": "application/x-www-form-urlencoded",
        "Authorization": token
      },
      success: function success(res) {
        var info = res.data;
        console.log(info);
        if(info.data.Avatar!=""){
          that.setData({
            userimg:'http://47.100.108.193:8080/'+info.data.Avatar
          })   
        }
        that.setData({
          user: info.data,
          
        })
      }
    })
  },
  upShopLogo: function () {
    var that = this;
    wx.showActionSheet({
      itemList: ['从相册中选择', '拍照'],
      itemColor: "#f7982a",
      success: function (res) {
        if (!res.cancel) {
          if (res.tapIndex == 0) {
            that.chooseWxImageShop('album'); //从相册中选择

          } else if (res.tapIndex == 1) {
            that.chooseWxImageShop('camera'); //手机拍照

          }

        }

      }

    })

  },
  upUserInfo: function () {
    var that = this;
    var token = "Bearer " + wx.getStorageSync('token');
    wx.request({
      url: 'http://47.100.108.193:8080/api/user/updata/name',
      data: {

      }
    })
  },
  updataName: function (name) {
    var that = this;
    var token = "Bearer " + wx.getStorageSync('token');
    wx.request({
      url: 'http://47.100.108.193:8080/api/user/update/name',
      data: {
        Name: name
      },
      method: "post",
      header: {
        "content-type": "application/x-www-form-urlencoded",
        "Authorization": token
      },
      success: function success(asd) {
        var info = asd.data;
        console.log(info);
        if (info.code == 200) {
          wx.showToast({
            title: '成功',
            duration: 1000
          })
        } else {
          wx.showToast({
            title: '失败',
            duration: 1000
          })
        }
      }
    })

  },
  showNameWindows: function () {
    var that = this;
    wx.showModal({
      title: "输入新名字",
      content: '请输入',
      editable: true,
      placeholderText: '数量',
      success(res) {
        if (res.confirm) {
          console.log('用户点击确定', res.content);
          if (res.content!=""){
            that.updataName(res.content);
            that.userInfo();
          }else{
            wx.showToast({
              title: '用户名为空',
            })
          }
        } else if (res.cancel) {
          console.log('用户点击取消');
        }
      }
    })
  },
  showPasswordWindows:function(){
    // wx.navigateTo({
    //   url: "/pages/login/login"
    // })
     
    
    var that = this;
    wx.showModal({
      title: "输入新密码",
      content: '请输入',
      editable: true,
      placeholderText: '数量',
      //password:ture,
      success(res) {
        if (res.confirm) {
          console.log('用户点击确定', res.content);
          if (res.content!=""){
            that.updatePassword(res.content);
            wx.navigateTo({
              url: "/pages/login/login"
            })
           // that.userInfo();
          }else{
            wx.showToast({
              title: '密码为空',
            })
          }
        } else if (res.cancel) {
          console.log('用户点击取消');
        }
      }
    })
  },
 
  updatePassword: function (newPassword) {
    var that = this;
    var token = "Bearer " + wx.getStorageSync('token');
    wx.request({
      url: 'http://47.100.108.193:8080/api/user/update/password',
      data: {
        Password: newPassword
      },
      method: "post",
      header: {
        "content-type": "application/x-www-form-urlencoded",
        "Authorization": token
      },
      success: function success(asd) {
        var info = asd.data;
        console.log(info);
        if (info.code == 200) {
          wx.showToast({
            title: '成功',
            duration: 1000
          })
        } else {
          wx.showToast({
            title: '失败',
            duration: 1000
          })
        }
      }
    })

  },
  chooseWxImageShop: function (type) {
    var that = this;
    var img;
    wx.chooseImage({
      sizeType: ['original', 'compressed'],

      sourceType: [type],

      success: function (res) {

        img = res.tempFilePaths[0],

          that.upload_file('http://47.100.108.193:8080/api/user/upload/avatar' + 'shop/shopIcon', res.tempFilePaths[0])

        img = res.tempFilePaths[0];

        that.setData({
          userimg: img

        })

      }

    })
    console.log(data.userimg)
  },
  upload_file: function (url, filePath) {
    var that = this;
    var token = "Bearer " + wx.getStorageSync('token');
    console.log(filePath)
    wx.uploadFile({
      url: 'http://47.100.108.193:8080/api/user/upload/avatar', //后台处理接口
      filePath: filePath,
      name: 'Avatar',
      header: {
        "Authorization": token
      }, // 设置请求的 header

      // formData: {//需要的参数
      //   Avatar:filePath
      // }, // HTTP 请求中其他额外的 form data

      success: function (res) {
        console.log(res.data)
        var data = JSON.parse(res.data);
        that.setData({
          //userimg: data.path,

        });
        that.userInfo();
        wx.showToast({
          title: "成功", // 提示的内容
          icon: "success", // 图标，默认success
          image: "", // 自定义图标的本地路径，image 的优先级高于 icon
          duration: 3000, // 提示的延迟时间，默认1500
          mask: false, // 是否显示透明蒙层，防止触摸穿透
        })

      },

      fail: function (res) {}

    })

  },



  /**
   * 生命周期函数--监听页面加载
   */


  onLoad: function (options) {
    var that = this;
    this.userInfo();

  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady: function () {
    var that = this;
    this.userInfo();
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