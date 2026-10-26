<?php

namespace App\Http\Controllers\V1\MaterialChannel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\EnterpriseAdmin;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use App\Exceptions\Error;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\Log;
use Throwable;

class LogListController extends BaseController
{
    /**
     * 活动日志
     */
    public const SYSTEM_OPERATE_DEFAULT = 0;
    public const SYSTEM_OPERATE_ACTIVITY = 10000;
    public const SYSTEM_OPERATE_COLLECT = [
        self::SYSTEM_OPERATE_DEFAULT    => '系统',
        self::SYSTEM_OPERATE_ACTIVITY   => '系统[活动]',
    ];

    /**
     * 渠道日志列表
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->validateInput();

        try {
            $data = (new MaterialChannel())->logList($input);
            $enterpriseAdminSdk = new EnterpriseAdmin();

            // 匹配门店名称
            $storeList = (new Store())->lists([
                'enterprise_id' => $this->enterpriseId,
                'store_ids'     => array_column($data['list'], 'source_id'),
                'size'          => 999999
            ])['list'];
            $storeList = array_column($storeList, 'name', 'id');
            foreach ($data['list'] as $key => $value) {
                $data['list'][$key]['store_name'] = $storeList[$value['source_id']] ?? '商城或其他';
                if (in_array((int) $value['operate_id'], self::SYSTEM_OPERATE_COLLECT, true)) {
                    $adminInfo['username'] = self::SYSTEM_OPERATE_COLLECT[$value['operate_id']];
                } else {
                    try {
                        // 查操作人名称
                        $adminInfo = $enterpriseAdminSdk->detail([
                            'enterprise_id' => $this->enterpriseId,
                            'admin_id'      => $value['operate_id']
                        ])['data'];
                    } catch (Throwable $exception) {
                        $adminInfo['username'] = self::SYSTEM_OPERATE_COLLECT[self::SYSTEM_OPERATE_DEFAULT];
                    }

                }
                $data['list'][$key]['operate_name'] = $adminInfo['username'] ?? '';
            }
        } catch (Error $error) {

            Log::error('获取渠道列表失败', $error->data);
            notice(100000, $error->getMessage());

        } catch (\Exception $e) {

            Log::error('获取渠道列表失败'.$e->getMessage(), $input);
            notice(100000, $e->getMessage());

        }

        return success($data);
    }

    /**
     * 参数验证
     *
     * @param $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    protected function validateInput()
    {
        $this->validate($this->request, [
            'material_id' => 'required|integer'
        ], [
            'material_id.*' => "物料ID不存在",
        ]);

        $input = $this->request->all();
        $input['page'] = $this->request->input('page', 1);
        $input['size'] = $this->request->input('size', 10);
        $input['enterprise_id'] = $this->enterpriseId;

        return $input;
    }
}
