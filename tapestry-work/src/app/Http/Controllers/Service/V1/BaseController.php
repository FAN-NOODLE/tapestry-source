<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/16
 * Time: 9:37 AM
 */

namespace App\Http\Controllers\Service\V1;

use App\Http\Controllers\Controller;
use CrsSdk\Utils\Jaeger;

abstract class BaseController extends Controller
{

    /**
     *  数据权限 ： 自然数据
     */
    const DATA_PERMISSION_TYPE_NATURE = 1;
    /**
     *  数据权限： 门店数据
     */
    const DATA_PERMISSION_TYPE_STORE = 2;
    /**
     *  数据权限： 自然+门店
     */
    const DATA_PERMISSION_TYPE_ALL = 3;

    /**
     * 企业Hash值
     *
     * @var null
     */
    protected $enterpriseHash = null;

    /**
     * BaseController constructor.
     *
     * @throws \Exception
     */
    public function __construct()
    {
        parent::__construct();

        $this->enterpriseHash = request()->header('Enterprise-Hash');

        /**
         * job（按企业）切换db时用到
         */
        $_SERVER['enterprise_hash'] = $this->enterpriseHash;

    }

    /**
     * @throws \Exception
     */
    public function __destruct()
    {

    }
}