<!-- @author Chen Zikang  -->
<!-- @date   2022/02/23 -->
<!-- Created By Goland -->

<template>
  <div class="editContainer">
    <!-- 面包屑导航区 -->
    <el-breadcrumb separator-class="el-icon-arrow-right">
      <el-breadcrumb-item :to="{ path: '/home' }">首页</el-breadcrumb-item>
      <el-breadcrumb-item>密码修改</el-breadcrumb-item>
    </el-breadcrumb>
    <!-- 卡片视图区域 -->
    <el-card>
      <!-- 修改表单区域 -->
      <el-form
        ref="ruleForm"
        :model="ruleForm"
        label-width="100px"
        :rules="rules"
      >
        <!-- 用户信息区域 -->
        <el-form-item prop="oldPassword" label="原密码">
          <el-input
            v-model="ruleForm.oldPassword"
            type="password"
            placeholder="请输入旧密码"
            show-password
            id=oldPassword
          ></el-input>
        </el-form-item>

        <el-form-item prop="newPassword" label="新密码">
          <el-input
            v-model="ruleForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
            id=newPassword
          ></el-input>
        </el-form-item>

        <el-form-item prop="confirmPwd" label="确认新密码">
          <el-input
            v-model="ruleForm.confirmPwd"
            type="password"
            placeholder="请再次输入新密码"
            show-password
            id=confirmPwd
          ></el-input>
        </el-form-item>

        <el-form-item class="btn">
          <el-button type="warning" @click="updatePassword">修改密码</el-button>
          <el-button @click="return_page">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
export default {
  data () {
    var validatePass2 = (rule, value, cb) => {
      const npw = document.getElementById('newPassword').value
      if (npw !== value) { cb(new Error('两次新密码输入不一致')) }
      const opw = document.getElementById('oldPassword').value
      if (opw !== value) { return cb() } cb(new Error('新密码不能和旧密码输入一致'))
    }
    return {
      // 图片上传的请求头对象
      headerObj: {
        Authorization: window.sessionStorage.getItem('entry-token')
      },
      ruleForm: {
        oldPassword: '',
        newPassword: '',
        confirmPwd: ''
      },
      rules: {
        oldPassword: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, max: 20, message: '长度在6到20个字符', trigger: 'blur' }
        ],
        newPassword: [
          { required: true, message: '请输入新密码', trigger: 'blur' },
          { min: 6, max: 20, message: '长度在6到20个字符', trigger: 'blur' }
        ],
        confirmPwd: [
          { required: true, message: '请再次输入新密码', trigger: 'blur' },
          { min: 6, max: 20, message: '长度在6到20个字符', trigger: 'blur' },
          { validator: validatePass2, trigger: 'blur' }
        ]
      }
    }
  },
  created () {
  },
  methods: {
    async updatePassword () {
      this.$refs.ruleForm.validate(async valid => {
        if (valid) {
          const newPW = {
            oldPassword: this.ruleForm.oldPassword,
            newPassword: this.ruleForm.newPassword
          }
          const { data: res } = await this.$http.patch('examPlaceInvigilator/update', newPW)
          if (res.code !== 10000) {
            // resetFields用不了，直接一个个更新了
            this.ruleForm.oldPassword = ''
            this.ruleForm.newPassword = ''
            this.ruleForm.confirmPwd = ''
            return this.$message.error(res.data)
          }
          if (res.code === 10000) {
            this.$message.success(res.data)
            this.return_page()
          }
        }
      })
    },

    // 返回主页
    return_page () {
      this.$router.push({ path: '/home' })
    }
  }
}
</script>

<style lang="less" scoped>
.el-form {
  width: 460px;
  margin: 0 150px;
}

.btn {
  text-align: center;
}
.el-card {
  background: url("../assets/img/background.svg");
  width: 800px;
  height: 800px;
  margin: 0 auto;
}
.editContainer {
  background: url("../assets/img/background.svg");
}
.cropper-content {
  width: 500px;
  height: 500px;
  background: pink;
}
.cropper {
  width: 500px;
  height: 500px;
  background: yellow;
}
</style>
