<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/8/20
 * Time: 10:45
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;
use App\Micros\MicroScrm\GoMicroScrm;
use App\Models\StoreMemberViewTab;
use CrsSdk\GoMicroScrm\Guide;
use CrsSdk\MicroAuthority\Account;
use CrsSdk\Sdk;
use CrsSdk\Work\DataCompass;
use CrsSdk\Work\StoreWorker;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class StatisticController extends BaseController
{
    const OVERVIEW_COLUMN_NAME_MAP = [
        'data1'  => 'cumulative_customer_num',
        'data2'  => 'today_customer_num',
        'data3'  => 'customer_access',
        'data4'  => 'customer_deal',
        'data5'  => 'cumulative_contact_card_way_customer_num',
        'data6'  => 'cumulative_other_way_customer_num',
        'data7'  => 'today_contact_card_way_customer_num',
        'data8'  => 'today_other_way_customer_num',
        'data9'  => 'chat_cnt',
        'data10' => 'reply_percentage',
        'data11' => 'cumulative_binding_customer_num',
        'data12' => 'today_binding_customer_num',
        'data13' => 'cumulative_member_num',
        'data14' => 'today_member_num',
        'data15' => 'cumulative_non_member_num',
        'data16' => 'today_non_member_num',
        'data17' => 'cumulative_customer_join_group_rate',
        'data18' => 'cumulative_vip_external_user_rate',
    ];

    /**
     * 顾客数据概览
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request)
    {
        $params                  = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->loginMobile();
        $params['qy_domain']     = $this->loginAdminEnterprise['domain'];
        $params['admin_id']      = $this->adminId;

        $overview      = (new DataCompass())->customerStatisticOverviewForApp($params);
        $overviewData  = $overview['overview'];
        $statisticList = $overview['detail_list'];

        $result = [];
        foreach (self::OVERVIEW_COLUMN_NAME_MAP as $key => $value) {
            $result[$key] = $overviewData[$value] ?? 0;
        }
        $x = $y = [];
        if (!empty($statisticList)) {
            $x           = array_map(function ($statTime) {
                return date('Y-m-d', $statTime);
            }, array_column($statisticList, 'stat_time'));
            $defaultData = array_fill(0, count($x), 0);
            foreach (self::OVERVIEW_COLUMN_NAME_MAP as $key => $value) {
                if (in_array($key, ['data3', 'data4'])) {
                    $y[$key] = $defaultData;
                } else {
                    $y[$key] = array_column($statisticList, $value);
                }
            }
        }

        $result['detail_data'] = [
            'x' => $x,
            'y' => $y
        ];
        return success($result);
    }


    /**
     * 记录顾客数据提醒
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function setViewTab(Request $request)
    {
        $params             = $request->input();
        $params['store_id'] = $params['store_id'] ?? 0;
        //获取店员id
        $workerDetail = $this->storeWorkerDetail($this->enterpriseId, $params['store_id']);
        if (!empty($workerDetail['id'])) {
            $this->saveViewTab($this->enterpriseId, $params['store_id'], $workerDetail['id'] ?? 0, StoreMemberViewTab::CUSTOMER_TOTAL_TAB_TYPE);
        }
        return success();
    }


    public function getCompassSummary(Request $request)
    {
        $input                  = $request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['mobile']        = $this->mobile;
        $searchAll              = 0;

        if (isset($input['search_all']) && $input['search_all'] == 1) {
            $adminId    = $this->loginAdminEnterprise['admin_id'];
            $dataRange  = (new Account())->dataRangeList(['admin_id' => $adminId]);
            $dpIds      = array_column($dataRange, 'object_id');
            $memberList = Sdk::work()->setApi("service/v1/department/member-list-by-dp-ids")->post([
                'enterprise_id' => $this->enterpriseId,
                'dp_ids'        => $dpIds,
            ]);
            $useridArr  = array_column($memberList, 'userid');
            $memberInfo = GoMicroScrm::getMemberInfo([
                'enterprise_id' => $this->enterpriseId,
                'mobile'        => $this->mobile,
            ]);
            if (!empty($memberInfo)) {
                $useridArr[] = $memberInfo['userid'];
            }
            $input['user_ids'] = $useridArr;
        }

        if ($input['search_all'] == 0) {
            // 后续3.3版本之后，shopping_guide的参数判断验证可以去掉
            if (isset($input['userids']) && !empty($input['userids'])) { // 此处的验证为coach的导购查看自己的数据罗盘，userids的值为[null]的处理
                $userids = array_filter($input['userids'], function ($item) {
                    return !empty($item);
                });
                if (!empty($userids)) {
                    $input['user_ids'] = $userids;
                }
            }
            if (empty($input['user_ids'])){  // 这里是针对导购访问自己的数据罗盘，都不传userids 和shopping_guide的情况
                $memberInfo = (new StoreWorker())->detail([
                    'mobile'        => $this->mobile,
                    'enterprise_id' => $this->enterpriseId,
                    'use_scrm_db'   => true
                ]);
                $input['user_ids'] = [$memberInfo['userid']];
            }
        }

        $redisKeyMd5 = md5(json_encode($input));
        $redisKey    = 'DATA_COMPASS_SUMMARY' . $redisKeyMd5;
        $cache       = getOrSetData($redisKey);
        if ($cache) {
            return success($cache);
        }
        try {
            $input['search_all'] = $searchAll;
            $result              = (new Guide())->getCompassSummary($input);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        $return = $result['data'] ?? $result;
        getOrSetData($redisKey, $return, 300);
        return success($return);
    }


    /**
     * Ck查询数据罗盘统计
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     */
    public function getCompassSummaryCk(Request $request)
    {
        $input                  = $request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['mobile']        = $this->mobile;
        $searchAll              = 0;
        $dpIds                  = [];

        if (isset($input['search_all']) && $input['search_all'] == 1) {
            $adminId    = $this->loginAdminEnterprise['admin_id'];
            $dataRange  = (new Account())->dataRangeList(['admin_id' => $adminId]);
            $dpIds      = array_column($dataRange, 'object_id');
            $memberList = Sdk::work()->setApi("service/v1/department/member-list-by-dp-ids")->post([
                'enterprise_id' => $this->enterpriseId,
                'dp_ids'        => $dpIds,
            ]);
            $useridArr  = array_column($memberList, 'userid');
            $memberInfo = GoMicroScrm::getMemberInfo([
                'enterprise_id' => $this->enterpriseId,
                'mobile'        => $this->mobile,
            ]);
            if (!empty($memberInfo)) {
                $useridArr[] = $memberInfo['userid'];
            }
            $input['user_ids'] = $useridArr;
        }

        if ($input['search_all'] == 0) {
            if (isset($input['shopping_guide']) && !empty($input['shopping_guide'])) {
                $shopping_guide    = array_filter($input['shopping_guide'], function ($item) {
                    return !empty($item);
                });
                $res               = (new StoreWorker())->getList(['enterprise_id' => $this->enterpriseId, 'mem_uids' => $shopping_guide, "is_page" => false]);
                $useridArr         = array_column($res['list'], 'mem_userid');
                $input['user_ids'] = $useridArr;
            }
            // 后续3.3版本之后，shopping_guide的参数判断验证可以去掉
            if (isset($input['userids']) && !empty($input['userids'])) {
                $userids = array_filter($input['userids'], function ($item) {
                    return !empty($item);
                });
                if (!empty($userids)) {
                    $input['user_ids'] = $userids;
                }
            }
        }

        try {
            $input['search_all'] = $searchAll;
            $result = Common::microScrm()->setApi("micro-scrm/v1/data-com-pass/get_compass_summary_ck")->post($input);

        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }

        $return = $result['data'] ?? $result;
        return success($return);
    }
}

