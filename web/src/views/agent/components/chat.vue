<template>
  <div class="full-content flex">
    <el-main class="scroll">
      <div class="smart-center" style="padding: 0">
        <!--开场白设置-->
        <div v-show="echo" class="session rl echo">
          <streamGreetingField
            :editForm="editForm"
            @setProloguePrompt="setProloguePrompt"
          />
        </div>
        <!--对话-->
        <div v-show="!echo" class="center-session">
          <streamMessageField
            ref="session-com"
            class="component"
            :chatType="'agent'"
            :sessionStatus="sessionStatus"
            :recommendConfig="recommendConfig"
            @clearHistory="clearHistory"
            @refresh="refresh"
            @queryCopy="queryCopy"
            @handleRecommendClick="handleRecommendClick"
            :defaultUrl="editForm.avatar.path"
          />
        </div>
        <!--停止生成-重新生成-->
        <div class="center-editable">
          <div v-show="stopBtShow" class="stop-box">
            <span v-show="sessionStatus === 0" class="stop" @click="preStop">
              <img
                class="stop-icon mdl"
                :src="require('@/assets/imgs/stop.png')"
              />
              <span class="mdl">{{ $t('agent.stop') }}</span>
            </span>
            <span v-show="sessionStatus !== 0" class="stop" @click="refresh">
              <img
                class="stop-icon mdl"
                :src="require('@/assets/imgs/refresh.png')"
              />
              <span class="mdl">{{ $t('agent.refresh') }}</span>
            </span>
          </div>
          <!-- 输入框 -->
          <streamInputField
            ref="editable"
            source="perfectReminder"
            :fileTypeArr="fileTypeArr"
            :type="type"
            @preSend="preSend"
            @setSessionStatus="setSessionStatus"
          />
          <!-- 版权信息 -->
          <div v-if="appUrlInfo" class="appUrlInfo">
            <span v-if="appUrlInfo.copyrightEnable">
              {{ $t('app.copyright') }}: {{ appUrlInfo.copyright }}
            </span>
            <span v-if="appUrlInfo.privacyPolicyEnable">
              {{ $t('app.privacyPolicy') }}:
              <a
                :href="appUrlInfo.privacyPolicy"
                target="_blank"
                style="color: var(--color)"
              >
                {{ appUrlInfo.privacyPolicy }}
              </a>
            </span>
            <span v-if="appUrlInfo.disclaimerEnable">
              {{ $t('app.disclaimer') }}: {{ appUrlInfo.disclaimer }}
            </span>
          </div>
        </div>
      </div>
    </el-main>
  </div>
</template>

<script>
// import SessionComponentSe from './SessionComponentSe';
// import EditableDivV3 from './EditableDivV3';
// import Prologue from './Prologue';
import streamMessageField from '@/components/stream/streamMessageField';
import streamInputField from '@/components/stream/streamInputField';
import streamGreetingField from '@/components/stream/streamGreetingField';
import { parseSub, convertLatexSyntax } from '@/utils/util.js';
import {
  delConversation,
  createConversation,
  getConversationHistory,
  delOpenurlConversation,
  openurlConversation,
  OpenurlConverHistory,
} from '@/api/agent';
import sseMethod from '@/mixins/sseMethod';
import { md } from '@/mixins/markdown-it';
import { mapGetters, mapState } from 'vuex';
import { getRecommendQuestionUrl } from '@/api/agent';
import { fetchEventSource } from '@microsoft/fetch-event-source';

