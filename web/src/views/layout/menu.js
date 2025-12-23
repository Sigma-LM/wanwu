import { PERMS } from '@/router/permission';
import { i18n } from '@/lang';

/**
 *  关于菜单中 line 的 perm: 控制 line 是否展示
 *  1. 正常 line 由其上方的菜单权限控制，若 line 上方菜单有一个则展示
 *  2. 最后的 line 由上下两方的菜单控制，因为最后的 line 若下方菜单无权限则最后一个 line 也不会展示
 *  3. 特殊情况如：模板广场 templateSquare 无权限，常显，所以最后一个 line 目前不受上方权限的控制，
 *  只受下方菜单的控制，若不是最后的 line，则无需配置 perm 即表示 line 常显
 *  加菜单时根据以上规则配置 line 的权限
 */
export const menuList = [
  {
    name: i18n.t('menu.modelAccess'),
    key: 'modelAccess',
    img: require('@/assets/imgs/model.svg'),
    imgActive: require('@/assets/imgs/model_active.svg'),
    path: '/modelAccess',
    perm: PERMS.MODEL,
  },
  {
    name: i18n.t('menu.knowledge'),
    key: 'knowledge',
    img: require('@/assets/imgs/knowledge.svg'),
    imgActive: require('@/assets/imgs/knowledge_active.svg'),
    path: '/knowledge',
    perm: PERMS.KNOWLEDGE,
  },
  {
    name: i18n.t('menu.tool'),
    key: 'tool',
    img: require('@/assets/imgs/tool.svg'),
    imgActive: require('@/assets/imgs/tool_active.svg'),
    path: '/tool',
    perm: PERMS.TOOL,
  },
  {
    name: i18n.t('menu.safetyGuard'),
    key: 'safetyGuard',
    img: require('@/assets/imgs/safety.svg'),
    imgActive: require('@/assets/imgs/safety_active.svg'),
    path: '/safety',
    perm: PERMS.SAFETY,
  },
  {
    key: 'line',
    perm: [PERMS.MODEL, PERMS.KNOWLEDGE, PERMS.TOOL, PERMS.SAFETY],
  },
  {
    name: i18n.t('menu.app.rag'),
    key: 'rag',
    img: require('@/assets/imgs/rag.svg'),
    imgActive: require('@/assets/imgs/rag_active.svg'),
    path: '/appSpace/rag',
    perm: PERMS.RAG,
  },
  {
    name: i18n.t('menu.app.workflow'),
    key: 'workflow',
    img: require('@/assets/imgs/workflow_icon.svg'),
    imgActive: require('@/assets/imgs/workflow_icon_active.svg'),
    path: '/appSpace/workflow',
    perm: PERMS.WORKFLOW,
  },
  {
    name: i18n.t('menu.app.agent'),
    key: 'agent',
    img: require('@/assets/imgs/agent.svg'),
    imgActive: require('@/assets/imgs/agent_active.svg'),
    path: '/appSpace/agent',
    perm: PERMS.AGENT,
  },
  {
    key: 'line',
    perm: [PERMS.RAG, PERMS.WORKFLOW, PERMS.AGENT],
  },
  {
    name: i18n.t('menu.mcp'),
    key: 'mcpManage',
    img: require('@/assets/imgs/mcp_menu.svg'),
    imgActive: require('@/assets/imgs/mcp_menu_active.svg'),
    path: '/mcp',
    perm: PERMS.MCP,
  },
  {
    name: i18n.t('menu.explore'),
    key: 'explore',
    img: require('@/assets/imgs/explore.svg'),
    imgActive: require('@/assets/imgs/explore_active.svg'),
    path: '/explore',
    perm: PERMS.EXPLORE,
  },
  {
    name: i18n.t('menu.templateSquare'),
    key: 'templateSquare',
    img: require('@/assets/imgs/template_square.svg'),
    imgActive: require('@/assets/imgs/template_square_active.svg'),
    path: '/templateSquare',
  },
  {
    key: 'line',
    perm: [PERMS.API_KEY],
  },
  {
    name: 'API Key',
    key: 'openApiKey',
    img: require('@/assets/imgs/api_key_management.svg'),
    imgActive: require('@/assets/imgs/api_key_management_active.svg'),
    path: '/openApiKey',
    perm: PERMS.API_KEY,
  },
];
