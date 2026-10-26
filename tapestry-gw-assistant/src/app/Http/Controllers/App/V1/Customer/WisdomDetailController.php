<?php
/**
 * Created by PhpStorm.
 * User: BaiZhongYang
 * Date: 2020/7/17
 * Time: 17:59
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Common\Constant;
use CrsSdk\Shop\User;
use CrsSdk\Work\Customer;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;
use Illuminate\Support\Facades\Log;

class WisdomDetailController extends BaseController
{
    /**
     * 有效类型 1:永久有效 2:有效期
     */
    const BINDING_TYPE_ONE = 1;
    const BINDING_TYPE_TWO = 2;
    /**
     * 导购是否自己 0:否 1:是
     */
    const BINDING_UN_OWN_STATUS = 0;
    const BINDING_OWN_STATUS = 1;

    /**
     * 智慧导购顾客详情
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/3 11:05
     */
    public function index(Request $request)
    {
        $param['union_id']      = $request->input('union_id');
        $param['relation_id']   = $request->input('relation_id');
        $param['user_id']       = $request->input('user_id');
        $param['store_id']       = $request->input('store_id',0);
        $param['mem_uid']       = $request->input('mem_uid');
        $param['enterprise_id'] = $this->enterpriseId;
        $param['qy_domain']     = $this->loginAdminEnterprise['domain'];
        $param['qy_hash']       = $this->loginAdminEnterprise['hash'];
        Log::info('顾客详情请求参数', [$param]);

        //获取店员id
        $workerDetail = $this->storeWorkerDetail($this->enterpriseId, $param['store_id']);

        $result = (new Customer())->customerWisdomDetail($param);

        // 添加方式
        $result['customer_info']['add_type'] = Constant::ADD_TYPE_BUSINESS_CARD_OTHER;
        if ($result['customer_info']['source_type'] == $workerDetail['contact_qr_code_id']) {
            $result['customer_info']['add_type'] = Constant::ADD_TYPE_BUSINESS_CARD;
        }

        $userId = $result['customer_info']['user_id'] ?? 0;

        // 获取最后登录时间以及购买时间
        $statisticsDetail = (new User())->statisticsDetail(['user_id' => $userId]);

        $userStatistics=$statisticsDetail['data'];

        // 计算平均订单价格
        $orderNum = $userStatistics['total_order_num'];
        if ($orderNum == 0) {
            $average = 0;
        } else {
            $average = $userStatistics['total_order_price'] / ($orderNum * 100);
            $average = sprintf("%.2f", $average) * 100;
        }

        // 拼接订单相关数据
        $result['customer_info']                        = array_merge($result['customer_info'], $userStatistics);
        $result['customer_info']['online_order_num']    = $orderNum;
        $result['customer_info']['online_order_amount'] = $userStatistics['total_order_amount'];
        $result['customer_info']['average']             = $average;
        $result['worker_info'] = [];
        if (!empty($param['user_id'])) {
            // 组装传递参数
            $bindingParam = [
                'enterprise_id' => $this->enterpriseId,
                'user_id'       => $userId
            ];
            // 获取导购数据
            $result['worker_info'] = (new StoreWorker())->customerBindingWorkerDetail($bindingParam);
        }

        $flage=false;
        $endTime=0;
        // 修改对应导购详情
        if (!empty($result['worker_info'])) {
            $result['worker_info']['end_time'] = 0;
            $settingType                       = self::BINDING_TYPE_TWO;

            // 当前导购和返回导购是同一个人
            $workerUserId = $workerDetail['user_id']??'';

            // 导购是自己
            if ($result['worker_info']['user_id'] == $param['user_id'] || $result['guide_achievement_setting']['binding_type'] == self::BINDING_TYPE_ONE) {
                $settingType = self::BINDING_TYPE_ONE;
            }
            // 如果是有有效期的
            if ($settingType == self::BINDING_TYPE_TWO) {
                $endTime=$result['worker_info']['bing_time'] + $result['guide_achievement_setting']['duration_time'] * 86400;
                $result['worker_info']['end_time'] =$endTime;
                $flage=true;
            }

            $isOwn = self::BINDING_UN_OWN_STATUS;
            // 如果是导购自己
            if ($result['worker_info']['worker_id'] == $workerDetail['id']) {
                $isOwn = self::BINDING_OWN_STATUS;
            }
            $result['worker_info']['is_own'] = $isOwn;
            $result['worker_info']['setting_type'] = $settingType;
        }
        Log::info('顾客详情返回参数', [$result]);
        // 有效期过期问题
        if ($endTime < time() && $flage) {
            $result['worker_info'] = [];
        }

        return success($result);
    }

}
