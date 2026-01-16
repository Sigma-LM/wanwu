<template>
  <el-select
    v-model="modelId"
    :placeholder="placeholder || $t('agent.form.modelSearchPlaceholder')"
    :loading-text="loadingText || $t('agent.toolDetail.modelLoadingText')"
    class="cover-input-icon model-select model-select-custom"
    :loading="modelLoading"
    filterable
    value-key="modelId"
    @change="handleModelChange"
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
        <div class="model-select-tags" v-if="item.tags && item.tags.length > 0">
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
</template>

<script>
import { selectModelList } from '@/api/modelAccess';

export default {
  name: 'ModelSelector',
  props: {
    value: {
      type: [String, Number],
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    loadingText: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      modelOptions: [],
      modelLoading: false,
    };
  },
  computed: {
    modelId: {
      get() {
        return this.value;
      },
      set(val) {
        this.$emit('input', val);
      },
    },
  },
  created() {
    this.getModelData();
  },
  methods: {
    async getModelData() {
      this.modelLoading = true;
      try {
        const res = await selectModelList();
        if (res.code === 0) {
          this.modelOptions = res.data.list || [];
        }
      } catch (error) {
        console.error('Failed to fetch model list:', error);
      } finally {
        this.modelLoading = false;
      }
    },
    handleModelChange(val) {
      const selectedModel = this.modelOptions.find(
        item => item.modelId === val,
      );
      this.$emit('change', val, selectedModel);
    },
  },
};
</script>

<style lang="scss" scoped>
@import '@/style/draft.scss';
.model-select-custom {
  width: 100%;
}
</style>
