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
        <el-radio-group
          @change="changeSelectRadio"
          v-model="radio"
          style="width: 100%"
        >
          <el-radio label="0">{{ $t('apiKeyManage.table.custom') }}</el-radio>
          <el-radio label="1">
            {{ $t('apiKeyManage.table.permanent') }}
          </el-radio>
        </el-radio-group>
        <div v-if="radio === '0'">
          <el-date-picker
            style="width: 100%"
            type="date"
            v-model="form.expiredAt"
            value-format="yyyy-MM-dd"
            :placeholder="$t('common.select.placeholder')"
            :picker-options="pickerOptions"
          ></el-date-picker>
        </div>
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

const PERMANENT_TIME = '9999-12-31';

export default {
  components: { Pagination, SearchInput },
  data() {
    return {
      isEdit: false,
      radio: '0',
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
        expiredAt: [
          {
            required: true,
            message: this.$t('common.select.placeholder'),
            trigger: 'blur',
          },
        ],
      },
      pickerOptions: {
        disabledDate: time => {
          return time.getTime() < Date.now();
        },
        shortcuts: [
          {
            text: '一周',
            onClick(picker) {
              const date = new Date();
              date.setTime(date.getTime() + 3600 * 1000 * 24 * 7);
              picker.$emit('pick', date);
            },
          },
          {
            text: '一个月',
            onClick(picker) {
              const date = new Date();
              date.setTime(date.getTime() + 3600 * 1000 * 24 * 30);
              picker.$emit('pick', date);
            },
          },
          {
            text: '三个月',
            onClick(picker) {
              const date = new Date();
              date.setTime(date.getTime() + 3600 * 1000 * 24 * 90);
              picker.$emit('pick', date);
            },
          },
        ],
      },
      dialogVisible: false,
      submitLoading: false,
    };
  },
  methods: {
    changeSelectRadio(value) {
      if (value === '0') {
        this.form.expiredAt = '';
      } else {
        this.form.expiredAt = PERMANENT_TIME;
        this.$refs.form.clearValidate('expiredAt');
      }
    },
    setFormValue(row) {
      const obj = { ...this.form };
      for (let key in obj) {
        obj[key] = row ? row[key] : '';
      }
      this.form = obj;
    },
    handleClose() {
      this.setFormValue();
      this.$refs.form.resetFields();
      this.dialogVisible = false;
    },
    openDialog(row) {
      this.row = row;
      this.isEdit = Boolean(row && row.keyId);
      this.radio = row ? (row.expiredAt ? '0' : '1') : '0';
      this.setFormValue(
        row ? { ...row, expiredAt: row.expiredAt || PERMANENT_TIME } : row,
      );

      this.dialogVisible = true;
    },
    handleSubmit() {
      this.$refs.form.validate(async valid => {
        if (!valid) return;

        this.submitLoading = true;
        const params = {
          ...this.form,
          expiredAt: this.radio === '0' ? this.form.expiredAt : '',
        };
        if (this.isEdit) params.keyId = this.row.keyId;
        try {
          const res = this.isEdit
            ? await editApiKey(params)
            : await createApiKey(params);
          if (res.code === 0) {
            this.$message.success(this.$t('common.message.success'));
            this.handleClose();
            this.$emit('reloadData', this.isEdit ? {} : { pageNo: 1 });
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
