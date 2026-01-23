<template>
  <el-dialog
    :title="
      (form.externalApiId
        ? $t('common.button.edit')
        : $t('common.button.add')) + $t('knowledgeManage.externalAPI.title')
    "
    :visible.sync="visible"
    :before-close="handleClose"
    append-to-body
    width="50%"
  >
    <el-form :model="form" :rules="rules" ref="apiForm" label-width="120px">
      <el-form-item :label="$t('knowledgeManage.externalAPI.name')" prop="name">
        <el-input
          v-model="form.name"
          :placeholder="$t('common.input.placeholder')"
        ></el-input>
      </el-form-item>
      <el-form-item
        :label="$t('knowledgeManage.externalAPI.desc')"
        prop="description"
      >
        <el-input
          v-model="form.description"
          :placeholder="$t('common.input.placeholder')"
        ></el-input>
      </el-form-item>
      <el-row>
        <el-form-item label="Base Url" prop="baseUrl">
          <el-input
            v-model="form.baseUrl"
            :placeholder="$t('common.input.placeholder')"
          ></el-input>
        </el-form-item>
        <el-form-item label="API Key" prop="apiKey">
          <el-input
            v-model="form.apiKey"
            :placeholder="$t('common.input.placeholder')"
          ></el-input>
        </el-form-item>
      </el-row>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleClose">
        {{ $t('common.button.cancel') }}
      </el-button>
      <el-button type="primary" @click="submitForm">
        {{ $t('common.button.confirm') }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script>
import { importExternalAPI, editExternalAPI } from '@/api/knowledge';

export default {
  data() {
    return {
      visible: false,
      form: {
        externalApiId: '',
        name: '',
        description: '',
        baseUrl: '',
        apiKey: '',
      },
      rules: {
        name: [
          {
            required: true,
            message: this.$t('common.input.placeholder'),
            trigger: 'blur',
          },
        ],
        description: [
          {
            required: true,
            message: this.$t('common.input.placeholder'),
            trigger: 'blur',
          },
        ],
        baseUrl: [
          {
            required: true,
            message: this.$t('common.input.placeholder'),
            trigger: 'blur',
          },
        ],
        apiKey: [
          {
            required: true,
            message: this.$t('common.input.placeholder'),
            trigger: 'blur',
          },
        ],
      },
    };
  },
  methods: {
    showDialog() {
      this.visible = true;
    },
    handleClose() {
      this.visible = false;
      this.form = {
        externalApiId: '',
        name: '',
        description: '',
        baseUrl: '',
        apiKey: '',
      };
      this.$refs.apiForm.resetFields();
      this.$refs.apiForm.clearValidate();
    },
    submitForm() {
      this.$refs.apiForm.validate(valid => {
        if (valid) {
          if (this.form.externalApiId) {
            editExternalAPI({ ...this.form }).then(res => {
              if (res.code === 0) {
                this.$message.success(this.$t('common.info.edit'));
                this.$emit('confirm', this.form);
                this.handleClose();
              }
            });
          } else
            importExternalAPI({ ...this.form }).then(res => {
              if (res.code === 0) {
                this.$message.success(this.$t('common.info.create'));
                this.$emit('confirm', this.form);
                this.handleClose();
              }
            });
        } else {
          return false;
        }
      });
    },
  },
};
</script>

<style scoped lang="scss"></style>
