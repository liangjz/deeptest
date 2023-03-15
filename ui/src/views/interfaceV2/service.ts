import request from '@/utils/request';
import {requestToAgent} from '@/utils/request';
import {
    ApiKey,
    BasicAuth, BearerToken,
    BodyFormDataItem,
    BodyFormUrlEncodedItem, Checkpoint, Extractor,
    Header,
    Interface,
    OAuth20,
    Param
} from "@/views/interface/data";
import {isInArray} from "@/utils/array";
import {UsedBy} from "@/utils/enum";

const apiPath = 'interfaces';
const apiImport = 'import';
const apiSpec = 'spec';
const apiInvocation = 'invocations';
const apiAuth = 'auth';
const apiEnvironment = 'environments'
const apiEnvironmentVar = `${apiEnvironment}/vars`
const apiShareVar = `${apiEnvironment}/shareVars`
const apiSnippets = 'snippets'

const apiExtractor = 'extractors'
const apiCheckpoint = 'checkpoints'

const apiParser = 'parser'


interface InterfaceListReqParams {
    "prjectId"?: number,
    "page"?: number,
    "pageSize"?: number,
    "status"?: number,
    "userId"?: number,
    "title"?: string
}

//todo liguwe 待整理
interface SaveInterfaceReqParams {
    // project_id?: number,
    serveId?: number,
    title?: string,
    path?: string
}


/**
 * 接口列表
 * */
export async function getInterfaceList(data: any): Promise<any> {
    return request({
        url: `/endpoint/index`,
        method: 'post',
        data: data
    });
}

/**
 * 接口详情
 * */
export async function getInterfaceDetail(id: Number | String | any): Promise<any> {
    return request({
        url: `/endpoint/detail?id=${id}`,
        method: 'get',
    });
}

/**
 * 删除接口
 * */
export async function deleteInterface(id: Number): Promise<any> {
    return request({
        url: `/endpoint/delete?id=${id}`,
        method: 'delete',
    });
}


/**
 * 复制接口
 * */
export async function copyInterface(id: Number): Promise<any> {
    return request({
        url: `/endpoint/copy?id=${id}`,
        method: 'get',
    });
}


/**
 * 获取yaml展示
 * */
export async function getYaml(data: any): Promise<any> {
    return request({
        url: `/endpoint/yaml`,
        method: 'post',
        data: data
    });
}


/**
 * 接口过时
 * */
export async function expireInterface(id: Number): Promise<any> {
    return request({
        url: `/endpoint/expire?id=${id}`,
        method: 'put',
    });
}

/**
 * 保存接口
 * */
export async function saveInterface(data: SaveInterfaceReqParams): Promise<any> {
    return request({
        url: `/endpoint/save`,
        method: 'post',
        data: data
    });
}

/**
 * 创建分类
 * */
export async function newCategories(data: any): Promise<any> {
    return request({
        url: `/scenarios/categories`,
        method: 'post',
        data: data
    });
}

/**
 * 修改分类
 * */
export async function editCategories(data: any): Promise<any> {
    return request({
        url: `scenarios/categories/${data.id}/updateName`,
        method: 'put',
        data: data
    });
}

/**
 * 删除分类
 * */
export async function deleteCategories(data: any): Promise<any> {
    return request({
        url: `/scenarios/categories/${data.id}`,
        method: 'put',
    });
}

/**
 * 移动分类
 * */
export async function moveCategories(data: any): Promise<any> {
    return request({
        url: `/scenarios/categories/move`,
        method: 'post',
        data: data
    });
}

/**
 * 获取分类
 * */
export async function getCategories(params: any): Promise<any> {
    return request({
        url: `/scenarios/categories/load`,
        method: 'get',
        params,
    });
}