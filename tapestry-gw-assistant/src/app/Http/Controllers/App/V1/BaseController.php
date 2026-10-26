<?php


namespace App\Http\Controllers\App\V1;


use App\Common\Constant;
use App\Http\Controllers\Controller;
use App\Micros\MicroScrm\GoMicroScrm;
use App\Models\StoreMemberViewTab;
use CrsSdk\Work\PartnerMember;
use CrsSdk\Work\StoreWorker;
use Exception;
use Illuminate\Contracts\Auth\Authenticatable;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;

class BaseController extends Controller
{
    const  USERTYPE_WORKER = 1;
    const  USERTYPE_PARTNER = 2;

    const MEMBER_TYPE_SA = 3;//导购
    const MEMBER_TYPE_SM = 4;//管理
    /**
     * 分享商品类型
     */
    const  SHARE_TYPE = 'isg';
    /**
     * 是否必须登录
     *
     * @var bool
     */
    protected $isRequireLogin = true;
    /**
     * 当前的缓存的管理的企业信息数据
     *
     * @var array|Authenticatable|null
     */
    protected $loginAdminEnterprise = [];
    /**
     * 登录管理员ID
     *
     * @var int
     */
    protected $adminId = 0;
    /**
     * 登录管理员角色
     *
     * @var string
     */
    protected $role = '';
    /**
     * 登录SA是否企业管理
     *
     * @var string
     */
    protected $isManager = 0;
    /**
     * 登录SM的权限管理范围
     *
     * @var string
     */
    protected $dataRangeList = [];
    /**
     * 登录管理员手机号
     *
     * @var string
     */
    protected $mobile = '';
    /**
     * 登录导购信息
     *
     * @var string
     */
    protected $memberAuthInfo = [];
    /**
     * 当前所选登录企业ID
     *
     * @var int
     */
    protected $enterpriseId = 0;
    /**
     * 当前登陆人身份
     *  1导购 2合伙人
     * @var int
     */
    protected $userType = 0;
    /**
     * @var int 小程序类型 1 标准 2 灰度 3 爱点击主体
     */
    protected $identifier = 1;

    protected $memberInfo = [];

    /**
     * BaseController constructor.
     *
     * @throws Exception
     */
    public function __construct()
    {
        parent::__construct();

        $this->identifier = $this->request->header('identifier', Constant::STANDARD_IDENTIFIER);

        // 如果必须登录
        if ($this->isRequireLogin) {
            $loginCacheData = Auth::user();
            if (empty($loginCacheData)) {
                error(652004);
            }
            $this->adminId = $loginCacheData['choiceEnterprise']['admin_id'];
            $this->isManager = $loginCacheData['is_manager'];
            $this->mobile = $loginCacheData['mobile'];
            $this->enterpriseId = $loginCacheData['enterprise_id'];
            $this->loginAdminEnterprise = $loginCacheData['choiceEnterprise'];
            $this->dataRangeList = $loginCacheData['data_range_list'];
            $this->memberAuthInfo = $loginCacheData['member_auth_info'] ?? [];
            $this->userType = self::USERTYPE_WORKER;
        }
    }

    /**
     * 获取登录的手机号
     * @return string
     * author: zhaojie<zhaojie@vchangyi.com>
     * create: 2020/5/12 14:11
     */
    public function loginMobile(): string
    {
        try {
            $token = $this->request->header('token');
            if (empty($token)) {
                error(652006);
            }
            $tokenData = json_decode(Crypt::decrypt($token), true);
            return $tokenData['mobile'];
        } catch (Exception $e) {
            Log::error("获取登录员工手机号错误，错误信息为：{$e}");
            error(655001);
        }
        return '';
    }

