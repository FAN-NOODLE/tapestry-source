<?php
/**
 * 导购员工查询外部联系人类
 * Created by PhpStorm.
 * User: wangle@vchangyi.com
 * Date: 2021/04/26
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Services\GuideAchievementSettingService;
use Exception;
use App\Exceptions\Error;
use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\CustomerScrmService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class CustomerExternalController extends BaseController
{

    /**
     * 顾客列表筛选接口
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $param = $this->getParam();
        try {

            $result = CustomerScrmService::getNewCustomerList($param);

        } catch (Exception $e) {
            Log::info('顾客列表筛选接口失败' . $e);
            throw new Error($e->getCode(), $e->getMessage());
        }

        $result['page'] = $param['page'];
        $result['size'] = $param['size'];
        return success($result);
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'integer|required',
            'store_id'      => 'integer',
            'mobile'        => 'required',
            'page'          => 'integer',
            'size'          => 'integer',
        ], [
            'enterprise_id.*' => '企业id参数不合法',
            'store_id.*'      => '门店id参数不合法',
            'mobile.*'        => '用户手机号参数不合法',
            'page.integer'    => '分页不合法',
            'size.integer'    => '分页数量不合法',
        ]);
        $param = $this->request->only(['page', 'size', 'guide_tag_ids', 'enterprise_id', 'join_time_start', 'join_time_end', 'store_id', 'mobile', 'add_type', 'gender', 'search_key']);

        if (!isset($param['page'])) {
            $param['page'] = Constant::DEFAULT_PAGE;
        }
        if (!isset($param['size'])) {
            $param['size'] = Constant::DEFAULT_SIZE;
        }
        if (isset($param['add_type'])) {
            $param['add_way'] = $param['add_type'] == Constant::ADD_WAY_SCAN_QR_CODE ? Constant::ADD_WAY_SCAN_QR_CODE : Constant::ADD_WAY_OTHER;
        }
        if (isset($param['guide_tag_ids'])) {
            $param['tag_ids'] = explode(',', $param['guide_tag_ids']);
        }

        if (isset($param['join_time_start'])) {
            if (isset($param['join_time_end'])) {
                if ($param['join_time_start'] > $param['join_time_end']) {
                    list($param['join_time_start'], $param['join_time_end']) = [$param['join_time_end'], $param['join_time_start']];
                }
            } else {
                // 未设置结束时间则查询开始时间当天数据
                $param['join_time_end'] = $param['join_time_start'] + 86399;
            }
        } else if (isset($param['join_time_end'])) {
            // 未设置开始时间则查询结束时间当天数据
            $param['join_time_start'] = $param['join_time_end'] - 86399;
        }

        return $param;
    }


    /**
     * 查询外部联系人详情
     *
     * @param Request $request
     *
     * @return array|void
     */
    public function detail(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id', 0);
        if (!$enterpriseId) {
            return error(600000);
        }

        $mobile = $request->input('mobile', 0);
        if (empty($mobile)) {
            return error(600010);
        }
        // 列表返回的外部联系人主键relation_id
        $relationId = $request->input('relation_id', 0);
        if (empty($relationId)) {
            return error(600002);
        }

        $params = ['enterprise_id' => $enterpriseId];

        $guideAchievementSetting = (new GuideAchievementSettingService())->detail($params);
        if (empty($guideAchievementSetting)) {
            $guideAchievementSetting =
                [
                    'enterprise_id'  => $enterpriseId,
                    'binding_type'   => 1,
                    'duration_time'  => 0,
                    'change_owner'   => 1,
                    'protected_time' => 0
                ];
        }

        // 查询顾客基本信息
        $customerInfo = CustomerScrmService::userInfo($enterpriseId, $mobile, $relationId);

        return success([
            'customer_info'             => $customerInfo,
            'guide_achievement_setting' => $guideAchievementSetting,
        ]);

    }


}