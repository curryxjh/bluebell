<template>
  <div class="main">
    <div class="container">
      <h2 class="form-title">登录</h2>
      <div class="form-group">
        <label for="name">用户名</label>
        <input type="text" class="form-control" v-model="username" name="name" id="name" placeholder="用户名" />
      </div>
      <div class="form-group">
        <label for="pass">密码</label>
        <input type="password" class="form-control" v-model="password"  name="pass" id="pass" placeholder="密码" />
      </div>
      <div class="form-btn">
        <button type="button" class="btn btn-info" @click="submit">提交</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
	name: "Login",
	data() {
		return {
			username: "",
			password: "",
			submitted: false
		};
	},
	computed: {
	},
	created() {

	},
	methods: {
		submit() {
			this.$axios({
				method: 'post',
				url:'/login',
				data: JSON.stringify({
					username: this.username,
					password: this.password
				})
			}).then((res)=>{
				console.log(res.data)
				if (res.code == 1000) {
          localStorage.setItem("loginResult", JSON.stringify(res.data));
          this.$store.commit("login", res.data);
          this.$router.push({path: this.redirect || '/' })
				} else {
					console.log(res.msg)
				}
			}).catch((error)=>{
				console.log(error)
			})
		}
	}
};
</script>
<style lang="less" scoped>
.main {
  background: #f5f7fa;
  padding: 120px 0;
  min-height: 100vh;
  
  .container {
    width: 480px;
    background: #ffffff;
    margin: 0 auto;
    max-width: 1200px;
    padding: 48px;
    border-radius: 20px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08), 0 2px 8px rgba(0, 0, 0, 0.04);
    border: none;
    .form-title {
      margin-bottom: 40px;
      text-align: center;
      color: #1e293b;
      font-size: 28px;
      font-weight: 700;
    }
    .form-group {
      margin: 0 0 24px 0;
      label {
        display: inline-block;
        max-width: 100%;
        margin-bottom: 8px;
        font-weight: 600;
        color: #475569;
        font-size: 14px;
      }
      .form-control {
        display: block;
        width: 100%;
        height: 48px;
        padding: 12px 16px;
        font-size: 15px;
        line-height: 1.5;
        color: #1e293b;
        background-color: #f8fafc;
        background-image: none;
        border: 2px solid #e2e8f0;
        border-radius: 12px;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        &::placeholder {
          color: #94a3b8;
        }
        &:hover {
          border-color: #cbd5e1;
        }
        &:focus {
          outline: none;
          background-color: #ffffff;
          border-color: #3b82f6;
          box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
        }
      }
    }
    .form-btn {
      display: flex;
      justify-content: center;
      margin-top: 32px;
      .btn {
        width: 100%;
        padding: 14px 24px;
        font-size: 16px;
        line-height: 1.5;
        border-radius: 12px;
        display: inline-block;
        margin-bottom: 0;
        font-weight: 600;
        text-align: center;
        white-space: nowrap;
        vertical-align: middle;
        -ms-touch-action: manipulation;
        touch-action: manipulation;
        cursor: pointer;
        border: none;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      }
      .btn-info {
        color: #fff;
        background-color: #3b82f6;
        box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
        &:hover {
          background-color: #2563eb;
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
        }
        &:active {
          transform: translateY(0);
          box-shadow: 0 2px 4px rgba(59, 130, 246, 0.3);
        }
      }
    }
  }
}

/* 移动端适配 */
@media screen and (max-width: 768px) {
  .main {
    padding: 80px 10px;
    min-height: 100vh;
    box-sizing: border-box;
    .container {
      width: 100%;
      max-width: 100%;
      padding: 20px 15px;
      box-sizing: border-box;
      .form-title {
        margin-bottom: 25px;
        font-size: 20px;
      }
      .form-group {
        margin: 12px 0;
        label {
          margin-bottom: 4px;
          font-size: 14px;
        }
        .form-control {
          height: 36px;
          padding: 8px 12px;
          font-size: 14px;
          box-sizing: border-box;
          width: 100%;
        }
      }
      .form-btn {
        margin-top: 20px;
        .btn {
          padding: 8px 24px;
          font-size: 16px;
        }
      }
    }
  }
}
</style>