<template>
  <div class="content">
    <div class="left">
      <!-- <h4 class="c-l-title">热门帖子</h4> -->
      <div class="c-l-header">
        <div class="new btn-iconfont"
        :class="{ active: timeOrder }"
        @click="selectOrder('time')"
        >
          <i class="iconfont icon-polygonred"></i>New
        </div>
        <div class="top btn-iconfont"
         :class="{ active: scoreOrder }"
         @click="selectOrder('score')"
        >
          <i class="iconfont icon-top"></i>Score
        </div>
        <button class="btn-publish" @click="goPublish">发表</button>
      </div>
      <ul class="c-l-list">
        <li class="c-l-item"  v-for="post in postList" :key="post.id">
          <div class="post">
            <a class="vote">
              <span class="iconfont icon-up"
              @click.stop="vote(post.id, '1')"
              ></span>
            </a>
            <span class="text">{{post.vote_num}}</span>
            <a class="vote">
              <span class="iconfont icon-down"
              @click.stop="vote(post.id, '-1')"
              ></span>
            </a>
          </div>
          <div class="l-container" @click="goDetail(post.id)">
            <h4 class="con-title">{{post.title}}</h4>
            <div class="con-memo">
              <p>{{post.content}}</p>
            </div>
            <!-- <div class="user-btn">
              <span class="btn-item">
                <i class="iconfont icon-comment"></i>
                <span>{{post.comments}} comments</span>
              </span>
            </div> -->
          </div>
        </li>
      </ul>
    </div>
    <div class="right">
      <div class="communities">
        <h2 class="r-c-title">今日火热频道排行榜</h2>
        <ul class="r-c-content">
          <li class="r-c-item" v-for="(rank, index) in topCommunities" :key="rank.id" @click="goToCommunity(rank.id)">
            <span class="index">{{ index + 1 }}</span>
            <i class="icon"></i>
            {{ rank.name }}
          </li>
        </ul>
        <button class="view-all" @click="viewAllCommunities">查看所有</button>
      </div>
      <div class="r-trending">
        <h2 class="r-t-title">持续热门频道</h2>
        <ul class="rank">
          <li class="r-t-cell" v-for="community in trendingCommunities" :key="community.id">
            <div class="r-t-cell-info" @click="goToCommunity(community.id)">
              <div class="avatar"></div>
              <div class="info">
                <span class="info-title">{{ community.name }}</span>
                <p class="info-num">{{ community.members }} members</p>
              </div>
            </div>
            <button 
              class="join-btn" 
              :class="{ 'joined': community.joined }"
              @click.stop="joinCommunity(community.id)"
              :disabled="community.joined"
            >
              {{ community.joined ? '已加入' : 'JOIN' }}
            </button>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
// @ is an alias to /src

