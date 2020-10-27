<template>
  <el-container class="container">
    <!-- 旁边 -->
    <el-aside width="15%">
      <!-- brand -->
      <el-image :src="brand" class="brand"></el-image>
      <!-- 列表 -->
      <el-menu
          @open="handleOpen"
          @close="handleClose"
          :collapse="isCollapse"
          background-color="transparent"
          text-color="rgb(249,249,249)"
          active-text-color="rgb(249,249,249)"
          router
          unique-opened
      >
        <el-menu-item
            v-for="(item,i) in navList"
            :key="i"
            :index="item.name"
            @click="changeTitle(item.navItem)"
        >
          <i :class="item.icon"></i>
          <span slot="title">{{item.navItem}}</span>
        </el-menu-item>
        <el-menu-item @click="SignOut()">
          <i class="el-icon-switch-button"></i>
          <span slot="title">退出系统</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <!-- 头部 -->
      <!-- 237, 244, 250 -->
      <el-header height="70px">
        <el-row :gutter="20">
          <el-col :span="18">
            <div style="border-radius: 4px;min-height: 36px;font-size:15px;">
              当前页：
              <i class="el-icon-edit" style="font-size:20px;margin-right:10px;"></i>
              <span slot="title">{{title}}</span>
            </div>
          </el-col>
          <el-col :span="6">
            <div
                style="border-radius: 4px;min-height: 36px;display: flex;justify-content: space-evenly; align-items: center"
            >
              <el-image :src="avatar" class="avatar"></el-image>
              欢迎您：{{userName}}！
              <el-dropdown>
              <span class="el-dropdown-link">
               <i class="el-icon-arrow-down el-icon--right"></i>
                </span>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item @click.native="dialogVisible=true">修改名字</el-dropdown-item>

                  <el-dropdown-item @click.native="showPassword=true" >修改密码</el-dropdown-item>
                  <el-dropdown-item @click.native="uploadAvatar=true">上传头像</el-dropdown-item>
                  <el-dropdown-item  @click.native="SignOut()">退出系统</el-dropdown-item>

                </el-dropdown-menu>
              </el-dropdown>
            </div>



          </el-col>
        </el-row>
      </el-header>
      <!-- main -->
      <el-main>
        <router-view />
      </el-main>
    </el-container>
    <el-dialog :visible.sync="showPassword"  title="修改密码" width="360px">
      <div class="yy-step-content-1">输入旧密码</div>
      <div style="margin-top:20px">
        <i class="el-icon-lock yy-icon-1"></i>
        <input placeholder="请输入你的旧密码" class="yy-input" v-model="oldpassword" type="password" />
      </div>
      <div class="yy-step-content-1" style="margin-top:110px">密码修改</div>
      <div style="margin-top:30px">

        <input placeholder="请输入你的新密码" class="yy-input" v-model="password" type="password" />
      </div>
      <div style="margin-top:20px">

        <input placeholder="请再输入你的新密码" class="yy-input" type="password" v-model="repassword" />
      </div>
      <div class="dialog-footer" slot="footer">
        <el-button @click="showPassword=false">取 消</el-button>
        <el-button @click="savePassword" type="primary">确 定</el-button>
      </div>
    </el-dialog>


    <el-dialog :visible.sync="dialogVisible"  title="修改名字" width="360px">
      <div class="yy-step-content-1">请输入新名字</div>
      <div style="margin-top:20px">
        <i class="el-icon-user yy-icon-1"></i>
        <input placeholder="请输入新名字" class="yy-input" v-model="updateName" type="" />
      </div>
      <div class="dialog-footer" slot="footer">
        <el-button @click="dialogVisible=false">取 消</el-button>
        <el-button @click="saveName" type="primary">确 定</el-button>
      </div>
    </el-dialog>


    <el-dialog :visible.sync="uploadAvatar"  title="上传头像" width="360px">
      <el-upload
          class="avatar-uploader"
          action="http://47.100.108.193:8080/api/admin/upload/avatar"
          :headers="headers"
          :show-file-list="false"
          :limit="1"
          :on-success="handleAvatarSuccess"
          :name="Avatar"
          :before-upload="beforeAvatarUpload">

        <img v-if="imageUrl" :src="imageUrl" class="avatar">
        <i v-else class="el-icon-plus avatar-uploader-icon"></i>
      </el-upload>


      <div class="dialog-footer" slot="footer">
        <el-button @click="uploadAvatar=false">取 消</el-button>
        <el-button @click="saveAvatar" type="primary">确 定</el-button>
      </div>
    </el-dialog>

  </el-container>


