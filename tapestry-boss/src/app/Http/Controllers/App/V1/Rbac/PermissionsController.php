<?php

namespace App\Http\Controllers\App\V1\Rbac;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\MicroAuthority\Permissions;
use Illuminate\Support\Facades\Config;
use Illuminate\Support\Facades\Log;

class PermissionsController extends BaseController
{
    /**
     * 权限树接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function allTree()
    {
        $data = (new Permissions())->allTree();

        return success($data);
    }

    /**
     * 新增权限接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function add()
    {
        $input = $this->request->all();

        $data = (new Permissions())->add($input);

        return success($data);
    }

    /**
     * 编辑权限接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function edit()
    {
        $input = $this->request->all();

        $data = (new Permissions())->edit($input);

        return success($data);
    }

    /**
     * 发布超管权限接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function release()
    {
        $input = $this->request->all();

        $data = (new Permissions())->release($input);

        return success($data);
    }

    /**
     * 同步并发布权限接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function syncAndRelease()
    {
        $this->validate($this->request, [
            'env'           => 'required|string|in:dev,test,prod',
            'system_type'   => 'required|integer',
            'version_tag'   => 'required|string',
            'enterprise_id' => 'required|integer|min:0',
        ]);
        $input = $this->request->all();

        $key = 'crs:boss:permissions:sync-and-release';
        $redis = redis();
        if ($redis->get($key)) {
            return notice(100000, ['仅允许存在一个操作端']);
        }
        $redis->setex($key, 300, time());

        $permissions = new Permissions();

        try {
            // 1.获取当前环境权限列表
            $list = $permissions->lists([
                'system_type'  => $input['system_type'],
                'version_tag'  => $input['version_tag'],
                'need_deleted' => 1
            ])['list'];

            // 2.备份当前环境url
            Config::set('client.api_url', Config::get('client')["api_url_{$input['env']}"]);

            // 3.调用要发布环境的同步权限接口
            $permissions->sync([
                'list' => $list
            ]);

            // 4.调用要发布环境的发布超管权限接口
            $result = $permissions->release([
                'system_type'   => $input['system_type'],
                'version_tag'   => $input['version_tag'],
                'enterprise_id' => $input['enterprise_id']
            ]);

            // 5.恢复为当前环境
            Config::set('client.api_url', Config::get('client')['api_url_' . Config::get('app.env')]);
        } catch (\Exception $exception) {
            $redis->del([$key]);

            Log::error('同步并发布权限异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'input'   => $input
            ]);

            return error(100000, ["同步并发布权限异常: {$exception->getMessage()}"]);
        }
        $redis->del([$key]);

        return success($result);
    }

    /**
     * 系统类型列表
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function systemTypeList()
    {
        $data = (new Permissions())->systemTypeList();

        return success($data);
    }
}
