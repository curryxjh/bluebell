<template>
  <header class="header">
      <span class="logo" @click="goIndex">Forum</span>
    <div class="search">
      <label class="s-logo"></label>
      <input type="text" class="s-input" placeholder="搜索" />
    </div>
    <div class="btns">
      <div v-show="!isLogin">
        <a class="login-btn" @click="goLogin">登录</a>
        <a class="login-btn" @click="goSignUp">注册</a>
      </div>
      <div class="user-box" v-show="isLogin" @click.stop="toggleDropdown">
        <span class="user">{{ currUsername }}</span>
        <div class="dropdown-content" :class="{ 'show': showDropdown }">
          <a @click.stop="goLogout">登出</a>
        </div>
      </div>
    </div>
  </header>
</template>

<script>
export default {
  name: "HeadBar",
  data() {
    return {
      showDropdown: false
    };
  },
  created(){
    this.$store.commit("init");
    // 点击页面其他地方关闭下拉菜单
    document.addEventListener('click', this.closeDropdown);
  },
  beforeDestroy() {
    document.removeEventListener('click', this.closeDropdown);
  },
  computed: {
    isLogin() {
      return this.$store.getters.isLogin;
    },
    currUsername(){
      console.log(this.$store.getters.username);
      return this.$store.getters.username;
    }
  },
  methods: {
    goIndex(){
      this.$router.push({ name: "Home" });
    },
    goLogin() {
      this.$router.push({ name: "Login" });
    },
    goSignUp() {
      this.$router.push({ name: "SignUp" });
    },
    goLogout(){
      this.showDropdown = false;
      this.$store.commit("logout");
    },
    toggleDropdown() {
      this.showDropdown = !this.showDropdown;
    },
    closeDropdown() {
      this.showDropdown = false;
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="less">
.header {
  width: 100%;
  height: 56px;
  position: fixed;
  background: #ffffff;
  display: flex;
  display: -webkit-flex;
  align-items: center;
  top: 0;
  z-index: 100;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.06);
  border-bottom: 1px solid #e2e8f0;
  .logo {
    margin-left: 16px;
    height: 32px;
    background: url("../assets/images/logo.png") no-repeat;
    background-size: 32px 32px;
    background-position: left center;
    padding-left: 40px;
    line-height: 32px;
    flex-grow: 0;
    margin-right: 24px;
    cursor: pointer;
    color: #1e293b;
    font-weight: 700;
    font-size: 20px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    &:hover {
      color: #3b82f6;
    }
  }
  .search {
    flex-grow: 1;
    margin: 0 auto;
    max-width: 690px;
    position: relative;
    display: flex;
    display: -webkit-flex;
    .s-logo {
      width: 18px;
      height: 18px;
      background: url("../assets/images/search.png") no-repeat;
      background-size: cover;
      display: inline-block;
      position: absolute;
      top: 50%;
      margin-top: -9px;
      left: 15px;
    }
    .s-input {
      flex-grow: 1;
      -webkit-appearance: none;
      appearance: none;
      background-color: #f1f5f9;
      border-radius: 12px;
      border: 2px solid transparent;
      box-shadow: none;
      color: #1e293b;
      display: block;
      height: 40px;
      outline: none;
      padding: 0 16px 0 40px;
      width: 100%;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      &::placeholder {
        color: #94a3b8;
      }
      &:hover {
        background-color: #e2e8f0;
      }
      &:focus {
        border-color: #3b82f6;
        background-color: #ffffff;
        box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
      }
    }
  }
  .btns {
    flex-grow: 0;
    margin-left: 16px;
    margin-right: 10px;
    display: flex;
    display: -webkit-flex;
    align-items: center;
    .login-btn {
      border: 2px solid #e2e8f0;
      border-radius: 10px;
      box-sizing: border-box;
      text-align: center;
      letter-spacing: 0px;
      text-decoration: none;
      font-size: 14px;
      font-weight: 500;
      line-height: 20px;
      text-transform: none;
      padding: 8px 20px;
      color: #64748b;
      fill: #64748b;
      display: inline-block;
      cursor: pointer;
      background: transparent;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      &:hover {
        background: #f1f5f9;
        border-color: #cbd5e1;
        color: #1e293b;
        transform: translateY(-1px);
      }
      &:active {
        transform: translateY(0);
      }
      &:nth-child(1) {
        margin-right: 8px;
      }
      &:nth-child(2) {
        margin-right: 16px;
        background: #3b82f6;
        border-color: #3b82f6;
        color: #ffffff;
        box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
        &:hover {
          background: #2563eb;
          border-color: #2563eb;
          box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
        }
      }
    }
    .user-box {
      position: relative;
      cursor: pointer;
    }
    .user {
      width: auto;
      height: 24px;
      background: url("../assets/images/avatar.png") no-repeat;
      background-size: 24px 24px;
      background-position: left center;
      padding-left: 28px;
      display: flex;
      display: -webkit-flex;
      align-items: center;
      cursor: pointer;
      padding: 12px 16px 12px 32px;
      border-radius: 12px;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      color: #1e293b;
      font-weight: 500;
      font-size: 14px;
      &:hover {
        background: #f1f5f9;
      }
      &::after {
        content: "";
        width: 0;
        height: 0;
        border-top: 5px solid #64748b;
        border-right: 5px solid transparent;
        border-bottom: 5px solid transparent;
        border-left: 5px solid transparent;
        margin-top: 5px;
        margin-left: 10px;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      }
    }
    .user-box.active .user::after {
      transform: rotate(180deg);
      margin-top: -5px;
    }
    .dropdown-content {
      display: none;
      position: absolute;
      top: calc(100% + 8px);
      right: 0;
      background-color: #ffffff;
      min-width: 180px;
      box-shadow: 0 10px 20px rgba(0, 0, 0, 0.15), 0 3px 6px rgba(0, 0, 0, 0.1);
      border-radius: 12px;
      overflow: hidden;
      z-index: 1000;
      border: 1px solid #e2e8f0;
      &.show {
        display: block;
        animation: fadeIn 0.2s ease-in-out;
      }
      a {
        color: #1e293b;
        padding: 14px 20px;
        text-decoration: none;
        display: block;
        cursor: pointer;
        font-size: 14px;
        font-weight: 500;
        transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
        &:hover {
          background-color: #f1f5f9;
          color: #3b82f6;
          padding-left: 24px;
        }
        &:active {
          background-color: #e2e8f0;
        }
      }
    }
    
    @keyframes fadeIn {
      from {
        opacity: 0;
        transform: translateY(-10px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }
  }
}

/* 移动端适配 */
@media screen and (max-width: 768px) {
  .header {
    padding: 0 6px;
    .logo {
      margin-left: 0;
      margin-right: 6px;
      font-size: 13px;
      padding-left: 28px;
      background-size: 26px 26px;
      white-space: nowrap;
    }
    .search {
      flex-grow: 1;
      max-width: calc(100% - 180px);
      margin: 0 6px;
      .s-input {
        height: 30px;
        font-size: 13px;
        padding: 0 10px 0 32px;
      }
      .s-logo {
        width: 14px;
        height: 14px;
        left: 10px;
        margin-top: -7px;
      }
    }
    .btns {
      flex-shrink: 0;
      margin-left: 6px;
      margin-right: 8px;
      .login-btn {
        padding: 4px 8px;
        font-size: 12px;
        line-height: 18px;
        letter-spacing: 0.3px;
        white-space: nowrap;
        &:nth-child(1) {
          margin-right: 4px;
        }
        &:nth-child(2) {
          margin-right: 0;
        }
      }
      .user {
        font-size: 12px;
        padding: 6px 6px 6px 24px;
        background-size: 18px 18px;
        white-space: nowrap;
      }
    }
  }
}

/* 更小屏幕适配 */
@media screen and (max-width: 480px) {
  .header {
    padding: 0 4px;
    .logo {
      margin-right: 4px;
      font-size: 12px;
      padding-left: 26px;
      background-size: 24px 24px;
    }
    .search {
      max-width: calc(100% - 160px);
      margin: 0 4px;
      .s-input {
        height: 28px;
        font-size: 12px;
        padding: 0 8px 0 28px;
      }
      .s-logo {
        width: 12px;
        height: 12px;
        left: 8px;
        margin-top: -6px;
      }
    }
    .btns {
      margin-left: 4px;
      margin-right: 6px;
      .login-btn {
        padding: 3px 6px;
        font-size: 11px;
        line-height: 16px;
        &:nth-child(1) {
          margin-right: 3px;
        }
      }
      .user {
        font-size: 11px;
        padding: 5px 5px 5px 22px;
        background-size: 16px 16px;
      }
    }
  }
}
</style>
