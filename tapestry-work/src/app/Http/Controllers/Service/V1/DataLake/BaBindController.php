<?php


namespace App\Http\Controllers\Service\V1\DataLake;


use App\Http\Controllers\Service\V1\BaseController;
use App\Library\DataLake\UserRequestService;
use Illuminate\Validation\ValidationException;
use Throwable;

class BaBindController extends DataLakeBaseController
{

    protected $dataLakeService;
    protected $enterpriseId;
    private $customerUnionID;
    private $customerUserID;
    private $baUserID;
    private $bindStatus;
    private $bindType;
    private $opDate;
    private $baCode;
    private $memberCode;

    public function __construct()
    {
        parent::__construct();
        $this->checkParam();

        $this->enterpriseId    = request()->input('enterprise_id');
        $this->customerUnionID = request()->input('customer_union_id', '');
        $this->customerUserID = request()->input('customer_user_id', '');
        $this->baUserID = request()->input('ba_user_id', '');
        $this->bindStatus = request()->input('bind_status', 0);
        $this->bindType = request()->input('bind_type', 0);
        $this->opDate = request()->input('operation_date', '');
        $this->baCode = request()->input('ba_code', '');
        $this->memberCode = request()->input('member_code', '');
        $this->dataLakeService = new UserRequestService($this->enterpriseId);
    }

    /**
     * @throws ValidationException
     */
    public function checkParam()
    {
        $rule = [
            'customer_union_id' => 'required',
            'enterprise_id' => 'required',
            'customer_user_id' => 'required',
            'ba_user_id' => 'required',
            'bind_status' => 'required',
            'bind_type' => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业id缺失',
            'customer_union_id' => '用户union_id缺失',
            'customer_user_id.required' => '客户的UserID缺失',
            'ba_user_id.required' => 'BA的UserID缺失',
            'bind_status' => '操作类型缺失',
            'bind_type' => '绑定关系类型缺失',
        ];

        $this->validate($this->request, $rule, $message);
    }

    /**
     * @throws Throwable
     */
    public function requestApi()
    {
        return $this->dataLakeService->baBind($this->customerUnionID, $this->customerUserID, $this->baUserID, $this->bindStatus, $this->bindType, $this->opDate, $this->baCode, $this->memberCode);
    }

}