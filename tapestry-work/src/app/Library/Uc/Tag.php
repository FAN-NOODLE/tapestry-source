<?php


namespace App\Library\Uc;

/**
 * http://t-rst.vchangyi.com/doc.html 外部联系人标签管理
 * Class Tag
 * @package App\Library\Uc
 */
class Tag extends Abstracts
{
    /**
     * 同步企业标签库
     */
    const SYNC_TAG = '%s/b/%s/qy/external/enterprise-tag/sync';
    /**
     * 查询外部联系人标签关系信息
     *
     * @var string
     */
    const TAG_SEARCH = '%s/b/%s/qy/external/tag-relation/list';

    /**
     * 获取企业标签库
     *
     * @var string
     */
    const TAG_LIST = '%s/b/%s/qy/external/enterprise-tag/list';

    /**
     * 获取企业标签组列表
     *
     * @var string
     */
    const TAG_GROUP_LIST = '%s/b/%s/qy/external/tag-group/page-list';

    /**
     * 根据标签列表
     *
     * @var string
     */
    const TAG_LIST_PAGE = '%s/b/%s/qy/external/tag/page-list';

    /**
     * 添加企业客户标签
     *
     * @important 如果在原有标签组下新增标签则需要调用此新增接口
     *
     * @var string
     */
    const TAG_ADD = '%s/b/%s/qy/external/tag/add';

    /**
     * 编辑企业客户标签
     *
     * @important 只能对已有标签标签组进行编辑
     *
     * @var string
     */
    const TAG_UPDATE = '%s/b/%s/qy/external/tag/update';

    /**
     * 删除企业客户标签
     *
     * @var string
     */
    const TAG_DEL = '%s/b/%s/qy/external/tag/delete';

    /**
     * 同步企业标签库
     * @param $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     */
    public function syncTag($qyDomain, $params = [])
    {
        $url = vsprintf(self::SYNC_TAG, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }
    /**
     * 查询外部联系人标签关系信息
     * @param string $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     */
    public function tagDetail($qyDomain, $params = [])
    {
        $url = vsprintf(self::TAG_SEARCH, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

    /**
     * 添加标签组-标签
     *
     * @param $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   李坚  lijian@vchangyi.com
     * @date     2020/02/27 11:00 AM
     */
    public function add($qyDomain, $params = [])
    {
        $url = vsprintf(self::TAG_ADD, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

    /**
     * 获取企业标签库列表
     *
     * @param $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   李坚  lijian@vchangyi.com
     * @date     2020/02/28 19:00 PM
     */
    public function allList($qyDomain) {
        $url = vsprintf(self::TAG_LIST, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url);
    }

    /**
     * 获取企业标签组列表
     *
     * @param $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   李坚  lijian@vchangyi.com
     * @date     2020/02/27 11:00 AM
     */
    public function list($qyDomain, $params = [])
    {
        $url = vsprintf(self::TAG_GROUP_LIST, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

    /**
     * 根据组id获取标签列表
     *
     * @param $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   李坚  lijian@vchangyi.com
     * @date     2020/02/27 11:00 AM
     */
    public function pageList($qyDomain, $params = [])
    {
        $url = vsprintf(self::TAG_LIST_PAGE, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }


    /**
     * 更新标签组-标签
     *
     * @param $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   李坚  lijian@vchangyi.com
     * @date     2020/02/27 11:00 AM
     */
    public function update($qyDomain, $params = [])
    {
        $url = vsprintf(self::TAG_UPDATE, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

    /**
     * 删除标签组-标签
     *
     * @param $qyDomain
     * @param array $params
     * @return mixed
     * @throws \Exception
     * @author   李坚  lijian@vchangyi.com
     * @date     2020/02/27 11:00 AM
     * @request example
     * {
     *   "groupIds": [],
     *   "tagIds": [],
     *   "thirdGroupIds": [],
     *   "thirdTagIds": []
     * }
     */
    public function delete($qyDomain, $params = [])
    {
        $url = vsprintf(self::TAG_DEL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

}