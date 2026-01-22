
<template>
  <div class="content">
    <div class="left">
      <div class="post-name">创建新帖子</div>
      <div class="post-type">
    <!-- 关键修改：@click.stop -->
    <input 
      type="text" 
      class="post-type-value" 
      placeholder="选择一个频道" 
      v-model="selectCommunity.name" 
      @click.stop="showCommunity" 
      readonly
    />
    <ul class="post-type-options" v-show="showCommunityList">
      <li class="post-type-cell"
        v-for="(community, index) in communityList"
        :key="community.id"
        @click="selected(index)"
      >
        {{community.name}}
      </li>
    </ul>
    <!-- 这里的图标点击最好也加上开关功能 -->
    <i class="p-icon" @click.stop="showCommunity"></i>
  </div>
      <div class="post-content">
        <ul class="cat">
          <li class="cat-item active">
            <i class="iconfont icon-edit"></i>post
          </li>
        </ul>
        <div class="post-sub-container">
          <div class="post-sub-header">
            <textarea class="post-title" cols="30" rows="10" v-model="title" placeholder="标题" maxlength="100"></textarea>
            <span class="textarea-num">{{title.length}}/100</span>
          </div>
          <div class="post-text-con">
            <textarea
              class="post-content-t"
              cols="30"
              rows="10"
              v-model="content"
              placeholder="内容"
              maxlength="5000"
            ></textarea>
            <span class="content-num">{{content.length}}/5000</span>
          </div>
        </div>
        <div class="post-footer">
          <div class="btns">
            <button class="btn btn-cancel" @click="cancel">取消</button>
            <button class="btn btn-submit" @click="submit()" :disabled="isSubmitting">
              {{ isSubmitting ? '发布中...' : '发表' }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div class="right">
      <div class="post-rank">
        <h5 class="p-r-title">
          <i class="p-r-icon"></i>发帖规范
        </h5>
        <ul class="p-r-content">
          <li class="p-r-item">遵守法律法规，文明发言</li>
          <li class="p-r-item">禁止发布违禁、侵权内容</li>
          <li class="p-r-item">不得恶意灌水、刷屏</li>
          <li class="p-r-item">尊重他人，理性讨论</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Publish",
  data() {
    return {
      title: "",
      content: "",
      showCommunityList: false,
      selectCommunity: {},
      communityList: [],
      isSubmitting: false
    };
  },
  created() {
    document.addEventListener('click', this.closeDropdown);
  },
  beforeDestroy() {
    document.removeEventListener('click', this.closeDropdown);
  },
  methods: {
    async submit() {
      // 验证
      if (!this.selectCommunity.id) {
        this.$message?.error('请选择一个频道');
        return;
      }
      if (!this.title.trim() || this.title.trim().length < 5) {
        this.$message?.error('标题至少需要5个字符');
        return;
      }
      if (!this.content.trim() || this.content.trim().length < 10) {
        this.$message?.error('内容至少需要10个字符');
        return;
      }

      this.isSubmitting = true;

      try {
        const response = await this.$axios({
          method: "post",
          url: "/post",
          data: JSON.stringify({
            title: this.title.trim(),
            content: this.content.trim(),
            community_id: this.selectCommunity.id
          })
        });

        if (response.code == 1000) {
          this.$message?.success('发布成功！');
          setTimeout(() => {
            this.$router.push({ path: "/" });
          }, 500);
        } else {
          this.$message?.error(response.msg || '发布失败，请重试');
        }
      } catch (error) {
        console.log(error);
        this.$message?.error('网络错误，请稍后重试');
      } finally {
        this.isSubmitting = false;
      }
    },
    cancel() {
      if (this.title || this.content) {
        if (confirm('确定要取消吗？未保存的内容将会丢失。')) {
          this.$router.push({ path: "/" });
        }
      } else {
        this.$router.push({ path: "/" });
      }
    },
    getCommunityList() {
      this.$axios({
        method: "get",
        url: "/community"
      })
        .then(response => {
          if (response.code == 1000) {
            this.communityList = response.data || [];
          } else {
            this.$message?.error('获取社区列表失败');
          }
        })
        .catch(error => {
          console.log("获取社区失败:", error);
          this.$message?.error('网络错误，请稍后重试');
        });
    },
    showCommunity(){
      this.showCommunityList = !this.showCommunityList;
    },
    closeDropdown() {
      this.showCommunityList = false;
    },
    selected(index) {
      this.selectCommunity = this.communityList[index];
      this.showCommunityList = false;
    }
  },
  mounted: function() {
    this.getCommunityList();
  }
};
</script>

