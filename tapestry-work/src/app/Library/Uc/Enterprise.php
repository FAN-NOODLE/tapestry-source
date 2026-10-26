<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/13
 * Time: 10:12 AM
 */

namespace App\Library\Uc;

class Enterprise extends Abstracts
{
    /**
     * 企业详情接口
     *
     * @var string
     */
    const DETAIL_URL = '%s/b/%s/enterprise/detail';

    /**
     * 企业列表接口
     *
     * @var string
     */
    const LIST_URL = '%s/s/enterprise/page-list';

    /**
     * 企业信息详情
     *
     * @var string
     */
    const SYSTEM_DETAIL_URL = '%s/s/enterprise/detail';

    /**
     * 获取企业token和key
     *
     * @var string
     */
    const SYSTEM_TOKEN_KEY_URL = '%s/b/%s/ensetting/get';

    /**
     * 给UC回调地址
     *
     * @var string
     */
    const SYSTEM_MODIFY_URL = '%s/b/%s/ensetting/modify';

    /**
     * 企业信息添加
     *
     * @var string
     */
    const SYSTEM_ADD_URL = '%s/s/enterprise/base/register';

    /**
     * 企业信息编辑
     *
     * @var string
     */
    const SYSTEM_UPDATE_URL = '%s/b/%s/enterprise/modify';

    /**
     * 获取企业号登陆信息
     */
    const LOGIN_GET_ADMIN_URL = '%s/qy/login/get-admin';

    /**
     * 单点登录：获取微信套件授权安装登录信息接口
     */
    const WX_SUITE_LOGIN_INFO = '%s/s/qy/login/wx-suite-login-info';

    /**
     * 企业配置信息查询
     */
    const CONFIG_URL = '%s/b/%s/ensetting/query';

    /**
     * 获取企业详情
     *
     * @param string $qyDomain 企业域名
     * @param array $params 筛选条件
     * @return mixed
     * @throws \Exception
     * @author   陈朔  <chenshuo@vchangyi.com>
     * @date     2018/10/13 10:14 AM
     */
    public function detail($qyDomain, $params = [])
    {
        $url = vsprintf(self::DETAIL_URL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

    /**
     * 获取企业详情
     *
     * @param array $params 筛选条件
     * @return mixed
     * @throws \Exception
     * @author   陈朔  <chenshuo@vchangyi.com>
     * @date     2018/10/13 10:14 AM
     */
    public function list($params = [])
    {
        $url = vsprintf(self::LIST_URL, [$this->getUcApiUrl()]);
        return $this->post($url, $params);
    }

    /**
     * 根据crop_id 从uc 获取企业信息
     *
     * @param $corpId
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/15 3:37 PM
     */
    public function enterpriseDetail($corpId)
    {
        $url = vsprintf(self::SYSTEM_DETAIL_URL, [$this->getUcApiUrl()]);
        return $this->post($url, ['corpId' => $corpId]);
    }

    /**
     * 获取外部联系人回调配置参数
     *
     * @param $key
     * @param $qyDomain
     * @return mixed
     * @throws \Exception
     * @author   heshupeng
     */
    public function entSettingGet($qyDomain, $key)
    {
        $url = vsprintf(self::SYSTEM_TOKEN_KEY_URL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, ['key' => $key]);
    }

    /**
     * 设置业务回调
     *
     * @param $qyDomain
     * @param $callback
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/30 11:16 AM
     */
    public function entSettingModify($qyDomain, $callback)
    {
        $url = vsprintf(self::SYSTEM_MODIFY_URL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, ['wxExternalCallback' => $callback]);
    }

    /**
     * 添加企业信息
     *
     * @param $params
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/15 3:37 PM
     */
    public function enterpriseAdd($params)
    {
        $url = vsprintf(self::SYSTEM_ADD_URL, [$this->getUcApiUrl()]);
        return $this->post($url, $params);
    }

    /**
     * 更新企业信息
     *
     * @param $qyDomain
     * @param $params
     * @return mixed
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/15 3:37 PM
     */
    public function enterpriseUpdate($qyDomain, $params)
    {
        $url = vsprintf(self::SYSTEM_UPDATE_URL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

    /**
     * 单点登录：获取微信套件授权安装登录信息接口
     *
     * @param $params
     * @return mixed
     * @throws \Exception
     */
    public function wxSuiteLoginInfo($params)
    {
        $url = vsprintf(self::WX_SUITE_LOGIN_INFO, [$this->getUcApiUrl()]);

        return $this->post($url, $params);
    }

    /**
     * 获取企业号登陆信息接口
     *
     * @param $enterprise
     * @return mixed
     * @throws \Exception
     */
    public function configInfo($enterprise)
    {
        $url = vsprintf(self::CONFIG_URL, [$this->getUcApiUrl(), $enterprise]);

        return $this->post($url);
    }

    /**
     * 获取企业号登陆信息接口
     *
     * @param $params
     * @return mixed
     * @throws \Exception
     */
    public function authAdminInfo($params)
    {
        $url = vsprintf(self::LOGIN_GET_ADMIN_URL, [$this->getUcApiUrl()]);

        return $this->get($url, $params);
    }
}
