import request from "@/utils/request"
const BASE_URL = '/use/model/api/v1'

/*---工作流模板---*/
export const getWorkflowTempList = (data)=>{
    return request({
        url: `${BASE_URL}/mcp/square/list`,
        method: 'get',
        params: data
    })
};
export const getWorkflowTempInfo = (data)=>{
    return request({
        url: `${BASE_URL}/mcp/square`,
        method: 'get',
        params: data
    })
};
export const getWorkflowRecommendsList = (data)=>{
    return request({
        url: `${BASE_URL}/mcp/square/recommend`,
        method: 'get',
        params: data
    })
};