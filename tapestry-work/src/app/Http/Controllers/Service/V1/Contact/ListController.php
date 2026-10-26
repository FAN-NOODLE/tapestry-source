<?php

namespace App\Http\Controllers\Service\V1\Contact;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Department;
use App\Library\Uc\Member;
use App\Models\Channel;
use App\Models\ContactQrCode;
use App\Models\ContactQrCodeParams;
use App\Models\MemberUserRelation;
use App\Models\WelcomeSpeech;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class ListController extends BaseController
{
    const YES = 1;
    const NO = 2;

    const MAX_LIMIT = 100000; // UC数据查询返回最大条数

    /**
     * 联系我列表
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $size          = $request->post('size', 10);           // 每页个数
        $selectType    = $request->post('select_type', '1');    // 查询类型
        $channel       = $request->post('channel', '');        // 渠道或二维码名称(搜索)
        $enterpriseId  = $request->post('enterprise_id', 0);   // 企业id
        $qyDomain      = $request->post('qy_domain', '');      // 企业domain
        $appIdentifier = $request->post('app_identifier', ''); // 应用标识
        $hasWelcomeMsg = $request->post('has_welcome_msg', 1); // 是否包含已经添加欢迎语的联系我二维码：1、是，2、否
        // 查询联系我二维码列表
        $contactQrCodeList = ContactQrCode::query()
            ->where('enterprise_id', $enterpriseId)
            ->when($channel != '', function ($query) use ($selectType, $channel) {
                $column = 'channel';
                if ($selectType == 2) {
                    $column = 'name';
                }
                $query->where($column, 'like', '%' . $channel . '%');
            })
            ->orderByDesc('id');

        if ($size) {
            // 分页
            $contactQrCodeList = $contactQrCodeList
                ->paginate($size, [
                    'id',
                    'name',
                    'channel',
                    'enterprise_id',
                    'type',
                    'source_type',
                    'is_verify',
                    'official_url',
                    'transform_url',
                    'created_at',
                ])
                ->toArray();

            if ($contactQrCodeList['data']) {
                // 联系我二维码ID集合
                $contactQrCodeIds = array_column($contactQrCodeList['data'], 'id');
                // 联系我二维码ID集合对应的关联员工/部门列表(取第一条)[可能存在没有一个员工或部门]
                $contactQrCodeParamsList = ContactQrCodeParams::query()
                    ->whereIn('contact_qr_code_id', $contactQrCodeIds)
                    ->groupBy('contact_qr_code_id')
                    ->orderBy('id')
                    ->get(['id', 'enterprise_id', 'contact_qr_code_id', 'object', 'type'])
                    ->toArray();

                // 联系我二维码ID集合对应的关联员工
                $contactQrCodeParamsMember = ContactQrCodeParams::query()
                    ->select(DB::raw('object_extra, contact_qr_code_id'))
                    ->where('type', ContactQrCodeParams::MEMBER)
                    ->whereIn('contact_qr_code_id', $contactQrCodeIds)
                    ->get()
                    ->toArray();

                // 联系我二维码ID集合对应的关联部门
                $contactQrCodeParamsDepartment = ContactQrCodeParams::query()
                    ->select(DB::raw('object_extra, contact_qr_code_id'))
                    ->where('type', ContactQrCodeParams::DEPARTMENT)
                    ->whereIn('contact_qr_code_id', $contactQrCodeIds)
                    ->get()
                    ->toArray();
                // 联系我二维码累计添加人数
                $qrCodeAddedCount = MemberUserRelation::query()
                    ->selectRaw(
                        'source_type,count(id) as code_added_count'
                    )
                    ->whereIn('source_type', $contactQrCodeIds)
                    ->where('enterprise_id', $enterpriseId)
                    ->groupBy('source_type')
                    ->get()->toArray();
                // 联系我二维码累计减少人数
                $qrCodeDecreasedCount = MemberUserRelation::query()
                    ->selectRaw(
                        'source_type,count(id) as code_added_count'
                    )
                    ->whereIn('source_type', $contactQrCodeIds)
                    ->where('enterprise_id', $enterpriseId)
                    ->where('cancel_relation_at', '<>', 0)
                    ->where('delete_type', '<>', 0)
                    ->groupBy('source_type')
                    ->get()->toArray();
                $qrCodeAddedCountArr  = array_column($qrCodeAddedCount, 'code_added_count', 'source_type');
                $qrCodeDecreasedCount = array_column($qrCodeDecreasedCount, 'code_added_count', 'source_type');
                // 初始化员工/部门ID集合
                $memberParams     = [];
                $departmentParams = [];

                foreach ($contactQrCodeParamsList as $contactQrCodeParam) {
                    // 联系我二维码关联员工ID(memUid)
                    if (ContactQrCodeParams::MEMBER == $contactQrCodeParam['type']) {
                        $memberParams[] = $contactQrCodeParam['object'];
                    }

                    // 联系我二维码关联部门ID(dpId)
                    if (ContactQrCodeParams::DEPARTMENT == $contactQrCodeParam['type']) {
                        $departmentParams[] = $contactQrCodeParam['object'];
                    }
                }

                // 初始化员工/部门信息集合
                $members     = [];
                $departments = [];

                // 查询关联员工(memUid)的信息
                if ($memberParams) {
                    $memberResult = (new Member())->ignorePermissionsList($qyDomain, $appIdentifier, ['memUids' => $memberParams]);
                    $members      = array_column($memberResult['list'], null, 'memUid');
                }

                // 查询关联部门(dpId)的信息
                if ($departmentParams) {
                    $departmentResult = (new Department())->list($qyDomain, $appIdentifier, ['pageSize' => self::MAX_LIMIT], false);
                    $departments      = array_column($departmentResult['list'], null, 'dpId');
                }

                $contactQrCodeParamsList = array_column($contactQrCodeParamsList, null, 'contact_qr_code_id');

                // 联系我二维码ID集合对应的关联员工集合(用,分隔)
                $memberIdStringList = [];
                foreach ($contactQrCodeParamsMember as $value) {
                    $memberIdStringList[$value['contact_qr_code_id']] =
                        ($memberIdStringList[$value['contact_qr_code_id']] ?? '0') . ',' . $value['object_extra'];
                }

                // 联系我二维码ID集合对应的关联部门员工集合(用,分隔)
                $departmentMemberIdStringList = [];
                foreach ($contactQrCodeParamsDepartment as $value) {
                    $departmentMemberIdStringList[$value['contact_qr_code_id']] =
                        ($departmentMemberIdStringList[$value['contact_qr_code_id']] ?? '0') . ',' . $value['object_extra'];
                }

                // 组装数据
                foreach ($contactQrCodeList['data'] as &$value) {
                    // 如果联系我二维码ID集合对应的关联员工/部门列表(取第一条)存在
                    if ($contactQrCodeParamsList[$value['id']] ?? 0) {
                        // 如果联系我二维码ID集合对应的关联员工/部门列表(取第一条)是员工，则取员工信息
                        if (ContactQrCodeParams::MEMBER == $contactQrCodeParamsList[$value['id']]['type']) {
                            $value['object_name'] = ($members[$contactQrCodeParamsList[$value['id']]['object']] ?? [])['memUsername'] ?? '';
                            $value['object_type'] = ContactQrCodeParams::MEMBER;
                        }

                        // 如果联系我二维码ID集合对应的关联员工/部门列表(取第一条)是部门，则取部门信息
                        if (ContactQrCodeParams::DEPARTMENT == $contactQrCodeParamsList[$value['id']]['type']) {
                            $value['object_name'] = ($departments[$contactQrCodeParamsList[$value['id']]['object']] ?? [])['dpName'] ?? '';
                            $value['object_type'] = ContactQrCodeParams::DEPARTMENT;
                        }

                        $totalMemberIdStringList = '0';

                        if ($memberIdStringList[$value['id']] ?? 0) {
                            $totalMemberIdStringList = $totalMemberIdStringList . ',' . $memberIdStringList[$value['id']];
                        }

                        // 联系我二维码ID集合对应的关联部门员工数量
                        if ($departmentMemberIdStringList[$value['id']] ?? 0) {
                            $totalMemberIdStringList = $totalMemberIdStringList . ',' . $departmentMemberIdStringList[$value['id']];
                        }

                        $contactQrCodeParamsMemberTotalCount = count(array_filter(array_unique(explode(',', $totalMemberIdStringList))));

                        // 联系我二维码ID集合对应的关联总计员工数量
                        $value['count'] = $contactQrCodeParamsMemberTotalCount;
                    } else {
                        // 如果联系我二维码ID集合对应的关联员工/部门列表(取第一条)存在，即关联的最后一个员工离职或关联的最后一个部门删除
                        $value['object_name'] = '';
                        $value['object_type'] = ContactQrCodeParams::MEMBER;
                        $value['count']       = 0;
                    }

                    // 转换存储url(大)1417x1417
                    $value['transform_url_large'] = $value['transform_url'] . '?imageMogr2/thumbnail/1417x1417!';
                    // 转换存储url(中)430x430
                    $value['transform_url_middle'] = $value['transform_url'] . '?imageMogr2/thumbnail/430x430!';
                    // 转换存储url(小)258x258
                    $value['transform_url_small'] = $value['transform_url'] . '?imageMogr2/thumbnail/258x258!';
                    //联系人二维码累计添加人数
                    $qrCodeAddMemCount         = ($qrCodeAddedCountArr[$value['id']] ?? 0)
                        - ($qrCodeDecreasedCount[$value['id']] ?? 0);
                    $value['code_added_count'] = $qrCodeAddMemCount;
                }
                unset($value);
            }

            return success([
                'data' => $contactQrCodeList['data'],
                'meta' => [
                    'total' => $contactQrCodeList['total'],
                    'size'  => (int)$contactQrCodeList['per_page'],
                    'page'  => $contactQrCodeList['current_page']
                ]
            ]);
        } else {
            $welcomeSpeechContactQrCodeIds = [];

            if (self::NO == $hasWelcomeMsg) {
                $welcomeSpeechList = WelcomeSpeech::query()
                    ->where('enterprise_id', $enterpriseId)
                    ->where('is_default', WelcomeSpeech::IS_DEFAULT_NO)
                    ->get(['id', 'channel_id', 'contact_qr_code_id', 'enterprise_id', 'type'])
                    ->toArray();

                if ($welcomeSpeechList) {
                    $welcomeSpeechContactQrCodeIds = array_column($welcomeSpeechList, 'channel_id');
                }
            }
            //所有渠道
            $contactQrCodeList = Channel::query()
                ->where('enterprise_id', $enterpriseId)
                ->when($welcomeSpeechContactQrCodeIds, function ($query) use ($welcomeSpeechContactQrCodeIds) {
                    $query->whereNotIn('id', $welcomeSpeechContactQrCodeIds);
                })
                //->where('source_type', Channel::SOURCE_TYPE_DEFAULT)
                ->get(['id', 'channel'])->toArray();
            // 查出所有
//            $contactQrCodeList = $contactQrCodeList
//                ->when($welcomeSpeechContactQrCodeIds, function ($query) use ($welcomeSpeechContactQrCodeIds) {
//                    $query->whereNotIn('id', $welcomeSpeechContactQrCodeIds);
//                })
//                ->get(['id', 'name', 'channel', 'enterprise_id', 'type', 'is_verify', 'official_url', 'transform_url', 'created_at'])
//                ->toArray();

            return success([
                'data' => $contactQrCodeList,
            ]);
        }
    }
}
