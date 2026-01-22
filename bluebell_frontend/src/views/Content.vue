<template>
  <div class="content">
    <div class="left">
      <!-- 返回按钮 -->
      <div class="back-nav">
        <button class="back-btn" @click.stop="goBack">
          <i class="iconfont icon-left"></i>
          返回
        </button>
      </div>

      <!-- 帖子内容 -->
      <div class="post-container" v-if="post.post_id || post.id">
        <div class="vote-section">
          <!-- 这里的 vote('1') 和 vote('-1') 内部会自动转字符串 -->
          <button class="vote-btn" @click.stop="vote('1')">
            <span class="iconfont icon-up"></span>
          </button>
          <span class="vote-count">{{post.vote_num || 0}}</span>
          <button class="vote-btn" @click.stop="vote('-1')">
            <span class="iconfont icon-down"></span>
          </button>
        </div>

        <div class="post-content">
          <div class="post-header">
            <h1 class="post-title">{{post.title}}</h1>
            <div class="post-meta">
              <span class="community-tag" v-if="post.community_name" @click="goToCommunity">
                b/{{post.community_name}}
              </span>
              <span class="author">
                作者: {{post.author_name || '匿名'}}
              </span>
              <span class="time">
                {{formatDate(post.create_time)}}
              </span>
            </div>
          </div>

          <div class="post-body">
            <div class="post-text">{{post.content}}</div>
          </div>

          <div class="post-footer">
            <div class="post-stats">
              <span class="stat-item">
                <i class="iconfont icon-comment"></i>
                <span>{{comments.length}} 条评论</span>
              </span>
            </div>
          </div>
        </div>
      </div>
      <!-- 加载状态 -->
      <div class="post-container" v-else>
         <div class="post-content" style="text-align:center; padding: 40px;">加载中...</div>
      </div>

      <!-- 评论区域 -->
      <div class="comments-section">
        <div class="comments-header">
          <h3>评论</h3>
          <span class="comment-count">{{comments.length}} 条评论</span>
        </div>

        <!-- 评论输入框 -->
        <div class="comment-input-wrapper">
          <div class="user-avatar"></div>
          <div class="comment-input-box">
            <textarea 
              class="comment-input" 
              v-model="newComment"
              placeholder="写下你的评论..."
              maxlength="500"
            ></textarea>
            <div class="comment-actions">
              <span class="char-count">{{newComment.length}}/500</span>
              <button 
                class="submit-comment-btn" 
                @click="submitComment"
                :disabled="!newComment.trim() || isSubmittingComment"
              >
                {{isSubmittingComment ? '发送中...' : '发送'}}
              </button>
            </div>
          </div>
        </div>

        <!-- 评论列表 -->
        <div class="comments-list" v-if="comments.length > 0">
          <div class="comment-item" v-for="comment in comments" :key="comment.id">
            <div class="comment-avatar"></div>
            <div class="comment-content">
              <div class="comment-header">
                <span class="comment-author">{{comment.author_name}}</span>
                <span class="comment-time">{{formatDate(comment.create_time)}}</span>
              </div>
              <div class="comment-text">{{comment.content}}</div>
              <div class="comment-footer">
                <button class="reply-btn">
                  <i class="iconfont icon-comment"></i>
                  回复
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="no-comments" v-else>
          <p>还没有评论，快来抢沙发吧！</p>
        </div>
      </div>
    </div>

    <!-- 右侧边栏 -->
    <div class="right">
      <div class="community-card">
        <div class="community-header">
          <div class="header-bg"></div>
        </div>
        <div class="community-body">
          <div class="community-avatar"></div>
          <h3 class="community-name">b/{{communityDetail.name || '加载中...'}}</h3>
          <p class="community-desc">{{communityDetail.introduction || '暂无社区简介'}}</p>
          
          <div class="community-stats">
            <div class="stat-item">
              <div class="stat-number">{{communityDetail.member_count || 0}}</div>
              <div class="stat-label">成员</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">{{post.view_count || 0}}</div>
              <div class="stat-label">浏览</div>
            </div>
          </div>

          <div class="community-date" v-if="communityDetail.create_time">
            创建于 {{formatDate(communityDetail.create_time)}}
          </div>

          <button class="join-btn" @click="joinCommunity" :disabled="isJoined">
            {{isJoined ? '已加入' : '加入社区'}}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "PostDetail",
  data() {
    return {
      post: {},
      communityDetail: {},
      isJoined: false,
      newComment: '',
      isSubmittingComment: false,
      comments: []
    };
  },
  methods: {
    goBack() {
      this.$router.go(-1);
    },
    // 获取帖子 ID 的统一方法，解决精度丢失问题
    getRealPostId() {
      // 这里的 post_id 通常是后端返回的字符串，id 可能是路由里的
      return String(this.post.post_id || this.post.id || this.$route.params.id);
    },
    getPostDetail() {
      const id = String(this.$route.params.id); // 路由传参
      this.$axios({
        method: "get",
        url: "/post/" + id
      })
        .then(response => {
          if (response.code == 1000) {
            this.post = response.data;
            // 成功拿到帖子后，根据帖子的社区ID去拉取侧边栏社区详情
            if (this.post.community_id) {
              this.getCommunityDetail(this.post.community_id);
            }
          } else {
            this.$message?.error(response.msg || '获取帖子详情失败');
          }
        })
        .catch(error => {
          console.error("获取详情错误:", error);
          this.$message?.error('网络错误，请稍后重试');
        });
    },
    getCommunityDetail(communityId) {
      this.$axios({
        method: "get",
        url: "/community/" + communityId
      })
        .then(response => {
          if (response.code == 1000) {
            this.communityDetail = response.data;
          }
        })
        .catch(error => {
          console.error("获取社区详情失败:", error);
        });
    },
    vote(direction) {
      const token = window.localStorage.getItem('token');
      if (!token) {
        this.$message?.warning('请先登录后再投票');
        return;
      }

      const post_id = this.getRealPostId();
      
      this.$axios({
        method: "post",
        url: "/vote",
        data: {
          post_id: post_id, 
          direction: String(direction) // 后端要求 "1", "0", "-1"
        }
      })
        .then(response => {
          if (response.code == 1000) {
            this.$message?.success('投票成功');
            this.getPostDetail(); // 刷新点赞数
          } else {
            // 这里会拦截到“投票时间已过”等业务错误
            this.$message?.error(response.msg || '操作失败');
          }
        })
        .catch(error => {
          console.error("投票失败:", error);
          this.$message?.error('操作失败，请重试');
        });
    },
    submitComment() {
      const post_id = this.getRealPostId();
      const content = this.newComment.trim();

      if (!content) return;
      this.isSubmittingComment = true;

      this.$axios({
        method: "post",
        url: "/comment",
        data: {
          post_id: post_id,
          content: content,
          parent_id: "0" // 后端如果是 snowflake ID，parent_id 也要字符串
        }
      })
        .then(response => {
          if (response.code == 1000) {
            this.$message?.success('发表成功');
            this.newComment = '';
            this.getComments(); // 刷新列表
          } else {
            this.$message?.error(response.msg || '发表失败');
          }
        })
        .finally(() => {
          this.isSubmittingComment = false;
        });
    },
    getComments() {
      const id = String(this.$route.params.id);
      this.$axios({
        method: "get",
        url: "/comments/" + id
      })
        .then(response => {
          if (response.code == 1000) {
            this.comments = response.data || [];
          }
        });
    },
    joinCommunity() {
      const communityId = this.post.community_id || this.communityDetail.id;
      if (!communityId) return;

      this.$axios({
        method: "post",
        url: "/community/join",
        data: {
          community_id: String(communityId)
        }
      })
        .then(response => {
          if (response.code == 1000) {
            this.isJoined = true;
            this.$message?.success('欢迎加入！');
          }
        })
        .catch(() => {
          this.$message?.error('请先登录');
        });
    },
    formatDate(dateString) {
      if (!dateString) return '';
      const date = new Date(dateString);
      const now = new Date();
      const diff = now - date;
      const minutes = Math.floor(diff / 60000);
      const hours = Math.floor(diff / 3600000);
      
      if (minutes < 60) return minutes <= 0 ? '刚刚' : `${minutes}分钟前`;
      if (hours < 24) return `${hours}小时前`;
      
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
    }
  },
  mounted() {
    this.getPostDetail();
    this.getComments();
  }
};
</script>

