<?php
/**
 * Created by PhpStorm.
 * User: weihua
 * Date: 2020/9/8
 * Time: 11:34
 */

namespace App\Http\Controllers\V1\User;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\StoreWorker;
use Illuminate\Support\Facades\Log;

class GuideListByStoreOrDepartmentController extends BaseController
{
    public function index()
    {
        $input    = $this->request->input();
        $rules    = [
            'store_ids' => 'array',
            'dp_ids'    => 'array',
            'worker_name'=>'string',
            'worker_status'=>'integer',
        ];
        $messages = [
            'store_ids.array' => '门店ids为空必须为数组',
            'dp_ids.array'    => '部门ids必须为数组',
            'worker_name.string'    => '导购名称必须为字符串',
            'worker_status.integer'    => '导购状态必须为int类型',
        ];
        $this->validate($this->request, $rules, $messages);

        if(empty($input['store_ids']) && empty($input['dp_ids']) &&
            empty($input['worker_name']) ){
            error(320164);
        }

        $input['enterprise_id'] = $this->enterpriseId;
        $input['domain']        = $this->qyDomain;
        Log::info('根据门店获取导购列表接收参数', [$input]);
        // 获取导购列表
        $guideList = (new StoreWorker())->guideListByStore($input);

        return success(['list' => $guideList]);
    }
}