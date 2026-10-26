<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/12
 * Time: 10:35 AM
 */

namespace App\Library\Uc;

use Illuminate\Support\Facades\Log;

class Member extends Abstracts
{
    /**
     * 增加员工接口
     *
     * @var string
     */
    const ADD_URL = '%s/a/%s/%s/qy/member/add';

    /**
     * 编辑员工接口
     *
     * @var string
     */
    const EDIT_URL = '%s/a/%s/%s/qy/member/update';

    /**
     * 员工列表接口
     *
     * @var string
     */
    const LIST_URL = '%s/a/%s/%s/qy/member/list';

    /**
     * 员工详情接口
     *
     * @var string
     */
    const DETAIL_URL = '%s/a/%s/%s/qy/member/get';

    /**
     * 根据userId获取企业人员详情接口
     * @var string
     */
    const DETAIL_BY_USER_ID = '%s/a/%s/%s/qy/member/get-by-userid';

    /**
     * 删除员工
     *
     * @var string
     */
    const DELETE_URL = '%s/a/%s/%s/qy/member/delete';

    /**
     * 员工二维码
     *
     * @var string
     */
    const QR_CODE_URL = '%s/a/%s/%s/qy/member/qr-code/get';

    /**
     * 忽略权限查询人员列表
     *
     * @var string
     */
    const IGNORE_PERMISSIONS_LIST = '%s/a/%s/%s/qy/member/ignore-permissions/list';

    /**
     * 批量更新员工角色
     *
     * @var string
     */
    const UPDATE_ROLE_URL = '%s/a/%s/%s/qy/member/bat-update';

    /**
     * 员工身份校验
     *
     * @var string
     */
    const CHECK_MEMBER_URL = '%s/a/%s/%s/qy/member/check';

    /**
     * 获取企业微信素材 mediaId
     *
     * @var string
     */
    const IMAGE_TO_MEDIA_ID_URL = '%s/a/%s/%s/qy/media/upload';

    /**
     * 新增联系我二维码
     *
     * @var string
     */
    const CONTACT_QR_CODE_URL = '%s/b/%s/qy/contact-way/add';

    /**
     * 更新联系我二维码
     *
     * @var string
     */
    const CONTACT_QR_CODE_URL_UPDATE = '%s/b/%s/qy/contact-way/update';

    /**
     * 删除联系我二维码
     *
     * @var string
     */
    const CONTACT_QR_CODE_URL_DEL = '%s/b/%s/qy/contact-way/del';

    /**
     * 同步通讯录
     *
     * @var string
     */
    const CONTACT_SYNC = '%s/b/%s/qy/member/sync';

    /**
     * 通讯录同步状态
     *
     * @var string
     */
    const CONTACT_SYNC_STATUS = '%s/b/%s/qy/member/sync-status';

    /**
     * 没有部门人员列表
     *
     * @var string
     */
    const GET_NO_DEPARTMENT_MEMBER = '%s/a/%s/%s/qy/member/ignore-permissions/get-no-department-member';

    /**
     * 企业人员基本信息
     *
     * @var string
     */
    const BASIC_INFO_URL = '%s/a/%s/%s/qy/member/basic-info/get';

    /**
     * 企业人员基本信息分页列表接口
     */
    const BASIC_LIST_URL = '%s/a/%s/%s/qy/member/ignore-permissions/basic-list';

    /**
     * 获取员工联系客户统计数据
     */
    const BEHAVIOR_DATA = '%s/b/%s/qy/external/behavior-data/get';

    /**
     * 添加员工
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param array $condition
     * @return mixed
     * @throws \Exception
     */
    public function add($qyDomain, $appIdentifier, $condition = [])
    {
        $url = vsprintf(self::ADD_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $condition);
    }

    /**
     * 编辑员工
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param array $condition
     * @return mixed
     * @throws \Exception
     */
    public function edit($qyDomain, $appIdentifier, $condition = [])
    {
        $url = vsprintf(self::EDIT_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $condition);
    }

    /**
     * 员工列表
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param array $condition
     * @return mixed
     * @throws \Exception
     */
    public function list($qyDomain, $appIdentifier, $condition = [])
    {
        $appIdentifier = $this->getAppIdentifier($qyDomain);

        $url = vsprintf(self::LIST_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $condition);
    }

    /**
     * 员工详情
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param array $condition
     * @return mixed
     * @throws \Exception
     */
    public function detail($qyDomain, $appIdentifier, $condition = [])
    {
        $appIdentifier = $this->getAppIdentifier($qyDomain);

        $url = vsprintf(self::DETAIL_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $condition);
    }

    /**
     * 根据userId获取企业人员详情接口
     * @param $qyDomain
     * @param $appIdentifier
     * @param array $condition
     * @return mixed
     * @throws \Exception
     */
    public function detailByUserId($qyDomain, $appIdentifier, $condition = [])
    {
        $url = vsprintf(self::DETAIL_BY_USER_ID, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $condition);
    }

