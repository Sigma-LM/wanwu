<template>
  <div class="agent-from-content" :class="{ 'disable-clicks': disableClick }">
    <div class="form-header">
      <div class="header-left">
        <span class="el-icon-arrow-left btn" @click="goBack"></span>
        <div class="basicInfo">
          <div class="img">
            <img
              :src="
                editForm.avatar.path
                  ? `/user/api` + editForm.avatar.path
                  : '@/assets/imgs/bg-logo.png'
              "
            />
          </div>
          <div class="basicInfo-desc">
            <span class="basicInfo-title">
              {{
                (editForm.name || $t('agent.form.noInfo')).length > 12
                  ? (editForm.name || $t('agent.form.noInfo')).substring(
                      0,
                      12,
                    ) + '...'
                  : editForm.name || $t('agent.form.noInfo')
              }}
            </span>
            <span
              class="el-icon-edit-outline editIcon"
              @click="editAgent"
            ></span>
            <LinkIcon type="agent" />
            <p>{{ editForm.desc || $t('agent.form.noInfo') }}</p>
            <p>
              uuid: {{ this.editForm.uuid }}
              <copyIcon
                :text="this.editForm.uuid"
                :onlyIcon="true"
                size="mini"
              />
            </p>
          </div>
        </div>
      </div>
      <div class="header-right">
        <VersionPopover
          ref="versionPopover"
          v-if="publishType"
          style="pointer-events: auto"
          :appId="editForm.assistantId"
          :appType="AGENT"
          @reloadData="reloadData"
          @previewVersion="previewVersion"
        />
        <el-button
          v-if="publishType"
          size="small"
          type="primary"
          style="padding: 13px 12px"
          @click="handlePublishSet"
        >
          <span class="el-icon-setting"></span>
          {{ $t('agent.form.publishConfig') }}
        </el-button>
        <el-popover
          placement="bottom-end"
          trigger="click"
          style="margin-left: 13px"
        >
          <el-button
            slot="reference"
            size="small"
            type="primary"
            style="padding: 13px 12px"
          >
            {{ $t('common.button.publish') }}
            <span class="el-icon-arrow-down" style="margin-left: 5px"></span>
          </el-button>
          <el-form ref="publishForm" :model="publishForm" :rules="publishRules">
            <el-form-item :label="$t('list.version.no')" prop="version">
              <el-input
                v-model="publishForm.version"
                :placeholder="$t('list.version.noPlaceholder')"
              ></el-input>
            </el-form-item>
            <el-form-item :label="$t('list.version.desc')" prop="desc">
              <el-input
                v-model="publishForm.desc"
                :placeholder="$t('list.version.descPlaceholder')"
              ></el-input>
            </el-form-item>
            <el-form-item
              :label="$t('list.version.publishType')"
              prop="publishType"
            >
              <el-radio-group v-model="publishForm.publishType">
                <div>
                  <el-radio label="private">
                    {{ $t('agent.form.publishType') }}
                  </el-radio>
                </div>
                <div>
                  <el-radio label="organization">
                    {{ $t('agent.form.publishType1') }}
                  </el-radio>
                </div>
                <div>
                  <el-radio label="public">
                    {{ $t('agent.form.publishType2') }}
                  </el-radio>
                </div>
              </el-radio-group>
            </el-form-item>
            <div class="saveBtn">
              <el-button size="mini" type="primary" @click="savePublish">
                {{ $t('common.button.save') }}
              </el-button>
            </div>
          </el-form>
        </el-popover>
      </div>
    </div>
    <!-- 智能体配置 -->
    <div class="agent_form">
      <div class="block drawer-info">
        <div class="promptTitle">
          <h3>{{ $t('agent.form.systemPrompt') }}</h3>
          <div class="prompt-title-icon">
            <el-tooltip
              class="item"
              effect="dark"
              :content="$t('agent.form.submitToPrompt')"
              placement="top-start"
            >
              <span class="el-icon-folder-add" @click="handleShowPrompt"></span>
            </el-tooltip>
            <el-tooltip
              class="item"
              effect="dark"
              :content="$t('tempSquare.promptOptimize')"
              placement="top-start"
            >
              <span
                style="margin-left: 5px"
                class="el-icon-s-help"
                @click="showPromptOptimize"
              ></span>
            </el-tooltip>
            <el-tooltip
              class="item"
              effect="dark"
              :content="$t('tempSquare.promptCompare')"
              placement="top-start"
            >
              <span class="tool-icon" @click="showPromptCompare">
                <img :src="require('@/assets/imgs/temp-compare.png')" />
              </span>
            </el-tooltip>
          </div>
        </div>
        <div class="rl" style="padding-top: 10px">
          <el-input
            class="desc-input"
            v-model="editForm.instructions"
            :placeholder="$t('agent.form.promptTips')"
            type="textarea"
            show-word-limit
            :rows="12"
            @blur="handleInstructionsBlur"
          ></el-input>
        </div>
        <promptTemplate ref="promptTemplate" />
      </div>
      <div class="drawer-form">
        <div class="block">
          <h3 class="box labelTitle">{{ $t('agent.form.agentConfig') }}</h3>
          <div class="box">
            <p class="block-title common-set">
              <span class="common-set-label">
                <img
                  :src="require('@/assets/imgs/require.png')"
                  class="required-label"
                />
                {{ $t('agent.form.modelSelect') }}
              </span>
              <span class="common-add" @click="showModelSet">
                <el-tooltip
                  class="item"
                  effect="dark"
                  :content="$t('agent.form.modelSelectConfigTips')"
                  placement="top-start"
                >
                  <span class="el-icon-s-operation operation">
                    <span class="handleBtn">{{ $t('agent.form.config') }}</span>
                  </span>
                </el-tooltip>
              </span>
            </p>
            <div class="rl">
              <el-select
                v-model="editForm.modelParams"
                :placeholder="$t('agent.form.modelSearchPlaceholder')"
                @visible-change="visibleChange"
                :loading-text="$t('agent.toolDetail.modelLoadingText')"
                class="cover-input-icon model-select"
                :loading="modelLoading"
                filterable
                value-key="modelId"
                @change="handleModelChange($event)"
              >
                <el-option
                  class="model-option-item"
                  v-for="item in modelOptions"
                  :key="item.modelId"
                  :value="item.modelId"
                  :label="item.displayName"
                >
                  <div class="model-option-content">
                    <span class="model-name">{{ item.displayName }}</span>
                    <div
                      class="model-select-tags"
                      v-if="item.tags && item.tags.length > 0"
                    >
                      <span
                        v-for="(tag, tagIdx) in item.tags"
                        :key="tagIdx"
                        class="model-select-tag"
                      >
                        {{ tag.text }}
                      </span>
                    </div>
                  </div>
                </el-option>
              </el-select>
              <div
                class="model-select-tips"
                v-if="editForm.visionsupport === 'support'"
              >
                {{ $t('agent.form.visionModelTips') }}
              </div>
              <div
                class="model-select-tips"
                v-if="
                  editForm.functionCalling === 'noSupport' && editForm.newAgent
                "
              >
                {{ $t('agent.form.functionCallTips') }}
              </div>
            </div>
          </div>
          <div class="box">
            <p class="block-title">
              <img
                :src="require('@/assets/imgs/require.png')"
                class="required-label"
              />
              {{ $t('agent.form.prologue') }}
            </p>
            <div class="rl">
              <el-input
                class="desc-input"
                v-model="editForm.prologue"
                maxlength="100"
                :placeholder="$t('agent.form.prologuePlaceholder')"
                type="textarea"
              ></el-input>
              <span class="el-input__count">
                {{ editForm.prologue.length }}/100
              </span>
            </div>
          </div>
          <div class="recommend-box">
            <p class="block-title recommend-title">
              <span>{{ $t('agent.form.recommendQuestion') }}</span>
              <span @click="addRecommend" class="common-add">
                <span class="el-icon-plus"></span>
                <span class="handleBtn">{{ $t('agent.add') }}</span>
              </span>
            </p>
            <div
              class="recommend-item"
              v-for="(n, i) in editForm.recommendQuestion"
              @mouseenter="activeIndex = i"
              @mouseleave="activeIndex = -1"
              :key="`${i}rml`"
            >
              <el-input
                class="recommend--input"
                v-model.lazy="n.value"
                maxlength="50"
                :key="`${i}rml`"
              ></el-input>
              <span
                class="el-icon-delete recommend-del"
                @click="clearRecommend(n, i)"
                v-if="activeIndex === i"
              ></span>
            </div>
          </div>
        </div>
        <!-- 知识库库配置 -->
        <div class="block">
          <knowledgeDataField
            :knowledgeConfig="editForm.knowledgeBaseConfig"
            :category="0"
            @getSelectKnowledge="getSelectKnowledge"
            @knowledgeDelete="knowledgeDelete"
            @knowledgeRecallSet="knowledgeRecallSet"
            @updateMetaData="updateMetaData"
            :labelText="$t('agent.form.linkKnowledge')"
            :type="'knowledgeBaseConfig'"
            :appType="AGENT"
          />
        </div>

        <div class="block">
          <p class="block-title common-set">
            <span class="common-set-label">
              {{ $t('agent.form.tool') }}
              <template v-if="allTools.length">
                [{{ useToolNum }}/{{ allTools.length }}]
              </template>
            </span>
            <span @click="addTool" class="common-add">
              <span class="el-icon-plus"></span>
              <span class="handleBtn">{{ $t('agent.add') }}</span>
            </span>
          </p>
          <div class="rl tool-content">
            <div class="tool-right tool" v-show="allTools.length">
              <div class="action-list">
                <div
                  class="action-item"
                  v-for="(n, i) in allTools"
                  :key="`${i}ac`"
                >
                  <div class="name">
                    <div class="toolImg">
                      <img
                        :src="'/user/api/' + n.avatar.path"
                        v-show="n.avatar && n.avatar.path"
                      />
                    </div>
                    <el-tooltip
                      class="item"
                      effect="dark"
                      :content="displayName(n)"
                      placement="top-start"
                    >
                      <span>
                        {{
                          displayName(n).length > 20
                            ? displayName(n).substring(0, 20) + '...'
                            : displayName(n)
                        }}
                      </span>
                    </el-tooltip>
                    <el-tooltip
                      class="item"
                      effect="dark"
                      :content="n.mcpName || n.toolName"
                      placement="top-start"
                    >
                      <span
                        class="el-icon-info desc-info"
                        v-if="n.mcpName || n.toolName"
                      ></span>
                    </el-tooltip>
                  </div>
                  <div class="bt">
                    <span
                      class="el-icon-s-operation bt-operation"
                      @click="handleBuiltin(n)"
                      v-if="
                        n.type === 'action' &&
                        n.toolType &&
                        n.toolType === 'builtin'
                      "
                    ></span>
                    <el-switch
                      v-model="n.enable"
                      class="bt-switch"
                      @change="toolSwitch(n, n.type, n.enable)"
                    ></el-switch>
                    <span
                      @click="toolRemove(n, n.type)"
                      class="el-icon-delete del"
                    ></span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="block">
          <p class="block-title common-set">
            <span class="common-set-label">
              {{ $t('agent.form.safetyConfig') }}
              <el-tooltip
                class="item"
                effect="dark"
                :content="$t('agent.form.safetyConfigTips')"
                placement="top"
              >
                <span class="el-icon-question question-tips"></span>
              </el-tooltip>
            </span>
            <span class="common-add" @click="showSafety">
              <el-tooltip
                class="item"
                effect="dark"
                :content="$t('agent.form.safetyConfigTips')"
                placement="top-start"
              >
                <span class="el-icon-s-operation operation">
                  <span class="handleBtn">{{ $t('agent.form.config') }}</span>
                </span>
              </el-tooltip>
              <el-switch
                v-model="editForm.safetyConfig.enable"
                :disabled="!(editForm.safetyConfig.tables || []).length"
              ></el-switch>
            </span>
          </p>
        </div>
        <div class="block" v-if="editForm.visionsupport === 'support'">
          <p class="block-title common-set">
            <span class="common-set-label">
              {{ $t('agent.form.vision') }}
              <el-tooltip
                class="item"
                effect="dark"
                :content="$t('agent.form.visionTips1')"
                placement="top"
              >
                <span class="el-icon-question question-tips"></span>
              </el-tooltip>
            </span>
            <span class="common-add" @click="showVisualSet">
              <el-tooltip
                class="item"
                effect="dark"
                :content="$t('agent.form.visionTips')"
                placement="top-start"
              >
                <span class="el-icon-s-operation operation">
                  <span class="handleBtn">{{ $t('agent.form.config') }}</span>
                </span>
              </el-tooltip>
            </span>
          </p>
        </div>
      </div>
      <div class="block drawer-test">
        <Chat
          :editForm="editForm"
          :chatType="'test'"
          :disableClick="disableClick"
        />
      </div>
    </div>

    <!-- 编辑智能体 -->
    <CreateIntelligent
      ref="createIntelligentDialog"
      :type="'edit'"
      :editForm="editForm"
      @updateInfo="getAppDetail"
    />
    <!-- 模型设置 -->
    <ModelSet
      @setModelSet="setModelSet"
      ref="modelSetDialog"
      :modelform="editForm.modelConfig"
      :limitMaxTokens="limitMaxTokens"
    />
    <!-- 选择工具类型 -->
    <ToolDialog
      ref="toolDialog"
      @updateDetail="updateDetail"
      :assistantId="editForm.assistantId"
    />
    <!-- 敏感词设置 -->
    <setSafety ref="setSafety" @sendSafety="sendSafety" />
    <!-- 视图设置 -->
    <visualSet ref="visualSet" @sendVisual="sendVisual" />
    <!-- 内置工具详情 -->
    <ToolDetail ref="toolDetail" @updateDetail="updateDetail" />
    <!-- 提交至提示词 -->
    <createPrompt
      :isCustom="true"
      :type="promptType"
      ref="createPrompt"
      @reload="updatePrompt"
    />
    <!-- 提示词优化 -->
    <PromptOptimize ref="promptOptimize" @promptSubmit="promptSubmit" />
  </div>
