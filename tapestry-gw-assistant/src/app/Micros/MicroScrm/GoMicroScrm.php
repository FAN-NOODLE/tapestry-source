<?php

namespace App\Micros\MicroScrm;

use App\Micros\Common;

class GoMicroScrm
{
    /**
     * 排行榜
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function rank($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/datacompass/rank_data')->post($params);
    }

    /**
     * 外部联系人列表
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/10/20
     */
    public static function getExternalUserList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-user/list')->post($params);
    }

    /**
     * 外部联系人列表
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/10/20
     */
    public static function getExternalUserListCk($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-user/list-ck')->post($params);
    }


    /**
     * 外部联系人筛选条件标签
     *
     * @return array|mixed
     * @Time: 2022/10/21
     */
    public static function getChoiceOption($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-user/choice-option')->post($params);
    }

    /**
     * 导购信息
     *
     * @return array|mixed
     * @Time: 2022/10/21
     */
    public static function getMemberInfo($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/member/info')->post($params);
    }

    /**
     * 客户画像详情
     *
     * @return array|mixed
     * @Time: 2022/10/21
     */
    public static function portraitDetail($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/portrait-detail')->post($params);
    }

    /**
     * 侧边栏社群成员列表
     *
     * @return array|mixed
     * @Time: 2022/10/21
     */
    public static function groupUserChatMember($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/sidebar/group-chat/group-users-add-employees')->post($params);
    }

    /**
     * 顾客企微标签
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/10/25
     */
    public static function userTags($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/tag/detail/user-tags')->post($params);
    }

    /**
     * 顾客印象标签
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/10/25
     */
    public static function userImpressionTags($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/tag/impression/user-tags-list')->post($params);
    }

    /**
     * 绩效榜单活跃榜和顾客榜
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function performanceListActiveCustomers($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/member/performance-list/active-customers')->post($params);
    }

    /**
     * 绩效榜单 明星榜-活跃榜和顾客榜
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function performanceListStarActiveCustomers($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/member/performance-list/star-active-customers')->post($params);
    }

    /**
     * 绩效榜 邀约榜单
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function performanceListInvitation($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/member/performance-list/invitation')->post($params);
    }

    /**
     * 绩效榜 明星榜-邀约榜单
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function performanceListStarInvitation($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/member/performance-list/star-invitation')->post($params);
    }

    /**
     * 用户详情
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function externalUserDetail($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-user/detail')->post($params);
    }

    /**
     * 小程序端门店员工根据查看分组权限查询各素材列表
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/1
     */
    public static function selectList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/select-list')->post($params);
    }

    /**
     * 小程序端门店员工根据查看分组权限查询各素材上新数量
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/1
     */
    public static function getNewSourceNum($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/get-new-num')->post($params);
    }

    /**
     * 我的页面
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/11
     */
    public static function myCenter($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/my-center')->post($params);
    }

    /**
     * 我的页面 - 添加专属导购
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/16
     */
    public static function mySa($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/my-sa')->post($params);
    }


    /**
     * 导购或店长获取顾客邀约和播种数量
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/17
     */
    public static function getUserDetailTipsNum($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/user-detail-tips-num')->post($params);
    }

    /**
     * 顾客和导购的绑定关系
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/17
     */
    public static function getUserBindRelation($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/user-bind-relation')->post($params);
    }

    /**
     * 修改备注
     *
     * @param $params
     *
     * @return array|mixed
     * @author: wangle@vchangyi.com
     * @Time: 2022/11/18
     */
    public static function userUpdateRemark($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/update-remark')->post($params);
    }

    /**
     * 修改描述
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function userUpdateDescription($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/update-description')->post($params);
    }

    /**
     * 修改顾客地址尺码生日
     * @param $params
     *
     * @return array|mixed
     */
    public static function userSaveInfoAdd($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/save-info-add')->post($params);
    }

    /**
     * 操作顾客的印象标签
     *
     * @param $params
     *
     * @return array|mixed
     * @author: wangle@vchangyi.com
     * @Time: 2022/11/21
     */
    public static function userImpressionTagsOperate($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/tags/impression/operate')->post($params);
    }

    /**
     * 操作顾客的企微标签
     *
     * @param $params
     *
     * @return array|mixed
     * @author: wangle@vchangyi.com
     * @Time: 2022/11/21
     */
    public static function userWorkTagsOperate($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/tags/work/operate')->post($params);
    }


    /**
     * 邀请函模板详情
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function inviteInfo($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/invite/template/info')->post($params);
    }

    /**
     * 导购或管理的导购userid
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function getMemberUserIds($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-user/get-member-userid')->post($params);
    }

    /**
     * 根据cdp的会员卡号 查询顾客信息
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function getCustomerList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-user/customer-list')->post($params);
    }

    /**
     * 业绩榜单
     * @return array|mixed
     */
    public static function achievement($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/member/performance-list/achievement')->post($params);
    }


    /**
     * sw顾客榜单
     * @return void
     */
    public static function swCustomers($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/member/performance-list/sw-customers')->post($params);
    }

    /**
     * sw活跃榜单
     * @return void
     */
    public static function swActive($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/member/performance-list/sw-active')->post($params);
    }

    /**
     * 小程序登录
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function code2session($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/members/code-session')->post($params);
    }

    /**
     * 素材分组归属门店列表
     *
     * @param $params
     * @return array|mixed
     */
    public static function getSourceRelationList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/group-relation/list')->post($params);
    }

    /**
     * 根据门店获取素材分组列表
     *
     * @param $params
     * @return array|mixed
     */
    public static function getListForStore($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/group-relation/get-list-for-store')->post($params);
    }

    /**
     * 获取全选的外部联系人总数
     *
     * @param $params
     * @return array|mixed
     */
    public static function getExternalUserTotal($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-user/select-all-num')->post($params);
    }

    /**
     * 获取全选的外部联系人的external_user_id，最多返回20000个
     *
     * @param $params
     * @return array|mixed
     */
    public static function getExternalUserIds($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-user/select-all-list')->post($params);
    }

    /**
     * @param $params
     * @return array|mixed
     */
    public static function getTakeoverMemberList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/on-inherit/takeover-member-list')->post($params);
    }

    public static function getOnInheritSelectAllNum($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/on-inherit/select-all-num')->post($params);
    }

    /**
     * 用户确认邀请函时，SW和KS企业需要验证到店时间限制
     * @param array $params
     * @return array|mixed
     */
    public function checkReserveCreate(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/reserve/order/check-reserve-time')->post($params);
    }

    /**
     * 素材产品类别列表
     *
     * @param array $params
     * @return array|mixed
     */
    public function getSourceCategoryLists(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/category/list')->post($params);
    }

    /**
     * 素材上报计次
     *
     * @param array $params
     * @return array|mixed
     */
    public function reportTimes(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/report-times')->post($params);
    }

    /**
     * 验证是否用户首次打开素材交互提示
     *
     * @param array $params
     * @return array|mixed
     */
    public function checkFirst(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/check-first')->post($params);
    }


    /**
     * 门店列表
     *
     * @param $params
     *
     * @return array|mixed
     */
    public function storeList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/department/store-list')->post($params);
    }

    /**
     * 首次打开某个页面的引导标记
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/10/20
     */
    public function firstOpen($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/member/first-open')->post($params);
    }

    /**
     * 小品牌素材库查询商品素材
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2024/5/7
     */
    public static function brandSearchList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/brand/search-list')->post($params);
    }

    /**
     * 小品牌素材库查询商品素材查询sku
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2024/5/7
     */
    public static function brandSearchListGroupSku($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/brand/search-list-group-sku')->post($params);
    }

    /**
     * 小品牌素材库查询商品素材查询sku
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2024/5/7
     */
    public static function talkieImageSearch($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/brand/talkie-image-search')->post($params);
    }

    /**
     * 根据upcCode获取商品素材缓存
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2024/5/7
     */
    public static function getCacheInfo($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/brand/get-cache-info')->post($params);
    }

    /**
     * 小程序端门店员工收藏查询各素材列表
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/1
     */
    public static function collectSelectList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/collect-select-list')->post($params);
    }

    /**
     * 收藏素材
     * @param $params
     * @return array|mixed
     */
    public static function collect($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/collect/edit')->post($params);
    }


    /**
     * 小品牌素材库查询商品素材
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2024/5/7
     */
    public static function collectSearchList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/source/brand/collect/search-list')->post($params);
    }
}

