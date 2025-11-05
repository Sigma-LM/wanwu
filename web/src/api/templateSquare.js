import request from "@/utils/request"
const BASE_URL = '/user/api/v1'

/*---工作流模板---*/
export const getWorkflowTempList = (data)=>{
    return request({
        url: `${BASE_URL}/workflow/template/list`,
        method: 'get',
        params: data
    })
};
export const getWorkflowTempInfo = (data)=>{
    return request({
        url: `${BASE_URL}/workflow/template/detail`,
        method: 'get',
        params: data
    })
};
export const getWorkflowRecommendsList = (data)=>{
    return request({
        url: `${BASE_URL}/workflow/template/recommend`,
        method: 'get',
        params: data
    })
};
export const downloadWorkflow = (params) => {
    return request({
        url: `${BASE_URL}/workflow/template/download`,
        method: "get",
        params,
        responseType: 'blob'
    });
};
export const copyWorkflowTemplate = (data)=>{
    return request({
        url: `${BASE_URL}/workflow/template`,
        method: 'post',
        data
    })
};

/*---提示词模板---*/
export const getPromptTempList = (data)=>{
    return request({
        url: `${BASE_URL}/prompt/template/list`,
        method: 'get',
        params: data
    })
};

export const copyPromptTemplate = (data)=>{
    return request({
        url: `${BASE_URL}/prompt/template`,
        method: 'post',
        data
    })
};

/*---自定义提示词---*/
export const getCustomPromptList = (data)=>{
    return request({
        url: `${BASE_URL}/prompt/custom/list`,
        method: 'get',
        params: data
    })
};

export const createCustomPrompt = (data)=>{
    return request({
        url: `${BASE_URL}/prompt/custom`,
        method: 'post',
        data
    })
};

export const editCustomPrompt = (data)=>{
    return request({
        url: `${BASE_URL}/prompt/custom`,
        method: 'put',
        data
    })
};

export const copyCustomPrompt = (data)=>{
    return request({
        url: `${BASE_URL}/prompt/custom/copy`,
        method: 'post',
        data
    })
};

export const deleteCustomPrompt = (data)=>{
    return request({
        url: `${BASE_URL}/prompt/custom`,
        method: 'delete',
        data
    })
};