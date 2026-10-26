<?php
/**
 * 核对并批量修改商品价格
 * Created by PhpStorm.
 * User: Leon Liu
 * Date: 2020/6/28
 * Time: 16:01
 */

namespace App\Http\Controllers\V1\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\MultiChangePriceJob;
use App\Services\GoodsService;
use Illuminate\Http\Request;

class MultiChangePriceController extends BaseController
{
    /**
     *  改价列表模板的表头
     */
    const EXCEL_HEADER = [
        'SKUid',
        '售卖价',
        '上传失败原因'
    ];

    /**
     * 改价任务发布成功
     */
    const UPLOAD_SUCC = 1;

    /**
     * 改价操作已经进行中，无法添加新的操作任务
     */
    const UPLOAD_OCCUPIED = 3;

    /**
     * 改价操作失败
     */
    const UPLOAD_FAIL = 2;

    /**
     * 可以上传改价任务
     */
    const UPLOAD_YES = 1;

    /**
     * 不能上传改价任务
     */
    const UPLOAD_NO = 2;

    /**
     * 商品相关Service
     *
     * @var GoodsService
     */
    protected $goodsService;

    public function __construct(GoodsService $goodsService)
    {
        parent::__construct();
        $this->goodsService = $goodsService;
    }

    /**
     * 商品价格批量修改
     *
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @throws \PHPExcel_Exception
     * @throws \PHPExcel_Reader_Exception
     */
    public function index(Request $request)
    {
        //先确认是否有任务在进行中
        $res = $this->goodsService->checkTask($this->enterpriseId);
        if ($res) {
            return success([
                'status' => self::UPLOAD_OCCUPIED
            ]);
        }

        // 验证上传文件有效性
        $input = $request->only(['excel']);
        $rules = ['excel' => 'required|file'];
        $message = ['excel.*' => '请上传excel文件'];
        $this->validate($request, $rules, $message);

        // 获取服务器copy上传文件的路径
        $fileName = $this->goodsService->uploadFile($input['excel']);

        // 获取文件内容
        $fileInfo = $this->goodsService->excelImport($fileName);

        // 获取原始文件名
        $fileName = $input['excel']->getClientOriginalName();

        // 数据验证，金钱单位转换为分
        unset($fileInfo[0]);
        $priceList = [];
        foreach ($fileInfo as $info) {
            $sku = $info[0];
            $price = $info[1] * 100;
            $priceList[$sku] = $price;
        }
        if (empty($priceList)) {
            return success(['status' => self::UPLOAD_FAIL]);
        }

        // 初始化改价任务任务
        $this->initChangePriceWork($priceList, $fileName);

        // 分发任务
        try {
            $hash = $this->enterprise['hash'];
            dispatch(
                new MultiChangePriceJob(
                    $priceList,
                    $this->enterpriseId,
                    $this->enterpriseAdminId,
                    $hash,
                    $fileName
                )
            );
        } catch (\Exception $e) {
            // 如果发布任务失败，则解锁以便用户再次上传改价文件
            $this->goodsService->unTaskLock($this->enterpriseId);
            throw new Error($e->getCode(), $e->getMessage());
        }

        $res = [
            'status' => self::UPLOAD_SUCC
        ];
        return success($res);
    }

    /**
     * 初始化改价任务的一系列前置操作
     *
     * @param $priceList
     * @param $fileName
     * @return bool
     */
    public function initChangePriceWork($priceList, $fileName)
    {
        // 改价业务上锁
        $res = $this->goodsService->setTaskLock($this->enterpriseId);
        if (!$res) {
            return success([
                'status' => self::UPLOAD_FAIL
            ]);
        }

        // 设置redis中的总任务数
        $this->goodsService->setTaskNum($this->enterpriseId, count($priceList));

        // 更新文件名
        $this->goodsService->setfileName($this->enterpriseId, $fileName);

        // 删除之前进行的任务,重启新任务
        $this->goodsService->delOld($this->enterpriseId);

        // 初始化成功和失败数
        $this->goodsService->setSuccessNum($this->enterpriseId, 0);
        $this->goodsService->setFailNum($this->enterpriseId, 0);

        return true;
    }

    /**
     * 商品批量改价进度接口
     *
     * @return mixed
     */
    public function process()
    {
        // 获取企业所有进行中的改价任务id
        $successNum = $this->goodsService->getSuccessNum($this->enterpriseId);
        $failNum = $this->goodsService->getFailNum($this->enterpriseId);
        $totalNum = $this->goodsService->getTotalTaskNum($this->enterpriseId);
        $fileName = $this->goodsService->getFileName($this->enterpriseId);

        // 没有进行中的改价任务时
        if ($totalNum <= 0) {
            $data = [
                'total' => 0,
                'success' => 0,
                'fail' => 0,
                'percent' => 0,
                'file_name' => $fileName
            ];
            return success($data);
        }

        // 计算完成百分比
        $finPercent = intval(100 * ($successNum + $failNum) / $totalNum);

        // 返回数据给客户端
        $data = [
            'total' => $totalNum,
            'success' => $successNum,
            'fail' => $failNum,
            'percent' => $finPercent,
            'file_name' => $fileName
        ];
        return success($data);
    }

    /**
     * 商品批量改价失败数据列表下载
     *
     * @return mixed
     */
    public function failList()
    {
        // 获取所有的改价失败数据，原数据存在redis list中
        $arrr = [];
        while (true) {
            $arr = $this->goodsService->popFailList($this->enterpriseId);
            if (empty($arr)) {
                break;
            }
            $arrr[] = $arr;
        }

        // 数据处理为可以导出的形式
        $res = [];
        foreach ($arrr as $list) {
            $arr = explode('_', $list);
            $arr[1] = $arr[1] / 100;
            // 保证下载失败列表顺序与上传的列表相同
            array_unshift($res, $arr);
        }

        // 再次放入redis，以便下次再次下载失败列表
        foreach ($res as $re) {
            $this->goodsService->pushFailList($this->enterpriseId, $re[0], $re[1] * 100, $re[2]);
        }

        // 返回数据
        $data['title'] = self::EXCEL_HEADER;
        $data['list'] = $res;

        return success($data);
    }

    /**
     * 获取改价模板
     *
     * @return mixed
     */
    public function getTemplate()
    {
        $data['title'] = self::EXCEL_HEADER;
        $data['list'] = [['', '']];

        return success($data);
    }

    /**
     * 检测是否可以进行批量改价操作
     *
     * @return mixed
     */
    public function operCheck()
    {
        $res = $this->goodsService->checkTask($this->enterpriseId);

        // 检查目前是否有进行中任务，
        if ($res) {
            return success(['status' => self::UPLOAD_NO]);
        } else {
            return success(['status' => self::UPLOAD_YES]);
        }
    }

    /////////////////////////////////  调试用函数  /////////////////////////////////

    /**
     * 重置相关redis数据
     *
     * @return mixed
     */
    public function reset()
    {
        return $this->goodsService->reset($this->enterpriseId);
    }

    /**
     * 获取相关redis数据
     *
     * @return mixed
     */
    public function showredis()
    {
        return $this->goodsService->showredis($this->enterpriseId);
    }


}