export default {
  name: "Home",
  components: {},
  data() {
    return {
      order: "time",
      page: 1,
      postList: [],
      topCommunities: [],
      trendingCommunities: [],
      userJoinedCommunityIds: [] // 用户已加入的社区ID列表
    };
  },
  methods: {
    selectOrder(order){
      this.order = order;
      this.getPostList()
    },
    goPublish(){
      this.$router.push({ name: "Publish" });
    },
    goDetail(id){
      this.$router.push({ name: "Content", params: { id: id }});
    },
    goToCommunity(community_id) {
      // 跳转到社区详情页
      this.$router.push({ name: "CommunityDetail", params: { id: community_id } });
    },
    getPostList() {
      this.$axios({
        method: "get",
        url: "/posts2",
        params: {
          page: this.page,
          order: this.order,
        }
      })
        .then(response => {
          console.log(response.data, 222);
          if (response.code == 1000) {
            this.postList = response.data;
          } else {
            console.log(response.msg);
          }
        })
        .catch(error => {
          console.log(error);
        });
    },
    vote(post_id, direction){
      this.$axios({
        method: "post",
        url: "/vote",
        data: JSON.stringify({
          post_id: String(post_id),
          direction: direction,
        })
      })
        .then(response => {
          if (response.code == 1000) {
            console.log("vote success");
          } else {
            console.log(response.msg);
          }
        })
        .catch(error => {
          console.log(error);
        });
    },
    // 查看所有社区
    viewAllCommunities() {
      this.$router.push({ name: "AllCommunities" });
    },
    // 加入社区
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
            // 更新社区状态
            const community = this.trendingCommunities.find(c => c.id === community_id);
            if (community) {
              community.joined = true;
            }
            this.$message?.success('成功加入社区！');
            console.log("成功加入社区");
          } else {
            this.$message?.error(response.msg || '加入社区失败');
            console.log(response.msg);
          }
        })
        .catch(error => {
          console.log("加入社区失败:", error);
          this.$message?.error('网络错误，请稍后重试');
        });
    },
    // 获取社区列表（包含成员数量）
    getCommunityList() {
      this.$axios({
        method: "get",
        url: "/community/list"
      })
        .then(response => {
          console.log("社区列表:", response);
          if (response.code == 1000 && response.data && response.data.length > 0) {
            const communities = response.data;
            
            // 设置今日火热频道排行榜（前3个，按成员数排序）
            this.topCommunities = communities.slice(0, 3);
            
            // 设置持续热门频道（前5个社区，避免显示太多）
            this.trendingCommunities = communities.slice(0, 5).map(community => ({
              id: community.id,
              name: community.name,
              members: this.formatMemberCount(community.member_count),
              joined: this.userJoinedCommunityIds.includes(community.id)
            }));
          } else {
            console.log("社区数据为空或格式不正确:", response);
          }
        })
        .catch(error => {
          console.log("获取社区列表失败:", error);
        });
    },
    // 格式化成员数量显示
    formatMemberCount(count) {
      if (count >= 1000) {
        return (count / 1000).toFixed(1) + 'k';
      }
      return count.toString();
    },
    // 获取热门社区列表
    async getTrendingCommunities() {
      // 先尝试获取用户已加入的社区（可能会失败，比如未登录）
      try {
        await this.getUserCommunities();
      } catch (error) {
        console.log("跳过获取用户社区（可能未登录）");
      }
      
      // 无论用户是否登录，都要获取社区列表
      this.getCommunityList();
    },
    // 获取用户已加入的社区
    getUserCommunities() {
      return this.$axios({
        method: "get",
        url: "/community/user"
      })
        .then(response => {
          console.log("用户已加入的社区:", response);
          if (response.code == 1000 && response.data) {
            // 保存用户已加入的社区ID列表
            this.userJoinedCommunityIds = response.data.map(c => c.id);
          } else {
            this.userJoinedCommunityIds = [];
          }
        })
        .catch(error => {
          console.log("获取用户社区失败（可能未登录）:", error);
          // 如果获取失败（比如未登录），使用空数组，但抛出错误让外层catch
          this.userJoinedCommunityIds = [];
          throw error;
        });
    }
  },
  mounted: function() {
    this.getPostList();
    this.getTrendingCommunities();
  },
  computed:{
    timeOrder(){
      return this.order == "time";
    },
    scoreOrder(){
      return this.order == "score";
    }
  }
};
</script>