<!-- 样式保持你原来的不变即可 -->

<style lang="less" scoped>
* {
  box-sizing: border-box;
}

.content {
  max-width: 1200px;
  margin: 56px auto 0;
  padding: 20px 24px;
  background: #f5f7fa;
  min-height: calc(100vh - 56px);
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: 24px;
  
  .left {
    min-width: 0;
    overflow: hidden;
    box-sizing: border-box;
    
    .back-nav {
      margin-bottom: 16px;
      
      .back-btn {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 8px 16px;
        background: #ffffff;
        border: 2px solid #e2e8f0;
        border-radius: 10px;
        color: #64748b;
        font-size: 14px;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.3s ease;
        
        &:hover {
          background: #f1f5f9;
          border-color: #cbd5e1;
          color: #1e293b;
        }
        
        .iconfont {
          font-size: 14px;
        }
      }
    }
    
    .post-container {
      background: #ffffff;
      border-radius: 16px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
      overflow: hidden;
      display: flex;
      padding: 24px;
      gap: 16px;
      margin-bottom: 20px;
      box-sizing: border-box;
      
      .vote-section {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 8px;
        padding: 8px;
        background: #f8fafc;
        border-radius: 12px;
        height: fit-content;
        
        .vote-btn {
          width: 32px;
          height: 32px;
          display: flex;
          align-items: center;
          justify-content: center;
          background: transparent;
          border: none;
          border-radius: 8px;
          cursor: pointer;
          transition: all 0.2s ease;
          
          .iconfont {
            font-size: 20px;
            color: #64748b;
          }
          
          &:hover {
            background: #e2e8f0;
            
            .iconfont {
              color: #3b82f6;
            }
          }
        }
        
        .vote-count {
          font-size: 14px;
          font-weight: 700;
          color: #1a1a1b;
          min-width: 32px;
          text-align: center;
        }
      }
      
      .post-content {
        flex: 1;
        min-width: 0;
        box-sizing: border-box;
        
        .post-header {
          margin-bottom: 20px;
          
          .post-title {
            font-size: 28px;
            font-weight: 700;
            color: #1a1a1b;
            margin: 0 0 12px 0;
            line-height: 1.3;
            word-break: break-word;
          }
          
          .post-meta {
            display: flex;
            flex-wrap: wrap;
            gap: 12px;
            font-size: 13px;
            color: #64748b;
            
            .community-tag {
              display: inline-flex;
              align-items: center;
              padding: 4px 12px;
              background: #eff6ff;
              color: #3b82f6;
              border-radius: 6px;
              font-weight: 600;
              cursor: pointer;
              transition: all 0.2s ease;
              
              &:hover {
                background: #dbeafe;
                color: #2563eb;
              }
            }
            
            .author, .time {
              display: flex;
              align-items: center;
            }
          }
        }
        
        .post-body {
          margin-bottom: 24px;
          
          .post-text {
            font-size: 16px;
            line-height: 1.8;
            color: #1a1a1b;
            word-break: break-word;
            white-space: pre-wrap;

            &.empty {
              color: #94a3b8;
              text-align: center;
              padding: 40px 0;
            }
          }
        }
        
        .post-footer {
          padding-top: 20px;
          border-top: 1px solid #e2e8f0;
          
          .post-stats {
            display: flex;
            gap: 16px;
            
            .stat-item {
              display: flex;
              align-items: center;
              gap: 6px;
              color: #64748b;
              font-size: 14px;
              
              .iconfont {
                font-size: 16px;
              }
            }
          }
        }
      }
    }
    
    .comments-section {
      background: #ffffff;
      border-radius: 16px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
      padding: 24px;
      overflow: hidden;
      box-sizing: border-box;
      
      .comments-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 24px;
        padding-bottom: 16px;
        border-bottom: 2px solid #f1f5f9;
        
        h3 {
          font-size: 20px;
          font-weight: 700;
          color: #1a1a1b;
          margin: 0;
        }
        
        .comment-count {
          font-size: 14px;
          color: #94a3b8;
          font-weight: 500;
        }
      }

      .comment-input-wrapper {
        display: flex;
        gap: 12px;
        margin-bottom: 24px;
        padding-bottom: 24px;
        border-bottom: 1px solid #e2e8f0;
        width: 100%;
        box-sizing: border-box;

        .user-avatar {
          width: 40px;
          height: 40px;
          background: url("../assets/images/avatar.png") no-repeat;
          background-size: cover;
          border-radius: 50%;
          flex-shrink: 0;
        }

        .comment-input-box {
          flex: 1;
          min-width: 0;
          box-sizing: border-box;

          .comment-input {
            width: 100%;
            min-height: 80px;
            padding: 12px 16px;
            border: 2px solid #e2e8f0;
            border-radius: 12px;
            font-size: 14px;
            line-height: 1.6;
            color: #1a1a1b;
            resize: vertical;
            outline: none;
            transition: all 0.3s ease;
            font-family: inherit;
            margin-bottom: 8px;
            box-sizing: border-box;

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

          .comment-actions {
            display: flex;
            justify-content: space-between;
            align-items: center;

            .char-count {
              font-size: 12px;
              color: #94a3b8;
              font-weight: 500;
            }

            .submit-comment-btn {
              padding: 8px 20px;
              background: #3b82f6;
              color: #ffffff;
              border: none;
              border-radius: 10px;
              font-size: 14px;
              font-weight: 600;
              cursor: pointer;
              transition: all 0.3s ease;
              box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);

              &:hover:not(:disabled) {
                background: #2563eb;
                transform: translateY(-2px);
                box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
              }

              &:disabled {
                background: #e2e8f0;
                color: #94a3b8;
                cursor: not-allowed;
                box-shadow: none;
              }
            }
          }
        }
      }

      .comments-list {
        display: flex;
        flex-direction: column;
        gap: 20px;

        .comment-item {
          display: flex;
          gap: 12px;

          .comment-avatar {
            width: 36px;
            height: 36px;
            background: url("../assets/images/avatar.png") no-repeat;
            background-size: cover;
            border-radius: 50%;
            flex-shrink: 0;
          }

          .comment-content {
            flex: 1;
            min-width: 0;

            .comment-header {
              display: flex;
              align-items: center;
              gap: 12px;
              margin-bottom: 8px;

              .comment-author {
                font-size: 14px;
                font-weight: 600;
                color: #1a1a1b;
              }

              .comment-time {
                font-size: 12px;
                color: #94a3b8;
              }
            }

            .comment-text {
              font-size: 14px;
              line-height: 1.6;
              color: #1a1a1b;
              margin-bottom: 8px;
              word-break: break-word;
              white-space: pre-wrap;
            }

            .comment-footer {
              .reply-btn {
                display: inline-flex;
                align-items: center;
                gap: 4px;
                padding: 4px 12px;
                background: transparent;
                border: none;
                color: #64748b;
                font-size: 13px;
                font-weight: 500;
                cursor: pointer;
                border-radius: 6px;
                transition: all 0.2s ease;

                &:hover {
                  background: #f1f5f9;
                  color: #3b82f6;
                }

                .iconfont {
                  font-size: 14px;
                }
              }
            }
          }
        }
      }
      
      .no-comments {
        text-align: center;
        padding: 60px 20px;
        
        p {
          font-size: 15px;
          color: #94a3b8;
          margin: 0;
        }
      }
    }
  }
  
  .right {
    .community-card {
      background: #ffffff;
      border-radius: 16px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
      overflow: hidden;
      position: sticky;
      top: 76px;
      
      .community-header {
        .header-bg {
          width: 100%;
          height: 80px;
          background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
        }
      }
      
      .community-body {
        padding: 0 20px 24px;
        
        .community-avatar {
          width: 72px;
          height: 72px;
          background: url("../assets/images/avatar.png") no-repeat;
          background-size: cover;
          border-radius: 50%;
          border: 4px solid #ffffff;
          margin: -36px 0 16px 0;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }
        
        .community-name {
          font-size: 20px;
          font-weight: 700;
          color: #1a1a1b;
          margin: 0 0 8px 0;
        }
        
        .community-desc {
          font-size: 14px;
          line-height: 1.6;
          color: #64748b;
          margin: 0 0 20px 0;
        }
        
        .community-stats {
          display: flex;
          gap: 20px;
          margin-bottom: 20px;
          padding-bottom: 20px;
          border-bottom: 1px solid #e2e8f0;
          
          .stat-item {
            flex: 1;
            
            .stat-number {
              font-size: 20px;
              font-weight: 700;
              color: #1a1a1b;
              margin-bottom: 4px;
            }
            
            .stat-label {
              font-size: 12px;
              color: #94a3b8;
              font-weight: 500;
            }
          }
        }
        
        .community-date {
          font-size: 13px;
          color: #64748b;
          margin-bottom: 16px;
        }
        
        .join-btn {
          width: 100%;
          height: 44px;
          background: #3b82f6;
          color: #ffffff;
          border: none;
          border-radius: 12px;
          font-size: 15px;
          font-weight: 600;
          cursor: pointer;
          box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
          transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
          
          &:hover:not(:disabled) {
            background: #2563eb;
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
          }
          
          &:active:not(:disabled) {
            transform: translateY(0);
            box-shadow: 0 2px 4px rgba(59, 130, 246, 0.3);
          }
          
          &:disabled {
            background: #e2e8f0;
            color: #94a3b8;
            cursor: not-allowed;
            box-shadow: none;
          }
        }
      }
    }
  }
}

