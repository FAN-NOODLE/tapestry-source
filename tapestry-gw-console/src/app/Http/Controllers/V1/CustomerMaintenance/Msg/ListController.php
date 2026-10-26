<?php

namespace App\Http\Controllers\V1\CustomerMaintenance\Msg;

use CrsSdk\Work\Channel;
use App\Models\QYwxEmsg;
use Illuminate\Http\Request;
use App\Services\UcService;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    public $ucService = null;//推送服务

    public function __construct(Request $request)
    {
        parent::__construct();
        $this->request   = $request;
        $this->ucService = new UcService($this->qyDomain);
    }

    /**
     * 获取消息群发列表
     * @return JsonResponse|void
     */
    public function index()
    {
        $input = $this->request->only(
            'keyword', 'send_object', "status", "sts", "ets", "page", "size"
        );

        $rules    = [
            'keyword'     => 'max:30',
            'send_type'   => 'in:0,1,2,3',
            'send_object' => 'array',
            'sts'         => 'integer',
            'ets'         => 'integer',
        ];
        $messages = [
            'keyword.max'       => '搜索关键字不能超过30个字符',
            'send_type.in'      => '推送类型数据不在范围内',
            'send_object.array' => '推送对象格式错误',
            'sts.integer'       => '时间区间应该为数值类型',
            'ets.integer'       => '时间区间应该为数值类型'
        ];

        $validator = Validator::make($input, $rules, $messages);
        if ($validator->fails()) {
            notice(1000010, $validator->errors()->all());
        }

        $sts         = $input['sts'] ?? 0;
        $ets         = $input['ets'] ?? 0;
        $status      = isset($input['status']) ? $input['status'] : "";
        $keyword     = isset($input['keyword']) && trim($input['keyword']) ? trim($input['keyword']) : "";
        $send_object = empty($input['send_object']) ? "" : json_encode($input['send_object']);

        $db     = QYwxEmsg::query();
        $dbCout = QYwxEmsg::query();


        $db->select(["tsid",
                     "send_type",//推送内容类型  0文本消息|1图片素材|2图文链接|3小程序卡片
                     "send_object",
                     "send_time",//0：立即推送|时间戳：定时推送时间
                     "name",
                     "text_content",
                     "num_est",//预计触达
                     "num_fact",//实际触达
                     "status",//0草稿状态 1推送中 2已完成 3定时未开始
                     "updated_at",
                     "created_at"]);

        $where = [];
        if ($this->qyDomain) {
            $where[] = ['uc_eid', $this->qyDomain];
        }
        if ($sts) {
            $where[] = ['updated_at', '>=', $sts];
        }
        if ($ets) {
            $where[] = ['updated_at', '<=', $ets];
        }

        if ($status !== "") {
            $where[] = ['status', $status];
        }

        $db->where($where);
        $dbCout->where($where);

        if ($keyword) {
            //$sql = DB::raw("(name like '%{$keyword}%' OR tsid='%{$keyword}%')");-//此处%与getQuerySql方法冲突
            $sql = function ($query) use ($keyword) {
                $query->where('name', 'like', '%' . $keyword . '%')
                    ->orWhere('tsid', 'like', '%' . $keyword . '%');
            };
            $db->where($sql);
            $dbCout->where($sql);
        }

        if ($send_object) {
            //-如果mysql版本大于等于5.7 就优先使用JSON_CONTAINS
            if (version_compare(DB::select("select version() AS V")[0]->V, '5.7', 'gt')) {
                //$sql = DB::raw("JSON_CONTAINS(send_object,'{$send_object}','$')");
                $s = [];
                foreach ($input['send_object']['list'] AS $item) {
                    $s[] = "JSON_CONTAINS(send_object->'$.list[*]','\"{$item}\"', '$')";
                }
                $sql = DB::raw("send_object->'$.type'={$input['send_object']['type']}" . (count($s) ? " AND (" . implode(" OR ", $s) . ")" : ""));
                $db->whereRaw($sql);
                $dbCout->whereRaw($sql);
            } else {
                //JSON_CONTAINS 替代写法
                $sql = function ($query) use ($input, &$test) {
                    $query->where(function ($query) use ($input) {
                        $query->where('send_object', 'like', '%"type": ' . $input['send_object']['type'] . '%');
                        $query->orWhere('send_object', 'like', '%"type":' . $input['send_object']['type'] . '%');
                    });
                    if (isset($input['send_object']['list']) && is_array($input['send_object']['list'])) {
                        $query->where(function ($query) use ($input) {
                            foreach ($input['send_object']['list'] AS $v) {
                                $query->orWhere('send_object', 'like', '%"' . $v . '"%');
                            }
                        });
                    }
                };
                $db->where($sql);
                $dbCout->where($sql);
            }
        }

        $page = isset($input['page']) && $input['page'] ? $input['page'] : 1;
        $size = isset($input['size']) && $input['size'] ? $input['size'] : 20;

        $db->offset(($page - 1) * $size)->limit($size);
        $db->orderBy('updated_at', 'desc');

        $list = $db->get()->toArray();
        foreach ($list as $i => $item) {
            $type = [];
            if ($item['text_content']) {
                $type[] = 0;
            }
            unset($list[$i]['text_content']);
            $item['send_type'] && $type[] = $item['send_type'];
            $list[$i]['send_type']   = $type;
            $list[$i]['send_object'] = json_decode($item['send_object'], 1);
        }

        $data = [
            'list'  => $list,
            'page'  => $page,
            'size'  => $size,
            'total' => $dbCout->count(),
        ];

        return success($data);
    }

    /**
     * 获取下拉数据
     * @param string $type 默认所有|1所属员工|2渠道|3标签
     * @return array|JsonResponse
     * @throws \Exception
     */
    public function getObjectSelect($type = "")
    {
        $data = [];
        if ($type === "") {
            $data[0] = 0;
            $data[1] = $this->getObjectSelect(1)['data'];//-员工列表
            $data[2] = $this->getObjectSelect(2)['data'];//-渠道列表
            $data[3] = $this->getObjectSelect(3)['data'];//-标签列表
        } else {
            if (Cache::has('uc_cache_' . $this->qyDomain . '_' . $type) && env('APP_ENV') == 'prod' && 0) {
                eval('$data = ' . Cache::get('uc_cache_' . $this->qyDomain . '_' . $type) . ";");
            } else {
                switch ($type) {
                    case 1:
                        if ($list = $this->ucService->send('userList', '{"pageSize":10000}')) {
                            $list = @json_decode($list, 1);
                            if (isset($list['data']['list'])) {
                                foreach ($list['data']['list'] as $item) {
                                    $data[] = [
                                        "name" => $item['memUsername'],
                                        "id"   => (string)$item['memUserid']
                                    ];
                                }
                                Cache::put('uc_cache_' . $this->qyDomain . '_' . $type, var_export($data, true), (strtotime(date('Y-m-d', strtotime('+1 day'))) - time()) / 60);
                            }
                        }
                        break;
                    case 2:
                        $list = (new Channel())->channelList([
                            'enterprise_id' => $this->enterpriseId,
                            "size"          => 0
                        ]);
                        //$list = @json_decode($list,1);
                        if (isset($list['data'])) {
                            foreach ($list['data'] AS $item) {
                                $data[] = ['name' => $item['channel'], 'id' => (string)$item['id']];
                            }
                        }
                        if (count($data))
                            Cache::put('uc_cache_' . $this->qyDomain . '_' . $type, var_export($data, true), (strtotime(date('Y-m-d', strtotime('+1 day'))) - time()) / 60);
                        break;
                    case 3:
                        if ($list = $this->ucService->send('tagList')) {
                            $list = @json_decode($list, 1);
                            if (isset($list['data']))
                                foreach ($list['data'] AS $value) {
                                    foreach ($value['tagList'] as $item) {
                                        $data[] = [
                                            "name" => $value['groupName'] . "[{$item['tagName']}]",
                                            "id"   => (string)$item['tagId']
                                        ];
                                    }
                                }
                            if (count($data)) {
                                Cache::put('uc_cache_' . $this->qyDomain . '_' . $type, var_export($data, true), (strtotime(date('Y-m-d', strtotime('+1 day'))) - time()) / 60);
                            }
                        }
                        break;
                }
            }

            return ['data' => $data];
        }
        return success(['list' => $data]);
    }

    /**
     * 获取推送消息的结果详情
     *
     * @param string $tsid
     *
     * @return array|JsonResponse
     */
    public function getInfo($tsid = "")
    {
        $input = $this->request->only(
            "type", "page", "size"
        );

        $rules    = [
            'type' => 'integer|in:0,1,2',//0:all 1:已发送，2未发送
            'page' => 'integer',
            'size' => 'integer',
        ];
        $messages = [
            'type.integer' => '筛选类型为数值类型',
            'type.in'      => '筛选类型传值不正确',
            'page.integer' => '分页参数必须为数值类型',
            'size.integer' => '每页显示数必须为数值类型'
        ];

        $validator = Validator::make($input, $rules, $messages);
        if ($validator->fails()) {
            notice(1000010, $validator->errors()->all());
        }

        $page     = isset($input['page']) && $input['page'] ? $input['page'] : 1;
        $size     = isset($input['size']) && $input['size'] ? $input['size'] : 10;
        $type     = isset($input['type']) ? $input['type'] : 0;
        $sendData = [
            'tsid'      => $tsid,
            'page_idx'  => $page,
            'page_size' => $size,
            'type'      => (int)$type
        ];

        if (($ret = $this->ucService->send('sendRet', json_encode($sendData))) && $ret = @json_decode($ret, 1)) {
            if ($ret['Code']) {
                notice(1000009, ["[{$ret['Code']}] {$ret['Msg']}"]);
            }
        } else {
            notice(600005);
        }

        $data = [
            'info'  => [
                'all_cnt'          => $ret['all_cnt'],// 应发送
                'sended_cnt'       => $ret['sended_cnt'],// 已完成
                'un_send_cnt'      => $ret['un_send_cnt'],// 未发送
                'un_contact_cnt'   => $ret['un_contact_cnt'],// 非联系人
                'already_send_cnt' => $ret['already_send_cnt'] // 已经接收其他消息
            ],
            'list'  => $ret['list'],
            'total' => !$type ? $ret['all_cnt'] : ($type == 1 ? $ret['sended_cnt'] : $ret['un_send_cnt']),
            'page'  => $page,
            'size'  => $size
        ];

        return success($data);
    }


    /**
     * 获取推送消息编辑详情
     *
     * @param string $tsid
     *
     * @return array|JsonResponse
     */
    public function getEditInfo($tsid = "")
    {
        if (!$tsid || !($data = QYwxEmsg::query()->where('tsid', $tsid)->first())) error(1000004);
        $data                = $data->toArray();
        $data['send_object'] = json_decode($data['send_object'], 1);
        return success($data);
    }

    /**
     * 整理数据并推送到发送接口
     * @param $data
     * @return bool
     */
    public function send_api($data = [])
    {
        if (!count($data)) {
            notice(1000001);
        }

        //print_r($data);
        $sendData = [
            "tsid"      => $data['tsid'],
            "send_obj"  => json_decode($data['send_object'], 1),
            "send_time" => $data['send_time'] ? $data['send_time'] : time(),
        ];
        // 0 文本  1图片  2链接  3小程序卡片
        !empty($data['text_content']) && $sendData['txt'] = $data['text_content'];
        switch ($data['send_type']) {
            case 1:
                $sendData['picUrl'] = $data['image_media_url'];
                break;
            case 2:
                preg_match("/^(http:\/\/|https:\/\/)?([^\/]+)/i", $this->request->url(), $matches);
                $sendData['link'] = [
                    'title'  => $data['link_title'],
                    'desc'   => $data['link_desc'],
                    'url'    => 'https://' . $matches[2] . "/gw-console/v1/customer-maintenance/msg/article/{$data['tsid']}",
                    'picUrl' => $data['image_media_url']
                ];
                break;
            case 3:
                $sendData['miniprogram'] = [
                    "title"  => $data['minip_title'],
                    "appid"  => $data['minip_appid'],
                    "page"   => $data['minip_page'],
                    "picUrl" => $data['image_media_url']
                ];
                break;
        }

        if (($ret = $this->ucService->send('send', json_encode($sendData))) && $ret = @json_decode($ret, 1)) {
            if ($ret['Code']) {
                notice(1000009, ["当前推送对象的可推送次数已用完[{$ret['Code']}]"]);
            }
        } else {
            notice(600005);
        }

        if ($data['send_time']) {
            $updateStatus = 3;//--定时未开始
        } else {
            $updateStatus = 1;//--推送中
        }
        QYwxEmsg::query()->where('tsid', $data['tsid'])->update(['status' => $updateStatus]);
        return true;
    }

    /**
     * 获取预计触达人数
     *
     * @return JsonResponse
     */
    public function getEst()
    {
        $input = $this->request->only('send_object');

        $rules = [
            'send_object' => 'required|array',
        ];

        $messages = [
            'send_object.required' => '发送对象必须选择',
            'send_object.array'    => '推送对象格式错误'
        ];

        $validator = Validator::make($input, $rules, $messages);
        if ($validator->fails()) {
            notice(1000010, $validator->errors()->all());
        }

        $send_object = json_encode($input['send_object']);

        if (($ret = $this->ucService->send('numEst', $send_object)) && $ret = @json_decode($ret, 1)) {
            if (!$ret['Code']) {
                $data = ['noSend' => $ret['already_send_cnt'], 'num' => $ret['all_cnt']];
            } else {
                notice(1000009, ["[{$ret['Code']}] {$ret['Msg']}"]);
            }
        } else {
            notice(600005);
        }

        return success($data);
    }

}
