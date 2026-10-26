<?php
/**
 * Created by PhpStorm.
 * User: 陈朔
 * Date: 2018/11/9
 * Time: 9:54 AM
 */

namespace App\Http\Controllers\App\V1;

use App\Exceptions\Error;
use App\Http\Controllers\Controller;
use App\Models\MemberExtend;
use Illuminate\Support\Facades\Crypt;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Utils\Jaeger;
use Illuminate\Support\Facades\Log;

abstract class BaseController extends Controller
{
    /**
     * 是否必须登录
     *
     * @var bool
     */
    protected $requireLogin = true;

    /**
     * 【企业微信】侧用户的唯一标识
     *
     * @var null
     */
    protected $userId = null;

    /**
     * 企业微信CRS小程序登录状态
     *
     * @var null
     */
    protected $sessionKey = null;

    /**
     * 当前登录的员工ID
     *
     * @var null
     */
    protected $memUid = null;

    /**
     * 企业在 => UC的唯一标识
     *
     * @var null
     */
    protected $qyDomain = null;

    /**
     * 企业在 => 企业微信侧的唯一标识
     *
     * @var null
     */
    protected $corpId = null;

    /**
     * 企业ID
     *
     * @var null
     */
    protected $enterpriseId = null;

    /**
     * 企业hash值
     *
     * @var array|null|string
     */
    protected $enterpriseHash = null;

    /**
     * 企业微信应用标识
     *
     * @var null
     */
    protected $appIdentifier = '';

    /**
     * 当前员工所在企业信息
     *
     * @var null
     */
    protected $enterprise = null;

    /**
     * crs小程序开通状态
     *
     * @var null
     */
    protected $crsStatus = null;

    /**
     * 商城小程序开通状态
     *
     * @var null
     */
    protected $storeStatus = null;

    // shop请求来源
    const  SOURCE_SHOP = 'shop';


    /**
     * 来源work默认为空
     *
     * @var null
     */
    protected $source = '';
    /**
     * BaseController constructor.
     */
    public function __construct()
    {
        parent::__construct();

        (new Jaeger())->serviceJaeger(
            'work',
            [
                'path' => request()->path(),
            ]
        );

        // CRS小程序端开发模拟用户身份
        if (config('app.env') === 'dev' && $this->mockLogin()) {
            return true;
        }

        $token = request()->header('Token');
        // 请求来源
        $this->source = request()->header('Source','');
        $this->enterpriseHash = request()->header('Enterprise-Hash', '');

        /**
         * job（按企业）切换db时用到
         */
        $_SERVER['enterprise_hash'] = $this->enterpriseHash;

        // 如果请求接口是shop则免登录
        if ($this->source == self::SOURCE_SHOP&&!empty($this->enterpriseHash)) {

            // 获取该企业的授权信息
            $enterpriseData     = (new Enterprise())->detail(['hash' => $this->enterpriseHash]);
            if (empty($enterpriseData['data'])) {
                throw new Error(110016, '企业信息有误');
            }

            $enterpriseData     = $enterpriseData['data'];
            $this->enterpriseId = $enterpriseData['id'];
            $this->corpId       = $enterpriseData['corp_id'];
            $this->qyDomain     = $enterpriseData['domain'];
            return true;
        }

        // 非必须登录
        if (!$this->requireLogin) {
            return true;
        }

        if (empty($token)) {
            throw new Error(200005, '无效的用户身份');
        }

        // 解密token
        try {
            $token = Crypt::decrypt($token);
        } catch (\Exception $exception) {
            // 一般是由于token 错误导致无法解密
            throw new Error(110004, 'Token解密失败');
        }

        $tokenData = json_decode($token, true);

        if (empty($tokenData['qyDomain']) || empty($tokenData['memUid'])) {
            throw new Error(110014, '员工信息有误');
        }

        $this->qyDomain   = $tokenData['qyDomain'];
        $this->userId     = $tokenData['userId'];
        $this->sessionKey = $tokenData['sessionKey'];
        $this->corpId     = $tokenData['corpId'];
        $this->memUid     = $tokenData['memUid'];
        // $memberExtend = MemberExtend::query()->where('mem_uid', $this->memUid)->first();

        // if (!$memberExtend || MemberExtend::GUIDE != $memberExtend->role) {
        //     throw new Error(110015, '该员工不是导购，无权限访问小程序');
        // }

        $enterpriseData = (new Enterprise())->detail(
            ['corp_id' => $this->corpId]
        );

        if (empty($enterpriseData['data'])) {
            throw new Error(110016, '企业信息有误');
        }

        $this->enterprise   = $enterpriseData['data'];
        $this->enterpriseId = $this->enterprise['id'];
        $this->crsStatus    = $this->enterprise['crs_status'];
        $this->storeStatus  = $this->enterprise['store_status'];

        return true;
    }

    /**
     * 模拟登录
     *
     * @return bool
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/4 9:52 PM
     */
    private function mockLogin()
    {
        // CRS小程序端开发模拟用户身份
        $mockUserId       = request()->header('Mock-UserId');
        $mockEnterpriseId = request()->header('Mock-EnterpriseId');
        $mockQyDomain     = request()->header('Mock-QyDomain');
        $mockMemUid       = request()->header('Mock-MemUid');

        if (!empty($mockUserId) && !empty($mockEnterpriseId)) {
            $this->userId       = $mockUserId;
            $this->enterpriseId = $mockEnterpriseId;
            $this->qyDomain     = $mockQyDomain;
            $this->memUid       = $mockMemUid;
            return true;
        }

        return false;
    }

    /**
     * @throws \Exception
     */
    public function __destruct()
    {
        (new Jaeger())->flush();
    }
}
