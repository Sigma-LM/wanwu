import service from '@/utils/request';
import { USER_API } from '@/utils/requestConstants';

// 获取列表
export const fetchApiKeyList = params => {
  return service({
    url: `${USER_API}/api/key/list`,
    method: 'get',
    params,
  });
};
// 创建
export const createApiKey = data => {
  return service({
    url: `${USER_API}/api/key`,
    method: 'post',
    data,
  });
};
// 编辑
export const editApiKey = data => {
  return service({
    url: `${USER_API}/api/key/update`,
    method: 'put',
    data,
  });
};
// 删除
export const deleteApiKey = data => {
  return service({
    url: `${USER_API}/api/key`,
    method: 'delete',
    data,
  });
};
// 修改状态
export const changeApiKeyStatus = data => {
  return service({
    url: `${USER_API}/api/key/status`,
    method: 'put',
    data,
  });
};
