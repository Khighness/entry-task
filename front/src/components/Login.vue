<!-- @author Chen Zikang  -->
<!-- @date   2022/02/23 -->
<!-- Created By Goland -->

<template>
  <div class="login_container">
    <!-- 头部信息 -->
    <div class="login_header">
      <span>
        <b class="login_header_left">CZK</b>
        <b class="login_header_right">  entry-task</b>
      </span>
    </div>
    <div class="login_box">
      <!-- Logo图标 -->
      <div class="logo_box">
        <img src="../assets/img/k.jpg">
      </div>
      <!-- 登录表单 -->
      <el-form ref="loginFormRef" :model="loginForm" :rules="loginFormRules" label-width="0px" class="login_form">
        <el-form-item prop="username">
          <el-input v-model="loginForm.username" prefix-icon="iconfont iconyouxiang" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="loginForm.password" prefix-icon="iconfont iconmima" type="password" placeholder="请输入密码" show-password></el-input>
        </el-form-item>
        <el-form-item class="login_button">
          <el-button type="primary" @click="submitLoginForm('loginForm')">登录</el-button>
          <el-button type="info" @click="resetLoginForm()">重置</el-button>
        </el-form-item>
      </el-form>
      <div class="register_link">
        <el-link type="primary" href="/#/register">没有账号？点此注册</el-link>
      </div>
    </div>
    <!-- 尾部信息 -->
    <div class="login_tail">
      <span><b>Copyright © 2022 <a style="color: #00a0ff" href="https://www.parak.top">Chen Zikang</a> <a style="color: orangered" href="https://www.shopee.cn">Shopee</a></b></span>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Login',
  id: '',
  data () {
    return {
      /* 系统名称以及性质 */
      system: {
        name: 'CZK',
        property: 'entry-task'
      },
      /* 登录表单的数据：用户名和密码 */
      loginForm: {
        username: '',
        password: ''
      },
      /* 登录表单的校验规则 */
      loginFormRules: {
        username: [
          { required: 'true', message: '请输入用户名', trigger: 'blur' },
          { min: 3, max: 18, message: '长度在3到18个字符', trigger: 'blur' }
        ],
        password: [
          { required: 'true', message: '请输入密码', trigger: 'blur' },
          { min: 6, max: 20, message: '长度在6到20个字符', trigger: 'blur' }
        ]
      }
    }
  },
  created () {
    var _self = this
    document.onkeydown = function (e) {
      var key = window.event.keyCode
      if (key === 13 || key === 100) {
        _self.submitLoginForm('formName')
      }
    }
  },
  methods: {
    /**
     * 重置输入框
     */
    resetLoginForm () {
      /* 如果表单有默认值，则恢复到默认值 */
      this.$refs.loginFormRef.resetFields()
    },
    /**
     * 提交登录
     */
    submitLoginForm (formName) {
      const _this = this
      this.$refs.loginFormRef.validate(valid => {
        if (valid) {
          const loginRequest = { username: this.loginForm.username, password: this.loginForm.password }
          this.$http.post('login', loginRequest).then(function (response) {
            if (response.data.code === 10000) {
              _this.$message({
                message: '登录成功',
                type: 'success',
                duration: 1000,
                showClose: true
              })
              window.sessionStorage.setItem('entry-token', response.data.data.token)
              _this.$router.push('/home')
            } else {
              _this.$message({
                message: response.data.message,
                type: 'error',
                showClose: true
              })
            }
          })
        }
      })
    }
  }
}
</script>

<style lang="less" scope>

  .login_container {
    height: 100%;
    background:#f0f2f5 url('../assets/img/background.svg') no-repeat center 110px;
    display: flex;
    flex-direction: column;
  }

  .login_header {
    user-select: none;
    text-align: center;
    overflow: hidden;
    position: absolute;
    top: 8%;
    height: 120px;
    width: 100%;
    line-height: 120px;
    font-family: "Trebuchet MS",Arial,Helvetica,sans-serif;
    transition:height .8s cubic-bezier(0.075,0.82,0.165,1);
  }
  .login_header_left {
    text-decoration: underline;
    font-size: 30px;
    color: #3491ff;
  }
  .login_header_right {
    text-decoration: none;
    font-size: 30px;
    color: #000000;
  }

  .login_box {
    width: 450px;
    height: 300px;
    background-color: #FFFFFF;
    border: 1px solid #00a0ff;
    border-radius: 3px;
    position: absolute;
    left: 50%;
    top: 55%;
    transform: translate(-50%, -50%);
    border: none;
    box-shadow: 0 2px 12px 0 #68abd1;
  }
  .logo_box {
    position: absolute;
    height: 100px;
    width: 100px;
    border-radius: 50%;
    padding: 10px;
    left: 50%;
    transform: translate(-50%, -50%);
    img {
      height: 100%;
      width: 100%;
      border-radius:50%;
      box-shadow: 0px 0px 3px rgba(63, 81, 181, 0.72);
      border: 5px solid #3491ff;
      transform-origin: center center;
    }
  }

  .login_form {
    position: absolute;
    top: 90px;
    width: 100%;
    padding: 0 50px;
    box-sizing: border-box;
  }

  .login_button {
    padding-top: 25px;
    display: flex;
    justify-content: center;
  }

  .register_link {
    position: absolute;
    left: 60%;
    top: 68%;
  }

  .login_tail {
    margin-top: 10px;
    user-select: none;
    text-align: center;
    position: absolute;
    bottom: 30px;
    left: 40%;
    right: 40%;
    color: black;
    font-family: "Trebuchet MS";
    font-size: 5px;
  }

  a {
    text-decoration: none;
  }

</style>
