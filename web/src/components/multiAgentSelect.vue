<template>
  <div>
    <el-dialog
      :visible.sync="dialogVisible"
      width="40%"
      :before-close="handleClose"
    >
      <template slot="title">
        <div class="dialog_title">
          <h3>
            {{ $t('multiAgentSelect.title') }}
          </h3>
          <el-input
            v-model="agentName"
            :placeholder="$t('multiAgentSelect.searchPlaceholder')"
            class="input"
            suffix-icon="el-icon-search"
            @keyup.enter.native="searchAgent"
            clearable
          ></el-input>
        </div>
      </template>
      <div class="content">
        <div
          v-for="item in agentData"
          :key="item['appId']"
          class="content_item"
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
              <span class="desc">{{ item.desc }}</span>
            </div>
          </div>
          <el-button
            type="primary"
            @click="bindAgent($event, item)"
            v-if="!item.checked"
            size="small"
          >
            {{ $t('multiAgentSelect.add') }}
          </el-button>
          <el-button type="primary" v-else size="small">
            {{ $t('multiAgentSelect.added') }}
          </el-button>
        </div>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="handleClose">
          {{ $t('common.button.confirm') }}
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
import { getMultiAgentList, bindMultiAgent } from '@/api/agent';
import { avatarSrc } from '@/utils/util';
export default {
  props: ['appId'],
  data() {
    return {
      dialogVisible: false,
      agentData: [],
      agentList: [],
      agentName: '',
    };
  },
  created() {
    this.getAgentList('');
  },
  methods: {
    avatarSrc,
    getAgentList(name) {
      getMultiAgentList({ name })
        .then(res => {
          if (res.code === 0) {
            this.agentData = (res.data.list || [])
              .filter(item => item.appId !== this.appId)
              .map(m => ({
                ...m,
                checked: this.agentList.some(item => item.agentId === m.appId),
              }));
          }
        })
        .catch(() => {});
    },
    bindAgent(e, item) {
      if (!e) return;
      item.checked = !item.checked;
      bindMultiAgent({
        agentId: item.appId,
        assistantId: this.appId,
      })
        .then(res => {
          if (res.code === 0) {
            this.$message.success(this.$t('common.message.success'));
            this.$emit('bindAgent', item);
          }
        })
        .catch(() => {});
    },
    searchAgent() {
      this.getAgentList(this.agentName);
    },
    showDialog(data) {
      this.dialogVisible = true;
      this.agentList = data || [];
      this.agentData = this.agentData.map(m => ({
        ...m,
        checked: this.agentList.some(item => item.agentId === m.appId),
      }));
    },
    handleClose() {
      this.dialogVisible = false;
    },
  },
};
</script>
<style lang="scss" scoped>
::v-deep {
  .el-dialog__body {
    padding: 10px 20px;
  }
  .el-dialog__header {
    display: flex;
    align-items: center;
    .el-dialog__headerbtn {
      top: unset !important;
    }
  }
}
.dialog_title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex: 1;
  h3 {
    font-size: 16px;
    font-weight: bold;
  }
  .input {
    width: 250px;
    margin-right: 28px;
  }
}
.content {
  padding: 10px 0;
  max-height: 300px;
  overflow-y: auto;
  .content_item {
    padding: 5px 20px;
    border-bottom: 1px solid $color_opacity;
    border-radius: 6px;
    margin-bottom: 10px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    ::v-deep {
      .el-button--primary {
        background: #fff !important;
        border: 1px solid #eee !important;
        padding: 8px 16px;
        border-radius: 6px;
        span {
          color: $color !important;
          font-size: 14px;
        }
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
      display: flex;
      flex-direction: column;
      gap: 4px;

      .name {
        font-size: 14px;
        font-weight: 600;
        color: #1c1d23;
      }
      .desc {
        font-size: 12px;
        color: rgba(28, 29, 35, 0.8);
      }
    }
  }
  .content_item:hover {
    background: $color_opacity;
  }
}
.active {
  border: 1px solid $color !important;
  color: #fff;
  background: $color;
}
</style>
