<?php
/**
 * Description:用户自定义标签控制器
 * User: lintong
 * Email: 18511245965@163.com
 * Date: 2020/8/3
 * Time: 16:32
 */

namespace App\Http\Controllers\App\V1\Customer;


use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Work\CustomerTags;
use CrsSdk\Work\StoreWorker;

class UserTagController extends BaseController
{
    /**
     * 解除自定义标签绑定
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function deleteRelation(Request $request)
    {
        $param['id'] = $request->post('id');
        $result      = (new CustomerTags())->deleteRelation($param);
        return success($result);
    }

    /**
     * 删除自定义标签
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function deleteTag(Request $request)
    {
        $param['id'] = $request->post('id');
        $result      = (new CustomerTags())->deleteTag($param);
        return success($result);
    }

    /**
     * 添加用户自定义标签
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function addTag(Request $request)
    {
        $post    = $request->only([
            'name',
            'store_id'
        ]);
        $storeId = $post['store_id'] ?? 0;
        $wokerId                 = $this->getWorkerId($storeId);
        $params['worker_id']     = $wokerId;
        $params['name']          = $post['name'] ?? '';
        $params['enterprise_id'] = $this->enterpriseId;
        $result                  = (new CustomerTags())->addTag($params);
        if(isset($result['code'])){
            throw new Error($result['code'],$result['message'],[],200);
        }
        return success($result);
    }

    /**
     * 绑定标签
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function addRelation(Request $request)
    {
        $post                    = $request->only([
            'relation_id',
            'store_id',
            'tag_names'
        ]);
        $storeId                 = $post['store_id'] ?? 0;
        $wokerId                 = $this->getWorkerId($storeId);
        $params['worker_id']     = $wokerId;
        $params['relation_id']   = $post['relation_id'] ?? '';
        $params['tag_names']     = $post['tag_names'] ?? [];
        $params['enterprise_id'] = $this->enterpriseId;
        $result                  = (new CustomerTags())->addRelation($params);
        if(isset($result['code'])){
            throw new Error($result['code'],$result['message'],[],200);
        }
        return success($result);
    }

    /**
     * 获取绑定的自定义标签列表
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function getRelation(Request $request)
    {
        $post                    = $request->only([
            'relation_id',
            'store_id',
        ]);
        $storeId                 = $post['store_id'] ?? 0;
        $wokerId                 = $this->getWorkerId($storeId);
        $params['relation_id']   = $post['relation_id'] ?? '';
        $params['worker_id']     = $wokerId;
        $params['enterprise_id'] = $this->enterpriseId;
        $result                  = (new CustomerTags())->getRelation($params);
        return success($result);

    }

    /**
     * 获取未绑定的自定义订单列表
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function getTags(Request $request)
    {
        $post                    = $request->only([
            'relation_id',
            'store_id',
            'page',
            'size'
        ]);
        $storeId                 = $post['store_id'] ?? 0;
        $wokerId                 = $this->getWorkerId($storeId);
        $params['relation_id']   = $post['relation_id'] ?? '';
        $params['worker_id']     = $wokerId;
        $params['enterprise_id'] = $this->enterpriseId;
        $params['page']          = $post['page'] ?? 1;
        $params['size']          = $post['size'] ?? 10;
        $result                  = (new CustomerTags())->getTags($params);
        return success($result);
    }

    /**
     * 获取店员所有的自定义标签
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function getList(Request $request)
    {
        $post                    = $request->only([
            'store_id',
            'page',
            'size'
        ]);
        $storeId                 = $post['store_id'] ?? 0;
        $wokerId                 = $this->getWorkerId($storeId);
        $params['worker_id']     = $wokerId;
        $params['enterprise_id'] = $this->enterpriseId;
        $params['page']          = $post['page'] ?? 1;
        $params['size']          = $post['size'] ?? 10;
        $result                  = (new CustomerTags())->getList($params);
        return success($result);
    }

    /**
     * 搜索自定义标签
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function searchTags(Request $request)
    {
        $post                    = $request->only([
            'store_id',
            'keyword'
        ]);
        $storeId                 = $post['store_id'] ?? 0;
        $wokerId                 = $this->getWorkerId($storeId);
        $params['worker_id']     = $wokerId;
        $params['enterprise_id'] = $this->enterpriseId;
        $params['keyword']       = $post['keyword'];
        $result                  = (new CustomerTags())->searchTags($params);
        return success($result);
    }

    /**
     * 根据登陆信息获取店员worker_id
     * @param $storeId
     * @return mixed
     * @throws \Exception
     */
    public function getWorkerId($storeId)
    {
        // 获取店员信息
        $memberInfo = (new StoreWorker())->detail([
            'enterprise_id' => $this->enterpriseId,
            'store_id'      => $storeId,
            'mobile'        => $this->mobile,
        ]);
        if (empty($memberInfo)) {
            throw new Error(820000, '店员信息有误');
        }
        return $memberInfo['id'];
    }
}