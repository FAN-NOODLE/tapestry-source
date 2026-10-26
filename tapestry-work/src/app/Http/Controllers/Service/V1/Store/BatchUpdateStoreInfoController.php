<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StoreDepartment;
use App\Services\StoreService;

class BatchUpdateStoreInfoController extends BaseController
{

    /**
     * 批量更新门店信息
     * @return array
     * @throws \Exception
     */
    public function index(): array
    {
        $params = $this->request->only([
            'enterprise_id',
            'data'
        ]);
        if (empty($params['enterprise_id'])) {
            throw new Error(600000, '企业id不能为空');
        }
        $fields = ['id', 'wechat_department_id', 'third_store_code', 'store_name', 'store_type',
            'address', 'sa_share_order_url', 'sa_check_order_url', 'popup_store_keys','email',
            'city', 'phone', 'inventory_name'];
        $list = StoreService::departmentList(['enterprise_id' => $params['enterprise_id']], false, $fields)['list'];
        $newList = [];
        foreach ($list as $val) {
            $newList[] = array_values($val);
        }
        foreach ($params['data'] as $key => $val) {
            if ($key == 0) {
                if ($val[0] != '系统id' || $val[1] != '企微部门id' || $val[2] != '门店编号'
                    || $val[3] != '门店对外名称' || $val[4] != '门店属性:1奥莱 2零售' || $val[5] != '门店地址'
                    || $val[6] != '库存查询门店名' || $val[7] != '门店所属城市' || $val[8] != '门店电话'
                    || $val[9] != 'SA分享订单链接' || $val[10] != '店长查看门店链接' || $val[11] != '签名密钥' || $val[12] != '邮箱'
                    ) {
                    throw new Error(100000, '模板文件异常', 200);
                }
                continue;
            }
            //处理excel数据自动转为字符串的问题
            $val[0] = intval($val[0]); //系统id
            $val[1] = intval($val[1]); //企微部门id
            $val[4] = intval($val[4]); //门店属性
            $updateArr = [];
            if (!in_array($val, $newList)) {
                //如果有差异且新数据不为空则更新
                if (!empty($val[3])) {
                    $updateArr['store_name'] = $val[3];
                }
                if (!empty($val[4])) {
                    $updateArr['store_type'] = $val[4];
                }
                if (!empty($val[5])) {
                    $updateArr['address'] = $val[5];
                }
                if (!empty($val[12])) {
                    $updateArr['email'] = $val[12];
                }
                $updateArr['inventory_name'] = $val[6];
                $updateArr['city'] = $val[7];
                $updateArr['phone'] = $val[8];
                $updateArr['sa_share_order_url'] = $val[9];
                $updateArr['sa_check_order_url'] = $val[10];
                $updateArr['popup_store_keys'] = $val[11];
                if (empty($updateArr['sa_share_order_url']) || empty($updateArr['sa_check_order_url']) || empty($updateArr['popup_store_keys'])) {
                    //如果popup store其中一个信息为空则直接关闭该门店popup store功能
                    $updateArr['popup_store_status'] = StoreDepartment::POPUP_STORE_CLOSE;
                }
                if (!empty($updateArr)) {
                    $updateArr['sync_status'] = 0;
                    StoreDepartment::query()->where('id', $val[0])->update($updateArr);
                }
            }

        }

        return success();
    }


}