<style lang="less" scoped>
.content {
  max-width: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: row;
  justify-content: center;
  margin: 0 auto;
  padding: 20px 24px;
  margin-top: 56px;
  background: #f5f7fa;
  min-height: calc(100vh - 56px);
  
  .left {
    flex-grow: 1;
    max-width: 740px;
    word-break: break-word;
    flex: 1;
    margin: 32px;
    margin-right: 12px;
    padding-bottom: 30px;
    position: relative;
    
    .post-name {
      font-size: 24px;
      font-weight: 700;
      color: #1a1a1b;
      padding: 12px 4px;
      margin: 16px 0;
      border-bottom: 2px solid #3b82f6;
    }
    
    .post-type {
      position: relative;
      box-sizing: border-box;
      width: 100%;
      max-width: 400px;
      height: 48px;
      border-radius: 12px;
      transition: all 0.3s ease;
      border: 2px solid #e2e8f0;
      background-color: #ffffff;
      padding-left: 16px;
      margin-bottom: 20px;
      
      &:hover {
        border-color: #cbd5e1;
      }
      
      .post-type-value {
        font-size: 15px;
        font-weight: 500;
        line-height: 48px;
        width: calc(100% - 40px);
        vertical-align: middle;
        color: #1a1a1b;
        background-color: transparent;
        cursor: pointer;
        border: none;
        outline: none;
        
        &::placeholder {
          color: #94a3b8;
        }
      }
      
      .post-type-options {
        position: absolute;
        width: 100%;
        background-color: white;
        left: 0;
        top: calc(100% + 8px);
        z-index: 100;
        border-radius: 12px;
        box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
        border: 1px solid #e2e8f0;
        max-height: 280px;
        overflow-y: auto;
        
        .post-type-cell {
          margin: 0;
          padding: 14px 16px;
          font-size: 14px;
          list-style: none;
          border-bottom: 1px solid #f1f5f9;
          color: #1a1a1b;
          cursor: pointer;
          transition: all 0.2s ease;
          
          &:last-child {
            border-bottom: none;
          }
          
          &:hover {
            background: #f8fafc;
            padding-left: 20px;
          }
        }
      }
      
      .p-icon {
        width: 0;
        height: 0;
        border-left: 6px solid transparent;
        border-right: 6px solid transparent;
        border-top: 6px solid #64748b;
        position: absolute;
        top: 50%;
        right: 16px;
        margin-top: -3px;
        cursor: pointer;
        transition: transform 0.3s ease;
      }
    }
    
    .post-content {
      background: #ffffff;
      margin: 10px 0;
      padding-bottom: 15px;
      border-radius: 16px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
      border: none;
      overflow: hidden;
      
      .cat {
        display: flex;
        align-items: center;
        width: 100%;
        height: 60px;
        padding: 0 24px;
        border-bottom: 2px solid #f1f5f9;
        
        .cat-item {
          display: flex;
          align-items: center;
          gap: 8px;
          height: 40px;
          line-height: 40px;
          list-style: none;
          color: #3b82f6;
          font-weight: 600;
          font-size: 15px;
          
          .iconfont {
            font-size: 18px;
          }
        }
      }
      
      .post-sub-container {
        padding: 24px;
        
        .post-sub-header {
          position: relative;
          margin-bottom: 20px;
          
          .post-title {
            resize: none;
            box-sizing: border-box;
            overflow: hidden;
            display: block;
            width: 100%;
            height: 56px;
            padding: 16px;
            outline: none;
            border: 2px solid #e2e8f0;
            border-radius: 12px;
            color: #1a1a1b;
            font-size: 16px;
            font-weight: 500;
            line-height: 24px;
            transition: all 0.3s ease;
            font-family: inherit;
            
            &::placeholder {
              color: #94a3b8;
            }
            
            &:hover {
              border-color: #cbd5e1;
            }
            
            &:focus {
              border-color: #3b82f6;
              box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
            }
          }
          
          .textarea-num {
            font-size: 12px;
            font-weight: 500;
            color: #94a3b8;
            position: absolute;
            bottom: 12px;
            right: 16px;
            pointer-events: none;
          }
        }
        
        .post-text-con {
          width: 100%;
          position: relative;
          
          .post-content-t {
            resize: vertical;
            box-sizing: border-box;
            display: block;
            width: 100%;
            min-height: 280px;
            padding: 16px;
            outline: none;
            border: 2px solid #e2e8f0;
            border-radius: 12px;
            color: #1a1a1b;
            font-size: 15px;
            font-weight: 400;
            line-height: 1.6;
            transition: all 0.3s ease;
            font-family: inherit;
            
            &::placeholder {
              color: #94a3b8;
            }
            
            &:hover {
              border-color: #cbd5e1;
            }
            
            &:focus {
              border-color: #3b82f6;
              box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
            }
          }
          
          .content-num {
            font-size: 12px;
            font-weight: 500;
            color: #94a3b8;
            position: absolute;
            bottom: 12px;
            right: 16px;
            pointer-events: none;
          }
        }
      }
      
      .post-footer {
        display: flex;
        margin: 0 24px;
        padding-top: 20px;
        border-top: 1px solid #f1f5f9;
        justify-content: flex-end;
        
        .btns {
          display: flex;
          gap: 12px;
          
          .btn {
            min-width: 100px;
            height: 44px;
            border: none;
            border-radius: 12px;
            box-sizing: border-box;
            text-align: center;
            font-size: 14px;
            font-weight: 600;
            padding: 0 24px;
            cursor: pointer;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            
            &:disabled {
              opacity: 0.6;
              cursor: not-allowed;
            }
            
            &.btn-cancel {
              background: transparent;
              border: 2px solid #e2e8f0;
              color: #64748b;
              
              &:hover:not(:disabled) {
                background: #f1f5f9;
                border-color: #cbd5e1;
                color: #1e293b;
              }
            }
            
            &.btn-submit {
              background: #3b82f6;
              color: #ffffff;
              box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
              
              &:hover:not(:disabled) {
                background: #2563eb;
                transform: translateY(-2px);
                box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
              }
              
              &:active:not(:disabled) {
                transform: translateY(0);
                box-shadow: 0 2px 4px rgba(59, 130, 246, 0.3);
              }
            }
          }
        }
      }
    }
  }
  
  .right {
    flex-grow: 0;
    width: 312px;
    margin-top: 32px;
    
    .post-rank {
      background: #ffffff;
      border-radius: 16px;
      margin-top: 15px;
      padding: 24px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
      border: none;
      
      .p-r-title {
        display: flex;
        align-items: center;
        font-size: 16px;
        font-weight: 600;
        color: #1a1a1b;
        border-bottom: 2px solid #f1f5f9;
        padding-bottom: 12px;
        margin-bottom: 16px;
        
        .p-r-icon {
          width: 36px;
          height: 36px;
          background: url("../assets/images/avatar.png") no-repeat;
          background-size: cover;
          margin-right: 10px;
        }
      }
      
      .p-r-content {
        display: flex;
        flex-direction: column;
        
        .p-r-item {
          list-style: none;
          border-bottom: 1px solid #f1f5f9;
          color: #64748b;
          padding: 12px 8px;
          font-size: 14px;
          line-height: 1.5;
          
          &:last-child {
            border-bottom: none;
          }
        }
      }
    }
  }
}

