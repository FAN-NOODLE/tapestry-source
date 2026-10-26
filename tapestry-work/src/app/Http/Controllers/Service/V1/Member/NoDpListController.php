<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:47
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Library\Uc\Member;
use Illuminate\Http\Request;
use App\Services\MemberUserRelationService;
use App\Http\Controllers\Service\V1\BaseController;

class NoDpListController extends BaseController
{
    /**
     * @var MemberUserRelationService
     */
    protected $memberUserRelationService;

    /**
     * CountController constructor.
     *
     */
    public function __construct() {
        parent::__construct();
    }

    /**
     * 没有部门员工列表
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', '');
        $params        = $request->post('params', []);

        $result = (new Member())->noDpList($qyDomain, $appIdentifier, $params);

        return success($result);
    }
}
