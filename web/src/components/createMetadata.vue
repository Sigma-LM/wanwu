<template>
  <el-dialog
    :title="$t('knowledgeManage.docList.metaDataManagement')"
    :visible.sync="metaVisible"
    width="550px"
    append-to-body
    :before-close="handleClose"
  >
    <metadata
      ref="metadata"
      @updateMeta="updateMeta"
      type="create"
      :knowledgeId="knowledgeId"
      :disableOld="true"
      class="metadata"
    />
    <span slot="footer" class="dialog-footer">
      <el-button type="primary" @click="createMeta">
        {{ $t('common.button.create') }}
      </el-button>
      <el-button type="primary" @click="submitMeta" :disabled="isDisabled">
        {{ $t('common.button.confirm') }}
      </el-button>
    </span>
  </el-dialog>
</template>

<script>
import metadata from '@/views/knowledge/component/metadata.vue';
import { updateDocMeta } from '@/api/knowledge';

export default {
  name: 'createMetadata',
  components: {
    metadata,
  },
  data() {
    return {
      metaVisible: false,
      knowledgeId: '',
      metaData: [],
      isDisabled: false,
    };
  },
  watch: {
    metaData: {
      handler(val) {
        this.isDisabled = !!(
          val.some(item => !item.metaKey || !item.metaValueType) || !val.length
        );
      },
    },
  },
  methods: {
    showDialog(knowledgeId) {
      this.metaVisible = true;
      this.knowledgeId = knowledgeId;
      this.$nextTick(() => {
        this.$refs.metadata.getList();
      });
    },
    createMeta() {
      this.$refs.metadata.createMetaData();
      this.scrollToBottom();
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.metadata;
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    },
    submitMeta() {
      this.isDisabled = true;
      const metaList = this.metaData
        .filter(item => item.option !== '')
        .map(({ metaId, metaKey, metaValueType, option }) => ({
          metaKey,
          ...(option === 'add' ? { metaValueType } : {}),
          option,
          ...(option === 'update' || option === 'delete' ? { metaId } : {}),
        }));
      const data = {
        docId: '',
        knowledgeId: this.knowledgeId,
        metaDataList: metaList,
      };
      updateDocMeta(data)
        .then(res => {
          if (res.code === 0) {
            this.$message.success(this.$t('common.message.success'));
            this.$emit('updateMeta');
            this.metaVisible = false;
            this.isDisabled = false;
          }
        })
        .catch(() => {
          this.isDisabled = false;
        });
    },
    updateMeta(data) {
      this.metaData = data;
    },
    handleClose() {
      this.metaVisible = false;
    },
  },
};
</script>

<style lang="scss" scoped>
.metadata {
  max-height: 400px;
  overflow-y: auto;
}
</style>
