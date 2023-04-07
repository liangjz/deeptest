/**
 * 可参考：https://json-schema.apifox.cn/#%E6%95%B0%E6%8D%AE%E7%B1%BB%E5%9E%8B
 * */
export const JSONSchemaDataTypes = [
    {
        label: "string",
        value: "string",
        color: 'pink',
        active: false,
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'enum',
                    name: 'enum',
                    component: 'selectTag',
                    type: 'array',
                    placeholder: '输入文本后按回车添加',
                    value: [],
                },
                {
                    label: 'format',
                    name: 'format',
                    type: 'string',
                    component: 'select',
                    placeholder: 'select a value',
                    value: null,
                    options: [
                        {
                            label: 'data-time',
                            value: 'data-time',
                        },
                        {
                            label: 'time',
                            value: 'time',
                        },
                        {
                            label: 'email',
                            value: 'email',
                        },
                        {
                            label: 'idn-email',
                            value: 'idn-email',
                        },
                        {
                            label: 'hostname',
                            value: 'hostname',
                        },
                        {
                            label: 'idn-hostname',
                            value: 'idn-hostname',
                        },
                        {
                            label: 'ipv4',
                            value: 'ipv4',
                        },
                        {
                            label: 'ipv6',
                            value: 'ipv6',
                        },
                        {
                            label: 'uri',
                            value: 'uri',
                        },
                        {
                            label: 'uri-reference',
                            value: 'uri-reference',
                        },
                        {
                            label: 'iri',
                            value: 'iri',
                        },
                        {
                            label: 'iri-reference',
                            value: 'iri-reference',
                        },
                        {
                            label: 'uri-template',
                            value: 'uri-template',
                        },
                        {
                            label: 'json-pointer',
                            value: 'json-pointer',
                        },
                        {
                            label: 'regex',
                            value: 'regex',
                        },
                        {
                            label: 'uuid',
                            value: 'uuid',
                        },
                        {
                            label: 'password',
                            value: 'password',
                        },
                        {
                            label: 'byte',
                            value: 'byte',
                        },
                    ],
                },
                {
                    label: 'behavior',
                    name: 'behavior',
                    type: 'boolean',
                    component: 'select',
                    placeholder: 'select a value',
                    value: null,
                    options: [
                        {
                            label: 'Read/Write',
                            value: 'Read/Write',
                        },
                        {
                            label: 'Read Only',
                            value: 'Read Only',
                        },
                        {
                            label: 'Write Only',
                            value: 'Write Only',
                        },
                    ]
                },
                {
                    label: 'default',
                    name: 'default',
                    component: 'input',
                    placeholder: 'default',
                    type: 'string',
                    value: '',
                },
                {
                    label: 'example',
                    name: 'example',
                    type: 'string',
                    component: 'input',
                    placeholder: 'example',
                    value: '',
                },
                {
                    label: 'pattern',
                    name: 'pattern',
                    type: 'string',
                    component: 'input',
                    placeholder: 'pattern',
                    value: '',
                },
                {
                    label: 'minLength',
                    name: 'minLength',
                    component: 'inputNumber',
                    placeholder: '>=0',
                    type: 'integer',
                    value: '',
                },
                {
                    label: 'maxLength',
                    name: 'maxLength',
                    type: 'integer',
                    component: 'inputNumber',
                    placeholder: '>=0',
                    value: '',
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                    value: false,
                },
            ]
        }
    },
    {
        label: "number",
        value: "number",
        color: 'cyan',
        active: false,
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'enum',
                    name: 'enum',
                    component: 'selectTag',
                    type: 'array',
                    placeholder: '输入文本后按回车添加',
                    value: [],
                },
                {
                    label: 'format',
                    name: 'format',
                    type: 'string',
                    component: 'select',
                    placeholder: 'select a value',
                    value: null,
                    options: [
                        {
                            label: 'float',
                            value: 'float',
                        },
                        {
                            label: 'double',
                            value: 'double',
                        },
                    ]
                },
                {
                    label: 'behavior',
                    name: 'behavior',
                    type: 'boolean',
                    component: 'select',
                    placeholder: 'select a value',
                    value: null,
                    options: [
                        {
                            label: 'Read/Write',
                            value: 'Read/Write',
                        },
                        {
                            label: 'Read Only',
                            value: 'Read Only',
                        },
                        {
                            label: 'Write Only',
                            value: 'Write Only',
                        },
                    ]
                },
                {
                    label: 'default',
                    name: 'default',
                    component: 'input',
                    placeholder: 'default',
                    type: 'string',
                    value: '',
                },
                {
                    label: 'example',
                    name: 'example',
                    type: 'string',
                    component: 'input',
                    placeholder: 'example',
                    value: '',
                },
                {
                    label: 'minimum',
                    name: 'minimum',
                    type: 'number',
                    component: 'inputNumber',
                    placeholder: '>=0',
                    value: '',
                },
                {
                    label: 'maximum',
                    name: 'maximum',
                    type: 'number',
                    value: '',
                    component: 'inputNumber',
                    placeholder: '>=0',
                },
                {
                    label: 'maxLength',
                    name: 'maxLength',
                    type: 'integer',
                    component: 'inputNumber',
                    placeholder: '>=0',
                    value: '',
                },
                {
                    label: 'multipleOf',
                    name: 'multipleOf',
                    type: 'number',
                    component: 'inputNumber',
                    placeholder: '>=0',
                    value: '',
                },
                {
                    label: 'exclusiveMin',
                    name: 'exclusiveMin',
                    type: 'boolean',
                    component: 'switch',
                    value: false,
                },
                {
                    label: 'exclusiveMax',
                    name: 'exclusiveMax',
                    type: 'boolean',
                    component: 'switch',
                    value: false,
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                    value: false,
                },
            ]
        }
    },
    {
        label: "integer",
        value: "integer",
        color: 'green',
        active: false,
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'enum',
                    name: 'enum',
                    component: 'selectTag',
                    type: 'array',
                    placeholder: '输入文本后按回车添加',
                    value: [],
                },
                {
                    label: 'format',
                    name: 'format',
                    type: 'string',
                    component: 'select',
                    placeholder: 'select a value',
                    value: null,
                    options: [
                        {
                            label: 'int32',
                            value: 'int32',
                        },
                        {
                            label: 'int64',
                            value: 'int64',
                        },
                    ]
                },
                {
                    label: 'behavior',
                    name: 'behavior',
                    type: 'boolean',
                    component: 'select',
                    placeholder: 'select a value',
                    value: null,
                    options: [
                        {
                            label: 'Read/Write',
                            value: 'Read/Write',
                        },
                        {
                            label: 'Read Only',
                            value: 'Read Only',
                        },
                        {
                            label: 'Write Only',
                            value: 'Write Only',
                        },
                    ]
                },
                {
                    label: 'default',
                    name: 'default',
                    component: 'input',
                    placeholder: 'default',
                    type: 'string',
                    value: '',
                },
                {
                    label: 'example',
                    name: 'example',
                    type: 'string',
                    component: 'input',
                    placeholder: 'example',
                    value: '',
                },
                {
                    label: 'minimum',
                    name: 'minimum',
                    type: 'number',
                    component: 'inputNumber',
                    placeholder: '>=0',
                    value: '',
                },
                {
                    label: 'maximum',
                    name: 'maximum',
                    type: 'number',
                    value: '',
                    component: 'inputNumber',
                    placeholder: '>=0',
                },
                {
                    label: 'maxLength',
                    name: 'maxLength',
                    type: 'integer',
                    component: 'inputNumber',
                    placeholder: '>=0',
                    value: '',
                },
                {
                    label: 'multipleOf',
                    name: 'multipleOf',
                    type: 'number',
                    component: 'inputNumber',
                    placeholder: '>=0',
                    value: '',
                },
                {
                    label: 'exclusiveMin',
                    name: 'exclusiveMin',
                    type: 'boolean',
                    component: 'switch',
                    value: false,
                },
                {
                    label: 'exclusiveMax',
                    name: 'exclusiveMax',
                    type: 'boolean',
                    component: 'switch',
                    value: false,
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                    value: false,
                },
            ]
        }
    },
    {
        label: "object",
        value: "object",
        color: 'blue',
        active: false,
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'minProperties',
                    name: 'minProperties',
                    type: 'integer',
                    placeholder: '>=0',
                    component: 'inputNumber',
                    value: null,
                },
                {
                    label: 'maxProperties',
                    name: 'maxProperties',
                    type: 'integer',
                    component: 'inputNumber',
                    placeholder: '>=0',
                    value: null,
                },
                {
                    label: 'allow additional Properties',
                    name: 'additionalProperties',
                    type: 'boolean',
                    component: 'switch',
                    value: false,
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                    value: false,
                },
            ]
        },
    },
    {
        label: "array",
        value: "array",
        color: 'orange',
        active: false,
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'minItems',
                    name: 'minItems',
                    type: 'integer',
                    placeholder: '>=0',
                    component: 'inputNumber',
                    value: null,
                },
                {
                    label: 'maxItems',
                    name: 'maxItems',
                    type: 'integer',
                    placeholder: '>=0',
                    component: 'inputNumber',
                    value: null,
                },
                {
                    label: 'uniqueItems',
                    name: 'additionalProperties',
                    component: 'switch',
                    type: 'boolean',
                    value: false,
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                    value: false,
                },
            ],
            subTypes: []
        },
    },
    {
        label: "boolean",
        value: "boolean",
        color: 'red',
        active: false,
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'behavior',
                    name: 'behavior',
                    type: 'string',
                    component: 'select',
                    value: '',
                    options: [
                        {
                            label: 'Read/Write',
                            value: 'Read/Write',
                        },
                        {
                            label: 'Read Only',
                            value: 'Read Only',
                        },
                        {
                            label: 'Write Only',
                            value: 'Write Only',
                        },
                    ]
                },
                {
                    label: 'default',
                    name: 'default',
                    type: 'boolean',
                    component: 'select',
                    placeholder: 'select a value',
                    value: '',
                    options: [
                        {
                            label: 'true',
                            value: true,
                        },
                        {
                            label: 'false',
                            value: false,
                        },
                    ]
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                    value: false,
                },
            ],
        }
    },
];
/**
 * 设置schema模块数据类型
 * */
