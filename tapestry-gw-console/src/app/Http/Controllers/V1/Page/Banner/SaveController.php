<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/23
 * Time: 18:07
 */

namespace App\Http\Controllers\V1\Page\Banner;

use App\Exceptions\Error;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\PageBannerRequest;
use CrsSdk\Shop\PageBanner;

class SaveController extends BaseController
{
    /**
     * 轮播图验证器
     *
     * @var PageBannerRequest
     */
    protected $pageBannerRequest;

    /**
     * SaveController constructor.
     * @param PageBannerRequest $pageBannerRequest
     */
    public function __construct(PageBannerRequest $pageBannerRequest)
    {
        parent::__construct();

        // 轮播图验证器赋值
        $this->pageBannerRequest = $pageBannerRequest;
    }

    /**
     * 轮播图保存
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/25 11:16
     * @param Request $request
     * @return mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['data']);

        // 请求数据验证和重组
        $haveDatas = $editData = $addData = $deleteIds = [];
        if (isset($params['data']) && !empty($params['data'])) {
            foreach ($params['data'] as $value) {
                $validator = $this->pageBannerRequest->common($value);
                if ($validator->fails()) {
                    throw new Error(310000, $validator->errors()->first());
                }

                //数据特殊处理
                if (!empty($value['redirect_value'])) {
                    $value['redirect_value'] = json_encode($value['redirect_value']);
                }else{
                    $value['redirect_value']='';
                }

                // 有ID的和要新添加的数据区分
                if ($value['id'] > 0) {
                    $haveDatas[] = $value;
                } else {
                    $addData[] = $value;
                }
            }
        }

        // 区分编辑数据和删除的数据
        $bannerlist = $this->getList();
        if (!empty($haveDatas) && !empty($bannerlist['data'])) {
            $haveIds = array_column($haveDatas, 'id');
            $haveDatas = array_column($haveDatas, NULL, 'id');
            foreach ($bannerlist['data'] as $banner) {
                if (in_array($banner['id'], $haveIds)) {
                    $editData[] = $haveDatas[$banner['id']];
                } else {
                    $deleteIds[] = $banner['id'];
                }
            }
        } else if (!empty($bannerlist['data'])) {
            $deleteIds = array_column($bannerlist['data'], 'id');
        } else {
            $editData = $haveDatas;
        }

        //数据删除
        !empty($deleteIds) && $this->deleteData($deleteIds);
        //数据编辑
        !empty($editData) && $this->editData($editData);
        //数据添加
        !empty($addData) && $this->addData($addData);

        $this->writeOperatorLog('编辑轮播图');
        return success();
    }

    /**
     * 获取轮播图列表
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/28 11:16
     * @return mixed
     */
    public function getList()
    {
        try {
            return (new PageBanner())->lists(['enterprise_id' => $this->enterpriseId]);
        } catch (\Exception $e) {
            return error(700001);
        }
    }

    /**
     * 删除轮播图
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/28 11:16
     * @param array $deleteIds 要删除的数据
     * @return mixed
     */
    public function deleteData($deleteIds)
    {
        try {
            return (new PageBanner())->delete(['banner_id' => $deleteIds]);
        } catch (\Exception $e) {
            return error(700002);
        }
    }

    /**
     * 编辑轮播图
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/28 11:16
     * @param array $editData 要编辑的数据
     * @return bool
     */
    public function editData($editData)
    {
        try {
            foreach ($editData as $edit) {
                $edit['banner_id'] = $edit['id'];
                unset($edit['id']);
                (new PageBanner())->save($edit);
            }
            return true;
        } catch (\Exception $e) {
            return error(700003);
        }
    }

    /**
     * 添加轮播图
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/28 11:16
     * @param array $addData 要添加的数据
     * @return bool
     */
    public function addData($addData)
    {
        try {
            foreach ($addData as $add) {
                unset($add['id']);
                $add['enterprise_id'] = $this->enterpriseId;
                (new PageBanner())->add($add);
            }
            return true;
        } catch (\Exception $e) {
            return error(700004);
        }
    }
}