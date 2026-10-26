<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/8/20
 * Time: 10:37
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\DataCompass;
use CrsSdk\Work\Member;
use Illuminate\Http\Request;

class RankController extends BaseController {

    const RANK_LENGTH = 50;

    const RANK_TYPE_COLUMN_MAP = [
        1 => 'cumulative_customer_num',
        2 => 'today_customer_num',
        5 => 'chat_cnt',
        6 => 'reply_percentage',
        7 => 'cumulative_binding_customer_num',
        8 => 'today_binding_customer_num'
    ];

    /**
     * 顾客数据排行榜
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();
        $params['qy_domain'] = $this->loginAdminEnterprise['domain'];
        $params['admin_id'] = $this->adminId;
        $params['rank_column'] = !empty($params['type']) && isset(self::RANK_TYPE_COLUMN_MAP[$params['type']]) ? self::RANK_TYPE_COLUMN_MAP[$params['type']] : '';
        $params['rank_length'] = 50;
        if (empty($params['rank_column'])) {
            return success(['list' => []]);
        }

        $rankList = (new DataCompass())->customerStatisticRankForApp($params);

        return success(['list' => $rankList['list']]);
    }
}