export default {
  inject: {
    getHeaderConfig: {
      default: () => null,
    },
  },
  props: {
    editForm: {
      type: Object,
      default: null,
    },
    chatType: {
      type: String,
      default: '',
    },
    type: {
      type: String,
      default: 'agentChat',
    },
    appUrlInfo: {
      type: Object,
      default: null,
    },
  },
  components: {
    // SessionComponentSe,
    // EditableDivV3,
    streamMessageField,
    streamInputField,
    streamGreetingField,
    // Prologue,
  },
  mixins: [sseMethod],
  computed: {
    ...mapGetters('app', ['sessionStatus']),
    ...mapGetters('menu', ['basicInfo']),
    ...mapGetters('user', ['commonInfo']),
    ...mapState('user', ['userInfo']),
  },
  data() {
    return {
      echo: true,
      fileTypeArr: ['doc/*', 'image/*'],
      hasDrawer: false,
      drawer: true,
      fileId: [],
      recommendConfig: {
        reqController: new AbortController(),
        list: [],
        loading: false,
      },
      recommendTimer: null,
    };
  },
  methods: {
    createConversion() {
      if (this.echo) {
        this.$message({
          type: 'info',
          message: this.$t('app.switchSession'),
          customClass: 'dark-message',
          iconClass: 'none',
          duration: 1500,
        });
        return;
      }
      this.conversationId = '';
      this.echo = true;
      this.clearPageHistory();
      this.$emit('setHistoryStatus');
    },
    //切换对话
    conversionClick(n) {
      if (this.sessionStatus === 0) {
        return;
      } else {
        this.stopBtShow = false;
      }

      this.$emit('setHistoryStatus');
      this.amswerNum = 0;
      n.active = true;
      this.clearPageHistory();
      this.echo = false;
      this.conversationId = n.conversationId;
      this.getConversationDetail(this.conversationId, true);
    },
    async getConversationDetail(id, loading) {
      loading && this.$refs['session-com'].doLoading();
      let res = null;
      if (this.type === 'agentChat') {
        res = await getConversationHistory({
          conversationId: id,
          pageSize: 1000,
          pageNo: 1,
        });
      } else {
        const config = this.getHeaderConfig();
        res = await OpenurlConverHistory(
          { conversationId: id },
          this.editForm.assistantId,
          config,
        );
      }

      if (res.code === 0) {
        let history = res.data.list
          ? res.data.list.map((n, index) => {
              return {
                ...n,
                query: n.prompt,
                finish: 1, //兼容流式问答
                response: md.render(
                  parseSub(convertLatexSyntax(n.response), index),
                ),
                oriResponse: n.response,
                searchList: n.searchList || [],
                fileList: n.requestFiles,
                gen_file_url_list: n.responseFileUrls || [],
                isOpen: true,
                toolText: this.$t('agent.tooled'),
                thinkText: this.$t('agent.thinked'),
                showScrollBtn: null,
              };
            })
          : [];
        this.$refs['session-com'].replaceHistory(history);
        this.$nextTick(() => {
          this.addCopyClick();
        });
      }
    },
    //删除对话
    async preDelConversation(n) {
      if (this.sessionStatus === 0) {
        return;
      }
      let res = null;
      if (this.type === 'agentChat') {
        res = await delConversation({ conversationId: n.conversationId });
      } else {
        const config = this.getHeaderConfig();
        res = await delOpenurlConversation(
          { conversationId: n.conversationId },
          this.editForm.assistantId,
          config,
        );
      }

      if (res.code === 0) {
        this.$emit('reloadList');
        if (this.conversationId === n.conversationId) {
          this.conversationId = '';
          this.$refs['session-com'].clearData();
        }
        this.echo = true;
      }
    },
    /*------会话------*/
    async preSend(val, fileList, fileInfo) {
      if (this.recommendTimer) {
        clearInterval(this.recommendTimer);
        this.recommendTimer = null;
      }
      if (this.recommendConfig.loading) {
        this.recommendConfig.reqController.abort();
        this.recommendConfig.reqController = new AbortController();
      }
      this.recommendConfig.list = [];
      this.recommendConfig.loading = false;
      this.inputVal = val || this.$refs['editable'].getPrompt();
      this.fileId = fileInfo || [];
      this.isTestChat = this.chatType === 'test' ? true : false;
      this.fileList = fileList || this.$refs['editable'].getFileList();
      if (!this.inputVal) {
        this.$message.warning(this.$t('agent.inputContent'));
        return;
      }
      if (!this.verifiyFormParams()) {
        return;
      }
      //如果是新会话，先创建
      if (!this.conversationId && this.chatType === 'chat') {
        let res = null;
        if (this.type === 'agentChat') {
          res = await createConversation({
            prompt: this.inputVal,
            assistantId: this.editForm.assistantId,
          });
        } else {
          const config = this.getHeaderConfig();
          res = await openurlConversation(
            { prompt: this.inputVal },
            this.editForm.assistantId,
            config,
          );
        }

        if (res.code === 0) {
          this.conversationId = res.data.conversationId;
          this.$emit('reloadList', true);
          this.setParams();
        }
      } else {
        this.setParams();
      }
    },
    verifiyFormParams() {
      if (this.chatType === 'chat') return true;
      const { matchType, priorityMatch, rerankModelId } =
        this.editForm.knowledgeBaseConfig.config;
      const isMixPriorityMatch = matchType === 'mix' && priorityMatch;
      const conditions = [
        {
          check: !this.editForm.modelParams,
          message: this.$t('agent.form.selectModel'),
        },
        {
          check: !isMixPriorityMatch && !rerankModelId,
          message: this.$t('knowledgeManage.hitTest.selectRerankModel'),
        },
        {
          check: !this.editForm.prologue,
          message: this.$t('agent.form.inputPrologue'),
        },
      ];
      for (const condition of conditions) {
        if (condition.check) {
          this.$message.warning(condition.message);
          return false;
        }
      }
      return true;
    },
    setParams() {
      const fileInfo = this.$refs['editable'].getFileIdList();
      let fileId = !fileInfo.length ? this.fileId : fileInfo;
      // this.useSearch = this.$refs['editable'].sendUseSearch();
      this.setSseParams({
        conversationId: this.conversationId,
        fileInfo: fileId,
        assistantId: this.editForm.assistantId,
      });
      this.doSend();
      this.echo = false;
    },
    /*--右侧提示词--*/
    showDrawer() {
      this.drawer = true;
    },
    hideDrawer() {
      this.drawer = false;
    },
    async getReminderList(cb) {
      let res = await getTemplateList({ pageNo: 0, pageSize: 0, title: '' });
      if (res.code === 0) {
        this.reminderList = res.data.list || [];
        cb && cb();
      }
    },
    reminderClick(n) {
      this.$refs['editable'].setPrompt(n.prompt);
    },
    // 打印结束回调
    onMainPrintEnd() {
      const history = this.$refs['session-com'].getSessionData().history;
      const lastMessage = history[history.length - 1];

      // 只有当最后一条消息存在且 finish 状态为 1 (真正结束) 时才触发推荐
      if (
        lastMessage &&
        lastMessage.finish === 1 &&
        this.editForm.recommendConfig &&
        this.editForm.recommendConfig.recommendEnable &&
        this.editForm.recommendConfig.modelConfig.modelId
      ) {
        this.recommendConfig.list = [];
        this.getRecommendQuestion();
      }
    },
    handleRecommendClick(val) {
      this.preSend(val);
    },
    // 请求推荐问题
    getRecommendQuestion() {
      const history = this.$refs['session-com'].getSessionData().history;
      const lastUserMessage = history
        .slice()
        .reverse()
        .find(item => item.query);
      const query = lastUserMessage ? lastUserMessage.query : '';
      const signal = this.recommendConfig.reqController.signal;

      class RetriableError extends Error {}
      class FatalError extends Error {}

      const params = {
        query: query,
        assistantId: this.editForm.assistantId,
        conversationId: this.conversationId,
        trial: this.chatType === 'test' ? true : false,
      };

      this.recommendConfig.loading = true;

      let currentBuffer = ''; // 用于暂存当前正在拼接的问题片段
      let baseList = []; // 用于存储已经确认完成的问题
      let contentQueue = []; // 字符队列，用于模拟打字机效果
      let isFinished = false; // 标记 SSE 是否已结束接收

      if (this.recommendTimer) {
        clearInterval(this.recommendTimer);
        this.recommendTimer = null;
      }

      // 核心处理逻辑：从队列中取字符并更新 UI
      const processQueue = () => {
        if (contentQueue.length > 0) {
          const item = contentQueue.shift();
          const { char, type } = item;
          currentBuffer += char;
          const delimiter = currentBuffer.includes('\\n')
            ? '\\n'
            : currentBuffer.includes('\n')
              ? '\n'
              : null;

          if (delimiter) {
            // 使用分隔符拆分内容
            const parts = currentBuffer.split(delimiter);
            // 除了最后一部分外，前面的部分都是已经接收完整的
            for (let i = 0; i < parts.length - 1; i++) {
              const finishedContent = parts[i].trim();
              if (finishedContent) {
                baseList.push({
                  content: finishedContent,
                  type: type,
                });
              }
            }
            // 将最后一部分（可能还不完整）留回缓冲区
            currentBuffer = parts[parts.length - 1];
          }

          // 实时渲染展示列表（已完成列表 + 当前正在输入的问题）
          const displayList = [...baseList];
          if (currentBuffer.trim()) {
            displayList.push({
              content: currentBuffer.trim(),
              type: type,
            });
          }
          this.recommendConfig.list = displayList;
        } else if (isFinished) {
          // 如果数据接收完毕且队列已空，执行最后收尾
          clearInterval(this.recommendTimer);
          this.recommendTimer = null;

          // 处理缓冲区剩余的内容
          const finalContent = currentBuffer.trim();
          if (finalContent) {
            // 获取最后一个元素的类型，如果没有则默认为 answer
            const lastType =
              this.recommendConfig.list.length > 0
                ? this.recommendConfig.list[
                    this.recommendConfig.list.length - 1
                  ].type
                : 'answer';

            baseList.push({
              content: finalContent,
              type: lastType,
            });
          }

          this.recommendConfig.list = [...baseList];
          this.recommendConfig.loading = false;
          currentBuffer = '';
        }
      };

      const _this = this;

      fetchEventSource(`${getRecommendQuestionUrl}`, {
        method: 'POST',
        signal,
        openWhenHidden: true,
        headers: {
          'Content-Type': 'text/event-stream; charset=utf-8',
          Authorization: 'Bearer ' + this.token,
          'x-user-id': this.userInfo.uid,
          'x-org-id': this.userInfo.orgId,
        },
        body: JSON.stringify(params),
        async onopen(response) {
          if (
            response.ok &&
            response.headers.get('content-type').includes('text/event-stream')
          ) {
            console.log('连接成功，开始获取数据...');
          } else if (
            response.status >= 400 &&
            response.status < 500 &&
            response.status !== 429
          ) {
            _this.recommendConfig.loading = false;
            throw new FatalError();
          } else {
            throw new RetriableError();
          }
        },

        onmessage: msgData => {
          if (msgData.data) {
            try {
              const _data = JSON.parse(msgData.data);
              const choice = _data.choices && _data.choices[0];
              if (choice) {
                const content = choice.delta && choice.delta.content;
                const contentType = choice.contentType || 'answer';

                if (content) {
                  // 将内容拆分为带类型信息的字符对象存入队列
                  const items = content.split('').map(char => ({
                    char,
                    type: contentType,
                  }));
                  contentQueue.push(...items);

                  if (!this.recommendTimer) {
                    this.recommendTimer = setInterval(processQueue, 30);
                  }
                }

                if (['stop', 'accidentStop'].includes(choice.finish_reason)) {
                  isFinished = true;
                  if (!this.recommendTimer) {
                    processQueue();
                  }
                }
              }
            } catch (e) {
              console.error('解析推荐问题失败', e);
            }
          }
          if (msgData.event === 'FatalError') {
            isFinished = true;
            throw new FatalError(msgData.data);
          }
        },
        async onclose() {
          console.log('连接关闭...');
          isFinished = true;
          if (!_this.recommendTimer) {
            processQueue();
          }
          return false;
        },
        onerror(event) {
          console.log('连接错误:', event);
          isFinished = true;
          _this.recommendConfig.loading = false;
          throw event;
        },
      });
    },
  },
  beforeDestroy() {
    if (this.recommendTimer) {
      clearInterval(this.recommendTimer);
      this.recommendTimer = null;
    }
  },
};
</script>

<style lang="scss" scoped>
@import '@/style/chat.scss';
.appUrlInfo {
  margin-top: 10px;
  display: flex;
  justify-content: center;
  span {
    cursor: pointer;
    color: #bbb;
    margin-right: 15px;
  }
}
</style>
