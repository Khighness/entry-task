<!-- @author Chen Zikang  -->
<!-- @date   2022/02/23 -->
<!-- Created By Goland -->

<template>
  <div class="editContainer">
    <!-- 面包屑导航区 -->
    <el-breadcrumb separator-class="el-icon-arrow-right">
      <el-breadcrumb-item :to="{ path: '/home' }">首页</el-breadcrumb-item>
      <el-breadcrumb-item>账号设置</el-breadcrumb-item>
    </el-breadcrumb>
    <!-- 卡片视图区域 -->
    <el-card>
      <!-- 修改表单区域 -->
      <el-form ref="editFormRef" :model="editForm" label-width="100px" :rules="editFormRules">
        <!-- 用户头像区域 -->
        <el-form-item label="上传头像" style="padding-top: 150px">
          <el-upload
            ref="uploadRef"
            class="avatar-uploader"
            :action="uploadUrl"
            :on-change="changePhotoFile"
            :show-file-list="false"
            :file-list="photoList"
            :on-success="handleAvatarSuccess"
            :before-upload="beforeAvatarUpload"
            :auto-upload="false"
          >
            <img v-if="avatarUrl" :src="avatarUrl" class="avatar">
            <i v-else class="el-icon-plus avatar-uploader-icon"></i>
          </el-upload>
        </el-form-item>
        <!-- 用户信息区域 -->
        <el-form-item label="用户名" prop="username" style="padding-top: 50px">
          <el-input v-model="editForm.username" prefix-icon="el-icon-user" clearable></el-input>
        </el-form-item>
        <el-form-item label="每日一言">
          <el-input type="textarea" v-model="hitokoto" maxlength="200" show-word-limit :rows="4" :disabled="true"></el-input>
        </el-form-item>
        <el-form-item class="btn" style="padding-top: 50px">
          <el-button type="warning" @click="onSubmit">更新</el-button>
          <el-button @click="return_page">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    <!-- 裁剪头像对话框 -->
     <el-dialog title="图片剪裁" :visible.sync="dialogVisible" append-to-body>
      <div class="cropper-content">
        <div class="cropper" style="text-align:center">
          <vueCropper
            ref="cropper"
            :img="option.img"
            :outputSize="option.outputSize"
            :outputType="option.outputType"
            :info="option.info"
            :canScale="option.canScale"
            :autoCrop="option.autoCrop"
            :autoCropWidth="option.autoCropWidth"
            :autoCropHeight="option.autoCropHeight"
            :fixed="option.fixed"
            :fixedBox="option.fixedBox"
            :fixedNumber="option.fixedNumber"
          ></vueCropper>
        </div>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="finish" :loading="loading">确认</el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import { VueCropper } from 'vue-cropper'