/* 移动端适配 */
@media screen and (max-width: 1024px) {
  .content {
    grid-template-columns: 1fr;
    padding: 16px 12px;
    
    .right {
      .community-card {
        position: static;
      }
    }
  }
}

@media screen and (max-width: 768px) {
  .content {
    margin-top: 48px;
    padding: 12px 8px;
    gap: 16px;
    
    .left {
      .back-nav {
        margin-bottom: 12px;
        
        .back-btn {
          padding: 6px 12px;
          font-size: 13px;
        }
      }
      
      .post-container {
        padding: 16px;
        gap: 12px;
        
        .vote-section {
          padding: 6px;
          
          .vote-btn {
            width: 28px;
            height: 28px;
            
            .iconfont {
              font-size: 18px;
            }
          }
          
          .vote-count {
            font-size: 13px;
          }
        }
        
        .post-content {
          .post-header {
            margin-bottom: 16px;
            
            .post-title {
              font-size: 22px;
            }
            
            .post-meta {
              font-size: 12px;
              gap: 8px;
            }
          }
          
          .post-body {
            .post-text {
              font-size: 15px;
              line-height: 1.7;
            }
          }
          
          .post-footer {
            padding-top: 16px;
            
            .post-stats {
              gap: 12px;
              
              .stat-item {
                font-size: 13px;
                
                .iconfont {
                  font-size: 15px;
                }
              }
            }
          }
        }
      }
      
      .comments-section {
        padding: 16px;
        
        .comments-header {
          margin-bottom: 16px;
          padding-bottom: 12px;
          
          h3 {
            font-size: 18px;
          }
          
          .comment-count {
            font-size: 13px;
          }
        }

        .comment-input-wrapper {
          gap: 10px;
          margin-bottom: 20px;
          padding-bottom: 20px;

          .user-avatar {
            width: 32px;
            height: 32px;
          }

          .comment-input-box {
            .comment-input {
              min-height: 60px;
              font-size: 13px;
              padding: 10px 12px;
            }

            .comment-actions {
              .char-count {
                font-size: 11px;
              }

              .submit-comment-btn {
                padding: 6px 16px;
                font-size: 13px;
              }
            }
          }
        }

        .comments-list {
          gap: 16px;

          .comment-item {
            gap: 10px;

            .comment-avatar {
              width: 32px;
              height: 32px;
            }

            .comment-content {
              .comment-header {
                gap: 8px;
                margin-bottom: 6px;

                .comment-author {
                  font-size: 13px;
                }

                .comment-time {
                  font-size: 11px;
                }
              }

              .comment-text {
                font-size: 13px;
                margin-bottom: 6px;
              }

              .comment-footer {
                .reply-btn {
                  padding: 3px 10px;
                  font-size: 12px;

                  .iconfont {
                    font-size: 13px;
                  }
                }
              }
            }
          }
        }
        
        .no-comments {
          padding: 40px 16px;
          
          p {
            font-size: 14px;
          }
        }
      }
    }
    
    .right {
      .community-card {
        .community-body {
          padding: 0 16px 20px;
          
          .community-avatar {
            width: 64px;
            height: 64px;
            margin: -32px 0 12px 0;
          }
          
          .community-name {
            font-size: 18px;
          }
          
          .community-desc {
            font-size: 13px;
            margin-bottom: 16px;
          }
          
          .community-stats {
            gap: 16px;
            margin-bottom: 16px;
            padding-bottom: 16px;
            
            .stat-item {
              .stat-number {
                font-size: 18px;
              }
              
              .stat-label {
                font-size: 11px;
              }
            }
          }
          
          .community-date {
            font-size: 12px;
            margin-bottom: 12px;
          }
          
          .join-btn {
            height: 40px;
            font-size: 14px;
          }
        }
      }
    }
  }
}
</style>
