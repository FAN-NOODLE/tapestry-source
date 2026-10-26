<?php
/**
 * 商品相关业务
 * Created by PhpStorm.
 * User: Leon Liu
 * Date: 2020/6/28
 * Time: 16:29
 */

namespace App\Services;

use Illuminate\Support\Facades\Redis;

class GoodsService extends BaseService
{

    /**
     * 允许上传的文件类型
     *
     * @var array
     */
    const ALLOW_FILE = ['xls', 'xlsx'];

    /**
     * 允许上传的文件大小
     *
     * @var array
     */
    const FILE_SIZE = 1024 * 1024 * 5;

    /**
     * 文件中文表头
     *
     * @var array
     */
    public $changePriceValueHand = ['SKUid', '售卖价'];

    /**
     * 文件英文表头
     *
     * @var array
     */
    public $changePriceKeyHand = ['sku_id', 'price'];

    /**
     *  客户上传的文件名
     */
    const FILENAME_MARK = '_MultiPriceChangeFileName';

    /**
     *  标记有改价任务在进行中的锁
     */
    const LOCK_MARK = '_MultiPriceChangeLock';

    /**
     *  标记目前总的任务数
     */
    const TOTAL_MARK = '_MultiPriceChangeTotal';

    /**
     *  标记目前总的任务数
     */
    const SUCCESS_NUM = '_MultiPriceChangeSuccessTotal';

    /**
     *  标记目前总的任务数
     */
    const FAIL_NUM = '_MultiPriceChangeFailTotal';

    /**
     *  失败任务的链表
     */
    const FAIL_LIST = '_FailTaskList';

    /**
     * excel文件上件
     * @param array $file 客户端的文件对象
     * @return array
     */
    public function uploadFile($file)
    {
        // 判断文件大小
        if ($file->getClientSize() > self::FILE_SIZE) {
            notice('320260');
        }

        // 生成随机的文件名
        $filename = date('Ymd', time()) . mt_rand(10000, 99999) . '.xls';

        // 校验文件扩展名
        if (!in_array($file->getClientOriginalExtension(), self::ALLOW_FILE)) {
            notice('451506');
        }

        // 创建路径
        if (!file_exists(resource_path('goods'))) {
            mkdir(resource_path('goods'));
        }

        // 生成文件的拷贝
        $path = $file->move(resource_path('change_price'), $filename);

        if (empty($path)) {
            notice('320261');
        }
        return $path;
    }

    /**
     * excel 导入 获取上传的数据
     *
     * @param array $filename 要读取的文件路径类
     * @param int $num 最大行数
     * @param array $changePriceValueHand 对应的表头
     * @return array $import_data 上传表格的数据
     * @throws \PHPExcel_Exception
     * @throws \PHPExcel_Reader_Exception
     */
    public function excelImport($filename, $num = 5000, $changePriceValueHand = [])
    {
        if (!file_exists($filename)) {
            notice(320258);
        }

        $objPHPExcel = \PHPExcel_IOFactory::load($filename);

        $sheet = $objPHPExcel->getSheet(0);

        // 总行数
        $highestRow = $sheet->getHighestRow();

        // 最后一列
        $highestColumn = $sheet->getHighestColumn();

        $char = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
            'U', 'V', 'W', 'X', 'Y', 'Z'];

        // 删除上传的文件
        unlink($filename);

        if ($highestRow >= $num) {
            notice(320262);
        }

        // 超过26列 最大列为两个以上字母
        if (strlen($highestColumn) > 1) {
            notice(320256);
        }

        // 默认表头对应的数据
        if (empty($changePriceValueHand)) {
            $changePriceValueHand = $this->changePriceValueHand;
        }

        // 验证表头
        foreach ($changePriceValueHand as $key => $val) {
            $excelName = trim($sheet->getCell($char[$key] . '1')->getValue());
            if ($excelName != $val) {
                notice(451500);
            }
        }

        // 获取上传文件的数据
        $import_data = [];
        for ($row = 1; $row <= $highestRow; $row++) {
            $data = [];
            for ($column = 0; $column < count($changePriceValueHand); $column++) {
                $val = $sheet->getCellByColumnAndRow($column, $row)->getValue();
                $val = trim($val);
                $data[] = $val;
            }
            $import_data[] = $data;
        }

