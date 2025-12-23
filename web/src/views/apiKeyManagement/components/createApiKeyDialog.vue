<template>
  <el-dialog
    :title="
      isEdit ? $t('apiKeyManage.dialog.edit') : $t('apiKeyManage.dialog.create')
    "
    :visible.sync="dialogVisible"
    width="580px"
    append-to-body
    :close-on-click-modal="false"
    :before-close="handleClose"
  >
    <el-form :model="form" :rules="rules" ref="form" style="margin-top: -16px">
      <el-form-item :label="$t('apiKeyManage.table.name')" prop="name">
        <el-input
          v-model="form.name"
          maxlength="20"
          show-word-limit
          :placeholder="$t('common.hint.nameHint')"
        />
      </el-form-item>
      <el-form-item :label="$t('apiKeyManage.table.desc')" prop="desc">
        <el-input
          v-model="form.desc"
          maxlength="100"
          show-word-limit
          :placeholder="$t('common.input.placeholder')"
        />
      </el-form-item>
      <el-form-item
        :label="$t('apiKeyManage.table.expiredAt')"
        prop="expiredAt"
      >
        <el-date-picker
          style="width: 100%"
          type="date"
          v-model="form.expiredAt"
          value-format="yyyy-MM-dd"
          :placeholder="$t('common.select.placeholder')"
          :picker-options="pickerOptions"
        ></el-date-picker>
      </el-form-item>
    </el-form>
    <span slot="footer" class="dialog-footer">
      <el-button size="small" @click="handleClose">
        {{ $t('common.button.cancel') }}
      </el-button>
      <el-button
        size="small"
        type="primary"
        :loading="submitLoading"
        @click="handleSubmit"
      >
        {{ $t('common.button.confirm') }}
      </el-button>
    </span>
  </el-dialog>
</template>

<script>
import Pagination from '@/components/pagination.vue';
import SearchInput from '@/components/searchInput.vue';
import { createApiKey, editApiKey } from '@/api/apiKeyManagement';
export default {
  components: { Pagination, SearchInput },
  data() {
    return {
      isEdit: false,
      form: {
        name: '',
        desc: '',
        expiredAt: '',
      },
      rules: {
        name: [
          {
            required: true,
            message: this.$t('common.input.placeholder'),
            trigger: 'change',
          },
          {
            pattern: /^[A-Za-z0-9.\u4e00-\u9fa5_-]+$/,
            message: this.$t('common.hint.nameHint'),
            trigger: 'change',
          },
        ],
      },
      pickerOptions: {
        disabledDate: time => {
          return time.getTime() < Date.now();
        },
      },
      dialogVisible: false,
      submitLoading: false,
    };
  },
  methods: {
    setFormValue(row) {
      const obj = { ...this.form };
      for (let key in obj) {
        obj[key] = row ? row[key] : '';
      }
      this.form = obj;
    },
    handleClose() {
      this.$refs.form.resetFields();
      this.dialogVisible = false;
    },
    openDialog(row = {}) {
      this.row = row;
      this.isEdit = Boolean(row.keyId);
      this.setFormValue(row);

      this.dialogVisible = true;
    },
    handleSubmit() {
      this.$refs.form.validate(async valid => {
        if (!valid) return;

        this.submitLoading = true;
        const params = { ...this.form };
        if (this.isEdit) params.keyId = this.row.keyId;
        try {
          const res = this.isEdit
            ? await editApiKey(params)
            : await createApiKey(params);
          if (res.code === 0) {
            this.$message.success(this.$t('common.message.success'));
            this.dialogVisible = false;
            this.$emit('reloadData');
          }
        } finally {
          this.submitLoading = false;
        }
      });
    },
  },
};
</script>

<style lang="scss" scoped></style>
