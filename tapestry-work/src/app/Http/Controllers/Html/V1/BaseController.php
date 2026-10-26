<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/13
 * Time: 10:54 AM
 */

namespace App\Http\Controllers\Html\V1;

use App\Exceptions\Error;
use App\Http\Controllers\Controller;
use App\Library\Uc\Member;
use App\Services\EnterpriseService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;

abstract class BaseController extends Controller
{
    /**
     * 是否必须登录
     *
     * @var bool
     */
    protected $requireLogin = true;

    /**
     * 企业微信的唯一ID
     *
     * @var mixed|null
     */
    protected $corpId = null;

    /**
     * 企业微信的domain
     *
     * @var mixed|null
     */
    protected $domain = null;


    /**
     * 企业ID
     *
     * @var null
     */
    protected $enterpriseId = null;

    /**
     * 企业数据 [当前登录人员所在企业]
     *
     * 有 UC epId
     * 有 小程序端的 enterprise_id
     *
     * @var array
     */
    protected $enterprise = [];

    /**
     * 企业微信应用标识
     *
     * @var mixed|null
     */
    protected $agentId = null;

    /**
     * 应用标识符
     *
     * @var array|null|string
     */
    protected $appIdentifier = null;

    /**
     * 企业内部人员ID
     *
     * @var null
     */
    protected $memUid = null;

    /**
     * 内部人员 [当前登录]
     *
     * @var array
     */
    protected $member = [];

    /**
     * 外部人员ID
     *
     * @var null
     */
    protected $extUserId = null;

    /**
     * 外部人员 [当前登录]
     *
     * @var array
     */
    protected $external = [];

    /**
     * 内部人员Cookie Name
     *
     * @var string
     */
    const COOKIE_MEMBER_NAME = 'MEMBER_USER';

    /**
     * Cookie 有效期
     *
     * @var int
     */
    const COOKIE_EXPIRE = 7 * 24 * 60 * 60; // 7天

    /**
     * BaseController constructor.
     *
     * @param Request $request
     * @param int     $more 1代表不查UC多余数据, 0代表查UC多余数据
     *
     * @throws BaseException
     */
    public function __construct(Request $request)
    {

        // 以下两个参数所有接口请求必须以GET形式传递
        $corpId        = $request->get('corp_id');
        $appIdentifier = $request->get('app_identifier');

        if (empty($corpId) || empty($appIdentifier)) {
            throw new Error(10000, '企业标识或者应用标识参数错误');
        }

        $this->appIdentifier = $appIdentifier;
        $this->corpId        = $corpId;

        // 根据 corpId 获取 enterpriseId 和 企业信息
        $this->getEnterprise();

        if (empty($this->enterpriseId) || !is_integer($this->enterpriseId)) {
            throw new Error(10002, '企业ID错误');
        }
        // 增加企业hash
        $_SERVER["Enterprise-Hash"] = $this->enterprise["hash"];

        // 分别读取内部、外部人员的cookie 信息
        $memberCookie   = $this->getMemberCookie();
        // 非强制登录 -> 暂时返回处理
        if ($this->requireLogin === false) {
            return true;
        }

        // 既不存在内部人员也不存在外部人员，身份失效
        if (empty($memberCookie)) {
            throw new Error(10001, '身份信息错误');
        }

        // 解析内部人员cookie
        if (!empty($memberCookie)) {
            $this->parseMemberCookie($memberCookie);
        }

        return true;
    }

    /**
     * 获取enterprise数据
     *
     * @return bool
     * @throws BaseException
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/10/15 4:09 PM
     */
    protected function getEnterprise()
    {
        try {
            $data               = (new EnterpriseService())->getCorpIdInfo($this->corpId);
            $this->enterpriseId = $data['id'];
            $this->domain   = $data['domain'];
            $this->enterprise   = $data;
        } catch (\Exception $exception) {
            throw new Error(1003, $exception->getMessage());
        }

        return true;
    }

    /**
     * 响应成功的结果
     *
     * @param array $data
     * @param int   $httpStatus
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public static function success($data = [], $httpStatus = 200)
    {
        $response = [
            'id'        => md5(uniqid()),
            'status'    => 200,
            'code'      => 0,
            'data'      => $data,
            'message'   => 'success',
            'timestamp' => time()
        ];

        return response()->json($response, $httpStatus);
    }

    /**
     * 设置企业内人员cookie
     *
     * @param $value
     *
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/10/13 1:00 PM
     */
    protected function setMemberCookie($value)
    {
        $expire = (int)self::COOKIE_EXPIRE;
        setcookie(self::COOKIE_MEMBER_NAME, $value, time() + $expire, '/');
    }

    /**
     * 读取企业内人员cookie
     *
     * @return bool
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/10/13 1:00 PM
     */
    protected function getMemberCookie()
    {
        if (isset($_COOKIE[self::COOKIE_MEMBER_NAME])) {
            return $_COOKIE[self::COOKIE_MEMBER_NAME];
        }

        return false;
    }

    /**
     * 解析内部人员cookie
     *
     * @param $cookie
     * @param $more
     *
     * @return bool
     * @throws BaseException
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/10/15 9:52 AM
     */
    protected function parseMemberCookie($cookie)
    {
        $cookieData = json_decode(Crypt::decrypt($cookie), true);

        if (!is_array($cookieData)) {
            throw new Error(1004, '用户信息解析错误');
        }

        $this->memUid = $cookieData['memUid'];

        // 获取用户信息
        try {
            // 调用uc方法获取数据
            $member = (new Member())->detail($this->domain, $this->appIdentifier, [
                'memUid' => $this->memUid,
                'app_identifier' => 'contact',
            ]);

            $this->member = $member;
        } catch (\Exception $exception) {
            throw new Error(10005, $exception->getMessage());
        }

        return true;
    }

    /**
     * 解析外部人员cookie
     *
     * @param $cookie
     *
     * @return bool
     * @throws BaseException
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/10/16 10:30 AM
     */
    protected function parseExternalCookie($cookie)
    {
        $cookieData = json_decode(Crypt::decrypt($cookie), true);

        if (!is_array($cookieData)) {
            throw new BaseException(10006, '外部人员信息解析错误');
        }

        $this->extUserId = $cookieData['memUserId'];
        $this->external  = $cookieData;

        return true;
    }
}
