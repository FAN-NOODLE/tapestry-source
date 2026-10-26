<?php
/**
 * 新增消息推送
 * User: tony
 * Date: 20200220
 * Time: 11:20
 */

namespace App\Http\Controllers\V1\CustomerMaintenance\Msg;

use App\Models\QYwxEmsg;
use App\Models\QYwxMedia;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class AddController extends BaseController
{
    protected $send = null;//--消息推送接口地址

    public function __construct(Request $request)
    {
        parent::__construct();
        $this->request = $request;
        $this->send    = new ListController($request);
    }

    /**
     * 添加群发消息
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $input = $this->request->only(
            'send_type', 'send_object', 'name', 'text_content', "num_est", "send_time",
            "media_id"
        );

        $rules     = [
            'name'        => 'required|max:30',
            'send_type'   => 'required|in:0,1,2,3',
            'send_object' => 'required|array',
            'send_time'   => 'integer',
            'num_est'     => 'integer',
        ];
        $messages  = [
            'name.required'        => '推送名称不能为空',
            'name.max'             => '推送名称不能超过30个字符',
            'send_type.required'   => '推送类型不能为空',
            'send_type.in'         => '推送类型数据不在范围内',
            'send_object.required' => '推送对象不能为空',
            'send_object.array'    => '推送对象格式错误',
            'send_time.integer'    => '定时推送时间参数错误',
            'num_est.integer'      => '预计触达人数参数错误'
        ];
        $validator = Validator::make($input, $rules, $messages);
        if ($validator->fails()) {
            notice(1000010, $validator->errors()->all());
        }

        $data = $input;

        if (!empty($data['text_content'])) {
            $data["text_content"] = trim($data['text_content']);
        }


        if (isset($input['send_time']) && $input['send_time'] < 0) {
            $data['send_time'] = 0;
            $data['status']    = 0;
        }

        // - 文本类型的 内容不能为空
        if (!$data['send_type'] && !$data["text_content"]) {
            notice(1000002);
        } else {
            if ($input['send_type'])
                if (!empty($input['media_id']) && ($media = QYwxMedia::query()->where(['scid' => $input['media_id']])->first())) {
                    $data['image_media_url'] = $media->image_media_url;
                    switch ($input['send_type']) {
                        case 1:
                            $data['minip_page'] = $media->name;//字段复用-存储图片名称
                        case 2:
                            $data['link_title']   = $media->name;
                            $data['link_desc']    = $media->link_desc;
                            $data['link_content'] = $media->link_content;
                            $data['minip_page']   = $media->user_name;//字段复用-存储作者名称
                            break;
                        case 3:
                            $data['minip_title'] = $media->name;
                            $data['minip_appid'] = $media->minip_appid;
                            $data['minip_page']  = $media->minip_page;
                            break;
                    }
                } else {
                    notice(1000003);
                }
        }

        //--去掉无用参数
        if (isset($data['media_id'])) {
            unset($data['media_id']);
        }
        $data['send_object'] = json_encode($data['send_object']);
        $data['user_id']     = $this->enterpriseAdminId ? $this->enterpriseAdminId : "";
        $data['tsid']        = $this->createID();
        $data['uc_eid']      = $this->qyDomain ? $this->qyDomain : "";
        $data['created_at']  = time();
        $data['updated_at']  = $data['created_at'];
        $id                  = '';
        //---此处不能使用事务处理-中间有数据更新如果失败删除数据即可
        try {
            if (QYwxEmsg::query()->insert($data)) {
                $id = $data['tsid'];
                if (isset($input['send_time']) && $input['send_time'] >= 0) {
                    //--立即群发和定时群发 -- 直接推送到接口
                    $this->send->send_api($data);
                }
            }
        } catch (\Exception $e) {
            //--流水号减一
            //$this->createID('-');
            QYwxEmsg::query()->where('tsid', $data['tsid'])->delete();
            Log::info("添加推送失败UC2：[{$e->getCode()}] " . $e->getMessage());
            notice(1000010, [$e->getMessage()]);
        }

        return success(['id' => $id]);
    }

    /**
     * 最新的id流水号
     * @return bool|string
     */
    protected function createID($jj = '+')
    {
        $ling = "000000";
        //--不同环境修改前缀以免回调信息互相影响
        $d = '';
        if (env('APP_ENV') == 'test')
            $d = 'T';
        if (env('APP_ENV') == 'dev')
            $d = 'D';
        $d .= 'TS' . date('ymd');
        $i = 0;
        if (Cache::has($d . '_liushui')) {
            $i = Cache::get($d . '_liushui');
        }
        if ($jj == '+') $new_num = (int)$i + 1;
        else $new_num = (int)$i - 1;
        Cache::put($d . '_liushui', $new_num, 1440);//--仅保留一天
        return $d . substr($ling . $new_num, -3);
    }
}