export const schemaSettingInfo = [
    {
        label: 'Type',
        subLabel: 'SubType',
        type: 'type',
        value: 'string',
        active: true,
        props: JSONSchemaDataTypes
    },
    {
        label: 'Components',
        type: '$ref',
        value: '',
        active: false,
        subLabel: 'Components',
    },
]
export const typeOpts = ['string', 'number', 'boolean', 'array', 'object', 'integer'];
// 树形结构的层级递进宽度
export const treeLevelWidth = 24;

/**
 * 是否是对象类型
 * */
export function isObject(type: string) {
    return type === 'object';
}

/**
 * 是否是数组类型
 * */
export function isArray(type: string) {
    return type === 'array';
}


/**
 * 根据传入的 schema 结构信息，添加需要额外的渲染属性
 * */
export function addExtraInfo(val: any) {
    if (!val) {
        return null
    }
    val.extraViewInfo = {
        "isExpand": true,
        "name": "root",
        "depth": 1,
        "type": val.type,
        "parent": null,
    };
    function fn(obj: any, depth: number) {
        if (obj.properties && obj.type === 'object') {
            Object.entries(obj.properties).forEach(([key, value]: any) => {
                value.extraViewInfo = {
                    "isExpand": true,
                    "name": key,
                    "depth": depth,
                    "type": obj.type,
                    "parent": obj,
                }
                if (value.type === 'object') {
                    fn(value, depth + 1);
                }
            })
        }
        if (obj.type === 'array' && obj.items) {
            Object.entries(obj?.items?.properties).forEach(([key, value]: any) => {
                value.extraViewInfo = {
                    "isExpand": true,
                    "name": key,
                    "depth": depth,
                    "type": obj.type,
                    "items": obj?.items,
                    "parent": obj,
                }
                if (value.type === 'object') {
                    fn(value, depth + 1);
                }
            })
        }
    }
    fn(val, 2);
    console.log(222, val);
    return val;
}

