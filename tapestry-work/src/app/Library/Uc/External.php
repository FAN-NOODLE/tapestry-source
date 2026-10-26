<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2018/10/12
 * Time: 5:57 PM
 */

namespace App\Library\Uc;

class External extends Abstracts
{
    /**
     * 外部联系人列表
     *
     * @var string
     */
    const LIST_URL = '%s/b/%s/qy/external/list';

    /**
     * 通过memUid查询外部联系人列表
     *
     * @var string
     */
    const LIST_FOR_MEMBER_URL = '%s/b/%s/qy/external/list-for-member';

    /**
     * 通过外部联系人id获取外部联系人详情
     *
     * @var string
     */
    const DETAIL_URL = '%s/b/%s/qy/external/detail';

    /**
     * 通过外部联系人userid获取外部联系人列表
     *
     * @var string
     */
    const LIST_BY_USER_ID_URL = '%s/b/%s/qy/external/list-by-userid';

    /**
     * 通过code获取员工外部联系人-关联企业客户库
     *
     * @var string
     */
    const CODE_URL = '%s/qy/external/get-by-code';

    /**
     * 同步某个员工的外部联系人
     *
     * @var string
     */
    const SYNC_URL = '%s/qy/external/sync';

    /**
     * 更新外部人员信息
     *
     * @var string
     */
    const UPD_URL = '%s/b/%s/qy/external/update-by-unionid';

    /**
     * 获取企业微信userid/sessionkey
     *
     * @var string
     */
    const INF_URL = '%s/miniprogram/jscode2session';

    /**
     * 导入外部人员
     *
     * @var string
     */
    const IMPORT_URL = '%s/b/%s/qy/external/import';

    /**
     * 查询外部联系人总数
     *
     * @var string
     */
    const CONSUMER_URL = '%s/b/%s/qy/external/number';

    /**
     * 离职成员的外部联系人再分配
     *
     * @var string
     */
    const TRANSFER_URL = '%s/b/%s/qy/external/transfer_external_contact';

    /**
     * 通过外部联系人userId获取外部联系人平台userId
     *
     * @var string
     */
    const EXCHANGE_INFO = '%s/a/%s/%s/qy/external/exchange-info';

    /**
     * 发送欢迎语
     *
     * @var string
     */
    const WELCOME_SEND = '%s/b/%s/qy/external/send/welcome-msg';

    /**
     * 外部联系人user_id转openid
     *
     * @var string
     */
    const CONVERT_TO_OPENID = '%s/b/%s/qy/external/convert_to_openid';

    /**
     * 通过外部联系人userId列表获取外部联系人平台userId列表
     */
    const EXCHANGE_INFO_LIST = '%s/a/%s/%s/qy/external/exchange-info-list';

    /**
     * 通过外部联系人userid获取外部联系人列表
     *
     * @var string
     */
    const CUSTOMER_LIST_URL = '%s/b/%s/qy/external/relation-list';

    /**
     * 同步缺少的客户基础信息及关系数据
     */
    const SYNC_MISSING = '%s/b/%s/qy/external/sync/missing';

    /**
     * 获取外部联系人列表
     *
     * @param $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/15 3:33 PM
     */
    public function list($qyDomain, $params = [])
    {
        $url = vsprintf(self::LIST_URL, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $params);
    }

    /**
     * 添加外部联系人
     *
     * @param $qyDomain
     * @param array $condition
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/15 3:33 PM
     */
    public function import($qyDomain, $condition = [])
    {
        $url = vsprintf(self::IMPORT_URL, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $condition);
    }

    /**
     * 获取外部联系人详情
     *
     * @param $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/15 3:33 PM
     */
    public function detail($qyDomain, $params = [])
    {
        $url = vsprintf(self::DETAIL_URL, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $params);
    }

    /**
     * 更新外部联系人信息
     *
     * @param       $qyDomain
     * @param array $params
     *
     * @return mixed
     * @throws \Exception
     */
    public function update($qyDomain, $params = [])
    {
        $url = vsprintf(self::UPD_URL, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $params);
    }

    /**
     * 离职成员的外部联系人再分配
     *
     * @param       $qyDomain
     * @param array $condition
     *
     * @return mixed
     * @throws \Exception
     */
    public function transfer($qyDomain, $condition = [])
    {
        $url = vsprintf(self::TRANSFER_URL, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $condition);
    }

    /**
     * 通过memUid查询外部联系人列表
     *
     * @param       $qyDomain
     * @param array $condition
     *
     * @return mixed
     * @throws \Exception
     */
    public function listForMember($qyDomain, $condition = [])
    {
        $url = vsprintf(self::LIST_FOR_MEMBER_URL, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $condition);
    }

    /**
     * 通过userid查询外部联系人列表
     *
     * @param       $qyDomain
     * @param array $condition
     *
     * @return mixed
     * @throws \Exception
     */
    public function listByUserId($qyDomain, $condition = [])
    {
        $url = vsprintf(self::LIST_BY_USER_ID_URL, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $condition);
    }

    /**
     * 通过外部联系人userId获取外部联系人平台userId
     * 1-忽略 0-查当前应用套件，默认0
     *
     * @param        $qyDomain
     * @param string $appIdentifier 应用标识
     * @param array  $params
     *
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 9:43 AM
     */
    public function exchangeInfo($qyDomain, $appIdentifier, $params = [])
    {
        $url = vsprintf(self::EXCHANGE_INFO, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $params);
    }

    /**
     * 发送欢迎语
     *
     * @param       $qyDomain
     * @param array $condition
     *
     * @return mixed
     * @throws \Exception
     */
    public function welcomeSend($qyDomain, $condition = [])
    {
        $url = vsprintf(self::WELCOME_SEND, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $condition);
    }

    /**
     * 外部联系人user_id转openid
     *
     * @param       $qyDomain
     * @param array $condition
     *
     * @return mixed
     * @throws \Exception
     */
    public function convertToOpenid($qyDomain, $condition = [])
    {
        $url = vsprintf(self::CONVERT_TO_OPENID, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $condition);
    }

    /**
     * exchangeInfo
     * @desc 通过外部联系人userId列表获取外部联系人平台userId列表
     * @param       $qyDomain
     * @param       $appIdentifier
     * @param array $params
     *
     * @return mixed
     * @date 2020-07-26
     * @throws \Exception
     */
    public function exchangeInfoList($qyDomain, $appIdentifier, $params = [])
    {
        $url = vsprintf(self::EXCHANGE_INFO_LIST, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $params);
    }

    /**
     * 外部联系人关系列表
     *
     * @param       $qyDomain
     * @param array $condition
     *
     * @return mixed
     * @throws \Exception
     */
    public function relationList($qyDomain, $condition = [])
    {
        $url = vsprintf(self::CUSTOMER_LIST_URL, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $condition);
    }

    /**
     * 同步缺少的客户基础信息及关系数据
     *
     * @param       $qyDomain
     * @param array $condition
     *
     * @return mixed
     * @throws \Exception
     */
    public function syncMissing($qyDomain, $condition = [])
    {
        $url = vsprintf(self::SYNC_MISSING, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $condition);
    }
}
