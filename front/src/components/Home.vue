<!-- @author Chen Zikang  -->
<!-- @date   2022/02/23 -->
<!-- Created By Goland -->

<template>
  <el-container class="home-container">
    <!-- 头部区域 -->
    <el-header>
      <div>
        <span>entry-task</span>
      </div>
      <div class="head-right">
        <img v-if="avatarUrl" :src="avatarUrl" class="avatar" />
        <el-dropdown trigger="click" size="small">
          <p class="el-dropdown-link">
            <span v-if="username">{{ username }}</span>
            <i class="el-icon-arrow-down el-icon--right"></i>
          </p>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item icon="el-icon-user" @click.native="editInfo"
              >账号设置</el-dropdown-item
            >
            <el-dropdown-item icon="el-icon-lock" @click.native="updatePass"
              >密码设置</el-dropdown-item
            ><el-dropdown-item icon="el-icon-chat-dot-round" @click.native="feedback"
              >问题反馈</el-dropdown-item
            >
            <el-dropdown-item icon="el-icon-circle-close" @click.native="logout"
              >退出登录</el-dropdown-item
            >
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </el-header>
    <el-container>
      <!-- 左侧边栏 -->
      <el-aside :width="isCollapse ? '64px' : '200px'">
        <div class="iconfont iconzhedie" @click="toggleCollapse"></div>
        <!-- 侧边栏菜单区域 -->
        <el-menu
          background-color="#333744"
          text-color="#fff"
          active-text-color="#409fff"
          :unique-opened="true"
          :router="true"
          :collapse="isCollapse"
          :collapse-transition="false"
          :default-active="activePach"
        >
          <!-- 一级菜单 -->
          <el-submenu
            :index="item.id + ''"
            v-for="item in menulist"
            :key="item.id"
          >
            <!-- 一级菜单模板区 -->
            <template slot="title">
              <!-- 图标 -->
              <i :class="item.icon"></i>
              <!-- 文本 -->
              <span>{{ item.data }}</span>
            </template>
            <!-- 一级菜单模板区域范围 -->

            <el-menu-item
              :index="'/' + subitem.path"
              v-for="subitem in item.chiledren"
              :key="subitem.id"
              @click="saveNavState('/' + subitem.path)"
            >
              <template slot="title">
                <!-- 图标 -->
                <i :class="subitem.icon"></i>
                <!-- 文本 -->
                <span>{{ subitem.data }}</span>
              </template>
            </el-menu-item>
          </el-submenu>
        </el-menu>
      </el-aside>
      <!-- 右侧主体 -->
      <el-main>
        <!-- 路由占位符 -->
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
export default {
  name: 'home',
  data () {
    return {
      // 菜单图标
      iconsObject: 'iconfont icon-zhedie',
      // 是否折叠
      isCollapse: false,
      menulist: {
        joinMonitor: {
          id: '1',
          data: '个人中心',
          icon: 'el-icon-s-help',
          chiledren: [
            { data: '我的信息', path: 'profile', icon: 'el-icon-menu' }
          ]
        }
      },
      activePach: '',
      username: '',
      avatarUrl: ''
    }
  },
  methods: {
    /**
     * 获取用户信息
     */
    async getUserProfile () {
      const { data: res } = await this.$http.get('user/profile')
      if (res.code === 10000) {
        this.username = res.data.username
        this.avatarUrl = res.data.profilePicture
      }
    },
    // 点击按钮，折叠/展开菜单
    toggleCollapse () {
      this.isCollapse = !this.isCollapse
    },
    saveNavState (activePath) {
      // 保存连接的激活状态
      window.sessionStorage.setItem('activePath', activePath)
      this.activePach = activePath
    },
    // 账户设置
    editInfo () {
      this.$router.push({ path: '/profile' })
    },
    // 退出登录
    logout () {
      const _this = this
      this.$http.get('logout').then(function (res) {
        if (res.data.code === 10000) {
          window.sessionStorage.clear()
          _this.$router.push('/login')
        }
      })
    },
    // 修改密码
    updatePass () {
      this.$router.push({ path: '/updatePass' })
    },
    // 系统反馈
    feedback () {
      this.$router.push({ path: '/feedback' })
    }
  },
  created () {
    this.getUserProfile()
  }
}
</script>

<style lang="less" scoped>
  .home-container {
    height: 100%;
  }

  .el-header {
    background-color: #2b303b;
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: #ffffff;
    font-size: 20px;
    font-family: "Trebuchet MS", Arial, Helvetica, sans-serif;
    > div {
      display: flex;
      align-items: center;
      span {
        margin-left: 15px;
      }
    }
    img {
      width: 9%;
      height: 9%;
    }
  }

  .el-aside {
    background-color: #181e2b;

    .el-menu {
      border-right: none;
    }
  }

  .el-main {
    background-color: #f9f9f9;
  }
  .home-container {
    height: 100%;
  }
  .iconfont {
    margin-right: 10px;
    color: #90939a;
    background-color: #1b2a3f;
    width: 100%;
    cursor: pointer;
    justify-content: center;
    text-align: center;
  }
  .head-right img {
    width: 50px;
    height: 50px;
    border-radius: 44px;
  }
  .head-right {
    display: flex;
    align-items: center;
    color: #808080;
    font-weight: 400;
    font-size: 15.07px;
    margin-right: 20px;
  }
  .head-right p {
    margin-left: 12px;
    cursor: pointer;
  }
</style>
