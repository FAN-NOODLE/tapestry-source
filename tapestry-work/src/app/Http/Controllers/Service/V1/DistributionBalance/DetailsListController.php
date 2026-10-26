<?php
/**
 * Created by :PhpStorm
 * User: zhaoyang
 * Date: 2020/10/16
 * Time: 2:47 下午
 */

namespace App\Http\Controllers\Service\V1\DistributionBalance;


use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;
use App\Services\DistributionBalanceService;
use App\Library\Uc\Member;


class DetailsListController extends BaseController
{
    // 接受参数
    private $input;

    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {
        // 参数验证
        $this->getParams();
        // 切换域名
        switchAppHost(['enterprise_id' => $this->input['enterprise_id']]);
        $data = (new DistributionBalanceService())->balanceDetailList($this->input);

        // 门店
        $workerCondition = [
            'enterprise_id' =>  $this->input['enterprise_id'],
            'ids'   => [$this->input['worker_id']]
        ];

        $workerData = StoreWorkerService::list($workerCondition)['list'];
        $storeMember = !empty($workerData) ? current($workerData) : [];

        $memDetail = (new Member())->detail($this->input['domain'], 'contact', [
            'memUid' => $storeMember['mem_uid'],
        ]);

        $dpName = '';
        if(isset($memDetail['dpName'])){
            $dpName    = join(',', array_column($memDetail['dpName'], 'dpName'));
        }
        $storeMember['dpName'] = $dpName;

        return success([
            'list'   => $data['list'],
            'member' => $storeMember,
            'page'   => intval($this->input['page']),
            'size'   => intval($this->input['size']),
            'total'  => $data['total']
        ]);
    }

    /**
     * 检验参数
     *
     * @return bool
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams()
    {
        $rules    = [
            'enterprise_id' => 'required',
            'domain'        => 'required',
            'worker_id'     => 'required'
        ];
        $messages = [
            'enterprise_id.*' => '请确认企业',
            'domain.*'        => '请确认domain',
            'worker_id.*'     => '导购信息不能为空'
        ];

        $this->validate($this->request, $rules, $messages);

        $input = $this->request->all();
        $input['page'] = $input['page'] ?? Constant::PAGE;
        $input['size'] = $input['size'] ?? 10;
        // 接值
        $this->input = $input;
        return true;
    }
}