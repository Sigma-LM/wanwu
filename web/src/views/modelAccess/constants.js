import { i18n } from '@/lang';

export const LLM = 'llm';
export const RERANK = 'rerank';
export const EMBEDDING = 'embedding';
export const OCR = 'ocr';
export const OCR_DS = 'ocr-deepseek';
export const OCR_PADDLE = 'ocr-paddle';
export const GUI = 'gui';
export const PDF_PARSER = 'pdf-parser';
export const ASR = 'asr';

export const MODEL_TYPE_OBJ = {
  [LLM]: 'LLM',
  [RERANK]: 'Rerank',
  [EMBEDDING]: 'Embedding',
  [OCR]: 'OCR',
  [GUI]: 'GUI',
  [PDF_PARSER]: i18n.t('modelAccess.type.pdfParser'),
  /*[OCR_DS]: 'OCR-DeepSeek',
  [OCR_PADDLE]: 'OCR-PaddleOCR',*/
  // [ASR]: i18n.t('modelAccess.type.asr')
};

export const MODEL_TYPE = Object.keys(MODEL_TYPE_OBJ).map(key => ({
  key,
  name: MODEL_TYPE_OBJ[key],
}));

export const YUAN_JING = 'YuanJing';
export const OPENAI_API = 'OpenAI-API-compatible';
export const OLLAMA = 'Ollama';
export const QWEN = 'Qwen';
export const HUOSHAN = 'HuoShan';
export const INFINI = 'Infini';
export const DEEPSEEK = 'DeepSeek';
export const QIANFAN = 'QianFan';

export const PROVIDER_OBJ = {
  [OPENAI_API]: 'OpenAI-API-compatible',
  [YUAN_JING]: i18n.t('modelAccess.type.yuanjing'),
  [OLLAMA]: 'Ollama',
  [QWEN]: i18n.t('modelAccess.type.qwen'),
  [HUOSHAN]: i18n.t('modelAccess.type.huoshan'),
  [INFINI]: i18n.t('modelAccess.type.infini'),
  [DEEPSEEK]: 'DeepSeek',
  [QIANFAN]: i18n.t('modelAccess.type.qianfan'),
};

export const PROVIDER_IMG_OBJ = {
  [OPENAI_API]: require('@/assets/imgs/openAI.png'),
  [YUAN_JING]: require('@/assets/imgs/yuanjing.png'),
  [OLLAMA]: require('@/assets/imgs/ollama.png'),
  [QWEN]: require('@/assets/imgs/qwen.png'),
  [HUOSHAN]: require('@/assets/imgs/volcano.png'),
  [INFINI]: require('@/assets/imgs/infini.png'),
  [DEEPSEEK]: require('@/assets/imgs/deepseek.png'),
  [QIANFAN]: require('@/assets/imgs/qianfan.png'),
};

const COMMON_MODEL_KEY = [LLM, RERANK, EMBEDDING];
const OLL_MODEL_KEY = [LLM, EMBEDDING];
export const PROVIDER_MODEL_KEY = {
  [OPENAI_API]: COMMON_MODEL_KEY,
  [YUAN_JING]: [...COMMON_MODEL_KEY, OCR, GUI, PDF_PARSER],
  [OLLAMA]: OLL_MODEL_KEY,
  [QWEN]: COMMON_MODEL_KEY,
  [HUOSHAN]: OLL_MODEL_KEY,
  [INFINI]: COMMON_MODEL_KEY,
  [DEEPSEEK]: [LLM],
  [QIANFAN]: [...COMMON_MODEL_KEY], // OCR_DS, OCR_PADDLE
};

export const PROVIDER_TYPE = Object.keys(PROVIDER_OBJ).map(key => {
  return {
    key,
    name: PROVIDER_OBJ[key],
    children: MODEL_TYPE.filter(item =>
      PROVIDER_MODEL_KEY[key]
        ? PROVIDER_MODEL_KEY[key].includes(item.key)
        : false,
    ),
  };
});

export const DEFAULT_CALLING = 'noSupport';
export const FUNC_CALLING = [
  { key: 'noSupport', name: i18n.t('modelAccess.noSupport') },
  { key: 'toolCall', name: 'Tool call' },
  /*{key: 'functionCall', name: 'Function call'},*/
];

export const DEFAULT_SUPPORT = 'noSupport';
export const SUPPORT_LIST = [
  { key: 'noSupport', name: i18n.t('modelAccess.noSupport') },
  { key: 'support', name: i18n.t('modelAccess.support') },
];

export const TYPE_OBJ = {
  apiKey: {
    [YUAN_JING]: 'sk-abc********************xyz',
    [OPENAI_API]: 'sk_7e4*************4s-BpI1l',
    [OLLAMA]: '',
    [QWEN]: 'sk-b************c70d',
    [HUOSHAN]: 'd8008ac0-****-****-****-**************',
    [INFINI]: 'sk-nw****gzjb6',
    [DEEPSEEK]: 'sk-14082***********************5e95',
    [QIANFAN]: 'bce-v3/ALTAK******82d1',
  },
  inferUrl: {
    [OCR]: 'https://maas-api.ai-yuanjing.com/openapi/v1',
    [GUI]: 'https://maas-api.ai-yuanjing.com/openapi/v1',
    [PDF_PARSER]: 'https://maas-api.ai-yuanjing.com/openapi/v1',
    [YUAN_JING]: 'https://maas.ai-yuanjing.com/openapi/compatible-mode/v1',
    [OPENAI_API]: 'https://api.siliconflow.cn/v1',
    [OLLAMA]: 'https://192.168.21.100:11434',
    [QWEN]: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
    [HUOSHAN]: 'https://ark.cn-beijing.volces.com/api/v3',
    [INFINI]: 'https://cloud.infini-ai.com/maas/v1',
    [DEEPSEEK]: 'https://api.deepseek.com/v1',
    [QIANFAN]: 'https://qianfan.baidubce.com/v2',
    [ASR]:
      'https://maas-api.ai-yuanjing.com/openapi/synchronous/asr/audio/file/transfer/unicom/sync/file/asr',
  },
};