</template>

<script>
export default {
  data: function () {
    return {
      brand: require("../assets/images/brand.png"),
      msg: "Welcome to Your Vue.js App",
      avatar: "http://47.100.108.193:8080/"+localStorage.getItem("avatar"),
      Avatar:"Avatar",
      down: require("../assets/images/down.png"),
      isCollapse: false,
      navList: [],
      title: "首页",
      userName: "管理员",
      dialogVisible: false,
      showPassword: false,
      oldpassword: "",
      password: "",
      repassword: "",
      updateName:"",
      uploadAvatar:false,
      imageUrl:"",
      headers:{
        Authorization:"Bearer "+localStorage.getItem("token")
      }
    };
  },
  methods: {
    handleOpen(key, keyPath) {
      console.log(key, keyPath);
    },
    handleClose(key, keyPath) {
      console.log(key, keyPath);
    },
    changeTitle(name) {
      this.title = name;
    },
    SignOut: function () {
      //删除本地token，
      localStorage.clear();
      //跳到登录界面
      this.$router.push("/");
    },

    savePassword() {
      var that = this;
      if (that.password == that.repassword) {
        var parma = {
          Password: that.password,
        };
        that.$axios.adminUpdatePassword(parma).then(function (res) {
          if (res.code == 200) {
            that.$message({
              message: "修改密码成功",
              type: "success",
            });
            //删除本地token，
            localStorage.clear();
            //跳到登录界面
            that.$router.push("/");
          }
        });
      } else {
        this.$alert("两次密码输入不一致", "提示", {
          confirmButtonText: "确定",
          type: "warning",
        });
        that.password = "";
        that.repassword = "";
      }
    },
    saveName(){
      var that = this;
      if (that.updateName!=""){
        var parma={
          Name:that.updateName,
        };
        that.$axios.adminUpdateName(parma).then(function (res) {
          if (res.code == 200) {
            that.$message({
              message: "修改名字成功",
              type: "success",
            });
            //删除本地token，
          //  localStorage.clear();
            //跳到登录界面
           // that.$router.push("/Home");
            that.userName=that.updateName;
            that.updateName="";
            that.dialogVisible=false;
          }
      });
    }else{
        this.$alert("名字为空", "提示", {
          confirmButtonText: "确定",
          type: "warning",
        });

      }
    },
    saveAvatar(){

    },
    handleAvatarSuccess(res, file) {
      this.imageUrl = URL.createObjectURL(file.raw);
      this.$alert("上传成功","提示",{
        confirmButtonText: "确定",
        type: "info",
      })
      this.uploadAvatar=false;
      //this.$router.push("/Home");
      this.$axios.adminInfo().then(function (res){
        if(res.code==200){
          // this.$message({
          //   message: "查询成功",
          //   type: "success",
          // });
          this.avatar="http://47.100.108.193:8080/"+res.data.avatar;
          console.log(this.avatar);
        }
      })
    },
    beforeAvatarUpload(file) {
      const isJPG = file.type === 'image/jpeg';
      const isLt2M = file.size / 1024 / 1024 < 2;

      if (!isJPG) {
        this.$message.error('上传头像图片只能是 JPG 格式!');
      }
      if (!isLt2M) {
        this.$message.error('上传头像图片大小不能超过 2MB!');
      }
      return isJPG && isLt2M;
    }
  },
  mounted() {
    //管理员名称
    const Name = localStorage.getItem("name");
    this.userName = Name;
    const Avatar =localStorage.getItem("avatar");
    this.avatar="http://47.100.108.193:8080/"+Avatar;
    //console.log(this.avatar)

    //功能显示
    this.navList = [
      { name: "/Home", navItem: "首页", icon: "el-icon-discount" },
      { name: "/UserList", navItem: "用户管理", icon: "el-icon-user" },
      { name: "/Pi", navItem: "气象站管理", icon: "el-icon-umbrella" },
      {name:"/AdminList",navItem:"管理员管理",icon:"el-icon-user"},
      {name:"/Register",navItem:"添加管理员",icon:"el-icon-user"},
    ];
    //标题的显示
    switch (this.$route.name) {
      case "UserList":
        this.title = "用户管理";
        break;
      case "Home":
        this.title = "首页";
        break;
      case "Pi":
        this.title = "气象站管理";
        break;
      case "AdminList":
        this.title="管理员管理";
        break;
      case "Register":
        this.title="管理员管理";
        break;
    }
  },
};
</script>

