<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/25 0025
 * Time: 19:59
 */

namespace App\Http\Controllers\Service\V1\DataLake;


use App\Http\Controllers\Service\V1\BaseController;
use App\Library\DataLake\UserRequestService;
use Illuminate\Http\Request;


abstract class DataLakeBaseController extends BaseController
{
    protected $dataLakeService;
    protected $type;
    protected $value;
    protected $enterpriseId;

    public function __construct()
    {
        parent::__construct();
        $this->checkParam();
        $this->enterpriseId    = request()->input('enterprise_id');
        $this->type            = request()->input('type', '');
        $this->value           = request()->input('value', '');
        $this->dataLakeService = new UserRequestService($this->enterpriseId);
    }

    public function index()
    {
        return success(
            [
                'id'      => defined('LOG_ID') ? LOG_ID : '',
                'code'    => '',
                'status'  => 200,
                'message' => '',
                'error'   => '',
                'data'    => [], //$this->requestApi()
            ]
        );
    }

    public function checkParam()
    {
        $rules = [
            'enterprise_id' => 'required',
            'type'          => 'required',
            'value'         => 'required',
        ];

        $message = [
            'enterprise_id.required' => '企业id缺失',
            'type.required'          => '查询方式缺失',
            'value.required'         => '查询值缺失',
        ];
        $this->validate($this->request, $rules, $message);
    }

    abstract function requestApi();

    public function validate(Request $request, array $rules, array $messages = [], array $customAttributes = [])
    {
        
    }

}