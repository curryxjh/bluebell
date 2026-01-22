<template>
  <div class="community-detail">
    <div class="container">
      <!-- 社区信息头部 -->
      <div class="community-header">
        <div class="header-content">
          <div class="avatar"></div>
          <div class="info">
            <h1 class="community-name">{{ community.name }}</h1>
            <p class="community-intro">{{ community.introduction }}</p>
            <div class="stats">
              <span class="stat-item">
                <i class="iconfont icon-user"></i>
                <span>成员数: {{ memberCount }}</span>
              </span>
              <span class="stat-item">
                <i class="iconfont icon-time"></i>
                <span>创建于: {{ formatDate(community.create_time) }}</span>
              </span>
            </div>
          </div>
          <div class="actions">
            <button 
              class="join-btn" 
              :class="{ 'joined': isJoined }"
              @click="toggleJoin"
            >
              {{ isJoined ? '已加入' : '加入社区' }}
            </button>
            <button class="back-btn" @click="goBack">返回</button>
          </div>
        </div>
      </div>

      <!-- 社区帖子列表 -->
      <div class="community-posts">
        <div class="posts-header">
          <h2>社区帖子</h2>
          <div class="sort-btns">
            <button 
              class="sort-btn" 
              :class="{ active: order === 'time' }"
              @click="changeOrder('time')"
            >
              <i class="iconfont icon-polygonred"></i>最新
            </button>
            <button 
              class="sort-btn" 
              :class="{ active: order === 'score' }"
              @click="changeOrder('score')"
            >
              <i class="iconfont icon-top"></i>热门
            </button>
          </div>
        </div>

        <ul class="post-list" v-if="posts.length > 0">
          <li class="post-item" v-for="post in posts" :key="post.id" @click="goToPost(post.id)">
            <div class="vote-section">
              <span class="iconfont icon-up" @click.stop="vote(post.id, '1')"></span>
              <span class="vote-count">{{ post.vote_num }}</span>
              <span class="iconfont icon-down" @click.stop="vote(post.id, '-1')"></span>
            </div>
            <div class="post-content">
              <h3 class="post-title">{{ post.title }}</h3>
              <p class="post-excerpt">{{ post.content }}</p>
              <div class="post-meta">
                <span class="author">作者: {{ post.author_name || '匿名' }}</span>
                <span class="time">{{ formatDate(post.create_time) }}</span>
              </div>
            </div>
          </li>
        </ul>

        <div class="empty-state" v-else>
          <p>该社区暂无帖子</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "CommunityDetail",
  data() {
    return {
      community: {
        id: 0,
        name: '',
        introduction: '',
        create_time: ''
      },
      posts: [],
      isJoined: false,
      memberCount: 0,
      order: 'time',
      page: 1,
      size: 20
    };
  },
  methods: {
    getCommunityDetail() {
      const communityId = this.$route.params.id;
      this.$axios({
        method: "get",
        url: `/community/${communityId}`
      })
        .then(response => {
          console.log("社区详情:", response.data);
          if (response.code == 1000) {
            this.community = response.data;
          } else {
            this.$message?.error(response.msg || '获取社区详情失败');
          }
        })
        .catch(error => {
          console.log("获取社区详情失败:", error);
          this.$message?.error('网络错误，请稍后重试');
        });
    },
    getCommunityPosts() {
      const communityId = this.$route.params.id;
      this.$axios({
        method: "get",
        url: "/posts2",
        params: {
          page: this.page,
          size: this.size,
          order: this.order,
          community_id: communityId
        }
      })
        .then(response => {
          console.log("社区帖子:", response.data);
          if (response.code == 1000) {
            this.posts = response.data || [];
          } else {
            console.log(response.msg);
          }
        })
        .catch(error => {
          console.log("获取社区帖子失败:", error);
        });
    },
    checkIfJoined() {
      // 检查用户是否已加入该社区
      this.$axios({
        method: "get",
        url: "/community/user"
      })
        .then(response => {
          if (response.code == 1000 && response.data) {
            const communityId = parseInt(this.$route.params.id);
            this.isJoined = response.data.some(c => c.id === communityId);
          }
        })
        .catch(error => {
          console.log("检查加入状态失败:", error);
        });
    },
    toggleJoin() {
      const communityId = this.$route.params.id;
      const url = this.isJoined ? "/community/leave" : "/community/join";
      
      this.$axios({
        method: "post",
        url: url,
        data: JSON.stringify({
          community_id: parseInt(communityId)
        })
      })
        .then(response => {
          if (response.code == 1000) {
            this.isJoined = !this.isJoined;
            this.$message?.success(this.isJoined ? '成功加入社区！' : '已退出社区');
          } else {
            this.$message?.error(response.msg || '操作失败');
          }
        })
        .catch(error => {
          console.log("操作失败:", error);
          this.$message?.error('网络错误，请稍后重试');
        });
    },
    changeOrder(order) {
      this.order = order;
      this.getCommunityPosts();
    },
    vote(post_id, direction) {
      this.$axios({
        method: "post",
        url: "/vote",
        data: JSON.stringify({
          post_id: post_id,
          direction: direction
        })
      })
        .then(response => {
          if (response.code == 1000) {
            console.log("投票成功");
            // 重新获取帖子列表以更新投票数
            this.getCommunityPosts();
          } else {
            this.$message?.error(response.msg || '投票失败');
          }
        })
        .catch(error => {
          console.log("投票失败:", error);
        });
    },
    goToPost(postId) {
      this.$router.push({ name: "Content", params: { id: postId } });
    },
    goBack() {
      this.$router.go(-1);
    },
    formatDate(dateString) {
      if (!dateString) return '';
      const date = new Date(dateString);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
    }
  },
  mounted() {
    this.getCommunityDetail();
    this.getCommunityPosts();
    this.checkIfJoined();
  }
};
</script>

