<?php

namespace App\Http\Controllers\Service\V1\Customer;

use App\Library\Uc\Tag;
use App\Models\Channel;
use App\Library\Uc\Member;
use App\Models\ContactQrCode;
use Illuminate\Http\Request;
use App\Models\CustomerGroup;
use App\Models\MemberUserRelation;
use App\Models\CustomerGroupRelation;
use App\Http\Controllers\Service\V1\BaseController;

class BaseDetailController extends BaseController
{
    /**
     * 性别映射
     */
    const GENDER = [
        MemberUserRelation::EC_GENDER_MAN        => '男',
        MemberUserRelation::EC_GENDER_WOMAN      => '女',
        MemberUserRelation::EC_GENDER_DEFAULT    => '未知',
        MemberUserRelation::EC_AS_GENDER_DEFAULT => '未知',
    ];

    public function userDetail(Request $request)
    {
        $id            = $request->input('id');
        $qyDomain      = $request->input('qy_domain');
        $enterpriseId  = $request->input('enterprise_id');
        $appIdentifier = $request->input('app_identifier');

        if (empty($id)) {
            notice(600002);
        }

        if (empty($enterpriseId)) {
            notice(600000);
        }

        if (empty($qyDomain) || empty($appIdentifier)) {
            notice(110006);
        }

        // 找到客户
        $member = MemberUserRelation::query()->where('id', '=', $id)->first();

        if (empty($member)) {
            notice(900002);
        }

        $data             = [];
        $data['name']     = $member->ec_name;
        $data['union_id'] = $member->union_id;
        $data['avatar']   = $member->ec_avatar;
        $data['add_at']   = date('Y-m-d H:i:s', $member->relation_at);
        $data['gender']   = self::GENDER[$member->ec_gender];

        // 获取用户所属员工以及员工的部门
        $memUser          = (new Member())->detail($qyDomain, $appIdentifier, ['memUid' => $member->mem_uid]);
        $data['mem_name'] = $memUser['memUsername'] ?? '';
        $data['dep_name'] = $memUser['dpName'][0]['dpName'] ?? '未知部门';

        // 获取渠道
        $data['channel_name'] = '';
        $contactQrcode = ContactQrCode::query()->where('id', $member->source_type)->first();
        if ($contactQrcode) {
            $channelByContact = Channel::query()->where('id', $contactQrcode->channel_id)->first();
            $data['channel_name'] = $channelByContact->channel;
        }

        // 获取分组
        $groupIds           = CustomerGroupRelation::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('ec_userid', '=', $member->ec_userid)
            ->pluck('customer_group_id')
            ->toArray();
        $data['group_name'] = CustomerGroup::query()
            ->whereIn('id', $groupIds)
            ->pluck('name')
            ->toArray();


        $tagNames = [];
        $tagList  = (new Tag)->tagDetail($qyDomain, ['ecUserid' => $member['ec_userid']]);
        unset($tagList['requestId']);

        // 获取这个员工给这个客户打的标签
        if (!empty($tagList)) {
            $tagCollect = [];
            foreach ($tagList as $key => $item) {
                if ($item['memUid'] == $member['mem_uid']) {
                    $tagCollect[$item['tagId']] = $item['tagName'];
                }
            }
            $tagNames = array_values($tagCollect);
        }

        // 客户员工关系正常, 显示该员工给客户打的标签
        if ($member['delete_type'] > 0) {
            // 查询是否与其他客户关系正常
            $existsOtherRelation = MemberUserRelation::query()
                ->where('ec_id', '=', $member['ec_id'])
                ->where('delete_type', '=', 0)
                ->where('mem_uid', '!=', $member['mem_uid'])
                ->exists();

            // 存在就啥都不做, 返回空
            if ($existsOtherRelation) {
                $tagNames = [];
            }
        }

        $data['tag_name'] = $tagNames;

        return success($data);
    }
}