<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/12/6
 * Time: 16:49
 */

namespace App\Http\Controllers\App\V1\User;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\MemberUserRelation;

class StarsController extends BaseController
{

    /**
     * @var MemberUserRelation
     */
    protected $memberUserRelation;

    public function __construct(MemberUserRelation $memberUserRelation)
    {
        parent::__construct();
        $this->memberUserRelation = $memberUserRelation;
    }

    /**
     * 星标客户
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $params = $this->validates();

        // 获取关系表信息
        $userRelation = $this->memberUserRelation->query()->where([
            'enterprise_id' => $this->enterpriseId,
            'ec_userid' => $params['ec_userid'],
            'mem_uid' => $this->memUid
        ])->first();
        if (!$userRelation) {
            error(200013);
        }

        $update = ['stars' => MemberUserRelation::STARS_FALSE];
        if ($userRelation->stars == MemberUserRelation::STARS_FALSE) {
            $update['stars'] = MemberUserRelation::STARS_TRUE;
        }

        $this->memberUserRelation->query()->where('id', $userRelation->id)->update($update);

        return success();
    }

    /**
     * 数据基础验证
     *
     * @return array|string|null
     * @throws \Illuminate\Validation\ValidationException
     */
    private function validates()
    {
        $rules = [
            'ec_userid' => 'required'
        ];

        $message = [
            'ec_userid.required' => '客户ID不存在'
        ];

        $this->validate(request(), $rules, $message);

        $params = request()->input();

        return $params;
    }
}