export default {
  components: {
    VueCropper
  },
  data () {
    return {
      // 上传的图片信息
      uploadUrl: 'http://127.0.0.1:10000/avatar/upload',
      // 图片地址
      avatarUrl: '',
      // 一言
      hitokoto: '',
      // 要修改的用户信息
      editForm: {
        username: ''
      },
      editFormRules: {
        username: [
          { required: 'true', message: '请输入用户名', trigger: 'blur' },
          { min: 3, max: 18, message: '用户名为3～18个字符', trigger: 'blur' }
        ]
      },
      // 控制头像裁剪对话框的显示与隐藏
      dialogVisible: false,
      // 裁剪组件的基础配置option
      option: {
        img: '', // 裁剪图片的地址
        info: true, // 裁剪框的大小信息
        outputSize: 0.8, // 裁剪生成图片的质量
        outputType: 'jpg', // 裁剪生成图片的格式
        canScale: true, // 图片是否允许滚轮缩放
        autoCrop: true, // 是否默认生成截图框
        autoCropWidth: 300, // 默认生成截图框宽度
        autoCropHeight: 300, // 默认生成截图框高度
        fixedBox: false, // 固定截图框大小 不允许改变
        fixed: true, // 是否开启截图框宽高固定比例
        fixedNumber: [1, 1], // 截图框的宽高比例
        full: true, // 是否输出原图比例的截图
        canMoveBox: false, // 截图框能否拖动
        original: false, // 上传图片按照原始比例渲染
        centerBox: false, // 截图框是否被限制在图片里面
        infoTrue: true, // true 为展示真实输出图片宽高 false 展示看到的截图框宽高
        canMove: true
      },
      loading: false,
      fileinfo: {},
      photoList: [],
      isChange: ''
    }
  },
  created () {
    this.getUserProfile()
    this.getHitokoto()
  },
  methods: {
    /**
     * 获取用户信息
     */
    async getUserProfile () {
      const { data: res } = await this.$http.get('user/profile')
      if (res.code !== 10000) {
        this.$message({
          message: '认证失败',
          type: 'error',
          duration: 1000,
          showClose: true
        })
        this.$router.push('/login')
      } else {
        this.editForm.username = res.data.username
        this.avatarUrl = res.data.profilePicture
      }
    },
    /**
     * 获取一言
     */
    getHitokoto () {
      const _this = this
      this.$http.get('https://v1.hitokoto.cn').then(function (response) {
        _this.hitokoto = response.data.hitokoto + '  -「' + response.data.creator + '」'
      })
    },
    changePhotoFile (file, fileList) {
      if (fileList.length > 0) {
        this.photoList = [fileList[fileList.length - 1]]
      }
      this.handleCrop(file)
    },
    handleCrop (file) {
      this.$nextTick(() => {
        this.option.img = URL.createObjectURL(file.raw)
        this.fileinfo = file.raw
        this.dialogVisible = true
      })
    },
    handleAvatarSuccess (res, file) {
      this.avatarUrl = URL.createObjectURL(file.raw)
    },
    beforeAvatarUpload (file) {
      const isIMG = file.type === 'image/jpeg' || 'image/jpg' || 'image/png'
      if (!isIMG) {
        this.$message.error(
          '请上传JPG/PNG/JPEG格式图片'
        )
      }
      const isLt2M = file.size / 1024 / 1024 < 2
      if (!isLt2M) {
        this.$message.error('上传头像图片大小不能超过2MB')
      }
      return isIMG && isLt2M
    },
    // base64转图片文件
    dataURLtoFile (dataurl, filename) {
      var arr = dataurl.split(',')
      var mime = arr[0].match(/:(.*?);/)[1]
      var bstr = atob(arr[1])
      var n = bstr.length
      var u8arr = new Uint8Array(n)
      while (n--) {
        u8arr[n] = bstr.charCodeAt(n)
      }
      return new File([u8arr], filename, { type: mime })
    },
    // 裁剪过后的头像上传
    finish () {
      this.$refs.cropper.getCropBlob(data => {
        this.avatarUrl = URL.createObjectURL(data)
      })
      this.$refs.cropper.getCropData(async data => {
        this.dialogVisible = false
        const file = this.dataURLtoFile(data, 'profilePicture.png')
        const formData = new FormData()
        formData.append('profile_picture', file)
        const { data: res } = await this.$http.post('/avatar/upload', formData)
        if (res.code !== 10000) {
          return this.$message({
            message: '头像上传失败',
            type: 'error',
            duration: 1000,
            showClose: true
          })
        }
        localStorage.setItem('avatar', res.data)
        this.$message({
          message: '头像上传成功',
          type: 'success',
          duration: 1000,
          showClose: true
        })
        this.dialogVisible = false
      })
    },
    onSubmit () {
      this.$refs.editFormRef.validate(async valid => {
        if (!valid) return
        const editInfo = {
          username: this.editForm.username
        }
        const { data: res } = await this.$http.put('user/update', editInfo)
        localStorage.setItem('name', this.editForm.name)
        if (res.code !== 10000) {
          return this.$message.error(res.message)
        } else {
          this.$message({
            message: '信息更新成功',
            type: 'success',
            duration: 1000,
            showClose: true
          })
          this.getUserProfile()
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
    background: url('../assets/img/background.svg');
    width: 800px;
    height: 800px;
    margin: 0 auto
  }

  .editContainer{
    background: url('../assets/img/background.svg');
  }

  .cropper-content{
    width:500px;height:500px;background: pink;
  }

  .cropper{
    width:500px;
    height:500px;
    background: yellow;
  }
</style>
