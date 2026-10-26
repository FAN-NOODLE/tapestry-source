<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/14
 * Time: 14:51
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Work\Tag;
use Illuminate\Http\Request;

class TagController extends BaseController {

    /**
     * allTag
     * @desc 查询企业所有标签信息
     * @return \Illuminate\Http\JsonResponse
     * @date 2020-07-14
     * @author XuHuitao
     * @throws \Exception
     */
    public function allTag() {
        $enterpriseId = $this->enterpriseId;
        $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];

        $list = (new Tag())->allList(['qy_domain' => $enterpriseInfo['domain']]);

        return success(['list' => $list['data']]);
    }



    /**
     * 标签分组列表
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function tagsGroupList(Request $request)
    {
        $params                  = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->loginMobile();

        $list = (new Tag())->allList($params);

        return success($list);
    }


}