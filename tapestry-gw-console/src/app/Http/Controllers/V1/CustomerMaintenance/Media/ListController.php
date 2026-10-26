<?php
/**
 * 素材列表
 * User: tony
 * Date: 20200220
 * Time: 11:20
 */

namespace App\Http\Controllers\V1\CustomerMaintenance\Media;

use CrsSdk\Boss\Auth;
use CrsSdk\Work\Common;
use App\Models\QYwxMedia;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    public function __construct(Request $request)
    {
        parent::__construct();
        $this->request = $request;
    }

    /**
     * 获取素材列表
     *
     * @param string $type 默认所有,0文本素材 1图片素材|2图文链接|3小程序卡片
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index($type = "")
    {
        $input    = $this->request->only(
            'keyword', "sts", "ets", "status", "page", "size"
        );
        $rules    = [
            'sts' => 'integer',
            'ets' => 'integer',
        ];
        $messages = [
            'sts.integer' => '时间区间应该为数值类型',
            'ets.integer' => '时间区间应该为数值类型'
        ];

        $validator = Validator::make($input, $rules, $messages);
        if ($validator->fails()) {
            notice(1000010, $validator->errors()->all());
        }

        $sts = $input['sts'] ?? 0;
        $ets = $input['ets'] ?? 0;

        $keyword = isset($input['keyword']) && trim($input['keyword']) ? trim($input['keyword']) : "";

        $db     = QYwxMedia::query();
        $dbCout = QYwxMedia::query();

        $where = [];
        if ($this->qyDomain) {
            $where[] = ['uc_eid', $this->qyDomain];
        }

        if (isset($input['status'])) {
            $where[] = ['status', $input['status']];
        } else {
            $where[] = ['status', ">=", 0];
        }

        if ($sts) {
            $where[] = ['updated_at', '>=', $sts];
        }
        if ($ets) {
            $where[] = ['updated_at', '<=', $ets];
        }

        if ($type !== "") {
            $where[] = ['type', $type];
        }

        $db->where($where);
        $dbCout->where($where);

        if ($keyword) {
            $sql = function ($query) use ($keyword, $type) {
                $query->where('name', 'like', '%' . $keyword . '%');
                if (!$type) {
                    $query->orWhere('text_content', 'like', '%' . $keyword . '%');
                }
                $query->orWhere('scid', 'like', '%' . $keyword . '%');

            };
            $db->where($sql);
            $dbCout->where($sql);
        }

        $page = isset($input['page']) && $input['page'] ? $input['page'] : 1;
        $size = isset($input['size']) && $input['size'] ? $input['size'] : 20;
        $db->offset(($page - 1) * $size)->limit($size);
        $db->orderBy('updated_at', 'desc');

        $list = $db->get()->toArray();

        $data = [
            'list'  => $list,
            'page'  => $page,
            'size'  => $size,
            'total' => $dbCout->count(),
        ];

        return success($data);
    }

    /**
     * 获取商城小程序appid
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function getShopAppid()
    {
        // 获取企业app_id
        try {
            $appIdArr = (new Auth())->getShopAppid(['enterprise_id' => $this->enterpriseId]);
        } catch (\Exception $exception) {
            $appIdArr['data']['shop_app_id'] = '';
        }

        // 获取绑定相关列表
        try {
            $params                   = [];
            $params['qy_domain']      = $this->qyDomain;
            $params['app_identifier'] = $this->appIdentifier;
            $params['enterprise_id']  = $this->enterpriseId;
            $list                     = (new Common())->appList($params);
        } catch (\Exception $exception) {
            $list = [];
        }

        // boss挂了
        if (empty($appIdArr['data']['shop_app_id'])) {
            notice(600005);
        }

        return success(['data' => $appIdArr['data']['shop_app_id'], 'list' => $list]);
    }
}