</template>

<script>
import { appPublish } from '@/api/appspace';
import { store } from '@/store/index';
import { mapGetters, mapActions } from 'vuex';
import CreateIntelligent from '@/components/createApp/createIntelligent';
import setSafety from '@/components/setSafety';
import visualSet from './visualSet';
import metaSet from '@/components/metaSet';
import ModelSet from './modelSetDialog';
import { selectModelList, getRerankList } from '@/api/modelAccess';
import { AGENT } from '@/utils/commonSet';
import {
  deleteMcp,
  enableMcp,
  getAgentInfo,
  delWorkFlowInfo,
  delActionInfo,
  putAgentInfo,
  enableWorkFlow,
  enableAction,
  delCustomBuiltIn,
  switchCustomBuiltIn,
  getAgentPublishedInfo,
} from '@/api/agent';
import ToolDialog from './toolDialog';
import ToolDetail from './toolDetail';
import { readWorkFlow } from '@/api/workflow';
import Chat from './chat';
import LinkIcon from '@/components/linkIcon.vue';
import promptTemplate from './prompt/index.vue';
import createPrompt from '@/components/createApp/createPrompt.vue';
import PromptOptimize from '@/components/promptOptimize.vue';
import knowledgeDataField from '@/components/app/knowledgeDataField.vue';
import VersionPopover from '@/components/versionPopover.vue';
import CopyIcon from '@/components/copyIcon.vue';
export default {
  components: {
    CopyIcon,
    VersionPopover,
    LinkIcon,
    Chat,
    CreateIntelligent,
    ModelSet,
    ToolDialog,
    setSafety,
    visualSet,
    metaSet,
    ToolDetail,
    promptTemplate,
    createPrompt,
    PromptOptimize,
    knowledgeDataField,
  },
  provide() {
    return {
      getPrompt: this.getPrompt,
    };
  },
  watch: {
    agentFormParams: {
      handler(newVal) {
        if (this.isSettingFromDetail) return;

        if (this.debounceTimer) {
          clearTimeout(this.debounceTimer);
        }

        this.debounceTimer = setTimeout(() => {
          if (!this.initialAutoSaveSnapshot) {
            this.initialAutoSaveSnapshot = JSON.parse(JSON.stringify(newVal));
            return;
          }

          const changed =
            JSON.stringify(newVal) !==
            JSON.stringify(this.initialAutoSaveSnapshot);

          if (changed) {
            if (
              this.editForm.modelParams !== '' &&
              this.editForm.prologue !== ''
            ) {
              this.updateInfo();
            }
          }
        }, 500);
      },
      deep: true,
    },
  },
  computed: {
    ...mapGetters('app', ['cacheData']),
    ...mapGetters('user', ['commonInfo']),
    agentFormParams() {
      const {
        modelParams,
        modelConfig,
        prologue,
        knowledgeBaseConfig,
        safetyConfig,
        recommendQuestion,
        visionConfig,
      } = this.editForm;

      return {
        modelParams,
        modelConfig,
        prologue,
        knowledgeBaseConfig,
        safetyConfig,
        recommendQuestion,
        visionConfig,
      };
    },
    useToolNum() {
      return this.allTools.filter(item => item.enable).length;
    },
  },
  data() {
    return {
      AGENT,
      disableClick: false,
      version: '',
      promptType: 'create',
      limitMaxTokens: 4096,
      knowledgeIndex: -1,
      currentKnowledgeId: '',
      currentMetaData: {},
      knowledgeCheckData: [],
      activeIndex: -1,
      rerankOptions: [],
      initialEditForm: null,
      publishType: this.$route.query.publishType,
      publishForm: {
        publishType: 'private',
        version: '',
        desc: '',
      },
      publishRules: {
        version: [
          {
            required: true,
            message: this.$t('list.version.noPlaceholder'),
            trigger: 'blur',
          },
          {
            pattern: /^v\d+\.\d+\.\d+$/,
            message: this.$t('list.version.versionMsg'),
            trigger: 'blur',
          },
        ],
        desc: [
          {
            required: true,
            message: this.$t('list.version.descPlaceholder'),
            trigger: 'blur',
          },
        ],
        publishType: [
          {
            required: true,
            message: this.$t('common.select.placeholder'),
            trigger: 'change',
          },
        ],
      },
      editForm: {
        newAgent: false,
        functionCalling: '',
        visionsupport: '',
        assistantId: '',
        uuid: '',
        avatar: {},
        name: '',
        desc: '',
        rerankParams: '',
        modelParams: '',
        prologue: '', //开场白
        instructions: '', //系统提示词
        visionConfig: {
          //视觉配置
          picNum: 3,
          maxPicNum: 6,
        },
        knowledgeBaseConfig: {
          config: {
            keywordPriority: 0.8, //关键词权重
            matchType: 'mix', //vector（向量检索）、text（文本检索）、mix（混合检索：向量+文本）
            priorityMatch: 1, //权重匹配，只有在混合检索模式下，选择权重设置后，这个才设置为1
            rerankModelId: '', //rerank模型id
            semanticsPriority: 0.2, //语义权重
            topK: 5, //topK 获取最高的几行
            threshold: 0.4, //过滤分数阈值
            maxHistory: 0, //最长上下文
            useGraph: false,
          },
          knowledgebases: [],
        },
        knowledgeConfig: {},
        recommendQuestion: [
          {
            value: '',
          },
        ],
        modelConfig: {
          temperature: 0.7,
          topP: 1,
          frequencyPenalty: 0,
          presencePenalty: 0,
          maxTokens: 512,
          maxTokensEnable: true,
          frequencyPenaltyEnable: true,
          temperatureEnable: true,
          topPEnable: true,
          presencePenaltyEnable: true,
        },
        safetyConfig: {
          enable: false,
          tables: [],
        },
      },
      hasPluginPermission: false,
      modelLoading: false,
      wfDialogVisible: false,
      workFlowInfos: [],
      actionInfos: [],
      mcpInfos: [],
      allTools: [], //所有的工具
      workflowList: [],
      modelParams: {},
      modelOptions: [],
      selectKnowledge: [],
      knowledgeData: [],
      loadingPercent: 10,
      nameStatus: '',
      saved: false, //按钮
      loading: false, //按钮
      t: null,
      logoFileList: [],
      imageUrl: '',
      defaultLogo: require('@/assets/imgs/bg-logo.png'),
      debounceTimer: null, //防抖计时器
      initialAutoSaveSnapshot: null,
      isSettingFromDetail: false, // 防止详情数据触发更新标记
      nameMap: {
        workflow: {
          displayName: this.$t('menu.app.workflow'),
          propName: 'name',
        },
        mcp: {
          displayName: 'MCP' + this.$t('tool.tool'),
          propName: 'actionName',
        },
        action: {
          displayName: this.$t('menu.app.custom'),
          propName: 'actionName',
        },
        // 可以继续添加其他类型
        default: {
          displayName: this.$t('knowledgeManage.docList.unknown'),
          propName: 'name', // 默认属性名
        },
      },
    };
  },
  mounted() {
    this.initialEditForm = JSON.parse(JSON.stringify(this.editForm));
  },
  created() {
    this.reloadData();
  },
  beforeDestroy() {
    store.dispatch('app/initState');
    this.clearMaxPicNum();
  },
  methods: {
    reloadData() {
      this.getModelData(); //获取模型列表
      this.getRerankData(); //获取rerank模型
      if (this.$route.query.id) {
        this.editForm.assistantId = this.$route.query.id;
        setTimeout(() => {
          this.getAppDetail();
        }, 500);
      }
      //判断是否有插件管理的权限
      const accessCert = localStorage.getItem('access_cert');
      const permission = accessCert
        ? JSON.parse(accessCert).user.permission.orgPermission
        : '';
      this.hasPluginPermission = permission.indexOf('plugin') !== -1;
    },
    previewVersion(item) {
      this.disableClick = !item.isCurrent;
      this.version = item.version || '';
      this.getAppDetail();
    },
    ...mapActions('app', ['setMaxPicNum', 'clearMaxPicNum']),
    //系统提示词失去焦点时，触发提示词更新
    handleInstructionsBlur(e) {
      this.updateInfo();
    },
    syncAutoSaveBaseline() {
      this.initialAutoSaveSnapshot = JSON.parse(
        JSON.stringify(this.agentFormParams),
      );
    },
    //获取知识库或问答库选中数据
    getSelectKnowledge(data, type) {
      this.editForm[type]['knowledgebases'] = data;
    },
    //删除知识库或问答库
    knowledgeDelete(index, type) {
      this.editForm[type]['knowledgebases'].splice(index, 1);
    },
    //设置知识库或问答库召回参数
    knowledgeRecallSet(data, type) {
      if (data) {
        this.editForm[type]['config'] = data;
      } else {
        this.editForm[type]['config'] = this.editForm[type]['config'];
      }
    },
    //更新知识库元数据
    updateMetaData(data, index, type) {
      this.$set(this.editForm[type]['knowledgebases'], index, {
        ...this.editForm[type]['knowledgebases'][index],
        ...data,
      });
    },
    showPromptCompare() {
      this.$router.push({
        path: `/agent/promptCompare/${this.editForm.assistantId}`,
      });
    },
    updatePrompt() {
      this.$refs.promptTemplate.getPromptTemplateList();
    },
    handleShowPrompt() {
      this.$refs.createPrompt.openDialog({
        prompt: this.editForm.instructions,
      });
    },
    showPromptOptimize() {
      if (!this.editForm.instructions) {
        this.$message.warning(this.$t('tempSquare.promptOptimizeHint'));
        return;
      }
      this.$refs.promptOptimize.openDialog({
        prompt: this.editForm.instructions,
      });
    },
    promptSubmit(prompt) {
      this.editForm.instructions = prompt;
      this.updateInfo();
    },
    getPrompt(prompt) {
      this.editForm.instructions = prompt;
      this.updateInfo();
    },
    handleBuiltin(n) {
      this.$refs.toolDetail.showDialog(n);
    },
    showVisualSet() {
      this.$refs.visualSet.showDialog(this.editForm.visionConfig);
    },
    sendVisual(data) {
      this.editForm.visionConfig.picNum = data.picNum;
    },
    handleModelChange(val) {
      this.setModelInfo(val);
    },
    setModelInfo(val) {
      if (!val) return;
      const selectedModel = this.modelOptions.find(
        item => item.modelId === val,
      );
      if (selectedModel) {
        this.editForm.modelParams = val;
        this.editForm.visionsupport = selectedModel.config.visionSupport;
        this.editForm.functionCalling = selectedModel.config.functionCalling;
        const maxTokens = selectedModel.config.maxTokens;
        this.limitMaxTokens = maxTokens && maxTokens > 0 ? maxTokens : 4096;
      } else {
        this.editForm.modelParams = '';
        if (val) this.$message.warning(this.$t('agent.form.modelNotSupport'));
      }
    },
    handlePublishSet() {
      this.$router.push({
        path: `/agent/publishSet`,
        query: {
          appId: this.editForm.assistantId,
          appType: AGENT,
          name: this.editForm.name,
        },
      });
    },
    displayName(item) {
      const config = this.nameMap[item.type] || this.nameMap['default'];
      return item[config.propName];
    },
    updateDetail() {
      this.getAppDetail();
    },
    showSafety() {
      this.$refs.setSafety.showDialog(this.editForm.safetyConfig.tables);
    },
    sendSafety(data) {
      const tablesData = data.map(({ tableId, tableName }) => ({
        tableId,
        tableName,
      }));
      this.editForm.safetyConfig.tables = tablesData;
    },
    actionSwitch(id) {
      enableAction({
        actionId: id,
      }).then(res => {
        if (res.code === 0) {
          this.getAppDetail();
        }
      });
    },
    toolSwitch(n, type, enable) {
      if (type === 'workflow') {
        this.workflowSwitch(n.workFlowId, enable);
      } else if (type === 'mcp') {
        this.mcpSwitch(n, enable);
      } else {
        this.customSwitch(n, enable);
      }
    },
    customSwitch(n, enable) {
      switchCustomBuiltIn({
        assistantId: this.editForm.assistantId,
        actionName: n.actionName,
        toolId: n.toolId,
        toolType: n.toolType,
        enable,
      })
        .then(res => {
          if (res.code === 0) {
            this.getAppDetail();
          }
        })
        .catch(() => {});
    },
    mcpSwitch(n, enable) {
      enableMcp({
        assistantId: this.editForm.assistantId,
        actionName: n.actionName,
        enable,
        mcpId: n.mcpId,
        mcpType: n.mcpType,
      })
        .then(res => {
          if (res.code === 0) {
            this.getAppDetail();
          }
        })
        .catch(() => {});
    },
    workflowSwitch(id, enable) {
      enableWorkFlow({
        assistantId: this.editForm.assistantId,
        workFlowId: id,
        enable,
      })
        .then(res => {
          if (res.code === 0) {
            this.getAppDetail();
          }
        })
        .catch(() => {});
    },
    addTool() {
      const data = {
        mcpInfos: this.mcpInfos,
        workFlowInfos: this.workFlowInfos,
        customInfos: this.actionInfos,
      };
      this.$refs.toolDialog.showDialog(data);
    },
    rerankVisible(val) {
      if (val) {
        this.getRerankData();
      }
    },
    getRerankData() {
      getRerankList().then(res => {
        if (res.code === 0) {
          this.rerankOptions = res.data.list || [];
        }
      });
    },
    goBack() {
      this.$router.push({
        path: '/appSpace/agent',
      });
    },
    savePublish() {
      if (this.editForm.modelParams === '') {
        this.$message.warning(this.$t('agent.form.selectModel'));
        return false;
      }
      if (this.editForm.prologue === '') {
        this.$message.warning(this.$t('agent.form.inputPrologue'));
        return false;
      }

      this.$refs.publishForm.validate(valid => {
        if (valid) {
          const data = {
            appId: this.editForm.assistantId,
            appType: AGENT,
            publishType: this.publishForm.publishType,
            desc: this.publishForm.desc,
            version: this.publishForm.version,
          };
          appPublish(data).then(res => {
            if (res.code === 0) {
              this.$router.push({ path: '/explore' });
            }
          });
        }
      });
    },
    setModelSet(data) {
      this.editForm.modelConfig = data;
    },
    showModelSet() {
      this.$refs.modelSetDialog.showDialog();
    },
    editAgent() {
      this.$refs.createIntelligentDialog.openDialog();
    },
    async getWorkFlowDetail(n, index) {
      let params = {
        workflowID: n.appId,
      };
      let res = await readWorkFlow(params);
      if (res.code === 0) {
        this.doCreateWorkFlow(n.appId, res.data.base64OpenAPISchema, index);
      }
    },
    preAddWorkflow() {
      this.wfDialogVisible = true;
    },
    toolRemove(n, type) {
      if (type === 'workflow') {
        this.doDeleteWorkflow(n.workFlowId);
      } else if (type === 'mcp') {
        this.mcpRemove(n);
      } else {
        this.customRemove(n);
      }
    },
    customRemove(n) {
      delCustomBuiltIn({
        assistantId: this.editForm.assistantId,
        toolId: n.toolId,
        toolType: n.toolType,
        actionName: n.actionName,
      })
        .then(res => {
          if (res.code === 0) {
            this.$message.success(this.$t('agent.form.deleteSuccess'));
            this.getAppDetail();
          }
        })
        .catch(err => {});
    },
    mcpRemove(n) {
      deleteMcp({
        assistantId: this.editForm.assistantId,
        actionName: n.actionName,
        mcpId: n.mcpId,
        mcpType: n.mcpType,
      })
        .then(res => {
          if (res.code === 0) {
            this.$message.success(this.$t('agent.form.deleteSuccess'));
            this.getAppDetail();
          }
        })
        .catch(err => {});
    },
    visibleChange(val) {
      //下拉框显示的时候请求模型列表
      if (val) {
        this.getModelData();
      }
    },
    async getModelData() {
      this.modelLoading = true;
      const res = await selectModelList();
      if (res.code === 0) {
        this.modelOptions = res.data.list || [];
        this.modelLoading = false;
      }
      this.modelLoading = false;
    },
    async updateInfo() {
      //模型数据
      let modeInfo;
      if (
        typeof this.editForm.modelParams === 'object' &&
        this.editForm.modelParams
      ) {
        modeInfo = this.editForm.modelParams;
      } else {
        modeInfo = this.modelOptions.find(
          item => item.modelId === this.editForm.modelParams,
        );
      }
      const rerankInfo = this.rerankOptions.find(
        item =>
          item.modelId ===
          this.editForm.knowledgeBaseConfig.config.rerankModelId,
      );
      const recommendQuestion = this.editForm.recommendQuestion.map(
        item => item.value,
      );
      const params = {
        assistantId: this.editForm.assistantId,
        prologue: this.editForm.prologue,
        recommendQuestion:
          recommendQuestion.length > 0 && recommendQuestion[0] !== ''
            ? recommendQuestion
            : [],
        instructions: this.editForm.instructions,
        knowledgeBaseConfig: this.editForm.knowledgeBaseConfig,
        modelConfig: {
          config: this.editForm.modelConfig,
          displayName: modeInfo.displayName,
          model: modeInfo.model,
          modelId: modeInfo.modelId,
          modelType: modeInfo.modelType,
          provider: modeInfo.provider,
        },
        safetyConfig: this.editForm.safetyConfig,
        visionConfig: {
          picNum: this.editForm.visionConfig.picNum,
        },
        rerankConfig: rerankInfo
          ? {
              displayName: rerankInfo.displayName,
              model: rerankInfo.model,
              modelId: rerankInfo.modelId,
              modelType: rerankInfo.modelType,
              provider: rerankInfo.provider,
            }
          : {},
      };
      let res = await putAgentInfo(params);
      if (res.code === 0) {
        this.getAppDetail();
      }
    },
    startLoading(val) {
      this.loadingPercent = val;
      if (val === 0) {
        this.loading = true;
      }
      if (val === 100) {
        setTimeout(() => {
          this.loading = false;
        }, 500);
      }
    },
    async getAppDetail() {
      this.startLoading(0);
      this.isSettingFromDetail = true;
      let res;
      if (this.version) {
        res = await getAgentPublishedInfo({
          assistantId: this.editForm.assistantId,
          version: this.version,
        });
      } else
        res = await getAgentInfo({
          assistantId: this.editForm.assistantId,
        });
      if (res.code === 0) {
        this.startLoading(100);
        let data = res.data;
        this.publishType = data.publishType;
        //兼容后端知识库数据返回null
        if (
          res.data.knowledgeBaseConfig &&
          res.data.knowledgeBaseConfig !== null
        ) {
          this.editForm.knowledgeBaseConfig.knowledgebases =
            res.data.knowledgeBaseConfig.knowledgebases;
          this.editForm.knowledgeBaseConfig.config =
            res.data.knowledgeBaseConfig.config === null ||
            !res.data.knowledgeBaseConfig.config.matchType
              ? this.editForm.knowledgeBaseConfig.config
              : res.data.knowledgeBaseConfig.config;
        }

        this.editForm = {
          ...this.editForm,
          uuid: data.uuid,
          newAgent: data.newAgent,
          avatar: data.avatar || {},
          prologue: data.prologue || '', //开场白
          name: data.name || '',
          desc: data.desc || '',
          instructions: data.instructions || '', //系统提示词
          rerankParams: data.rerankConfig.modelId || '',
          visionConfig: data.visionConfig, //图片配置
          modelConfig:
            data.modelConfig.config !== null
              ? data.modelConfig.config
              : this.editForm.modelConfig,
          recommendQuestion:
            data.recommendQuestion && data.recommendQuestion.length > 0
              ? data.recommendQuestion.map((n, index) => {
                  return {
                    value: n,
                  };
                })
              : [],
          safetyConfig:
            data.safetyConfig !== null
              ? data.safetyConfig
              : this.editForm.safetyConfig,
        };

        this.editForm.knowledgeBaseConfig.config.rerankModelId =
          res.data.rerankConfig.modelId;
        //设置模型信息
        this.setModelInfo(data.modelConfig.modelId);

        //回显自定义插件
        this.workFlowInfos = data.workFlowInfos || [];
        this.mcpInfos = data.mcpInfos || [];
        this.actionInfos = data.toolInfos || [];
        this.allTools = [
          ...this.workFlowInfos.map(item => ({
            ...item,
            type: 'workflow',
          })),
          ...this.mcpInfos.map(item => ({
            ...item,
            type: 'mcp',
          })),
          ...this.actionInfos.map(item => ({
            ...item,
            type: 'action',
          })),
        ];

        this.setMaxPicNum(this.editForm.visionConfig.picNum);

        this.$nextTick(() => {
          this.isSettingFromDetail = false;
          this.syncAutoSaveBaseline();
        });
      } else {
        this.isSettingFromDetail = false;
      }
    },
    async doDeleteWorkflow(workFlowId) {
      if (this.editForm.assistantId) {
        let res = await delWorkFlowInfo({
          workFlowId,
          assistantId: this.editForm.assistantId,
        });
        if (res.code === 0) {
          this.$message.success(this.$t('agent.delPluginTips'));
          this.getAppDetail();
        }
      } else {
        this.$message.error(this.$t('agent.otherTips'));
      }
    },
    //推荐问题
    addRecommend() {
      if (this.editForm.recommendQuestion.length > 3) {
        return;
      }
      this.editForm.recommendQuestion.push({
        value: '',
      });
    },
    clearRecommend(n, index) {
      if (this.editForm.recommendQuestion.length === 1) return;
      this.editForm.recommendQuestion.splice(index, 1);
      this.activeIndex = -1;
    },
    async preDelAction(actionId) {
      this.$confirm(
        this.$t('createApp.delActionTips'),
        this.$t('knowledgeManage.tip'),
        {
          confirmButtonText: this.$t('createApp.save'),
          cancelButtonText: this.$t('createApp.cancel'),
          type: 'warning',
        },
      )
        .then(async () => {
          let res = await delActionInfo({
            actionId,
          });
          if (res.code === 0) {
            this.$message.success(this.$t('createApp.delSuccess'));
            this.getAppDetail();
          }
        })
        .catch(() => {});
    },
  },
};
</script>

