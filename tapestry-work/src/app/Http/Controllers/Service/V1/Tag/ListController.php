<?php


namespace App\Http\Controllers\Service\V1\Tag;

use Exception;
use App\Exceptions\Error;
use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\TagService;
use CrsSdk\GoMicroScrm\Tag;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class ListController extends BaseController
{
    /**
     * 获取标签组列表
     *
     * @return array
     * @throws \Exception
     * @api /service/v1/tag/list
     */
    public function index()
    {
        $qyDomain = $this->request->input('qy_domain', '');

        if (empty($qyDomain)) {
            notice(200018);
        }

        // 获取列表
        $tagList = (new TagService())->listTag($qyDomain);

        return success(['data' => $tagList]);
    }

    /**
     * 获取标签组标签列表-根据标签组id
     *
     * @return array
     * @throws \Exception
     * @api /service/v1/tag/page-list
     */
    public function pageList()
    {
        $input = $this->request->only([
            'enterprise_id',
            'tag_group_id',
            'keyword',
            'custom_group_id',
            'qy_domain',
            'page',
            'size',
        ]);

        $qyDomain      = $input['qy_domain'] ?? '';
        $groupId       = $input['tag_group_id'] ?? 0;
        $input['page'] = $input['page'] ?? Constant::PAGE;
        $input['size'] = $input['size'] ?? Constant::SIZE;

        if (empty($qyDomain)) {
            notice(200018);
        }
        if (empty($groupId)) {
            notice(810002);
        }

        // 获取列表
        $tagList = (new TagService())->pageListTag($qyDomain, $input);

        return success([
            'data' => $tagList['list'] ?? [],
            'meta' => [
                'page'  => $tagList['page'] ?? $input['page'],
                'size'  => $tagList['size'] ?? $input['size'],
                'total' => $tagList['total'] ?? 0
            ]
        ]);
    }

    /**
     * 获取所有标签组标签列表
     *
     * @return array
     * @throws \Exception
     * @api /service/v1/tag/all-list
     */
    public function allList()
    {
        $input = $this->request->only([
            'qy_domain',
        ]);

        // 获取全部标签列表
        $list      = (new TagService())->getAllTagList($input['qy_domain']);
        $requestId = $list['requestId'];
        unset($list['requestId']);

        foreach ($list as &$value) {
            foreach ($value['tagList'] as &$v) {
                $v['tagId'] = (string)$v['tagId'];
            }
            $value['groupId'] = (string)$value['groupId'];
        }

        unset($v);
        unset($value);
        return success([
            'data'      => $list,
            'requestId' => $requestId
        ]);
    }

    /**
     * 获取企业标签库标签列表
     *
     * @return array
     * @throws \Exception
     * @api /service/v1/tag/select-list
     */
    public function selectList()
    {
        $input = $this->request->only([
            'qy_domain',
            'page',
            'size'
        ]);

        // 获取全部标签列表
        $list = (new TagService())->tagListForSelect($input);

        foreach ($list['list'] as $k => $v) {
            $list['list'][$k]['id'] = (string)$v['id'];
        }

        return success([
            'data' => $list['list'] ?? [],
            'meta' => [
                'page'  => $list['pageNum'] ?? $input['page'],
                'size'  => $list['pageSize'] ?? $input['size'],
                'total' => $list['total'] ?? 0
            ]
        ]);
    }


    /**
     * 获取标签分组列表
     *
     * @param Request $request
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @throws Exception
     */
    public function allGroupList(Request $request): array
    {
        $params = $request->only([
            'enterprise_id',
        ]);
        $rules  = [
            'enterprise_id' => 'required|integer|min:1',
        ];

        $message   = [
            'enterprise_id.required' => '企业ID不存在',
            'enterprise_id.integer'  => '企业ID数据类型有误',
            'enterprise_id.min'      => '企业ID数据类型有误',
        ];
        $validator = Validator::make($params, $rules, $message);
        if ($validator->fails()) {
            throw new \Illuminate\Validation\ValidationException($validator);
        }

        try {
            $result = (new Tag())->tagGroupList($params);

        } catch (Exception $e) {
            Log::info('获取标签分组列表失败' . $e);
            throw new Error($e->getCode(), $e->getMessage());
        }


//        $addWayTypeObject = [];
//        $temp             = [];
//        $addWayType       = Constant::ADD_WAY_TYPE;
//        foreach ($addWayType as $k => $v) {
//            $temp['add_way_key']  = $k;
//            $temp['add_way_name'] = $v;
//            $addWayTypeObject[]   = $temp;
//        }
        $addWayTypeObject = [
            ['add_way_key' => Constant::ADD_WAY_SCAN_QR_CODE, 'add_way_name' => '导购名片'],
            ['add_way_key' => Constant::ADD_WAY_OTHER, 'add_way_name' => '其他方式'],
        ];

        $data['list']    = $result;
        $data['add_way'] = $addWayTypeObject;

        return success(['code' => 200, 'data' => $data]);
    }

}