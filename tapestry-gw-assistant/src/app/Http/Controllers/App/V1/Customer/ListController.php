<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/14
 * Time: 10:09
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\ExternalSearch;
use App\Micros\MicroScrm\GoMicroScrm;
use App\Micros\Work\CustomerInherit\CustomerInherit;
use App\Micros\Work\Work;
use App\Models\MobileWhite;
use App\Models\StoreMemberViewTab;
use App\Services\MobileWhiteService;
use CrsSdk\Work\Customer;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

class ListController extends BaseController
{

    /**
     * index
     *
     * @desc 顾客列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @date 2020-07-17
     * @throws Exception
     * @author XuHuitao
     */
    public function index(Request $request)
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();
        $params['qy_domain'] = $this->loginAdminEnterprise['domain'];

        // 设置查看时间
        $this->setViewTime($params);

        $data = (new Customer())->filterCustomerList($params);
        return success($data);
    }


    /**
     * 设置查看时间
     *
     * @param $params
     *
     * @return bool
     */
    protected function setViewTime($params)
    {

        //获取店员id
        $workerDetail = $this->storeWorkerDetail($this->enterpriseId, $params['store_id'] ?? 0);
        if (!empty($workerDetail['id'])) {
            $this->saveViewTab($this->enterpriseId, $params['store_id'] ?? 0, $workerDetail['id'] ?? 0, StoreMemberViewTab::CUSTOMER_TAB_TYPE);
        }

        return true;
    }


    /**
     * 顾客列表(新)
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     */
    public function list(Request $request)
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();

        // 设置查看时间
        $this->setViewTime($params);

        $data = (new Customer())->filterCustomerList($params);

        return success($data);
    }

    /**
     * 外部联系人筛选
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     * @Time: 2022/10/21
     */
    public function getChoiceOption(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();

        $range = $this->getRoleRange();
        $params['member_type'] = $range['member_type'];  // 3 导购 4管理
        if ($range['member_type'] == 4) {
            $params['third_ids'] = $range['third_ids'];
        }

        $res = (new GoMicroScrm())->getChoiceOption($params);

        return success($res);
    }

    /**
     * 外部联系人列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     * @Time: 2022/10/21
     */
    public function userList(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $params['is_group'] = true;   // 是否分组 true 是
        $params['member_type'] = 3; // 3 导购 4管理
        if ($this->isManager == 1) {
            $params['member_type'] = 4;
            $dataRange = $this->dataRangeList;
            $third_id = array_column($dataRange, "object_id");
            $params['third_ids'] = $third_id;
        }
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $encodeParams = json_encode($params);
        $redisKey = md5($encodeParams);
        $listData = getOrSetData($redisKey);
        if (!$listData) {
            $listData = $this->getExternalUserList($params);
            getOrSetData($redisKey, $listData, 300);
        }
        return success($listData);

    }

    /**
     * @param       $params
     * @return array|JsonResponse|mixed
     * @throws Exception
     */
    public function getExternalUserList($params)
    {

        if (!in_array($params['enterprise_id'],[3363603667487424,3360610327843968])){
            $data = (new GoMicroScrm())->getExternalUserList($params);
            if (empty($data['list'])) {
                return $data;
            }
            $args["enterprise_id"]=$params['enterprise_id'];
            $args['external_ids'] = array_column($data['list'], 'work_external_user_id');
            $recordForExternal = (new CustomerInherit())->recordForExternalIds($args);

            foreach ($data['list'] as $users => $user) {

                $data['list'][$users]['is_new_inherit'] = 0;
                // 分配列表有记录并且接替的导购和现在关系的导购一样
                foreach ($recordForExternal as $item) {
                    if ($item['takeover_member_id'] == $user['member_id'] &&
                        $item['external_id'] == $user['work_external_user_id']) {
                        $data['list'][$users]['is_new_inherit'] = 1;
                    }
                }
            }
        }else{
            $params['is_search'] = 1;//这里默认需要检索，否则会查出所有的数据
            $data = (new ExternalSearch())->list($params);
            if (empty($data['list'])) {
                return $data;
            }

            foreach ($data['list'] as $users => $user) {
                $data['list'][$users]['is_new_inherit'] = 0;
            }
            $data['list'] = self::formatCoachCustomerSearch($data['list']);
        }
        return $data;
    }

    /**
     * 兼容老的前端数据结构
     *
     * @param $result
     * @return array
     */
    private static function formatCoachCustomerSearch($result): array
    {
        $cardDescriptionMap = [
            "" => "0",
            "PROSPECT" => "1",
            "C" => "2",
            "D" => "3",
            "GOLD" => "4",
            "SILVER" => "5",
            "DIAMOND" => "6"
        ];

        foreach ($result as $k => $v) {
            $result[$k]['customer_card'] = $cardDescriptionMap[$v['card_level']] ?? "0";
            $result[$k]['last_12m_purchase_amount'] = $v['amount_total'];
            $result[$k]['client_type'] = $v['clt_type'] ?? '';
            $result[$k]['is_relation'] = $v['self_relation_status'] == 0 ? 1 : 0;
            $result[$k]['key_sa_member_name'] = $v['key_sa_name'];
        }
        return $result;
    }

    /**
     * 外部联系人列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     * @Time: 2022/10/21
     */
    public function userListCk(Request $request): JsonResponse
    {
        $params = $request->input();
        $args['enterprise_id'] = $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();
        $params['is_group'] = true;   // 是否分组 true 是

        $range = $this->getRoleRange();
        $params['member_type'] = $range['member_type'];  // 3 导购 4管理
        if ($range['member_type'] == 4) {
            $params['third_ids'] = $range['third_ids'];
        }

        $data = (new GoMicroScrm())->getExternalUserListCk($params);

        if (empty($data['list'])) {
            return success($data);
        }
        $args['external_ids'] = array_column($data['list'], 'work_external_user_id');
        $recordForExternal = (new CustomerInherit())->recordForExternalIds($args);

        $cardDescription = ['O' => 0,
            'PROSPECT' => 1,
            'C' => 2,
            'D' => 3,
            'GOLD' => 4,
            'SILVER' => 5,
            'DIAMOND' => 6];

        foreach ($data['list'] as $users => $user) {
            //String	会员等级
            $data['list'][$users]['customer_card'] = $cardDescription[$user['card_level']] ?? 0;
            $data['list'][$users]['last_12m_purchase_amount'] = $user['amount_total'] / 100;
            $data['list'][$users]['client_type'] = $user['clt_type'];
            $data['list'][$users]['customer_code'] = $user['union_code'];

            $data['list'][$users]['is_new_inherit'] = 0;
            // 分配列表有记录并且接替的导购和现在关系的导购一样
            foreach ($recordForExternal as $item) {
                if ($item['takeover_member_id'] == $user['member_id'] &&
                    $item['external_id'] == $user['work_external_user_id']) {
                    $data['list'][$users]['is_new_inherit'] = 1;
                }
            }
        }

        return success($data);

    }

    /**
     * 三方接口会员聚合查询
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     */
    public function customerListForCDP(Request $request): JsonResponse
    {
        $params['params'] = $request->input();
        $param['enterprise_id'] = $params['enterprise_id'] = $this->enterpriseId;

        // 获取导购权限管理的导购userid
        $param['mobile'] = $this->loginMobile();

        $range = $this->getRoleRange();
        $param['member_type'] = $range['member_type'];  // 3 导购 4管理
        if ($range['member_type'] == 4) {
            $param['third_ids'] = $range['third_ids'];
        }
        $useridArr = (new GoMicroScrm())->getMemberUserIds($param);
        $useridString = implode(',', $useridArr);

        // 三方接口查询获取数据
        $params['params']['curPage'] = $params['params']['page'];
        $params['params']['pageSize'] = $params['params']['size'];
        $params['params']['saAccount'] = $useridString;

        $customerList = (new Work())->customerList($params);
        if ($customerList['totalSize'] == 0) {
            return success(
                ['list' => [],
                    'total' => $customerList['totalSize'],
                    'page' => $customerList['curPage'],
                    'size' => $customerList['pageSize'],]
            );
        }
        $param['union_codes'] = [];
        if (isset($customerList['wecustomerList']) && count($customerList['wecustomerList']) > 0) {
            // 判断是一维数组
            if (count($customerList['wecustomerList']) == count($customerList['wecustomerList'], 1)) {
                $customerList['wecustomerList'] = [$customerList['wecustomerList']];
            }
            foreach ($customerList['wecustomerList'] as $v) {
                $param['union_codes'][] = $v['memberId'];
            }
        }

        $data['list'] = $customerList['wecustomerList'] ?? [];
        $data['total'] = $customerList['totalSize'];
        $data['page'] = $customerList['curPage'];
        $data['size'] = $customerList['pageSize'];

        $userData = (new GoMicroScrm())->getCustomerList($param);

        $cardDescription = ['O' => 0,
            'PROSPECT' => 1,
            'C' => 2,
            'D' => 3,
            'GOLD' => 4,
            'SILVER' => 5,
            'DIAMOND' => 6];
        // 匹配
        $customerLists = array_column($userData, null, 'union_code');
        foreach ($data['list'] as $k => $v) {
            $data['list'][$k]['work_external_user_id'] = '';
            $data['list'][$k]['member_id'] = '';
            $data['list'][$k]['euid'] = '';
            $data['list'][$k]['remark'] = '';
            $data['list'][$k]['nick'] = '';
            $data['list'][$k]['name'] = $v['lastName'] . $v['firstName'];
            $data['list'][$k]['unionid'] = '';
            $data['list'][$k]['avatar'] = '';
            $data['list'][$k]['openid'] = '';
            $data['list'][$k]['birthday'] = '';
            $data['list'][$k]['union_code'] = '';
            $data['list'][$k]['mobile'] = '';
            $data['list'][$k]['member_name'] = '';
            $data['list'][$k]['is_key_sa'] = '';
            $data['list'][$k]['sa_user_id'] = '';
            $data['list'][$k]['key_sa_member_id'] = '';
            $data['list'][$k]['key_sa_member_name'] = '';
            $data['list'][$k]['is_relation'] = '';
            $data['list'][$k]['card_level'] = '';
            $data['list'][$k]['clt_type'] = '';
            $data['list'][$k]['amount_total'] = '';

            $data['list'][$k]['customer_card'] = 0;
            $data['list'][$k]['last_12m_purchase_amount'] = 0;
            $data['list'][$k]['client_type'] = '';
            $data['list'][$k]['customer_code'] = '';
            if (isset($customerLists[$v['memberId']])) {
                $temp = $customerLists[$v['memberId']];
                $data['list'][$k]['work_external_user_id'] = $temp['work_external_user_id'];
                $data['list'][$k]['member_id'] = $temp['member_id'];
                $data['list'][$k]['euid'] = $temp['euid'];
                $data['list'][$k]['remark'] = $temp['remark'];
                $data['list'][$k]['nick'] = $temp['nick'];
                $data['list'][$k]['name'] = $temp['name'];
                $data['list'][$k]['unionid'] = $temp['unionid'];
                $data['list'][$k]['avatar'] = $temp['avatar'];
                $data['list'][$k]['openid'] = $temp['openid'];
                $data['list'][$k]['birthday'] = $temp['birthday'];
                $data['list'][$k]['union_code'] = $temp['union_code'];
                $data['list'][$k]['mobile'] = $temp['mobile'];
                $data['list'][$k]['member_name'] = $temp['member_name'];
                $data['list'][$k]['is_key_sa'] = $temp['is_key_sa'];
                $data['list'][$k]['sa_user_id'] = $temp['sa_user_id'];
                $data['list'][$k]['key_sa_member_id'] = $temp['key_sa_member_id'];
                $data['list'][$k]['key_sa_member_name'] = $temp['key_sa_member_name'];
                $data['list'][$k]['is_relation'] = $temp['is_relation'];
                $data['list'][$k]['card_level'] = $temp['card_level'];
                $data['list'][$k]['clt_type'] = $temp['clt_type'];
                $data['list'][$k]['amount_total'] = $temp['amount_total'];

                //String	会员等级
                $data['list'][$k]['customer_card'] = $cardDescription[$temp['card_level']] ?? 0;
                $data['list'][$k]['last_12m_purchase_amount'] = $temp['amount_total'] / 100;
                $data['list'][$k]['client_type'] = $temp['clt_type'];
                $data['list'][$k]['customer_code'] = $temp['union_code'];
            }
        }

        return success($data);

    }

    /**
     * 会员聚合查询列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     */
    public function customerList(Request $request): JsonResponse
    {
        $enterpriseId = $this->enterpriseId;
        $param['enterprise_id'] = $enterpriseId;
        // 获取导购权限管理的导购userid
        $param['mobile'] = $this->mobile;
        $param['is_group'] = true;   // 是否分组 true 是
        $param['member_type'] = 3; // 3 导购 4管理
        if ($this->isManager == 1) {
            $param['member_type'] = 4;
            $dataRange = $this->dataRangeList;
            $third_id = array_column($dataRange, "object_id");
            $param['third_ids'] = $third_id;
        }
        $param['member_id'] = $this->memberAuthInfo['member_id'];
        $useridArr = (new GoMicroScrm())->getMemberUserIds($param);
        $useridString = implode(',', $useridArr);
        $size = $request->input('size', 10);
        $page = $request->input('page', 1);
        $res = $this->getCrmUnionCodes($request->input(), $enterpriseId, $useridString, 1, 2000);
        $param['union_codes'] = $res['union_codes'];
        $param['page'] = $page;
        $param['size'] = $size;
        // 未查询到结果
        if (count($param['union_codes']) == 0) {
            return success(['list' => [], 'total' => 0, 'page' => $param['page'], 'size' => $param['size'],]);
        }
        $customerList = $res['wecustomerList'];
        $customerLists = array_column($customerList['wecustomerList'], null, 'memberId');
        $param['search_type'] = 1;
        $param['is_self'] = (bool)$request->input('is_self', 0);
        $listData = $this->getExternalUserList($param);
        foreach ($listData['list'] as $k => $v) {
            if (!empty($v['union_code']) && isset($customerLists[$v['union_code']])) {
                $temp = $customerLists[$v['union_code']];
                $listData['list'][$k]['name'] = $temp['lastName'] . $temp['firstName'];
            }
        }

        return success($listData);

    }

    /**
     * 整理请求crm获取unionCodes
     *
     * @param $input
     * @param $enterpriseId
     * @param $userId
     * @param $page
     * @param $size
     * @return array
     */
    protected function getCrmUnionCodes($input, $enterpriseId, $userId, $page, $size): array
    {
        $params['params'] = $input;
        $params['params']['curPage'] = $page;
        $params['params']['pageSize'] = $size;
        $params['params']['saAccount'] = $userId;
        $params['enterprise_id'] = $enterpriseId;
        $customerList = (new Work())->customerList($params);
        $data = [
            'union_codes' => [],
            "wecustomerList" => [],
        ];
        if (isset($customerList['wecustomerList']) && count($customerList['wecustomerList']) > 0) {
            // 判断是一维数组
            if (count($customerList['wecustomerList']) == count($customerList['wecustomerList'], 1)) {
                $customerList['wecustomerList'] = [$customerList['wecustomerList']];
            }
            $data["wecustomerList"] = $customerList;
            foreach ($customerList['wecustomerList'] as $v) {
                $data["union_codes"][] = $v['memberId'];
            }
        }

        return $data;
    }

    /**
     * 客户详情
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/10/25
     */
    public function userDetail(Request $request): JsonResponse
    {
        $params = $request->input();
        // 验证输入参数
//        $rules = [
//            'external_userid' => 'required|string',
//            'customer_code'   => 'required|string',
//        ];
//        $message = [
//            'external_userid.required' => 'external_userid参数必传',  // openid
//        ];
//        $this->validate($request, $rules, $message);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $params['external_userid'] = $params['external_userid'] ?? "";
        $params['customer_code'] = $params['customer_code'] ?? "";
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params['member_type'] = 3; // 3 导购 4管理
        if ($this->isManager == 1) {
            $params['member_type'] = 4;
            $dataRange = $this->dataRangeList;
            $third_id = array_column($dataRange, "object_id");
            $params['third_ids'] = $third_id;
        }

        // redis 缓存用户数据
//        if ($params['external_userid']) {
//            $redisKey = "USER:CRM:" . $params['enterprise_id'] . ":info:" . $params['external_userid'];
//            $userCrm  = getOrSetData($redisKey);
//            if ($userCrm) {
//                return success($userCrm);
//            }
//        }
        // 用户详情
        $userDetail = (new GoMicroScrm())->externalUserDetail($params);
        $userCrm = (new Work())->getCrmUserInfo($params['enterprise_id'], $userDetail['union_id'], $userDetail['user_id'] ?? 0, $params['customer_code']);

        $userDetail['customer_card'] = $userDetail['card_level'] == "" ? "0" : $userDetail['card_level'];
        if ($userDetail['config_store_name'] != 'coach') {
            $userDetail['customer_card'] = $userDetail['card_level'] ?? "";
        }
        // 请求crm获取会员信息
        if ($userCrm['status']) {
            $userDetail['customer_code'] = $userCrm['customer_code'];
            $userDetail['customer_card'] = $userCrm['customer_card'] ?? "";
            $userDetail['customer_card_str'] = $userCrm['customer_card_str'] ?? "";
            $userDetail['membership'] = $userCrm['membership'];
            $userDetail['customer_type'] = $userCrm['customer_type'] ?: $userDetail['customer_type'];
            $userDetail['name'] = $userCrm['name'];
            $userDetail['birthday_str'] = $userCrm['birthday'];
            $userDetail['mobile'] = $userCrm['mobile'];
            $userDetail['expiry_date'] = $userCrm['expiry_date'];
            $userDetail['upgrade_date_time'] = $userCrm['upgrade_date_time'];
            $userDetail['upgrade_need_points'] = $userCrm['upgrade_need_points'] ?? 0;
            $userDetail['last_12m_purchase_amount'] = $userCrm['last_12m_purchase_amount'] ?? 0;
            $userDetail['last_36m_purchase_amount'] = $userCrm['last_36m_purchase_amount'] ?? 0;
            $userDetail['total_valid_points'] = $userCrm['total_valid_points'];
            $userDetail['client_type'] = $userCrm['client_type'] ?? "";
            $userDetail['vip_status'] = $userCrm['vipStatus'] ?? [];
            $userDetail['city'] = $userCrm['city'] ?? '';
            $userDetail['address'] = $userCrm['address'] ?? '';

            $userDetail['total_inta_valid_points'] = $userCrm['total_inta_valid_points'] ?? 0; // 可兑换积分
            $userDetail['pts_need_to_upgrade'] = $userCrm['pts_need_to_upgrade'] ?? 0; // 升级所需消费金额
            $userDetail['pts_need_to_renew'] = $userCrm['pts_need_to_renew'] ?? 0; // 保级所需消费金额
        }

        // 设置缓存
        if ($params['external_userid']) {
            $redisKey = "USER:CRM:" . $params['enterprise_id'] . ":info:" . $params['external_userid'];
            getOrSetData($redisKey, $userDetail, Work::RedisCacheTime);
        }
        return success($userDetail);

    }

    /**
     * 顾客手机号
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @Time: 2022/12/6
     */
    public function userMobile(Request $request): JsonResponse
    {
        $params = $request->input();
        // 验证输入参数
//        $rules   = [
//            'unionid' => 'required|string',
//        ];
//        $message = [
//            'unionid.required' => 'unionid参数必传',
//        ];
//        $this->validate($request, $rules, $message);

        $params['enterprise_id'] = $this->enterpriseId;
        $params['union_id'] = $params['unionid'] ?? "";
        $params['customer_code'] = $params['customer_code'] ?? "";

        $userCrm = (new Work())->getCrmUserInfo($params['enterprise_id'], $params['union_id'], 0, $params['customer_code']);

        if ($userCrm['status']) {
            return success(['mobile' => $userCrm['mobile']]);

        }
        return success(['mobile' => '']);

    }

    /**
     * 顾客订单
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/11
     */
    public function userCrmOrderList(Request $request): JsonResponse
    {
        $params = $request->input();
        // 验证输入参数
//        $rules   = [
//            'union_id' => 'required|string',
//            'customer_code'   => 'required|string',
//        ];
//        $message = [
//            'union_id.required' => 'union_id参数必传',
//        ];
//        $this->validate($request, $rules, $message);

        $params['enterprise_id'] = $this->enterpriseId;
        // 查询全部订单
        $params['startDateTime'] = 0;
        $params['endDateTime'] = time() * 1000;
        $params['union_id'] = $params['union_id'] ?? "";
        $params['customer_code'] = $params['customer_code'] ?? "";
        // 查询crm
        $userOrder = (new Work())->getCrmUserOrderList($params);

        return success($userOrder);
    }

    /**
     * 顾客卡劵
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/11
     */
    public function userCrmCouponsList(Request $request): JsonResponse
    {
        $params = $request->input();
        // 验证输入参数
//        $rules   = [
//            'union_id' => 'required|string',
//            'customer_code'   => 'required|string',
//        ];
//        $message = [
//            'union_id.required' => 'union_id参数必传',
//        ];
//        $this->validate($request, $rules, $message);

        $params['enterprise_id'] = $this->enterpriseId;
        $params['channel_type'] = ''; // 渠道 RETAIL: 零售店 OUTLET：奥莱店
        $params['status'] = 1; // 优惠券状态 1未使用 2已使用 3过期 0 all
        $params['source'] = 1; // 请求来源： 1会员中心 2 零售小程序3 奥莱小程序 4 企业微信
        $params['pageNum'] = isset($input['page']) && $input['page'] > 0 ? $input['page'] : 1;
        $params['pageSize'] = isset($input['size']) && $input['size'] > 0 ? $input['size'] : 99;
        $params['union_id'] = $params['union_id'] ?? "";
        $params['customer_code'] = $params['customer_code'] ?? "";
        // 查询crm
        $userCoupons = (new Work())->getCrmUserCouponsList($params);

        return success($userCoupons);
    }


    /**
     * 顾客卡劵
     *
     * @param Request $request
     *
     * @return array|mixed
     * @throws ValidationException
     * @Time: 2022/11/11
     */
    public function swUserCrmCouponsList(Request $request)
    {
        $params = $request->input();
        // 验证输入参数
        $rules   = [
            'union_id' => 'required|string',
        ];
        $message = [
            'union_id.required' => 'union_id参数必传',
        ];
        $this->validate($request, $rules, $message);
        // 查询crm
        return (new Work())->swCrmCouponsList($params);
    }

    /**
     * 获取顾客邀约和播种数量
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/17
     */
    public function userDetailTipsNum(Request $request): JsonResponse
    {
        $params = $request->input();
        // 验证输入参数
        $rules = [
            'external_userid' => 'required|string',
        ];
        $message = [
            'external_userid.required' => 'external_userid不能为空',
        ];
        $this->validate($request, $rules, $message);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();

        $params['member_type'] = 3;
        // 这里身份的验证以token为主,is_manager为店长
        if ($this->isManager) {
            $params['member_type'] = 4;
            $dataRange = $this->dataRangeList;
            $params['third_ids'] = array_column($dataRange, "object_id");
        }
//        $range                 = $this->getRoleRange();
//        $params['member_type'] = $range['member_type'];  // 3 导购 4管理
//        if ($range['member_type'] == 4) {
//            $params['third_ids'] = $range['third_ids'];
//        }

        $res = (new GoMicroScrm())->getUserDetailTipsNum($params);

        return success($res);

    }

    /**
     * 顾客和导购的绑定关系
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/17
     */
    public function userBindRelation(Request $request): JsonResponse
    {
        $params = $request->input();
        // 验证输入参数
        $rules = [
            'external_userid' => 'required|string',
        ];
        $message = [
            'external_userid.required' => 'external_userid不能为空',
        ];
        $this->validate($request, $rules, $message);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();

        $range = $this->getRoleRange();
        $params['member_type'] = $range['member_type'];  // 3 导购 4管理
        if ($range['member_type'] == 4) {
            $params['third_ids'] = $range['third_ids'];
        }

        $res = (new GoMicroScrm())->getUserBindRelation($params);

        return success($res);

    }

    /**
     * 修改备注
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/18
     */
    public function userUpdateRemark(Request $request): JsonResponse
    {
        $params = $request->input();
        // 验证输入参数
        $rules = [
            'external_userid' => 'required|string',
            'remark' => 'required|string',
        ];
        $message = [
            'external_userid.required' => 'external_userid不能为空',
            'remark.required' => 'remark不能为空',
        ];
        $this->validate($request, $rules, $message);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        // 删除缓存
        $redisKey = "USER:CRM:" . $params['enterprise_id'] . ":info:" . $params['external_userid'];
        delRedisKey($redisKey);

        $res = (new GoMicroScrm())->userUpdateRemark($params);

        return success($res);

    }


    /**
     * 修改用户的描述
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/18
     */
    public function userUpdateDescription(Request $request): JsonResponse
    {
        $params = $request->input();
        // 验证输入参数
        $rules = [
            'external_userid' => 'required|string',
        ];
        $message = [
            'external_userid.required' => 'external_userid不能为空',
        ];
        $this->validate($request, $rules, $message);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();
        $params['member_id'] = $this->memberAuthInfo['member_id'];

        // 删除缓存
        $redisKey = "USER:CRM:" . $params['enterprise_id'] . ":info:" . $params['external_userid'];
        delRedisKey($redisKey);

        $res = (new GoMicroScrm())->userUpdateDescription($params);

        return success($res);

    }


    /**
     * 操作客户印象 标签
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/21
     */
    public function userWorkTagsOperate(Request $request): JsonResponse
    {
        $params = $request->input();

        if (empty($params['add_tag_ids']) && empty($params['del_tag_ids'])) {
            error(66001, "操作标签参数不能为空");
        }
        $params['enterprise_id'] = $this->enterpriseId;
        $params['operate_mobile'] = $this->mobile;
        $params['member_id'] = $this->memberAuthInfo["member_id"];
        if (isset($params['external_userid'])) {
            // 删除缓存
            $redisKey = "USER:CRM:" . $params['enterprise_id'] . ":info:" . $params['external_userid'];
            delRedisKey($redisKey);
        }

        $res = (new GoMicroScrm())->userWorkTagsOperate($params);

        return success($res);

    }

    /**
     * 操作客户印象 标签
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/21
     */
    public function userImpressionTagsOperate(Request $request): JsonResponse
    {
        $params = $request->input();

        if (empty($params['add_tag_ids']) && empty($params['del_tag_ids'])) {
            error(66001, "操作标签参数不能为空");
        }
        $params['enterprise_id'] = $this->enterpriseId;
        $params['operate_mobile'] = $this->loginMobile();
        if (isset($params['external_userid'])) {
            // 删除缓存
            $redisKey = "USER:CRM:" . $params['enterprise_id'] . ":info:" . $params['external_userid'];
            delRedisKey($redisKey);
        }

        $res = (new GoMicroScrm())->userImpressionTagsOperate($params);

        return success($res);

    }

    /**
     * 会员列表聚合查询条件
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     */
    public function customerOption(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();

        $data = (new Work())->customerOption($params);

        return success($data);
    }

    /**
     * 会员查询白名单
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function mobileWhiteList(Request $request): JsonResponse
    {
        $params = $request->all();
        $type= $params['type'] ?? MobileWhite::CUSTOMER_SEARCH;

        if ($type == MobileWhite::IS_ESA) {
            return success(['show' => $this->memberAuthInfo['store_code'] == 'ESA' ? 0 : 1]);
        }
        $res = MobileWhiteService::WhiteList($this->enterpriseId,$this->mobile,$type);
        return success(['show' => $res ? 1 : 0]);

    }

    /**
     * 修改顾客地址生日鞋码
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     */
    public function saveInfoAdd(Request $request): JsonResponse
    {
        $params = $request->input();
        // 验证输入参数
        $rules = [
            'external_userid' => 'required|string',
            'enterprise_user_id' => 'required|int',
            'address' => 'required|string',
            'shoe_size' => 'required|string',
            'month' => 'required|int',
            'day' => 'required|int',
        ];
        $message = [
            'external_userid.required' => 'external_userid不能为空',
            'enterprise_user_id.required' => 'enterprise_user_id不能为空',
            'address.required' => 'address不能为空',
            'shoe_size.required' => 'shoe_size不能为空',
            'month.required' => 'month不能为空',
            'day.required' => 'day不能为空',
        ];
        $this->validate($request, $rules, $message);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();

        // 删除缓存
        $redisKey = "USER:CRM:" . $params['enterprise_id'] . ":info:" . $params['external_userid'];
        delRedisKey($redisKey);

        $memberInfo = $this->memberInfo;
        $params['operate_user_id'] = $memberInfo['userid'] ?? '';
        $params['operate_name'] = $memberInfo['name'] ?? '';

        $res = (new GoMicroScrm())->userSaveInfoAdd($params);

        return success($res);

    }

    /**
     * 全选的外部联系人总数
     * @param Request $request
     * @return JsonResponse
     */

    public function selectAllNum(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params['enterprise_id'] = $this->enterpriseId;
        $params['is_group'] = true;
        $params['member_type'] = 3; // 3 导购 4管理
        if ($this->isManager == 1) {
            $params['member_type'] = 4;
            $dataRange = $this->dataRangeList;
            $params['third_ids'] = array_column($dataRange, "object_id");
        }
        $params['userid'] = $this->memberAuthInfo['userid'];
        try {
            $res = (new GoMicroScrm())->getExternalUserTotal($params);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success(["num" => $res]);

    }

    /**
     * 全选的外部联系人ids
     * @param Request $request
     * @return JsonResponse
     */
    public function selectAllList(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_type'] = 3; // 3 导购 4管理
        if ($this->isManager == 1) {
            $params['member_type'] = 4;
            $dataRange = $this->dataRangeList;
            $params['third_ids'] = array_column($dataRange, "object_id");
        }
        $params['userid'] = $this->memberAuthInfo['userid'];
        try {
            $res = (new GoMicroScrm())->getExternalUserIds($params);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success(['external_ids' => $res]);
    }

    /**
     * 会员查询创建群发获取全选人数
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function customerAllNum(Request $request): JsonResponse
    {
        $param = $this->customerSearchCondition($request);

        if (count($param['union_codes']) == 0) {
            return success(['num' => 0]);
        }
        $num=0;
        try {
            $param['search_type'] = 1;
            $param['is_self'] = (bool)$request->input('is_self',1);
            $res =  (new ExternalSearch())->getOpenIds($param);
            $num=count($res);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success(['num' => $num]);
    }

    protected function customerSearchCondition(Request $request)
    {
        $condition = $request->input();
        $param = $request->input();
        $enterpriseId = $this->enterpriseId;
        $condition['enterprise_id'] = $param['enterprise_id'] = $enterpriseId;
        $param['member_type'] = 3; // 3 导购 4管理
        if ($this->isManager == 1) {
            $param['member_type'] = 4;
            $dataRange = $this->dataRangeList;
            $third_id = array_column($dataRange, "object_id");
            $param['third_ids'] = $third_id;
        }
        $useridArr = (new GoMicroScrm())->getMemberUserIds($param);
        $useridString = implode(',', $useridArr);
        $res = $this->getCrmUnionCodes($condition, $enterpriseId, $useridString, 1, 2000);
        $param['union_codes'] = $res['union_codes'];
        $param['member_id'] = $this->memberAuthInfo['member_id'];
        return $param;
    }

    /**
     * 会员查询的条件整理
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function customerAllList(Request $request): JsonResponse
    {
        $param = $this->customerSearchCondition($request);
        if (count($param['union_codes']) == 0) {
            return success(['external_ids' => []]);
        }
        try {
            $res = (new GoMicroScrm())->getExternalUserIds($param);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success(['external_ids' => $res]);
    }
}