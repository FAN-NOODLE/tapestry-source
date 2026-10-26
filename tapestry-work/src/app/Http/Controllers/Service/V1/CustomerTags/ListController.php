<?php
/**
 * Description:获取绑定/未绑定的自定义标签列表
 * User: lintong
 * Email: 18511245965@163.com
 * Date: 2020/8/3
 * Time: 11:06
 */

namespace App\Http\Controllers\Service\V1\CustomerTags;

use Exception;
use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\UserCustomizeTagsService;
use CrsSdk\GoMicroScrm\Tag;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class ListController extends BaseController
{

    /**
     * 获取绑定的自定义标签
     * @param Request $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getRelation(Request $request)
    {
        $post = $request->only([
            'enterprise_id',
            'relation_id',
            'worker_id',
        ]);
        $this->checkParams($post);
        if (!isset($post['relation_id']) || empty($post['relation_id'])) {
            notice(820000, ['relation_id有误']);
        }
        $result = (new UserCustomizeTagsService())->getRelation($post);
        return success(['code' => 200, 'data' => $result]);
    }

    /**
     * 获取未绑定的标签列表
     * @param Request $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getTags(Request $request)
    {
        $post = $request->only([
            'enterprise_id',
            'relation_id',
            'worker_id',
            'page',
            'size'
        ]);
        $this->checkParams($post);
        if (!isset($post['relation_id']) || empty($post['relation_id'])) {
            notice(820000, ['relation_id有误']);
        }
        $post['page'] = $post['page'] ?? 1;
        $post['size'] = $post['size'] ?? 10;
        $result       = (new UserCustomizeTagsService())->getTags($post);
        return success(['code' => 200, 'data' => $result]);
    }

    /**
     * 获取员工所有自定义标签
     * @param Request $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getList(Request $request)
    {
        $post = $request->only([
            'enterprise_id',
            'worker_id',
            'page',
            'size'
        ]);
        $this->checkParams($post);
        $post['page']        = $post['page'] ?? 1;
        $post['size']        = $post['size'] ?? 10;
        $result              = (new UserCustomizeTagsService())->getTagsList($post);
        return success(['code' => 200, 'data' => $result]);
    }

    /**
     * 搜索自定义标签接口
     * @param Request $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function searchTags(Request $request)
    {
        $post = $request->only([
            'enterprise_id',
            'worker_id',
            'keyword'
        ]);
        $this->checkParams($post);
        if (!isset($post['keyword']) || empty($post['keyword'])) {
            notice(820000, ['关键字不能为空']);
        }
        $result = (new UserCustomizeTagsService())->searchTags($post);
        return success(['code' => 200, 'data' => $result]);
    }

    /**
     * 检验参数
     * @param $post
     * @throws \Illuminate\Validation\ValidationException
     */
    protected function checkParams($post)
    {
        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'worker_id'     => 'required|integer',
        ];

        $message   = [
            'enterprise_id.required' => '企业ID不存在',
            'enterprise_id.integer'  => '企业ID数据类型有误',
            'enterprise_id.min'      => '企业ID数据类型有误',
            'worker_id.required'     => '员工ID不能为空',
            'worker_id.integer'      => '员工ID数据类型有误',
        ];
        $validator = Validator::make($post, $rules, $message);
        if ($validator->fails()) {
            throw new \Illuminate\Validation\ValidationException($validator);
        }
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
//    public function getGroupList(Request $request): array
//    {
//        $params = $request->only([
//            'enterprise_id',
//        ]);
//        $rules = [
//            'enterprise_id' => 'required|integer|min:1',
//        ];
//
//        $message   = [
//            'enterprise_id.required' => '企业ID不存在',
//            'enterprise_id.integer'  => '企业ID数据类型有误',
//            'enterprise_id.min'      => '企业ID数据类型有误',
//        ];
//        $validator = Validator::make($params, $rules, $message);
//        if ($validator->fails()) {
//            throw new \Illuminate\Validation\ValidationException($validator);
//        }
//
//        try {
//            $result = (new Tag())->tagGroupList($params);
//
//        } catch (Exception $e) {
//            Log::info('获取标签分组列表失败' . $e);
//            throw new Error($e->getCode(), $e->getMessage());
//        }
//
//        return success(['code' => 200, 'data' => $result]);
//    }


}