/* 移动端适配 */
@media screen and (max-width: 768px) {
  .content {
    flex-direction: column;
    padding: 10px 8px;
    margin-top: 48px;
    
    .left {
      max-width: 100%;
      margin: 15px 0;
      margin-right: 0;
      
      .post-name {
        margin: 10px 0;
        font-size: 20px;
      }
      
      .post-type {
        width: 100%;
        max-width: 100%;
        height: 44px;
        
        .post-type-value {
          line-height: 44px;
          font-size: 14px;
        }
      }
      
      .post-content {
        margin: 8px 0;
        padding-bottom: 12px;
        
        .cat {
          height: 50px;
          padding: 0 16px;
          
          .cat-item {
            font-size: 14px;
            
            .iconfont {
              font-size: 16px;
            }
          }
        }
        
        .post-sub-container {
          padding: 16px;
          
          .post-sub-header {
            margin-bottom: 16px;
            
            .post-title {
              height: 48px;
              font-size: 15px;
              padding: 12px;
            }
          }
          
          .post-text-con {
            .post-content-t {
              min-height: 200px;
              padding: 12px;
              font-size: 14px;
            }
          }
        }
        
        .post-footer {
          margin: 0 16px;
          padding-top: 16px;
          flex-direction: column-reverse;
          
          .btns {
            flex-direction: column-reverse;
            
            .btn {
              width: 100%;
              height: 42px;
              font-size: 13px;
            }
          }
        }
      }
    }
    
    .right {
      width: 100%;
      margin-top: 0;
      
      .post-rank {
        margin-top: 12px;
        padding: 16px;
        
        .p-r-title {
          font-size: 15px;
          padding-bottom: 10px;
          
          .p-r-icon {
            width: 32px;
            height: 32px;
            margin-right: 8px;
          }
        }
        
        .p-r-content {
          .p-r-item {
            padding: 10px 6px;
            font-size: 13px;
          }
        }
      }
    }
  }
}
</style>
