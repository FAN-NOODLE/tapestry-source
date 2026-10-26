# download


## 如何使用id分段来解决深度分页问题

首先介绍的是我们扩展的字段信息
```go

// 这个是我们的一些默认信息
var (
	defaultRequestSortField           = "sort_key"              // 默认的排序接收字段，就是我们下载中心反查micro-scrm时把排序字段放到那个field，micro-scrm拿到这个field从中解析真正要执行的key，下面好多都类似，我们就介绍前几个，下面的大家看名字就知道了
	defaultSortKey                    = "id"                    // 默认的排序字段，下载中心在查micro-scrm时实际使用到的排序字段，按现在这两行来说就是，micro-scrm拿到http请求的json数据中有sort_key这个field，里面的value时id，当然这是默认情况下的，下面会给大家说如何修改，目的是为了方便兼容大家的接口，减少冲突
	defaultRequestSortTypeField       = "sort_type"             // 默认的排序方式的接收字段
	defaultSortType                   = "ASC"                   // 默认使用的排序方式，当然你也可以在micro-scrm写死你的排序方式，就不用在这里写了，但是为了明确我们的接口是实现了排序，推荐大家在调用的时候给写上
	defaultRequestSelectField         = "will_select"           // 我们好多场景下其实是单表查询数据或者就是简单的关联，但是micro-scrm中确实全部数据返回，我们这块给定了一个预期，micro-scrm可以根据这个值来确定自己需要查询返回的数据，可以提高效率，减少网络io，这是一个可选项，当然这里只要是field结尾的都是在设置要回传给micro-scrm的字段
	defaultRequestScrollPageInfoField = "will_scroll_page_info" // 这个是一个状态字段，也是要回传给micro-scrm的，如果这个字段的值为true，证明我们在分页，我们要在业务代码里做专门的处理，最好是只查询且只返回一个id，这样能减少网络io，当然你可以按自己的想法来做
	defaultRequestMinID               = "will_min_id" // id分段的话，我们会根据页大小进行切片，划分出一个个id区间，这是那个最大值的字段，他不可修改，你们要再业务代码里进行兼容
	defaultRequestMaxID               = "will_max_id" // 与上同理
	defaultSplitSize                  = common.DefaultSplitSize // 这个是我们的默认一页数据大小，默认跟普通模式是一样的，后续可以调整
)

// 下面的时我们创建下载任务时可以指定的一些信息，当然你也可以通过上面的默认信息进行实现
type SourceRestfulIDPaginationType struct {
	SourceRestfulPaginationType      int8     `json:"source_restful_pagination_type"`       // =2代表是id分段，其他就是默认分页方式
	SourceRestfulReceiveSelectSwitch bool     `json:"source_restful_receive_select_switch"` // true代表我们要开启select的回传，我们的业务代码要对其进行优化，这是可选的
	SourceRestfulReceiveSelectField  string   `json:"source_restful_receive_select_field"`  // 我这个是我们接收所需数据区域的接受字段的
	SourceRestfulSortField           string   `json:"source_restful_sort_field"`            // 这个是我们修改默认分类接收字段的
	SourceRestfulSortKeyAs           string   `json:"source_restful_sort_key_as"`           // 我们的排序字段在结构体中的名字，这个不同于我们的sort_key，因为我们可能会使用到表关联，这种情况下，可能我们的sort_key 要传一个 a.id ，但是映射到结构体的名字是 id
	SourceRestfulSortKey             string   `json:"source_restful_sort_key"` // 我们创建下载任务时候指定我们应该使用的排序id
	SourceRestfulSortTypeField       string   `json:"source_restful_sort_type_field"` // 这个是我们的分类方式接收字段
	SourceRestfulSortType            string   `json:"source_restful_sort_type"` // 和上面同理
	SourceRestfulScrollPageInfoField string   `json:"source_restful_scroll_page_info_field"` // 修改我们的分页获取状态的接受字段
	SourceRestfulSelectInfo          []string `json:"source_restful_select_info"`            // 这个是我们可以在创建下载中心时指定我们需要获取的数据区域，如果表关联，应该与业务代码对应，例如 []string{"a.id", "b.nickname"} 这样形式，如果是单表查询，我们就更简单了
	SourceRestfulMinIDField          string   `json:"source_restful_min_id_field"` // 最小值的接收字段
	SourceRestfulMaxIDField          string   `json:"source_restful_max_id_field"` // 最大值的接收字段
}
```

