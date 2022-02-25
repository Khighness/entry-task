<!-- @author Chen Zikang  -->
<!-- @date   2022/02/23 -->
<!-- Created By Goland -->

<template>
  <div class="register_container">
    <!-- 头部信息 -->
    <div class="register_header">
      <span>
        <b class="register_header_left">CZK</b>
        <b class="register_header_right">  entry-task</b>
      </span>
    </div>
    <div class="register_box">
      <!-- 注册表单 -->
      <el-form ref="registerFormRef" :model="registerForm" :rules="registerFormRules" label-width="100px" label-position="left" class="register_form">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="registerForm.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="请输入您的密码" show-password></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="checkPassword">
          <el-input v-model="registerForm.checkPassword" type="password" placeholder="请确认您的密码" show-password></el-input>
        </el-form-item>
        <el-form-item class="register_button">
          <el-button type="primary" @click="submitRegisterForm('registerForm')" >注册</el-button>
          <el-button type="info" @click="resetRegisterForm()">重置</el-button>
        </el-form-item>
      </el-form>
      <diV class="login_link">
        <el-link type="primary" href="/#/login">已有账号? 点此登录</el-link>
      </diV>
    </div>
    <!-- 尾部信息 -->
    <div class="register_tail">
      <span><b>Copyright © 2022 <a style="color: #00a0ff" href="https://www.parak.top">Chen Zikang</a> <a style="color: orangered" href="https://www.shopee.cn">Shopee</a></b></span>
    </div>
  </div>
</template>

<script>

export default {
  name: 'Register',
  data () {
    var validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'))
      } else {
        if (this.registerForm.checkPassword !== '') {
          this.$refs.registerFormRef.validateField('checkPassword')
        }
        callback()
      }
    }
    var validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.registerForm.password) {
        callback(new Error('两次输入密码不一致!'))
      } else {
        callback()
      }
    }
    return {
      /* 系统名称以及性质 */
      system: {
        name: 'CZK',
        property: 'entry-task'
      },
      /* 注册表单的数据：邮箱，密码，确认密码，验证码 */
      registerForm: {
        username: '',
        password: '',
        checkPassword: '',
        code: ''
      },
      /* 注册表单的校验规则 */
      registerFormRules: {
        username: [
          { required: 'true', message: '请输入用户名', trigger: 'blur' },
          { min: 3, max: 18, message: '用户名为3～18个字符', trigger: 'blur' }
        ],
        password: [
          { required: 'true', message: '请输入密码', trigger: 'blur' },
          { min: 6, max: 20, message: '密码长度为6～20', trigger: 'blur' },
          { validator: validatePass, trigger: 'blur' }
        ],
        checkPassword: [
          { required: 'true', message: '请确认密码', trigger: 'blur' },
          { min: 6, max: 20, message: '密码长度为6～20', trigger: 'blur' },
          { validator: validatePass2, trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    /**
     * 重置输入框
     */
    resetRegisterForm () {
      /* 如果表单有默认值，则恢复到默认值 */
      this.$refs.registerFormRef.resetFields()
    },
    /**
     * 提交注册
     */
    submitRegisterForm (formName) {
      const _this = this
      this.$refs.registerFormRef.validate(valid => {
        if (valid) {
          const registerRequest = { username: this.registerForm.username, password: this.registerForm.password }
          this.$http.post('register', registerRequest).then(function (response) {
            if (response.data.code === 10000) {
              _this.$message({
                message: '注册成功',
                type: 'success',
                duration: 1000,
                showClose: true
              })
              _this.$router.push('login')
            } else {
              _this.$message({
                message: response.data.message,
                type: 'warning',
                duration: 1000,
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

  .register_container {
    height: 100%;
    background:#f0f2f5 url('../assets/img/background.svg') no-repeat center 110px;
    display: flex;
    flex-direction: column;
  }

  .register_header {
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
  .register_header_left {
    text-decoration: underline;
    font-size: 30px;
    color: #3491ff;
  }
  .register_header_right {
    text-decoration: none;
    font-size: 30px;
    color: #000000;
  }

  .register_box {
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

  .register_form {
    position: absolute;
    top: 30px;
    width: 100%;
    padding: 0 50px;
    box-sizing: border-box;
  }

  .register_button {
    padding-top: 25px;
    display: flex;
    justify-content: center;
  }

 .login_link {
    position: absolute;
    left: 60%;
    top: 67%;
  }

  .register_tail {
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
</style>
