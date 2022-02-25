<!-- @author Chen Zikang  -->
<!-- @date   2022/02/23 -->
<!-- Created By Goland -->

<template>
  <div class="editContainer">
    <!-- 面包屑导航区 -->
    <el-breadcrumb separator-class="el-icon-arrow-right">
      <el-breadcrumb-item :to="{ path: '/home' }"
        >首页</el-breadcrumb-item
      >
      <el-breadcrumb-item>问题反馈</el-breadcrumb-item>
    </el-breadcrumb>
    <!-- 卡片视图区域 -->
    <el-card>
      <!-- 修改表单区域 -->
      <el-form
        ref="ruleForm"
        label-width="100px"
        :rules="rules"
        :model='ruleForm'
      >
        <!-- 用户信息区域 -->
        <el-form-item>
          <h2>问题反馈</h2>
        </el-form-item>

        <el-form-item prop="feedback" >
          <el-input
            type="textarea"
            placeholder="请输入反馈内容"
            v-model="ruleForm.feedback"
            maxlength="200"
            show-word-limit
             :rows="12"
          >
          </el-input>
        </el-form-item>

        <el-form-item class="btn">
          <el-button type="warning" @click="FeedBack()">提交反馈</el-button>
          <el-button @click="return_page">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
export default {
  data () {
    return {
      // 文本域
      ruleForm: {
        feedback: ''
      },
      rules: {
        feedback: [
          // 验证是否为空
          { required: true, message: '请填写反馈内容' }
        ]
      }
    }
  },
  created () {},
  methods: {
    async FeedBack () {
      this.$refs.ruleForm.validate(async valid => {
        if (valid) {
          const { data: res } = await this.$http.post('sysFeedback/save', this.ruleForm)
          if (res.code !== 10000) {
            return this.$message.error(res.data)
          }
          if (res.code === 10000) {
            this.$message.success(res.data)
            this.$refs.ruleForm.resetFields()
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
.el-form{
    margin-left: 100px;
}

</style>