        return $import_data;
    }

    /**
     * 获取任务进行中的锁
     *
     * @param int $enterpriseId 企业id
     * @return mixed
     */
    public function checkTask($enterpriseId)
    {
        return Redis::connection()->get($enterpriseId . self::LOCK_MARK);
    }

    /**
     * 设置任务进行中的锁
     *
     * @param int $enterpriseId 企业id
     * @return bool
     */
    public function setTaskLock($enterpriseId)
    {
        // 确认是否已经上锁
        $get = Redis::connection()->get($enterpriseId . self::LOCK_MARK);
        if (!empty($get)) {
            return false;
        }

        // 未上锁则进行上锁标记，并设置1小时过期
        $res = Redis::connection()->set($enterpriseId . self::LOCK_MARK, 1, 'ex', 3600);
        if (empty($res)) {
            return false;
        }
        return true;
    }

    /**
     * 改价任务解锁
     *
     * @param int $enterpriseId 企业id
     * @return mixed
     */
    public function unTaskLock($enterpriseId)
    {
        // 企业锁，任务完成时删除
        $key = $enterpriseId . self::LOCK_MARK;

        return Redis::connection()->del([$key]);
    }

    /**
     * 获取提交的改价文件名
     *
     * @param int $enterpriseId 企业id
     * @return mixed
     */
    public function getFileName($enterpriseId)
    {
        return Redis::connection()->get($enterpriseId . self::FILENAME_MARK);
    }

    /**
     * 设置redis中的改价文件名
     *
     * @param int $enterpriseId 企业id
     * @param $filename
     * @return mixed
     */
    public function setFileName($enterpriseId, $filename)
    {
        return Redis::connection()->set($enterpriseId . self::FILENAME_MARK, $filename);
    }

    /**
     *  设置总改价的任务数量
     *
     * @param int $enterpriseId 企业id
     * @param $taskNum
     * @return mixed
     */
    public function setTaskNum($enterpriseId, $taskNum)
    {
        // 设置redis中的总任务数
        return Redis::connection()->set($enterpriseId . self::TOTAL_MARK, $taskNum);
    }

    /**
     * 获取总的改价的任务数
     *
     * @param int $enterpriseId 企业id
     * @return mixed
     */
    public function getTotalTaskNum($enterpriseId)
    {
        // 获取redis中的总任务数
        return Redis::connection()->get($enterpriseId . self::TOTAL_MARK);
    }

    /**
     * 存取失败任务list中弹出一条数据
     *
     * @param int $enterpriseId 企业id
     * @return mixed
     */
    public function popFailList($enterpriseId)
    {
        return Redis::connection()->lpop($enterpriseId . self::FAIL_LIST);
    }

    /**
     * 清空之前的任务数据
     *
     * @param int $enterpriseId 企业id
     * @return mixed
     */
    public function delOld($enterpriseId)
    {
        // 删除失败任务列表数据
        return Redis::connection()->del([$enterpriseId . self::FAIL_LIST]);
    }

    /**
     * 设置改价成功任务数
     *
     * @param int $enterpriseId 企业id
     * @param $num
     * @return mixed
     */
    public function setSuccessNum($enterpriseId, $num)
    {
        return Redis::connection()->set($enterpriseId . self::SUCCESS_NUM, $num);
    }

    /**
     * 获取改价成功任务数
     *
     * @param int $enterpriseId 企业id
     * @return mixed
     */
    public function getSuccessNum($enterpriseId)
    {
        return Redis::connection()->get($enterpriseId . self::SUCCESS_NUM);
    }

    /**
     * 获取改价失败任务数
     *
     * @param int $enterpriseId 企业id
     * @return mixed
     */
    public function getFailNum($enterpriseId)
    {
        return Redis::connection()->get($enterpriseId . self::FAIL_NUM);
    }

    /**
     * 设置改价失败任务数
     *
     * @param int $enterpriseId 企业id
     * @param $num
     * @return mixed
     */
    public function setFailNum($enterpriseId, $num)
    {
        return Redis::connection()->set($enterpriseId . self::FAIL_NUM, $num);
    }

    /**
     * 插入一条改价失败任务数据
     *
     * @param int $enterpriseId 企业id
     * @param string $sku sku数据
     * @param int $price 修改价格
     * @param string $errorMsg 失败报错信息
     * @return int
     */
    public function pushFailList($enterpriseId, $sku, $price, $errorMsg)
    {
        return Redis::connection()->lpush($enterpriseId . self::FAIL_LIST, [$sku . '_' . $price . '_' . $errorMsg]);
    }

    /**
     * 获取redis值
     *
     * @param $key
     * @return mixed
     */
    public function getValueByKey($key)
    {
        return Redis::connection()->get($key);
    }

    /**
     * 重置相关redis数据
     *
     * @param int $enterpriseId 企业id
     * @return \Illuminate\Http\JsonResponse
     */
    public function reset($enterpriseId)
    {
        //企业锁，任务完成时删除
        $key_lock = $enterpriseId . self::LOCK_MARK;
        $key_total = $enterpriseId . self::TOTAL_MARK;
        $key_file = $enterpriseId . self::FILENAME_MARK;
        $success_num = $enterpriseId . self::SUCCESS_NUM;
        $fail_num = $enterpriseId . self::FAIL_NUM;
        $fail_list = $enterpriseId . self::FAIL_LIST;
        $keys_all = array_merge([$key_lock, $key_total, $key_file, $success_num, $fail_num, $fail_list]);

        if (empty($keys_all)) {
            return success();
        } else {
            Redis::connection()->del($keys_all);
        }
        return success($keys_all);
    }

    /**
     * 获取以企业id开关的所有redis数据
     *
     * @param int $enterpriseId 企业id
     * @return array
     */
    public function allKeys($enterpriseId)
    {
        return [
            $enterpriseId . self::SUCCESS_NUM,
            $enterpriseId . self::FAIL_NUM,
            $enterpriseId . self::LOCK_MARK,
            $enterpriseId . self::TOTAL_MARK,
            $enterpriseId . self::FILENAME_MARK];
    }

    /**
     * 获取相关redis数据
     *
     * @param int $enterpriseId 企业id
     * @return \Illuminate\Http\JsonResponse
     */
    public function showredis($enterpriseId)
    {
        $data = [];
        $keys_all = $this->allKeys($enterpriseId);
        if (empty($keys_all)) {
            return success();
        }
        foreach ($keys_all as $k) {
            $data[$k] = $this->getValueByKey($k);
        }
        return success($data);
    }

    /**
     * 失败任务数 +1
     *
     * @param int $enterpriseId 企业id
     */
    public function incrFail($enterpriseId)
    {
        Redis::connection()->incr($enterpriseId . self::FAIL_NUM);
    }

    /**
     * 成功任务数 +1
     *
     * @param int $enterpriseId 企业id
     */
    public function incrSucc($enterpriseId)
    {
        Redis::connection()->incr($enterpriseId . self::SUCCESS_NUM);
    }

}