/**
 * 根据传入的 schema 结构信息，删除额外的渲染属性
 * */
export function removeExtraInfo(val: any): object | null {
    if (!val) {
        return null
    }
    if (val.extraViewInfo) {
        delete val.extraViewInfo;
    }

    function fn(obj: any) {
        if (obj.properties && obj.type === 'object') {
            Object.entries(obj.properties).forEach(([key, value]: any) => {
                if (value.extraViewInfo) {
                    delete value.extraViewInfo;
                }
                if (value.type === 'object') {
                    fn(value);
                }
            })
        }
    }

    fn(val);
    return val;
}

const obj = {
    "type": "object",
    "required": [
        "name"
    ],
    "properties": {
        "name": {
            "type": "string",
            "example": "8322222",
            "enum": ["8322222", '122', '2222'],
        },
        "age": {
            "type": "integer",
            "format": "int32",
            "minimum": 0,
        },
        'obj1': {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name1": {
                    "type": "string",
                },
                "age1": {
                    "type": "integer",
                    "format": "int32",
                    "minimum": 110,
                    "example": "8322222",
                    "enum": ["8322222", '122', '2222'],
                },
                'obj2': {
                    "type": "object",
                    "required": [
                        "name11"
                    ],
                    "properties": {
                        "name3232": {
                            "type": "string",
                        },
                        "age3333": {
                            "type": "integer",
                            "format": "int32",
                            "minimum": 0,

                        },
                        'obj4341': {
                            "type": "object",
                            "required": [
                                "name"
                            ],
                            "properties": {
                                "name1332323": {
                                    "type": "string",
                                },
                                "arr1": {
                                    "type": "array",
                                    "format": "int32",
                                    minItems: 13,
                                    "minimum": 0,
                                    "items": {
                                        "type": "string",
                                    }
                                },
                            }
                        }
                    }
                }
            }
        }
    }
};
const obj1 = {
    "type": "array",
    "items": {
        "type": "object",
        "required": [
            "name11"
        ],
        "properties": {
            "name3232": {
                "type": "string",
            },
            "age3333": {
                "type": "integer",
                "format": "int32",
                "minimum": 0,
            },
            'obj4341': {
                "type": "object",
                "required": [
                    "name"
                ],
                "properties": {
                    "name1332323": {
                        "type": "string",
                    },
                    "arr1": {
                        "type": "array",
                        "format": "int32",
                        minItems: 13,
                        "minimum": 0,
                        "items": {
                            "type": "string",
                        }
                    },
                }
            }
        }
    }
};