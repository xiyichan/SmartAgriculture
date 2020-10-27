Component({
  data: {
    selected: 0,
    color: "#7A7E83",
    selectedColor: "#3cc51f",
    list: [{
      pagePath: "/pages/device/device",
      iconPath: "/images/device_default.png",
      selectedIconPath: "/images/device_active.png",
      text: "设备"
    }, {
      pagePath: "/pages/user/user",
      iconPath: "/images/user_default.png",
      selectedIconPath: "/images/user_active.png",
      text: "用户"
    }]
  },
  attached() {
  },
  methods: {
    switchTab(e) {
      const data = e.currentTarget.dataset
      const url = data.path
      wx.switchTab({url})
      this.setData({
        selected: data.index
      })
    }
  }
})