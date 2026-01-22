<template>
  <el-drawer
    :title="$t('knowledgeManage.externalAPI.title')"
    :visible.sync="drawer"
    direction="rtl"
    :before-close="handleClose"
    size="70%"
  >
    <div style="margin: 20px">
      <div>{{ $t('knowledgeManage.externalAPI.hint') }}</div>
      <!--      <div>{{ $t('knowledgeManage.externalAPI.tips') }}</div>-->
      <el-button type="primary" style="margin-top: 15px" @click="handleAdd">
        {{ $t('common.button.add') + $t('knowledgeManage.externalAPI.title') }}
      </el-button>
      <el-table
        :data="tableData"
        style="width: 100%; margin-top: 15px"
        :header-cell-style="{ background: '#F9F9F9', color: '#999999' }"
      >
        <el-table-column
          prop="name"
          :label="$t('knowledgeManage.externalAPI.name')"
        />
        <el-table-column
          prop="description"
          :label="$t('knowledgeManage.externalAPI.desc')"
        />
        <el-table-column prop="baseUrl" label="Base Url" />
        <el-table-column prop="apiKey" label="API Key" />
        <el-table-column :label="$t('common.table.operation')">
          <template slot-scope="scope">
            <el-button
              class="operation"
              type="text"
              @click="handleEdit(scope.row)"
            >
              {{ $t('common.button.edit') }}
            </el-button>
            <el-button type="text" @click="handleDelete(scope.row)">
              {{ $t('common.button.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <externalAPIDialog
      ref="externalAPIDialog"
      @confirm="getExternalAPIList"
    ></externalAPIDialog>
  </el-drawer>
</template>

<script>
import externalAPIDialog from '@/components/externalAPIDialog.vue';
import { getExternalAPIList, delExternalAPI } from '@/api/knowledge';

export default {
  data() {
    return {
      drawer: false,
      tableData: [],
    };
  },
  components: {
    externalAPIDialog,
  },
  methods: {
    getExternalAPIList() {
      getExternalAPIList().then(res => {
        if (res.code === 0) {
          this.tableData = res.data.externalApiList;
          this.$emit('update', this.tableData);
        }
      });
    },
    showDrawer() {
      this.drawer = true;
      this.getExternalAPIList();
    },
    handleAdd() {
      this.$refs.externalAPIDialog.showDialog();
    },
    handleEdit(row) {
      this.$refs.externalAPIDialog.form = { ...row };
      this.$refs.externalAPIDialog.showDialog();
    },
    handleDelete(row) {
      delExternalAPI({ externalApiId: row.externalApiId }).then(res => {
        if (res.code === 0) {
          this.$message.success(this.$t('common.info.delete'));
          this.tableData.splice(this.tableData.indexOf(row), 1);
          this.$emit('update', this.tableData);
        }
      });
    },
    handleClose() {
      this.drawer = false;
    },
  },
};
</script>

<style scoped lang="scss">
::v-deep .el-drawer__header {
  margin-bottom: 0;
  span {
    color: #000;
    font-size: 22px !important;
  }
}
</style>