    /**
     * 数据格式转化
     *
     * @param $data
     * @param string $key1
     * @param string $key2
     * @param string $key3
     * @param string $key4
     *
     * @return array
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/21 19:22
     */
    public function transformDataFormat(
        $data,
        $key1 = 'extension_num',
        $key2 = 'extension_amount',
        $key3 = 'service_num',
        $key4 = 'service_amount',
        $key5 = 'staff_order_num',
        $key6 = 'partner_order_num'
    )
    {
        $x = array_keys($data['detail_data']);
        $y = [
            'data1' => array_column($data['detail_data'], $key1) ?? [],
            'data2' => array_column($data['detail_data'], $key2) ?? [],
            'data3' => array_column($data['detail_data'], $key3) ?? [],
            'data4' => array_column($data['detail_data'], $key4) ?? [],
            'data5' => array_column($data['detail_data'], $key5) ?? [], //纯属导购订单数
            'data6' => array_column($data['detail_data'], $key6) ?? [] //纯属合伙人订单数
        ];

        $result = [
            'data1' => $data[$key1],
            'data2' => $data[$key2],
            'data3' => $data[$key3],
            'data4' => $data[$key4],
            'data5' => $data[$key5],
            'data6' => $data[$key6],
            'detail_data' => [
                'x' => $x,
                'y' => $y,
            ],
        ];

        return $result;
    }

    /**
     * 数据格式转化
     *
     * @param $data
     * @param string $key1
     * @param string $key2
     * @param string $key3
     * @param string $key4
     *
     * @return array
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/21 19:22
     */
    public function transformCustomerDataFormat($data, $key1 = 'customer_total', $key2 = 'customer_new', $key3 = 'customer_access', $key4 = 'customer_deal', $key5 = 'guide_customer_total', $key6 = 'other_customer_total', $key7 = 'guide_customer_new', $key8 = 'other_customer_new')
    {
        $x = array_keys($data['detail_data']);
        $y = [
            'data1' => array_column($data['detail_data'], $key1) ?? [],
            'data2' => array_column($data['detail_data'], $key2) ?? [],
            'data3' => array_column($data['detail_data'], $key3) ?? [],
            'data4' => array_column($data['detail_data'], $key4) ?? [],
            'data5' => array_column($data['detail_data'], $key5) ?? [],
            'data6' => array_column($data['detail_data'], $key6) ?? [],
            'data7' => array_column($data['detail_data'], $key7) ?? [],
            'data8' => array_column($data['detail_data'], $key8) ?? [],
        ];

        $result = [
            'data1' => $data[$key1],
            'data2' => $data[$key2],
            'data3' => $data[$key3],
            'data4' => $data[$key4],
            'data5' => $data[$key5],
            'data6' => $data[$key6],
            'data7' => $data[$key7],
            'data8' => $data[$key8],
            'detail_data' => [
                'x' => $x,
                'y' => $y,
            ],
        ];

        return $result;
    }

    /**
     * 获取员工角色
     *
     * @param $mobile
     * @return JsonResponse|int
     * @is_boss 1、店长 2、店员
     *
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/24 17:45
     */
    public function getRole($mobile)
    {
        try {
            $role = 0;
            $guideInfo = GoMicroScrm::getMemberInfo([
                'enterprise_id' => $this->enterpriseId,
                'mobile' => $mobile,
            ]);
            if (empty($guideInfo)) {
                return $role;
            }
            $this->memberInfo = $guideInfo;
            // 企业导购
            if ($guideInfo['type'] == 3) {
                $role = 3;
            }
            // 企业管理
            if ($guideInfo['type'] == 4) {
                $role = 5;
            }
            return $role;
        } catch (Exception $e) {
            Log::error("获取登录员工角色错误，错误信息为：{$e}");
            error(-1, '获取登录员工角色错误');
        }
        return 0;
    }

