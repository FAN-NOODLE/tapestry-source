<?php
/**
 * 新增素材
 * User: tony
 * Date: 20200220
 * Time: 11:20
 */

namespace App\Http\Controllers\V1\CustomerMaintenance\Media;

use App\Models\QYwxMedia;
use CrsSdk\Work\Source;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class AddController extends BaseController
{
    public function __construct(Request $request)
    {
        parent::__construct();
        $this->request = $request;
    }

    /**
     * 添加素材
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $input = $this->request->only(
            'type', 'name', 'content', 'desc', 'img', 'user_name'
        );

        $rules    = [
            'type' => 'required|in:0,1,2,3',
            'img'  => 'url'
        ];
        $messages = [
            'type.required' => '推送类型不能为空',
            'type.in'       => '推送类型数据不在范围内',
            'img'           => '图像地址必须是一个合法的url'
        ];

        $validator = Validator::make($input, $rules, $messages);
        if ($validator->fails()) {
            notice(1000010, $validator->errors()->all());
        }

        $data = ['type' => $input['type']];

        if (isset($input['img']) && $input['img']) {
            $data['image_media_url'] = trim($input['img']);
        }
        if (isset($input['name']) && $input['name']) {
            $data['name'] = trim($input['name']);
        }

        $data['user_name'] = $this->enterpriseAdmin ? $this->enterpriseAdmin['username'] : "";
        //0文本素材|1图片素材|2图文链接消息|3小程序卡片素材
        switch ($input['type']) {
            case 0:
                $data['text_content'] = trim($input['content']);
                break;
            case 1:
                break;
            case 2:
                $data['link_desc']    = trim($input['desc']);
                $data['link_content'] = trim($input['content']);
                if (isset($input['user_name']))
                    $data['user_name'] = trim($input['user_name']);
                break;
            case 3:
                $data['minip_appid'] = trim($input['desc']);
                $data['minip_page']  = trim($input['content']);
                break;
        }

        $data['user_id']    = $this->enterpriseAdminId ? $this->enterpriseAdminId : "";
        $data['scid']       = $this->createID();
        $data['uc_eid']     = $this->qyDomain ? $this->qyDomain : "";
        $data['created_at'] = time();
        $data['updated_at'] = $data['created_at'];
        $id                 = '';
        if (QYwxMedia::query()->insert($data)) {
            $id = $data['scid'];
        }

        if (!empty($data['image_media_url'])) {
            try {
                $params = [
                    'uc_eid'          => $this->qyDomain ? $this->qyDomain : "",
                    'image_media_url' => $data['image_media_url'],
                    'scid'            => $id,
                ];
                (new Source())->update_media_id($params);
            } catch (\Exception $exception) {
                Log::error('生成mediaId失败params:' . json_encode($params) . '错误提示：' . $exception->getMessage());
            }
        }

        return success(['id' => $id]);
    }

    /**
     * 最新的id流水号
     * @return bool|string
     */
    protected function createID()
    {
        $ling = "00000000";
        $d    = 'SC' . date('ymd');
        $i    = 0;
        if (Cache::has($d . '_liushui')) {
            $i = Cache::get($d . '_liushui');
        }
        $new_num = (int)$i + 1;
        Cache::put($d . '_liushui', $new_num, 1440);//--仅保留一天
        return $d . substr($ling . $new_num, -3);
    }
}