<style scoped lang="less">
.content {
  max-width: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: row;
  justify-content: center;
  margin: 56px auto 0;
  padding: 20px 24px;
  background: #f5f7fa;
  min-height: calc(100vh - 56px);
  
  .left {
    width: 640px;
    padding-bottom: 10px;
    .c-l-title {
      font-size: 14px;
      font-weight: 500;
      line-height: 18px;
      color: #1a1a1b;
      text-transform: unset;
      padding-bottom: 10px;
    }
    .c-l-header {
      align-items: center;
      background: #ffffff;
      border: none;
      border-radius: 16px;
      box-sizing: border-box;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
      display: -ms-flexbox;
      display: flex;
      -ms-flex-flow: row nowrap;
      flex-flow: row nowrap;
      height: 56px;
      -ms-flex-pack: start;
      justify-content: flex-start;
      margin-bottom: 16px;
      padding: 0 16px;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08), 0 2px 4px rgba(0, 0, 0, 0.08);
        transform: translateY(-1px);
      }
      .iconfont {
        margin-right: 4px;
      }
      .btn-iconfont {
        display: flex;
        display: -webkit-flex;
        color: #64748b;
        cursor: pointer;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        padding: 6px 12px;
        border-radius: 12px;
        &:hover {
          color: #3b82f6;
          background: #f1f5f9;
          transform: translateY(-1px);
        }
      }
      .active {
        background: #3b82f6;
        color: #ffffff;
        fill: #ffffff;
        border-radius: 12px;
        height: 36px;
        line-height: 36px;
        margin-right: 8px;
        padding: 0 16px;
        box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
        font-weight: 500;
      }
      .new {
        font-size: 14px;
        margin-right: 12px;
      }
      .top {
        font-size: 14px;
      }
      .btn-publish {
        min-width: 72px;
        height: 36px;
        line-height: 36px;
        background: #3b82f6;
        color: #ffffff;
        border: none;
        border-radius: 12px;
        box-sizing: border-box;
        text-align: center;
        margin-left: auto;
        cursor: pointer;
        box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        font-weight: 500;
        padding: 0 16px;
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
      .sort {
        margin-left: 300px;
        display: flex;
        color: #0079d3;
        display: -webkit-flex;
        align-items: center;
        .sort-triangle {
          width: 0;
          height: 0;
          border-top: 5px solid #0079d3;
          border-right: 5px solid transparent;
          border-bottom: 5px solid transparent;
          border-left: 5px solid transparent;
          margin-top: 5px;
          margin-left: 10px;
        }
      }
    }
    .c-l-list {
      .c-l-item {
        list-style: none;
        border-radius: 16px;
        padding-left: 40px;
        cursor: pointer;
        border: none;
        margin-bottom: 12px;
        background: #ffffff;
        color: #64748b;
        position: relative;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        overflow: hidden;
        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08), 0 4px 8px rgba(0, 0, 0, 0.06);
        }
        &::after {
          content: '';
          position: absolute;
          top: 0;
          left: 0;
          width: 4px;
          height: 100%;
          background: #3b82f6;
          opacity: 0;
          transition: opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        }
        &:hover::after {
          opacity: 1;
        }
        .post {
          align-items: center;
          box-sizing: border-box;
          display: -ms-flexbox;
          display: flex;
          -ms-flex-direction: column;
          flex-direction: column;
          height: 100%;
          left: 0;
          padding: 8px 4px 8px 0;
          position: absolute;
          top: 0;
          width: 40px;
          border-left: 4px solid transparent;
          background: #f8fafc;
          border-radius: 16px 0 0 16px;
          z-index: 1;
          .iconfont {
            margin-right: 0;
          }
          .down {
            transform: scaleY(-1);
          }
          .text {
            color: #1a1a1b;
            font-size: 12px;
            font-weight: 700;
            line-height: 16px;
            pointer-events: none;
            word-break: normal;
          }
        }
        .l-container {
          padding: 15px;
          .con-title {
            color: #000000;
            font-size: 18px;
            font-weight: 500;
            line-height: 22px;
            text-decoration: none;
            word-break: break-word;
          }
          .con-memo {
            margin-top: 10px;
            margin-bottom: 10px;
          }
          .con-cover {
            height: 512px;
            width: 100%;
            background: url("https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1585999647247&di=7e9061211c23e3ed9f0c4375bb3822dc&imgtype=0&src=http%3A%2F%2Fi1.hdslb.com%2Fbfs%2Farchive%2F04d8cda08e170f4a58c18c45a93c539375c22162.jpg")
              no-repeat;
            background-size: cover;
            margin-top: 10px;
            margin-bottom: 10px;
          }
          .user-btn {
            font-size: 14px;
            display: flex;
            display: -webkit-flex;
            .btn-item {
              display: flex;
              display: -webkit-flex;
              margin-right: 10px;
              .iconfont {
                margin-right: 4px;
              }
            }
          }
        }
      }
    }
  }
  .right {
    width: 312px;
    margin-left: 24px;
    margin-top: 28px;
    .communities {
      background: #ffffff;
      color: #1e293b;
      border: none;
      border-radius: 16px;
      overflow: hidden;
      word-wrap: break-word;
      margin-bottom: 20px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      &:hover {
        box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08), 0 4px 8px rgba(0, 0, 0, 0.06);
        transform: translateY(-2px);
      }
      .r-c-title {
        background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
        height: 80px;
        width: 100%;
        color: #fff;
        font-size: 18px;
        line-height: 80px;
        padding-left: 10px;
        box-sizing: border-box;
        text-align: center;
        font-weight: 600;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
      }
      .r-c-content {
        .r-c-item {
          align-items: center;
          display: flex;
          display: -webkit-flex;
          height: 48px;
          padding: 0 12px;
          border-bottom: thin solid #edeff1;
          font-size: 14px;
          cursor: pointer;
          transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
          &:hover {
            background-color: #f8fafc;
            padding-left: 16px;
          }
          .index {
            width: 20px;
            color: #1c1c1c;
            font-size: 14px;
            font-weight: 500;
            line-height: 18px;
          }
          .icon {
            width: 32px;
            height: 32px;
            background-image: url("../assets/images/avatar.png");
            background-repeat: no-repeat;
            background-size: cover;
            margin-right: 20px;
          }
          &:last-child {
            border-bottom: none;
          }
        }
      }
      .view-all {
        background: #3b82f6;
        border: none;
        border-radius: 12px;
        box-sizing: border-box;
        text-align: center;
        letter-spacing: 0.5px;
        text-decoration: none;
        font-size: 13px;
        font-weight: 600;
        line-height: 20px;
        text-transform: none;
        padding: 10px 0;
        width: 280px;
        color: #fff;
        margin: 20px 0 20px 16px;
        box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
        cursor: pointer;
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
    .r-trending {
      padding-top: 16px;
      width: 312px;
      background: #ffffff;
      color: #1e293b;
      fill: #1e293b;
      border: none;
      border-radius: 16px;
      overflow: hidden;
      word-wrap: break-word;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      &:hover {
        box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08), 0 4px 8px rgba(0, 0, 0, 0.06);
        transform: translateY(-2px);
      }
      .r-t-title {
        font-size: 10px;
        font-weight: 700;
        letter-spacing: 0.5px;
        line-height: 12px;
        text-transform: uppercase;
        background-color: #ffffff;
        border-radius: 3px 3px 0 0;
        color: #1a1a1b;
        display: -ms-flexbox;
        display: flex;
        fill: #1a1a1b;
        padding: 0 12px 12px;
      }
      .rank {
        padding: 12px;
        .r-t-cell {
          display: flex;
          display: -webkit-flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 16px;
          .r-t-cell-info {
            display: flex;
            cursor: pointer;
            flex: 1;
            transition: all 0.2s ease;
            &:hover {
              opacity: 0.8;
            }
          }
          .avatar {
            width: 32px;
            height: 32px;
            background: url("../assets/images/avatar.png") no-repeat;
            background-size: cover;
            margin-right: 10px;
          }
          .info {
            margin-right: 10px;
            .info-title {
              font-size: 12px;
              font-weight: 500;
              line-height: 16px;
              text-overflow: ellipsis;
              width: 144px;
            }
            .info-num {
              font-size: 12px;
              font-weight: 400;
              line-height: 16px;
              padding-bottom: 4px;
            }
          }
          .join-btn {
            min-width: 80px;
            height: 34px;
            line-height: 34px;
            background: #3b82f6;
            color: #ffffff;
            border: none;
            border-radius: 10px;
            box-sizing: border-box;
            text-align: center;
            font-weight: 500;
            font-size: 13px;
            padding: 0 12px;
            box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
            cursor: pointer;
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
            &.joined {
              background: #e2e8f0;
              color: #64748b;
              box-shadow: none;
              cursor: default;
            }
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
      width: 100%;
      margin: 0;
      padding-bottom: 10px;
      .c-l-header {
        height: 48px;
        padding: 0 8px;
        margin-bottom: 12px;
        .new, .top {
          font-size: 13px;
          margin-right: 12px;
        }
        .active {
          padding: 0 8px;
        }
        .btn-publish {
          width: 56px;
          height: 28px;
          line-height: 28px;
          font-size: 12px;
        }
      }
      .c-l-list {
        .c-l-item {
          padding-left: 35px;
          margin-bottom: 8px;
          .post {
            width: 35px;
            padding: 6px 2px 6px 0;
            .text {
              font-size: 11px;
            }
          }
          .l-container {
            padding: 10px;
            .con-title {
              font-size: 16px;
              line-height: 20px;
            }
            .con-memo {
              margin-top: 8px;
              margin-bottom: 8px;
              font-size: 13px;
              p {
                display: -webkit-box;
                -webkit-line-clamp: 3;
                -webkit-box-orient: vertical;
                overflow: hidden;
                text-overflow: ellipsis;
              }
            }
          }
        }
      }
    }
    .right {
      width: 100%;
      margin-left: 0;
      margin-top: 12px;
      .communities {
        margin-bottom: 15px;
        .r-c-title {
          height: 60px;
          font-size: 18px;
          line-height: 80px;
          padding-left: 8px;
        }
        .r-c-content {
          .r-c-item {
            height: 40px;
            padding: 0 10px;
            font-size: 13px;
            .icon {
              width: 28px;
              height: 28px;
              margin-right: 15px;
            }
          }
        }
        .view-all {
          width: calc(100% - 20px);
          margin: 15px 10px 15px 10px;
        }
      }
      .r-trending {
        width: 100%;
        .r-t-title {
          font-size: 9px;
        }
        .rank {
          padding: 10px;
          .r-t-cell {
            margin-bottom: 12px;
            .avatar {
              width: 28px;
              height: 28px;
            }
            .info {
              .info-title {
                font-size: 11px;
                width: auto;
              }
              .info-num {
                font-size: 11px;
              }
            }
            .join-btn {
              width: 80px;
              height: 28px;
              line-height: 28px;
              font-size: 11px;
            }
          }
        }
      }
    }
  }
}
</style>
