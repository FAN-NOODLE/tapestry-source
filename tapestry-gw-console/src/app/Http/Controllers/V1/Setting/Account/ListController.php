<?php


namespace App\Http\Controllers\V1\Setting\Account;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Account;

/**
 * 账号列表
 * Class ListController
 * @package App\Http\Controllers\V1\account
 */
class ListController extends BaseController
{
    /**
     * 账号列表
     * @api {domain}/gw-console/v1/setting/account/list
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $input = $this->request->only([
            'search_type',
            'keyword',
            'object_id',
            'role_id',
            'page',
            'size'
        ]);

        $input['enterprise_id'] = $this->enterpriseId;
        $input['size']          = $input['size'] ?? Constant::DEFAULT_SIZE;
        $input['page']          = $input['page'] ?? Constant::PAGE;
        $input['qy_domain']     = $this->qyDomain;

        // 验证--都是非必传
        $result = (new Account())->listAccount($input);
        if (isset($result['list']) && !empty($result["list"])) {
            foreach ($result['list'] as &$item) {
                $item['mobile'] = preg_replace("/(.{3}).*(.{4})/", '$1****$2', $item["mobile"]);
            }
        }
        return success([
            'list'  => $result['list'] ?? [],
            'page'  => $result['page'] ?? $input['page'],
            'size'  => $result['size'] ?? $input['size'],
            'total' => $result['total'] ?? 0
        ]);
    }
}