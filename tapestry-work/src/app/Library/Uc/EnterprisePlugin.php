<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/12
 * Time: 10:33 AM
 */

namespace App\Library\Uc;

use Illuminate\Support\Facades\Log;

class EnterprisePlugin extends Abstracts
{
    /**
     * 手动的、特殊应用安装接口
     *
     * @var string
     */
    const APP_INSTALL_URL = '%s/s/enplugin/one-to-many-Install';

    /**
     * 外部联系人权限接口
     *
     * 业务端将 secret 提交到UC
     *
     * @var
     */
    const MODIFY_URL = '%s/b/%s/ensetting/modify';

    /**
     * 应用列表接口
     *
     * @var string
     */
    const LIST_URL = '%s/b/%s/enplugin/list';

    /**
     * 应用详情接口
     *
     * @var string
     */
    const DETAIL_URL = '%s/b/%s/enplugin/detail';

    /**
     * 应用安装接口
     *
     * @var string
     */
    const INSTALL_URL = '%s/b/%s/enplugin/install';

    /**
     * 应用卸载接口
     *
     * @var string
     */
    const UNINSTALL_URL = '%s/b/%s/enplugin/uninstall';

    /**
     * 应用更新接口
     *
     * @var string
     */
    const APP_UPDATE_URL = '%s/b/%s/enplugin/secret/modify';

    /**
     * 发送小程序消息
     */
    const SEND_MINI_MESSAGE = '%s/a/%s/%s/qy/send_mini';

    /**
     * 获取应用列表接口
     *
     * @param string $qyDomain 企业域名
     * @param array $params 筛选条件
     * @return array|\Psr\Http\Message\StreamInterface
     * @throws \Exception
     */
    public function list($qyDomain, $params = [])
    {
        $url = vsprintf(self::LIST_URL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

    /**
     * 获取应用详情接口
     *
     * @param string $qyDomain 企业域名
     * @param array $params 筛选条件
     * @return array|\Psr\Http\Message\StreamInterface
     * @throws \Exception
     */
    public function detail($qyDomain, $params = [])
    {
        $url = vsprintf(self::DETAIL_URL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

    /**
     * 安装应用接口
     *
     * @param string $qyDomain 企业域名
     * @param array 企业应用信息
     * @return array|\Psr\Http\Message\StreamInterface
     * @throws \Exception
     */
    public function install($qyDomain, $plugin)
    {
        $url = vsprintf(self::INSTALL_URL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $plugin);
    }

    /**
     * 卸载应用接口
     *
     * @param string $qyDomain 企业域名
     * @param array 企业应用信息
     * @return array|\Psr\Http\Message\StreamInterface
     * @throws \Exception
     */
    public function uninstall($qyDomain, $plugin)
    {
        $url = vsprintf(self::UNINSTALL_URL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $plugin);
    }

    /**
     * 特殊应用安装接口
     *
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 5:44 PM
     */
    public function appInstall($params)
    {
        $url = vsprintf(self::APP_INSTALL_URL, [$this->getUcApiUrl()]);
        return $this->post($url, $params);
    }

    /**
     * UC说的是外部联系人安装用这个
     *
     * 乱七八糟，鬼才知道
     *
     * @param $qyDomain
     * @param $params
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 6:11 PM
     */
    public function modify($qyDomain, $params)
    {
        $url = vsprintf(self::MODIFY_URL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

    /**
     * 特殊应用更新接口
     *
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   zhaoqi  zhaoqi@vchangyi.com
     * @date     2019/02/27 16:21 PM
     */
    public function appUpdate($qyDomain, $params)
    {
        $url = vsprintf(self::APP_UPDATE_URL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

    /**
     * 发送小程序消息
     * @param $qyDomain
     * @param $appIdentifier
     * @param $params
     * @return mixed
     * @throws \Exception
     */
    public function sendMiniMessage($qyDomain, $appIdentifier, $params)
    {
        $url = vsprintf(self::SEND_MINI_MESSAGE, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $params);
    }
}