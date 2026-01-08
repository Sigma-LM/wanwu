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
            @clearHistory="clearHistory"
            @refresh="refresh"
            @queryCopy="queryCopy"
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
import { md } from '@/mixins/marksown-it';
import { mapGetters } from 'vuex';
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
  },
  data() {
    return {
      echo: true,
      fileTypeArr: ['doc/*', 'image/*'],
      hasDrawer: false,
      drawer: true,
      fileId: [],
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
      const conditions = [
        {
          check: !this.editForm.modelParams,
          message: this.$t('agent.form.selectModel'),
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