    /**
     * 保存浏览时间
     *
     * @param     $enterpriseId
     * @param     $storeId
     * @param     $storeWorkerId
     * @param int $type
     *
     * @return bool
     */
    public function saveViewTab($enterpriseId, $storeId, $storeWorkerId, $type = StoreMemberViewTab::EXTENSION_TAB_TYPE)
    {
        if (empty($storeWorkerId)) {
            return true;
        }

        $data = [
            'extension_view_time' => time()
        ];
        // 拼接查询条件
        $where = [
            'enterprise_id' => $enterpriseId,
            'store_id' => $storeId,
            'store_worker_id' => $storeWorkerId,
            'type' => $type
        ];

        // 查询是否有记录
        $total = StoreMemberViewTab::query()->where($where)->count();
        if (!$total) {
            $createData = array_merge($where, $data);
            StoreMemberViewTab::create($createData);
        } else {
            StoreMemberViewTab::query()->where($where)->update($data);
        }

        return true;
    }

    /**
     * 获得当前登录这的信息
     * @date    2020/11/20 17:09
     */
    public function getCurrentUserInfo($params)
    {
        $userInfo = [];
        if ($this->userType == self::USERTYPE_WORKER) {//导购
            $userInfo = $this->storeWorkerDetail($this->enterpriseId, $params['store_id']);
        }

        if ($this->userType == self::USERTYPE_PARTNER) {//合伙人
            $userInfo = $this->partnerDetail($this->enterpriseId);
        }
        return $userInfo;
    }

    /**
     * 获取登录员工详情
     * @param $enterprise
     * @param $storeId
     *
     * @return array|mixed|void
     */
    public function storeWorkerDetail($enterprise, $storeId)
    {
        $loginType = $this->loginUserType();
        $info = [];
        try {
            if ($loginType == self::USERTYPE_PARTNER) {
                $info = $this->partnerDetail($enterprise);
                if (!empty($info)) {
                    $info['id'] = $info['partner_id'];
                } else {
                    error(655001);
                }
            } else {
                $params = [
                    'enterprise_id' => $enterprise,
                    'store_id' => $storeId,
                    'mobile' => $this->mobile,
                    'worker_status' => 1
                ];

                $info = (new StoreWorker())->detail($params);
            }
        } catch (Exception $e) {
            Log::error("获取登录员工信息错误，错误信息为：{$e}");
            error(655001);
        }
        return $info;
    }

    /**
     * 获取登录的角色
     * @return string
     * author: zhaojie<zhaojie@vchangyi.com>
     * create: 2020/5/12 14:11
     */
    public function loginUserType(): string
    {
        try {
            $token = $this->request->header('token');
            if (empty($token)) {
                error(652006);
            }
            $tokenData = json_decode(Crypt::decrypt($token), true);
            return $tokenData['user_type'] ?? self::USERTYPE_WORKER;
        } catch (Exception $e) {
            Log::error("获取登录员工手机号错误，错误信息为：{$e}");
            error(655001);
        }
        return self::USERTYPE_WORKER;
    }

    /**
     * 获取登录员工详情
     * @param $enterprise
     * @param $storeId
     *
     * @return array|mixed|void
     */
    public function partnerDetail($enterprise)
    {
        try {
            $params = [
                'enterprise_id' => $enterprise,
                'mobile' => $this->mobile
            ];

            $info = (new PartnerMember())->detailAll($params);
        } catch (Exception $e) {
            Log::error("获取登录合伙人信息错误，错误信息为：{$e}");
            error(655001);
        }

        return $info;
    }

    /**
     * @throws Exception
     */
    public function __destruct()
    {

    }

    /**
     * 获取企业管理者权限范围数据
     *
     * @return array
     * @Time: 2022/11/18
     */
    public function getRoleRange(): array
    {
        $params = [
            'member_type' => self::MEMBER_TYPE_SA,
            'third_ids' => [],
        ];
        if ($this->isManager == 1) {
            $params['member_type'] = self::MEMBER_TYPE_SM;
            $dataRange = $this->dataRangeList;
            $third_id = array_column($dataRange, "object_id");
            $params['third_ids'] = $third_id;
        }
        return $params;
    }
}