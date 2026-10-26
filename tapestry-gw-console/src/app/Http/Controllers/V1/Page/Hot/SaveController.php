<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/23
 * Time: 18:07
 */

namespace App\Http\Controllers\V1\Page\Hot;

use App\Exceptions\Error;
use CrsSdk\Shop\PageActivityBanner;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\PageHotRequest;
use CrsSdk\Shop\PageHot;
use Illuminate\Support\Facades\Redis;

class SaveController extends BaseController
{

    const TYPE_HOT = 0;
    const TYPE_RECOMMEND = 1;
    /**
     * 爆品推荐验证器
     *
     * @var PageHotRequest
     */
    protected $pageHotRequest;

    /**
     * SaveController constructor.
     * @param PageHotRequest $pageHotRequest
     */
    public function __construct(PageHotRequest $pageHotRequest)
    {
        parent::__construct();

        // 爆品推荐验证器赋值
        $this->pageHotRequest = $pageHotRequest;
    }

    /**
     * 爆品推荐保存
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/29 11:16
     * @param Request $request
     * @return mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['data','banner']);

        // 请求数据验证和重组
        $haveDatas = $editData = $addData = $deleteIds = [];
        if (isset($params['data']) && !empty($params['data'])) {
            foreach ($params['data'] as &$value) {
                $validator = $this->pageHotRequest->common($value);
                if ($validator->fails()) {
                    throw new Error(310000, $validator->errors()->first());
                }

                //数据特殊处理
                if (!empty($value['redirect_value'])) {
                    $value['redirect_value'] = json_encode($value['redirect_value']);
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
        $hotlist = $this->getList();
        if (!empty($haveDatas) && !empty($hotlist['data'])) {
            $haveIds = array_column($haveDatas, 'id');
            $haveDatas = array_column($haveDatas, NULL, 'id');
            foreach ($hotlist['data'] as $hot) {
                if (in_array($hot['id'], $haveIds)) {
                    $editData[] = $haveDatas[$hot['id']];
                } else {
                    $deleteIds[] = $hot['id'];
                }
            }
        } else if (!empty($hotlist['data'])) {
            $deleteIds = array_column($hotlist['data'], 'id');
        } else {
            $editData = $haveDatas;
        }

        //数据删除
        !empty($deleteIds) && $this->deleteData($deleteIds);
        //数据编辑
        !empty($editData) && $this->editData($editData);
        //数据添加
        !empty($addData) && $this->addData($addData);

        //活动banner图
        if(isset($params['banner']) && !empty($params['banner'])){
            if(!empty($params['banner']['id'])){
                $this->saveBanner($params['banner']);
            }else{
                $this->addBanner($params['banner']);
            }
        }

        $this->writeOperatorLog('编辑爆品推荐');

        return success(['data' => []]);
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
            return (new PageHot())->lists(['enterprise_id' => $this->enterpriseId]);
        } catch (\Exception $e) {
            return error(700101);
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
            return (new PageHot())->delete(['enterprise_id' => $this->enterpriseId, 'hot_id' => $deleteIds]);
        } catch (\Exception $e) {
            return error(700102);
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
                $edit['hot_id'] = $edit['id'];
                $edit['enterprise_id'] = $this->enterpriseId;
                unset($edit['id']);
                (new PageHot())->save($edit);
            }
            return true;
        } catch (\Exception $e) {
            return error(700103);
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
                (new PageHot())->add($add);
            }
            return true;
        } catch (\Exception $e) {
            return error(700104);
        }
    }


    /**
     * 编辑活动banner
     * @author guokai  guokai@vchangyi.com
     * @date 2019/04/18 10:45
     * @param array $saveBanner 要编辑的数据
     * @return bool
     */
    public function saveBanner($saveBanner)
    {
        try {

            $saveBanner['type'] = self::TYPE_HOT;
            $saveBanner['enterprise_id'] = $this->enterpriseId;
            (new PageActivityBanner())->save($saveBanner);
            return true;

        } catch (\Exception $e) {
            return error(700205);
        }
    }


    /**
     * 添加活动banner
     * @author guokai guokai@vchangyi.com
     * @date 2019/04/18 11:20
     * @param array $addBanner 要添加的数据
     * @return bool
     */
    public function addBanner($addBanner)
    {
        try {

            $addBanner['type'] = self::TYPE_HOT;
            $addBanner['enterprise_id'] = $this->enterpriseId;
            unset($addBanner['id']);
            (new PageActivityBanner())->add($addBanner);

            return true;
        } catch (\Exception $e) {
            return error(700206);
        }
    }
}