<template>
  <div class="ai-assistant-container">
    <!-- 隐藏的检测图片（用于检测服务是否运行） -->
    <img
      v-if="checking"
      ref="checkImage"
      :src="imageUrl"
      style="display: none"
      @load="onImageLoad"
      @error="onImageError"
    />

    <!-- 聊天界面（只在服务确认可用后显示） -->
    <iframe
      v-if="showChat"
      ref="chatIframe"
      :src="chatFrameUrl"
      class="chat-iframe"
      frameborder="0"
      allow="clipboard-write"
      @load="onChatLoad"
    ></iframe>

    <!-- 服务未启动：显示简化的静态页面 -->
    <div v-if="showError" class="error-page">
      <div class="error-content">
        <div class="error-icon">🤖</div>
        <div class="error-title">
          {{ $t('aiAssistant.serviceUnavailable') }}
        </div>
        <div class="error-message">{{ $t('aiAssistant.refreshMessage') }}</div>
      </div>
    </div>

    <!-- 检查中：显示加载状态 -->
    <div v-if="showLoading" class="loading-page">
      <el-spinner type="dots" :size="40"></el-spinner>
      <div class="loading-text">{{ $t('aiAssistant.connecting') }}</div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AIAssistant',
  data() {
    return {
      chatFrameUrl: 'http://localhost:8585',
      imageUrl: 'http://localhost:8585/user-avatar.png', // 使用public文件夹下的图片
      isChatLoaded: false,
      loadingTimer: null,
      checkTimer: null,
      styleElement: null,
      // 状态控制
      checking: false, // 是否正在检测
      showLoading: true, // 显示加载中
      showError: false, // 显示错误
      showChat: false, // 显示聊天界面
      checkStartTime: null, // 检测开始时间
    };
  },
  computed: {
    orgInfo() {
      return this.$store.state.user.orgInfo || {};
    },
  },
  mounted() {
    console.log('[AI Assistant] 组件已挂载');

    this.$nextTick(() => {
      this.disablePageScroll();
      this.startCheck(); // 开始检测服务
    });

    // 监听路由变化
    this.$watch(
      () => this.$route.path,
      () => {
        this.startCheck();
      },
    );

    // 监听组织变化
    this.$watch(
      () => this.$store.state.user.userInfo.orgId,
      (newVal, oldVal) => {
        if (newVal !== oldVal && this.showChat) {
          this.sendContextToIframe();
        }
      },
    );
  },
  beforeDestroy() {
    if (this.loadingTimer) {
      clearTimeout(this.loadingTimer);
    }
    if (this.checkTimer) {
      clearTimeout(this.checkTimer);
    }
    this.enablePageScroll();
    window.removeEventListener('message', this.handleIframeMessage);
  },
  methods: {
    // ========================================
    // 开始检测服务（使用图片检测）
    // ========================================
    async startCheck() {
      console.log('[AI Assistant] 🏥 开始检测服务...');
      console.log('[AI Assistant] chatFrameUrl:', this.chatFrameUrl);
      console.log('[AI Assistant] imageUrl:', this.imageUrl);

      // 重置状态
      this.checking = false;
      this.showLoading = true;
      this.showError = false;
      this.showChat = false;
      this.isChatLoaded = false;
      this.checkStartTime = null;

      // 清除之前的定时器
      if (this.checkTimer) {
        clearTimeout(this.checkTimer);
      }

      // 等待一帧，确保之前的img已移除
      await this.$nextTick();

      // 开始检测
      this.checking = true;
      this.checkStartTime = Date.now();
      console.log('[AI Assistant] 创建检测图片:', this.imageUrl);

      // 设置3秒超时
      this.checkTimer = setTimeout(() => {
        const elapsed = Date.now() - this.checkStartTime;
        console.log('[AI Assistant] ⏱️ 检测超时（3秒），耗时', elapsed, 'ms');

        this.checking = false; // 移除图片
        this.showLoading = false;
        this.showError = true;
        this.showChat = false;
        console.log('[AI Assistant] ❌ 服务不可用（超时）');
      }, 3000);
    },

    // 图片加载成功 → 服务可用
    onImageLoad(event) {
      const elapsed = Date.now() - this.checkStartTime;
      console.log('[AI Assistant] ✅ 图片加载成功，耗时', elapsed, 'ms');
      console.log(
        '[AI Assistant] 图片尺寸:',
        event.target.naturalWidth,
        'x',
        event.target.naturalHeight,
      );

      clearTimeout(this.checkTimer);
      this.checking = false; // 立即移除图片

      // 显示聊天界面
      this.showLoading = false;
      this.showError = false;
      this.showChat = true;

      console.log('[AI Assistant] showChat:', this.showChat);

      this.$nextTick(() => {
        setTimeout(() => this.sendContextToIframe(), 1000);
      });
    },

    // 图片加载失败 → 服务不可用
    onImageError(event) {
      const elapsed = Date.now() - this.checkStartTime;
      console.error('[AI Assistant] ❌ 图片加载失败，耗时', elapsed, 'ms');
      console.error('[AI Assistant] 错误信息:', event);

      clearTimeout(this.checkTimer);
      this.checking = false; // 立即移除图片

      this.showLoading = false;
      this.showError = true;
      this.showChat = false;
      console.log('[AI Assistant] ❌ 服务不可用（图片加载失败）');
    },

    // ========================================
    // 聊天iframe加载
    // ========================================
    onChatLoad() {
      console.log('[AI Assistant] ✅ 聊天iframe加载完成');
      clearTimeout(this.loadingTimer);
      this.isChatLoaded = true;

      setTimeout(() => this.sendContextToIframe(), 500);
    },

    // ========================================
    // 原有方法：禁用/恢复页面滚动
    // ========================================
    disablePageScroll() {
      const elMain = document.querySelector('.el-main');
      if (elMain) {
        elMain.style.overflow = 'hidden';
      }

      const rightPageContent = document.querySelector('.right-page-content');
      if (rightPageContent) {
        rightPageContent.style.overflow = 'hidden';
        rightPageContent.style.padding = '0';
        rightPageContent.style.height = '100%';
      }

      if (!this.styleElement) {
        this.styleElement = document.createElement('style');
        this.styleElement.id = 'ai-assistant-no-scroll';
        this.styleElement.innerHTML = `
          .el-main { overflow: hidden !important; }
          .right-page-content { overflow: hidden !important; padding: 0 !important; height: 100% !important; }
        `;
        document.head.appendChild(this.styleElement);
      }
    },

    enablePageScroll() {
      const elMain = document.querySelector('.el-main');
      if (elMain) {
        elMain.style.overflow = '';
      }

      const rightPageContent = document.querySelector('.right-page-content');
      if (rightPageContent) {
        rightPageContent.style.overflow = '';
        rightPageContent.style.padding = '';
        rightPageContent.style.height = '';
      }

      if (this.styleElement) {
        document.head.removeChild(this.styleElement);
        this.styleElement = null;
      }
    },

    // ========================================
    // 发送上下文到iframe
    // ========================================
    sendContextToIframe() {
      if (!this.showChat) {
        console.log('[AI Assistant] 跳过发送上下文（聊天界面未显示）');
        return;
      }

      const access_cert_str = localStorage.getItem('access_cert');
      let userInfo = {};
      let token = '';
      let orgId = '';

      if (access_cert_str) {
        try {
          const access_cert = JSON.parse(access_cert_str);
          userInfo = access_cert.user || {};
          token = userInfo.token || '';
          orgId = userInfo.orgId || '';

          if (
            !orgId &&
            this.orgInfo &&
            this.orgInfo.orgs &&
            this.orgInfo.orgs.length > 0
          ) {
            orgId = this.orgInfo.orgs[0].id;
          }
        } catch (e) {
          console.error('[AI Assistant] 解析失败:', e);
        }
      }

      const locale = localStorage.getItem('locale') || 'zh-CN';
      const orgName = this.getOrgName(orgId);
      let wanwuApiUrl = window.API_API_ORIGIN || '';

      if (!wanwuApiUrl) {
        wanwuApiUrl = window.location.origin.replace(':8080', ':8081');
      }

      const contextInfo = {
        type: 'INIT_CONTEXT',
        payload: {
          token: token,
          userId: userInfo.id,
          userName: userInfo.userName,
          orgId: orgId,
          orgName: orgName,
          locale: locale,
          wanwuApiUrl: wanwuApiUrl,
          timestamp: Date.now(),
        },
      };

      if (this.$refs.chatIframe && this.$refs.chatIframe.contentWindow) {
        this.$refs.chatIframe.contentWindow.postMessage(contextInfo, '*');
        window.addEventListener('message', this.handleIframeMessage);

        setTimeout(() => {
          if (!this.isChatLoaded) {
            this.sendContextToIframe();
          }
        }, 1000);
      }
    },

    handleIframeMessage(event) {
      const { type } = event.data;
      if (type === 'CLAWCHAT_CONTEXT_ACK') {
        console.log('[AI Assistant] ✅ 上下文已确认');
      } else if (type === 'CLAWCHAT_READY') {
        this.sendContextToIframe();
      }
    },

    getOrgName(orgId) {
      if (this.orgInfo && this.orgInfo.orgs) {
        const org = this.orgInfo.orgs.find(o => o.id === orgId);
        return org ? org.name : '';
      }
      return '';
    },
  },
};
</script>

<style lang="scss" scoped>
.ai-assistant-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: #f6f7fa;
}

.chat-iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border: none;
  display: block;
  overflow: hidden;
}

/* ==================== 错误页面 ==================== */

.error-page {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background: #f6f7fa;
}

.error-content {
  text-align: center;
}

.error-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.error-title {
  font-size: 18px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.error-message {
  font-size: 14px;
  color: #999;
}

/* ==================== 加载中页面 ==================== */

.loading-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background: #f6f7fa;
  gap: 16px;
}

.loading-text {
  font-size: 14px;
  color: #666;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>

<style lang="scss">
/* 全局样式：隐藏滚动条视觉效果（保持滚动功能） */
.el-main::-webkit-scrollbar,
.right-page-content::-webkit-scrollbar {
  display: none;
}
</style>
