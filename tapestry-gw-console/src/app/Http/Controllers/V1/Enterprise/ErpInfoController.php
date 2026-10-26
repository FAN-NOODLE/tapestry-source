<?php
/**
 * 查询企业对接erp的信息
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019年11月7日
 * Time: 21:49:44
 */
namespace App\Http\Controllers\V1\Enterprise;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Open\Enterprise;

class ErpInfoController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_hash'] = $this->enterprise['hash'];
        $data =  (new Enterprise())->detail($input)['data'];
        return success($data);
    }
}