<style scoped>
.container {
  height: 100%;
  background: rgb(237, 244, 250);
}
.el-aside {
  background: linear-gradient(
      rgb(67, 135, 199),
      rgb(58, 170, 186),
      rgb(50, 201, 174)
  );
  color: white;
  height: 100%;
}
.brand {
  width: 64%;
  margin-top: 40px;
  margin-left: 16%;
  margin-bottom: 20px;
}
.el_submenu__title:hover,
.el-submenu:hover,
.el-menu-item:hover {
  background-color: rgba(255, 255, 255, 0.4) !important ;
}
/* .el-menu-item.is_active {
  
  
} */
.el-menu {
  border-right: 0 !important;
}
.el-menu-item {
  margin-top: 10px;
  background-color: transparent !important;
}
.el-menu-item i {
  color: inherit;
  margin-right: 20px;
  margin-left: 10px;
}
.el-header {
  background: rgb(255, 255, 255);
  box-shadow: 0px 0px 10px rgb(67, 135, 199, 0.2);
  line-height: 70px;
}
.avatar {
  height: 40px;
  width: 40px;
  border-radius: 50%;
}
.el-main {
  padding: 0;
  margin: 20px;
  background: rgb(255, 255, 255);
}
/* background: -moz-linear-gradient(top, #000000 0%, #ffffff 100%);
    background: -webkit-gradient(linear, left top, left bottom, color-stop(0%,#000000), color-stop(100%,#ffffff));
    background: -webkit-linear-gradient(top, #000000 0%,#ffffff 100%);
    background: -o-linear-gradient(top, #000000 0%,#ffffff 100%);
    background: -ms-linear-gradient(top, #000000 0%,#ffffff 100%);
    background: linear-gradient(to bottom, #000000 0%,#ffffff 100%); */
.dropdown-group {
  min-width: 100px;
}
.dropdown-group {
  min-width: 100px;
}
.el-dropdown-link {
  cursor: pointer;
  color: #409EFF;
}
.el-icon-arrow-down {
  font-size: 12px;
}
.user-header{
  position: relative;
  display: inline-block;
}
.user-header-com{
  width: 144px;
  height: 144px;
  display: inline-block;
}
.header-upload-btn{
  position: absolute;
  left: 0;
  top: 0;
  opacity: 0;
  /* 通过定位把input放在img标签上面，通过不透明度隐藏 */
}
.tip{
  font-size: 14px;
  color: #666;
}
/* error是用于错误提示 */
.error{
  font-size: 12px;
  color: tomato;
  margin-left: 10px;
}

.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}
.avatar-uploader .el-upload:hover {
  border-color: #409EFF;
}
.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  line-height: 178px;
  text-align: center;
}
</style>