## 下面提供两个案例，对应的业务代码去micro-scrm库里面查看

### 1。 简单的单表查询
```json
请求地址：http://localhost:8080/download/v1/create-export

请求body：

{
    "filename": "标签预处理导出",
    "type": 1,
    "store_id": 1,
    "enterprise_id": 192104618298240,
    "operation_id": 192104618632064,
    "operation_name": "gongguoxin@vchangyi.com@qq.com",
    "source_type": 3,
    "source_title_alias": {
        "id": "序号",
        "tag_name": "标签名称",
        "union_id": "用户unionid",
        "user_id": "导购user_id"
    },
    "source_title": [
        "id",
        "tag_name",
        "union_id",
        "user_id"
    ],
    "source_header_prefix": null,
    "label": "tag-export",
    "source_data": null,
    "source_sql_database": "",
    "source_sql_count": 0,
    "source_sql": "",
    "source_restful_path": "https://ws-dev.vchangyi.com/micro-scrm/micro-scrm/v1/tag/future/export",
    "source_restful_method": "POST",
    "source_restful_headers": {},
    "source_restful_params": {
        "enterprise_id": 192104618298240,
        "file_id": 3633
    },
    "source_suffix": null,
    "source_restful_count_field": "data.total",
    "source_restful_data_field": "data.list",
    "source_restful_request_page_field": "page",
    "source_restful_request_size_field": "size",
    "source_restful_receive_select_switch": true,
    "source_restful_sort_key": "id",
    "source_restful_sort_type": "DESC",
    "source_restful_select_info": [],
    "source_excel_cell_merge": {
        "type": 0,
        "rule": "",
        "title": null,
        "premise": null
    },
    "source_restful_pagination_type": 2
}
```

### 2. 内部有表关联的形式
```json
{
    "filename": "标签预处理导出",
    "type": 1,
    "store_id": 1,
    "enterprise_id": 192104618298240,
    "operation_id": 192104618632064,
    "operation_name": "gongguoxin@vchangyi.com@qq.com",
    "source_type": 3,
    "source_title_alias": {
        "id": "序号",
        "tag_name": "标签名称",
        "uid": "用户id",
        "nickname": "用户昵称"
    },
    "source_title": [
        "id",
        "tag_name",
        "uid",
        "nickname"
    ],
    "source_header_prefix": null,
    "label": "tag-export",
    "source_data": null,
    "source_sql_database": "",
    "source_sql_count": 0,
    "source_sql": "",
    "source_restful_path": "https://ws-dev.vchangyi.com/micro-scrm/micro-scrm/v1/tag/future/export-join",
    "source_restful_method": "POST",
    "source_restful_headers": {},
    "source_restful_params": {
        "enterprise_id": 192104618298240,
        "file_id": 3633
    },
    "source_suffix": null,
    "source_restful_count_field": "data.total",
    "source_restful_data_field": "data.list",
    "source_restful_request_page_field": "page",
    "source_restful_request_size_field": "size",
    "source_restful_receive_select_switch": true,
    "source_restful_sort_key": "a.id",
    "source_restful_sort_key_as": "id",
    "source_restful_sort_type": "ASC",
    "source_restful_select_info": [
        "a.id",
        "a.tag_name",
        "b.id as uid",
        "b.nick as nickname"
    ],
    "source_excel_cell_merge": {
        "type": 0,
        "rule": "",
        "title": null,
        "premise": null
    },
    "source_restful_pagination_type": 2
}
```