<style lang="scss" scoped>
@import '@/style/draft.scss';

.agent_form {
  gap: 10px;

  .drawer-info {
    position: relative;
    width: 30%;
    margin: 10px 0;
  }

  .labelTitle {
    font-size: 18px;
    font-weight: 800;
  }

  .promptTitle {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .prompt-title-icon {
      display: flex;
      align-items: center;
    }

    h3 {
      font-size: 18px;
      font-weight: 800;
    }

    span {
      margin-left: 5px;
      font-size: 16px;
      color: $color;
      cursor: pointer;
      display: inline-block;
      padding: 8px;
      border-radius: 50%;
      background: #e0e7ff;
    }

    .tool-icon {
      display: inline-block;
      width: 32px;
      height: 32px;
      cursor: pointer;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }
  }

  .drawer-form {
    width: 30%;
    margin: 10px 0;

    /*通用*/
    .block {
      margin-bottom: 15px;

      .tool-content {
        display: flex;
        justify-content: space-between;
        gap: 10px;

        .tool {
          width: 100%;
          max-height: 300px;
          overflow-y: auto;

          .action-list {
            margin: 10px 0 15px 0;
            width: 100%;

            .action-item {
              display: flex;
              justify-content: space-between;
              align-items: center;
              border: 1px solid #ddd;
              border-radius: 6px;
              margin-bottom: 5px;
              width: 100%;

              .name {
                width: 80%;
                box-sizing: border-box;
                padding: 10px;
                cursor: pointer;
                display: flex;
                align-items: center;
                color: #333;

                .desc-info {
                  color: #ccc;
                  margin-left: 4px;
                }

                .toolImg {
                  width: 30px;
                  height: 30px;
                  border-radius: 50%;
                  background: #eee;
                  margin-right: 5px;

                  img {
                    width: 100%;
                    height: 100%;
                    border-radius: 50%;
                    object-fit: cover;
                  }
                }
              }

              .bt {
                text-align: center;
                width: 30%;
                display: flex;
                justify-content: flex-end;
                align-items: center;
                padding-right: 10px;
                box-sizing: border-box;
                cursor: pointer;

                .del {
                  color: $btn_bg;
                  font-size: 16px;
                  line-height: 20px;
                }

                .bt-switch {
                  margin: 0 6px 0 6px;
                }

                .bt-operation {
                  font-size: 16px;
                  line-height: 20px;
                }
              }
            }
          }
        }
      }

      .model-select-tips {
        margin-top: 10px;
        color: #dc6803;
      }
    }

    /*推荐问题*/
    .recommend-box {
      .recommend-title {
        display: flex;
        justify-content: space-between;

        span {
          font-size: 15px;
        }
      }

      .recommend-item {
        margin-bottom: 12px;
        display: flex;
        justify-content: space-between;
        position: relative;

        .recommend--input {
          width: calc(100% - 30px);
          margin-right: 30px;
        }

        .recommend-del {
          position: absolute;
          right: 10px;
          top: 10px;
          color: #595959;
          cursor: pointer;
        }
      }
    }
  }

  .drawer-test {
    width: 670px;
  }
}
</style>
