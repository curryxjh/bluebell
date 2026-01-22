<template>
  <div class="all-communities">
    <div class="container">
      <div class="header">
        <h1 class="title">所有社区</h1>
        <p class="subtitle">探索并加入你感兴趣的社区</p>
      </div>
      
      <div class="communities-grid" v-if="communities.length > 0">
        <div class="community-card" v-for="community in communities" :key="community.id">
          <div class="card-header">
            <div class="avatar"></div>
            <div class="info">
              <h3 class="community-name">{{ community.name }}</h3>
              <p class="community-id">ID: {{ community.id }}</p>
            </div>
          </div>
          <div class="card-footer">
            <button 
              class="join-btn" 
              :class="{ 'joined': community.joined }"
              @click="joinCommunity(community.id)"
              :disabled="community.joined"
            >
              {{ community.joined ? '已加入' : '加入社区' }}
            </button>
            <button class="view-btn" @click="viewCommunity(community.id)">
              查看详情
            </button>
          </div>
        </div>
      </div>
      
      <div class="empty-state" v-else>
        <p>暂无社区数据</p>
      </div>
      
      <div class="back-btn-container">
        <button class="back-btn" @click="goBack">返回首页</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "AllCommunities",
  data() {
    return {
      communities: []
    };
  },
  methods: {
    getCommunities() {
      this.$axios({
        method: "get",
        url: "/community"
      })
        .then(response => {
          console.log("获取所有社区:", response.data);
          if (response.code == 1000) {
            this.communities = response.data.map(community => ({
              ...community,
              joined: false
            }));
          } else {
            this.$message?.error(response.msg || '获取社区列表失败');
          }
        })
        .catch(error => {
          console.log("获取社区列表失败:", error);
          this.$message?.error('网络错误，请稍后重试');
        });
    },
    joinCommunity(community_id) {
      this.$axios({
        method: "post",
        url: "/community/join",
        data: JSON.stringify({
          community_id: community_id
        })
      })
        .then(response => {
          console.log("加入社区响应:", response.data);
          if (response.code == 1000) {
            const community = this.communities.find(c => c.id === community_id);
            if (community) {
              community.joined = true;
            }
            this.$message?.success('成功加入社区！');
          } else {
            this.$message?.error(response.msg || '加入社区失败');
          }
        })
        .catch(error => {
          console.log("加入社区失败:", error);
          this.$message?.error('网络错误，请稍后重试');
        });
    },
    viewCommunity(community_id) {
      this.$router.push({ name: "CommunityDetail", params: { id: community_id } });
    },
    goBack() {
      this.$router.push({ name: "Home" });
    }
  },
  mounted() {
    this.getCommunities();
  }
};
</script>

<style scoped lang="less">
.all-communities {
  min-height: calc(100vh - 56px);
  background: #f5f7fa;
  padding: 20px;
  margin-top: 56px;
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
    
    .header {
      text-align: center;
      padding: 40px 0;
      
      .title {
        font-size: 32px;
        font-weight: 700;
        color: #1a1a1b;
        margin-bottom: 10px;
      }
      
      .subtitle {
        font-size: 16px;
        color: #64748b;
      }
    }
    
    .communities-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
      gap: 20px;
      margin-bottom: 40px;
      
      .community-card {
        background: #ffffff;
        border-radius: 16px;
        padding: 24px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        
        &:hover {
          transform: translateY(-4px);
          box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08), 0 4px 8px rgba(0, 0, 0, 0.06);
        }
        
        .card-header {
          display: flex;
          align-items: center;
          margin-bottom: 20px;
          
          .avatar {
            width: 48px;
            height: 48px;
            background: url("../assets/images/avatar.png") no-repeat;
            background-size: cover;
            border-radius: 50%;
            margin-right: 16px;
          }
          
          .info {
            flex: 1;
            
            .community-name {
              font-size: 18px;
              font-weight: 600;
              color: #1a1a1b;
              margin-bottom: 4px;
            }
            
            .community-id {
              font-size: 12px;
              color: #64748b;
            }
          }
        }
        
        .card-footer {
          display: flex;
          gap: 10px;
          
          .join-btn, .view-btn {
            flex: 1;
            height: 36px;
            border: none;
            border-radius: 10px;
            font-size: 13px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
          }
          
          .join-btn {
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
            
            &.joined {
              background: #e2e8f0;
              color: #64748b;
              box-shadow: none;
              cursor: default;
            }
          }
          
          .view-btn {
            background: #ffffff;
            color: #3b82f6;
            border: 2px solid #3b82f6;
            
            &:hover {
              background: #f1f5f9;
              transform: translateY(-2px);
            }
            
            &:active {
              transform: translateY(0);
            }
          }
        }
      }
    }
    
    .empty-state {
      text-align: center;
      padding: 80px 0;
      color: #64748b;
      font-size: 18px;
    }
    
    .back-btn-container {
      text-align: center;
      padding: 20px 0;
      
      .back-btn {
        min-width: 160px;
        height: 44px;
        background: #3b82f6;
        color: #ffffff;
        border: none;
        border-radius: 12px;
        font-size: 14px;
        font-weight: 500;
        cursor: pointer;
        box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        
        &:hover {
          background: #2563eb;
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
  .all-communities {
    padding: 10px;
    margin-top: 48px;
    
    .container {
      .header {
        padding: 20px 0;
        
        .title {
          font-size: 24px;
        }
        
        .subtitle {
          font-size: 14px;
        }
      }
      
      .communities-grid {
        grid-template-columns: 1fr;
        gap: 15px;
        
        .community-card {
          padding: 16px;
          
          .card-header {
            margin-bottom: 15px;
            
            .avatar {
              width: 40px;
              height: 40px;
              margin-right: 12px;
            }
            
            .info {
              .community-name {
                font-size: 16px;
              }
              
              .community-id {
                font-size: 11px;
              }
            }
          }
          
          .card-footer {
            flex-direction: column;
            gap: 8px;
            
            .join-btn, .view-btn {
              width: 100%;
            }
          }
        }
      }
    }
  }
}
</style>
