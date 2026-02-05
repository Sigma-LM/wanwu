<template>
  <div>
    <el-dialog
      top="10vh"
      :title="title"
      :close-on-click-modal="false"
      :visible.sync="dialogVisible"
      width="40%"
      :before-close="handleClose"
    >
      <div style="display: flex; flex-direction: row; align-items: center">
        <div class="img">
          <img
            :src="avatarSrc(item.avatar.path)"
            v-if="item.avatar && item.avatar.path"
          />
        </div>
        <div class="info">
          <span class="name">{{ item.name }}</span>
        </div>
      </div>
      <el-input
        class="desc-input"
        style="margin-top: 10px"
        v-model="item.desc"
        :placeholder="$t('multiAgentEdit.descPlaceholder')"
        type="textarea"
        show-word-limit
        :rows="12"
      ></el-input>
      <span slot="footer" class="dialog-footer">
        <el-button @click="handleClose()">
          {{ $t('common.confirm.cancel') }}
        </el-button>
        <el-button type="primary" @click="submit()">
          {{ $t('common.confirm.confirm') }}
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
import { avatarSrc } from '@/utils/util';
export default {
  data() {
    return {
      title: this.$t('multiAgentEdit.title'),
      dialogVisible: false,
      item: {},
    };
  },
  methods: {
    avatarSrc,
    handleClose() {
      this.dialogVisible = false;
    },
    submit() {
      this.$emit('submit', this.item);
      this.handleClose();
    },
    showDialog(item) {
      this.dialogVisible = true;
      this.item = JSON.parse(JSON.stringify(item));
    },
  },
};
</script>
<style lang="scss" scoped>
::v-deep .desc-input {
  .el-textarea__inner {
    height: 90px !important;
  }
}
.img {
  width: 35px;
  height: 35px;
  background: #eee;
  border-radius: 50%;
  display: inline-block;
  margin-right: 5px;
  img {
    width: 100%;
    height: 100%;
    border-radius: 50%;
    object-fit: cover;
  }
}

.info {
  flex: 1;
  color: #333;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-right: 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;

  .name {
    font-size: 14px;
    font-weight: 600;
    color: #1c1d23;
  }
}
</style>