<style scoped lang="less">
.community-detail {
  min-height: calc(100vh - 56px);
  background: #f5f7fa;
  padding: 20px;
  margin-top: 56px;

  .container {
    max-width: 1200px;
    margin: 0 auto;

    .community-header {
      background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
      border-radius: 16px;
      padding: 40px;
      margin-bottom: 24px;
      box-shadow: 0 8px 16px rgba(59, 130, 246, 0.2);

      .header-content {
        display: flex;
        align-items: flex-start;
        gap: 24px;

        .avatar {
          width: 80px;
          height: 80px;
          background: url("../assets/images/avatar.png") no-repeat;
          background-size: cover;
          border-radius: 50%;
          border: 4px solid rgba(255, 255, 255, 0.3);
          flex-shrink: 0;
        }

        .info {
          flex: 1;
          color: #ffffff;

          .community-name {
            font-size: 32px;
            font-weight: 700;
            margin-bottom: 12px;
            text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          }

          .community-intro {
            font-size: 16px;
            line-height: 1.6;
            margin-bottom: 16px;
            opacity: 0.95;
          }

          .stats {
            display: flex;
            gap: 24px;
            font-size: 14px;

            .stat-item {
              display: flex;
              align-items: center;
              gap: 8px;
              opacity: 0.9;

              .iconfont {
                font-size: 16px;
              }
            }
          }
        }

        .actions {
          display: flex;
          flex-direction: column;
          gap: 12px;
          flex-shrink: 0;

          .join-btn, .back-btn {
            min-width: 120px;
            height: 44px;
            border: none;
            border-radius: 12px;
            font-size: 14px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
          }

          .join-btn {
            background: #ffffff;
            color: #3b82f6;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);

            &:hover:not(.joined) {
              transform: translateY(-2px);
              box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
            }

            &.joined {
              background: rgba(255, 255, 255, 0.2);
              color: #ffffff;
              border: 2px solid rgba(255, 255, 255, 0.5);
              cursor: default;
            }
          }

          .back-btn {
            background: rgba(255, 255, 255, 0.15);
            color: #ffffff;
            border: 2px solid rgba(255, 255, 255, 0.3);

            &:hover {
              background: rgba(255, 255, 255, 0.25);
              border-color: rgba(255, 255, 255, 0.5);
            }
          }
        }
      }
    }

    .community-posts {
      background: #ffffff;
      border-radius: 16px;
      padding: 24px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);

      .posts-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 24px;
        padding-bottom: 16px;
        border-bottom: 2px solid #e2e8f0;

        h2 {
          font-size: 24px;
          font-weight: 700;
          color: #1a1a1b;
        }

        .sort-btns {
          display: flex;
          gap: 8px;

          .sort-btn {
            display: flex;
            align-items: center;
            gap: 6px;
            padding: 8px 16px;
            border: none;
            border-radius: 10px;
            background: #f1f5f9;
            color: #64748b;
            font-size: 14px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

            &:hover {
              background: #e2e8f0;
              color: #1e293b;
            }

            &.active {
              background: #3b82f6;
              color: #ffffff;
              box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
            }

            .iconfont {
              font-size: 16px;
            }
          }
        }
      }

      .post-list {
        display: flex;
        flex-direction: column;
        gap: 12px;

        .post-item {
          display: flex;
          gap: 16px;
          padding: 16px;
          border-radius: 12px;
          background: #f8fafc;
          cursor: pointer;
          transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
          list-style: none;

          &:hover {
            background: #f1f5f9;
            transform: translateX(4px);
          }

          .vote-section {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 4px;
            min-width: 40px;

            .iconfont {
              font-size: 20px;
              color: #64748b;
              cursor: pointer;
              transition: all 0.2s;

              &:hover {
                color: #3b82f6;
                transform: scale(1.2);
              }

              &.icon-down {
                transform: scaleY(-1);

                &:hover {
                  transform: scaleY(-1) scale(1.2);
                }
              }
            }

            .vote-count {
              font-size: 14px;
              font-weight: 600;
              color: #1a1a1b;
            }
          }

          .post-content {
            flex: 1;

            .post-title {
              font-size: 18px;
              font-weight: 600;
              color: #1a1a1b;
              margin-bottom: 8px;
            }

            .post-excerpt {
              font-size: 14px;
              color: #64748b;
              line-height: 1.5;
              margin-bottom: 12px;
              display: -webkit-box;
              -webkit-line-clamp: 2;
              -webkit-box-orient: vertical;
              overflow: hidden;
            }

            .post-meta {
              display: flex;
              gap: 16px;
              font-size: 12px;
              color: #94a3b8;

              .author, .time {
                display: flex;
                align-items: center;
              }
            }
          }
        }
      }

      .empty-state {
        text-align: center;
        padding: 80px 0;
        color: #64748b;
        font-size: 16px;
      }
    }
  }
}

/* 移动端适配 */
@media screen and (max-width: 768px) {
  .community-detail {
    padding: 10px;
    margin-top: 48px;

    .container {
      .community-header {
        padding: 20px;

        .header-content {
          flex-direction: column;

          .avatar {
            width: 60px;
            height: 60px;
          }

          .info {
            .community-name {
              font-size: 24px;
            }

            .community-intro {
              font-size: 14px;
            }

            .stats {
              flex-direction: column;
              gap: 8px;
            }
          }

          .actions {
            width: 100%;
            flex-direction: row;

            .join-btn, .back-btn {
              flex: 1;
              min-width: 0;
            }
          }
        }
      }

      .community-posts {
        padding: 16px;

        .posts-header {
          flex-direction: column;
          align-items: flex-start;
          gap: 12px;

          h2 {
            font-size: 20px;
          }
        }

        .post-list {
          .post-item {
            padding: 12px;

            .vote-section {
              min-width: 32px;

              .iconfont {
                font-size: 18px;
              }
            }

            .post-content {
              .post-title {
                font-size: 16px;
              }

              .post-excerpt {
                font-size: 13px;
              }

              .post-meta {
                flex-direction: column;
                gap: 4px;
              }
            }
          }
        }
      }
    }
  }
}
</style>