    /**
     * 删除员工
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param array $condition
     * @return mixed
     * @throws \Exception
     */
    public function delete($qyDomain, $appIdentifier, $condition = [])
    {
        $url = vsprintf(self::DELETE_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $condition);
    }

    /**
     * 查询员工二维码
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/15 3:35 PM
     */
    public function qrCode($qyDomain, $appIdentifier, $params = [])
    {
        $appIdentifier = $this->getAppIdentifier($qyDomain);

        $url = vsprintf(self::QR_CODE_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $params);
    }

    /**
     * 忽略权限人员列表
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/15 3:35 PM
     */
    public function ignorePermissionsList($qyDomain, $appIdentifier, $params = [])
    {
        $appIdentifier = $this->getAppIdentifier($qyDomain);

        $url = vsprintf(self::IGNORE_PERMISSIONS_LIST, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $params);
    }

    /**
     * 批量更新员工角色
     *
     * @param       $qyDomain
     * @param       $appIdentifier
     * @param array $params
     *
     * @return mixed
     * @throws \Exception
     */
    public function updateRole($qyDomain, $appIdentifier, $params = [])
    {
        $url = vsprintf(self::UPDATE_ROLE_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);

        return $this->post($url, $params);
    }

    /**
     * 员工身份校验（手机号）
     *
     * @param       $qyDomain
     * @param       $appIdentifier
     * @param array $params
     *
     * @return mixed
     * @throws \Exception
     */
    public function check($qyDomain, $appIdentifier, $params = [])
    {
        $appIdentifier = $this->getAppIdentifier($qyDomain);

        $url = vsprintf(self::CHECK_MEMBER_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);

        return $this->post($url, $params);
    }

    /**
     * 获取企业微信素材 media_id
     *
     * @param       $corpId
     * @param       $appIdentifier
     * @param array $params
     *
     * @return mixed
     * @throws \Exception
     */
    public function imageToMediaId($corpId, $appIdentifier, $params = [])
    {
        $appIdentifier = $this->getAppIdentifier($corpId);
        $url           = vsprintf(self::IMAGE_TO_MEDIA_ID_URL, [$this->getUcApiUrl(), $corpId, $appIdentifier]);

        Log::debug('uc-media-api-url' . $url);

        return $this->postMultipartUpload($url, $params);
    }

    /**
     * 新增联系我二维码
     *
     * @param       $qyDomain
     * @param array $params
     *
     * @return mixed
     * @throws \Exception
     */
    public function getContactQrCode($qyDomain, $params = [])
    {
        $url = vsprintf(self::CONTACT_QR_CODE_URL, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $params);
    }

    /**
     * 删除联系我二维码
     *
     * @param       $qyDomain
     * @param array $params
     *
     * @return mixed
     * @throws \Exception
     */
    public function delContactQrCode($qyDomain, $params = [])
    {
        $url = vsprintf(self::CONTACT_QR_CODE_URL_DEL, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $params);
    }

    /**
     * 更新联系我二维码
     *
     * @param       $qyDomain
     * @param array $params
     *
     * @return mixed
     * @throws \Exception
     */
    public function updateContactQrCode($qyDomain, $params = [])
    {
        $url = vsprintf(self::CONTACT_QR_CODE_URL_UPDATE, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $params);
    }

    /**
     * 同步通讯录
     *
     * @param       $qyDomain
     * @param array $params
     *
     * @return mixed
     * @throws \Exception
     */
    public function syncContact($qyDomain, $params = [])
    {
        $url = vsprintf(self::CONTACT_SYNC, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $params);
    }

    /**
     * 获取同步通讯录状态
     *
     * @param       $qyDomain
     * @param array $params
     *
     * @return mixed
     * @throws \Exception
     */
    public function getSyncContactStatus($qyDomain, $params = [])
    {
        $url = vsprintf(self::CONTACT_SYNC_STATUS, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $params);
    }

    /**
     * 没有部门员工列表
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param array $condition
     * @return mixed
     * @throws \Exception
     */
    public function noDpList($qyDomain, $appIdentifier, $condition = [])
    {
        $appIdentifier = $this->getAppIdentifier($qyDomain);

        $url = vsprintf(self::GET_NO_DEPARTMENT_MEMBER, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $condition);
    }

    /**
     * 获取企业人员基本信息
     * @param $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     */
    public function basicInfo($qyDomain, $params = [])
    {
        $appIdentifier = $this->getAppIdentifier($qyDomain);

        $url = vsprintf(self::BASIC_INFO_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $params);
    }

    /**
     * 企业人员基本信息分页查询
     * @param $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     */
    public function basicList($qyDomain, $params = [])
    {
        $appIdentifier = $this->getAppIdentifier($qyDomain);

        $url = vsprintf(self::BASIC_LIST_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $params);
    }

    /**
     * 获取员工联系客户统计数据
     *
     * @param $qyDomain
     * @param array $params
     * @return array|mixed
     * @throws \Exception
     */
    public function getBehaviorData($qyDomain, $params = [])
    {
        $url = vsprintf(self::BEHAVIOR_DATA, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $params);
    }
}
