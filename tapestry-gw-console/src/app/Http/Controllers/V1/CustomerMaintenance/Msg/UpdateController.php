<?php
/**
 * 修改群发消息
 * User: tony
 * Date: 20200220
 * Time: 11:20
 */

namespace App\Http\Controllers\V1\CustomerMaintenance\Msg;

use App\Models\QYwxEmsg;
use App\Models\QYwxMedia;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class UpdateController extends BaseController
{
    protected $send = null;

    public function __construct(Request $request)
    {
        parent::__construct();
        $this->request = $request;
        $this->send    = new ListController($request);
    }

    /**
     * 修改群发消息 - 只有草稿状态下的消息才能被修改
     *
     * @param string $tsid
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index($tsid = "")
    {
        if (!$tsid) error(1000005);
        $input = $this->request->only(
            'send_type', 'send_object', 'name', 'text_content', "num_est", "send_time", "media_id"
        );

        $rules    = [
            'name'        => 'required|max:30',
            'send_type'   => 'required|in:0,1,2,3',
            'send_object' => 'required|array',
            'send_time'   => 'integer',
            'num_est'     => 'integer',
        ];
        $messages = [
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

        $data                 = $input;
        $data["text_content"] = !empty($data['text_content']) ? $data['text_content'] : "";

        if (isset($input['send_time']) && $input['send_time'] < 0) {
            $data['send_time'] = 0;
            $data['status']    = 0;
        }

        // - 文本类型的 内容不能为空
        if (!$data['send_type'] && !$data["text_content"]) {
            notice(1000002);
        } else {
            if ($input['send_type'] == 0) {
                $data['image_media_url'] = "";
                $data['link_title']      = "";
                $data['link_desc']       = "";
                $data['link_content']    = "";
                $data['minip_title']     = "";
                $data['minip_appid']     = "";
                $data['minip_page']      = "";
            } else {
                if (!empty($input['media_id'])) {
                    if (!($media = QYwxMedia::query()->where(['scid' => $input['media_id']])->first())) {
                        notice(1000003);
                    }
                    $data['image_media_url'] = $media->image_media_url;
                    $data['link_title']      = "";
                    $data['link_desc']       = "";
                    $data['link_content']    = "";
                    $data['minip_title']     = "";
                    $data['minip_appid']     = "";
                    $data['minip_page']      = "";
                    switch ($input['send_type']) {
                        case 2:
                            $data['link_title']   = $media->name;
                            $data['link_desc']    = $media->link_desc;
                            $data['link_content'] = $media->link_content;
                            break;
                        case 3:
                            $data['minip_title'] = $media->name;
                            $data['minip_appid'] = $media->minip_appid;
                            $data['minip_page']  = $media->minip_page;
                            break;
                    }
                }
            }

        }

        //--去掉无用参数
        if (isset($data['media_id'])) {
            unset($data['media_id']);
        }
        $data['send_object'] = json_encode($data['send_object']);
        $data['user_id']     = $this->enterpriseAdminId ? $this->enterpriseAdminId : "";
        $data['updated_at']  = time();
        $id                  = '';
        $where               = [['tsid', $tsid]];
        if ($this->qyDomain) {
            $where[] = ['uc_eid', $this->qyDomain];
        }
        DB::beginTransaction();
        try {
            if (QYwxEmsg::query()->where($where)->update($data)) {
                $id = $tsid;
                if (!isset($input['send_time']) || $input['send_time'] >= 0) {
                    //--立即群发和定时群发 -- 直接推送到接口
                    if ($data = QYwxEmsg::query()->where($where)->first()) {
                        $this->send->send_api($data->toArray());
                    }
                }
            }
            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            Log::info("更新推送失败：[{$e->getCode()}] " . $e->getMessage());
            notice($e->getCode());
        }
        return success(['id' => $id]);
    }
}
