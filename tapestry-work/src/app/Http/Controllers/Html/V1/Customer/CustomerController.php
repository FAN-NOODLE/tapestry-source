<?php
/*
* +------------------------------------+
* | Created: phpstorm
* +------------------------------------+
* | Author : wyg
* +------------------------------------+
* | Date   : 2019/12/31
* +------------------------------------+
* | Time   : 17:02
* +------------------------------------+
* | Email  : <wangyongguo@vchangyi.com>
* +------------------------------------+
*/


namespace App\Http\Controllers\Html\V1\Customer;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Html\V1\BaseController;
use App\Library\Uc\External;
use App\Library\Uc\Tag;
use App\Models\Address;
use App\Models\Channel;
use App\Models\ContactQrCode;
use App\Models\CustomerGroup;
use App\Models\CustomerGroupRelation;
use App\Models\MemberUserRelation;
use App\Services\CustomerService;
use CrsSdk\MicroUser\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class CustomerController extends BaseController
{
    const TYPE_ORDERS = 1; // 加载订单

    const TYPE_COUPONS = 2; // 加载优惠券信

    /**
     * 用户详情
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function detail(Request $request)
    {
        $userId   = $request->input('user_id', 0);
        $orderNo  = $request->input('order_no', '');
        $dataType = $request->input('type', 0);
        $limit    = $request->input('size', 20);
        $page     = $request->input('page', Constant::PAGE);

        if (empty($userId)) {
            throw new Error(10000, '用户[user_id]不能为空');
        }

        if (empty($dataType)) {
            throw new Error(11000, '用户[type]不能为空');
        }
        $userDetail = MemberUserRelation::query()
            ->where([
                'mem_uid'       => $this->memUid,
                'ec_userid'     => $userId,
                'enterprise_id' => $this->enterpriseId
            ])
            ->first();

        if (!$userDetail) {
            $externalRelations = (new External())->listByUserId($this->domain, [
                'ecUserid' => $userId,
                'memUid'   => $this->memUid,
            ]);

            if (empty($externalRelations['list'])) {
                throw new Error(12000, '未查询到客户信息');
            }

            $this->addExternalContact($this->enterpriseId, $this->memUid, $externalRelations['list'][0]);
            // 客户详情
            $userDetail = MemberUserRelation::query()
                ->where([
                    'mem_uid'     => $this->memUid,
                    'delete_type' => MemberUserRelation::NORMAL_STATUS,
                    'ec_userid'   => $userId,
                ])
                ->first();
        }
        switchAppHostByDomain($this->domain);
        // 新版的走新逻辑
        if ($this->enterprise['version'] != Constant::NORMAL_VERSION) {
            $list = $this->newVersionCustomerDetail($userDetail, $dataType, $page, $limit, $orderNo);
        } else {

            // 用户查询
            $user = CustomerService::user(['enterprise_id' => $this->enterpriseId, 'union_id' => $userDetail->union_id]);
            if (empty($user)) {
                $user = [];
            } else {
                $user = $user->toArray();
            }

            if (empty($user)) {
                //throw new Error(12000, '未查询到客户信息');
                $user = $this->getUserInfo($userDetail);
            }

            $list['user'] = $user;

            // 查询地址参数
            $where = [
                'user_id'       => $user['id'],
                'default'       => Address::DEF_ADDR,
                'enterprise_id' => $this->enterpriseId
            ];

            $addressInfo = CustomerService::address($where);

            if ($addressInfo) {
                $list['user']['address'] = $addressInfo->province . $addressInfo->city . $addressInfo->area . $addressInfo->address;
            }

            // 订单|优惠券查询参数
            $where = [
                ['enterprise_id', '=', $this->enterpriseId],
                ['union_id', '=', $userDetail->union_id],
            ];

            // 搜索条件追加
            ($orderNo || '0' == $orderNo) && $where[] = ['order_no', 'like', "%$orderNo%"];
            self::TYPE_ORDERS == $dataType && $list['info'] = CustomerService::orders($where, $limit);

            // 优惠券查询参数
            self::TYPE_COUPONS == $dataType && $list['info'] = CustomerService::coupons($where, $limit);
        }
        // 获取用户标签,分组,渠道
        // 获取渠道
        $list['user']['channel_name'] = '';
        $contactQrcode                = ContactQrCode::query()->where('id', $userDetail->source_type)->first();
        if ($contactQrcode) {
            $channelByContact             = Channel::query()->where('id', $contactQrcode->channel_id)->first();
            $list['user']['channel_name'] = $channelByContact->channel;
        }

        // 获取分组
        $groupIds                   = CustomerGroupRelation::query()
            ->where('enterprise_id', $this->enterpriseId)
            ->where('ec_userid', '=', $userDetail->ec_userid)
            ->pluck('customer_group_id')
            ->toArray();
        $list['user']['group_name'] = CustomerGroup::query()
            ->whereIn('id', $groupIds)
            ->pluck('name')
            ->toArray();


        $tagNames = [];
        $tagList  = (new Tag())->tagDetail($this->domain, ['ecUserid' => $userDetail->ec_userid]);
        unset($tagList['requestId']);

        // 获取这个员工给这个客户打的标签
        if (!empty($tagList)) {
            $tagCollect = [];
            foreach ($tagList as $key => $item) {
                if ($item['memUid'] == $this->memUid) {
                    $tagCollect[$item['tagId']] = $item['tagName'];
                    $tagNames[]                 = [
                        'tag_id'   => $item['tagId'],
                        'tag_name' => $item['tagName']
                    ];
                }
            }
        }

        $list['user']['tag_name'] = $tagNames;

        return success([
            'data' => $list
        ]);
    }

    /**
     * @param        $userDetail
     * @param int    $dataType
     * @param int    $page
     * @param int    $limit
     * @param string $orderNo
     *
     * @return array
     * @throws \Exception
     */
    private function newVersionCustomerDetail($userDetail, int $dataType, int $page, int $limit, string $orderNo = ''): array
    {
        $customer = [];
        $user = [];
        // 先查询用户信息
        if (!empty($userDetail->union_id)) {
            try {
                $user = (new User())->userDetail([
                    'enterprise_id' => $this->enterpriseId,
                    'union_id' => $userDetail->union_id,
                    'default_address' => 1
                ]);
            } catch (Error $error) {

            }
        }
        // 使用默认用户信息
        if (empty($user)) {
            $user = $this->getUserInfo($userDetail);
        }

        // 客户地址
        $user['address'] = '';
        if (!empty($user['default_address'])) {
            $user['address'] = ($user['default_address']['province'] ?? '') .
                ($user['default_address']['city'] ?? '') .
                ($user['default_address']['area'] ?? '') .
                ($user['default_address']['address'] ?? '');
        }
        if (self::TYPE_ORDERS == $dataType) {
            $where = [
                ['enterprise_id', '=', $this->enterpriseId],
                ['union_id', '=', $userDetail->union_id],
            ];
            if (!empty($orderNo)) {
                $where[] = ['order_no', 'like', "%$orderNo%"];
            }
            $customer['info'] = CustomerService::orders($where, $limit);
        } else if (self::TYPE_COUPONS == $dataType) {
            if ($user['id'] > 0) {
                $couponData = CustomerService::customerCoupons($this->enterpriseId, $user['id']);
            }
            $customer['info'] = $couponData ?? [];
        }
        $customer['user'] = $user;
        return $customer;
    }

    /**
     * 添加成员
     *
     * @param int    $enterpriseId
     * @param string $memUid
     * @param        $params
     *
     * @return bool|string
     * @throws \Exception
     */
    protected function addExternalContact(int $enterpriseId, string $memUid, $params)
    {
        $bindMemberData = [];
        foreach ($params['followMembers'] as $member) {
            if ($member['memUid'] == $memUid) {
                $bindMemberData = $member;
                break;
            }
        }
        $insertData = [
            'enterprise_id'    => $enterpriseId,
            'mem_uid'          => $memUid,
            'ec_id'            => $params['ecId'],
            'ec_userid'        => $params['ecUserid'],
            'ec_name'          => $params['ecName'],
            'ec_avatar'        => $params['ecAvatar'],
            'ec_gender'        => $params['ecGender'],
            'user_name_index'  => $params['ecIndex'],
            'user_name_remark' => $bindMemberData['emRemark'],
            'union_id'         => $params['ecUnionid'] ?? '',
            'relation_at'      => $bindMemberData['emCreatetime'],
            'source_type'      => $bindMemberData['emState'] ? intval($bindMemberData['emState']) : 1,
        ];
        MemberUserRelation::query()->create($insertData);
        return true;
    }

    /**
     * 获取默认用户信息
     *
     * @param $userRelationDetail
     * @return array
     */
    protected function getUserInfo($userRelationDetail)
    {
        return [
            'id'              => 0,
            'nick_name'       => $userRelationDetail->ec_name,
            'mobile'          => '',
            'avatar'          => $userRelationDetail->ec_avatar,
            'union_id'        => $userRelationDetail->union_id,
            'enterprise_id'   => 0,
            'open_id'         => '',
            'third_member_no' => '',
            'real_name'       => '',
            'real_name_index' => '',
            'birthday'        => '',
            'country'         => '',
            'province'        => '',
            'city'            => '',
            'register_at'     => 0,
            'sex'             => 3, // 性别（1:男，2:女，3:未知）
            'level_id'        => 0,
            'level_update_at' => 0,
            'user_amount'     => 0,
            'created_at'      => 0,
            'updated_at'      => 0,
            'auth_at'         => 0,
            'phone'           => '',
            'score_total'     => 0,
            'mem_uid'         => "",
            'is_crm_user'     => 0,
            'default_address' => [],
            'bind_guide'      => []
        ];
    }
}