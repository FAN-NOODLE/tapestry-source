<?php


namespace App\Jobs\Export;

use App\Exceptions\Error;
use App\Jobs\Job;
use Box\Spout\Common\Type;
use Box\Spout\Writer\Common\Creator\WriterEntityFactory;
use Box\Spout\Writer\Common\Creator\WriterFactory;
use CrsSdk\Shop\DownloadCenter;
use Illuminate\Support\Facades\Log;


class ExportJob extends Job
{
    /**
     * 导出类型 1-excel
     */
    const EXPORT_TYPE_EXCEL = 1;

    /**
     * 下载中心任务状态 1-任务准备 2-等待下载 3-任务失败 4-已删除
     */
    const JOB_DOWN = 2;
    const JOB_FAIL = 3;

    /**
     * 导出参数
     *
     * @var
     */
    protected $params;

    /**
     * 导出类型后缀
     */
    const  EXTENSION = 'xlsx';

    /**
     * 文件名称（列表显示和导出名称）
     *
     * @var
     */
    protected $name;

    /**
     * 生成文件名称
     * @var
     */
    protected $fileName;

    /**
     * 生成的文件全路径--下载路径
     * @var
     */
    protected $filePath;

    /**
     * 文件写入对象
     *
     * @var \Box\Spout\Writer\XLSX\Writer
     */
    protected $writer;

    /**
     * 下载中心任务唯一标识
     *
     * @var int
     */
    protected $downId;


    /**
     * ExportJob constructor.
     * @param array $params 用于处理业务得参数
     * @param array $paramJson 前端原始请求参数-传数组-下载中心转json
     * @param string $name 任务名称或导出文件sheet名称
     *
     * @throws \Box\Spout\Common\Exception\UnsupportedTypeException
     */
    protected function __construct(array $params, array $paramJson = [], string $name = '统计数据导出')
    {
        $this->name   = $name;
        $this->params = $params;
        $this->downId = self::jobReady($this->params, $paramJson);
        $this->writer = WriterFactory::createFromType(Type::XLSX);
    }

    /**
     * 通知下载中心下载正在准备-返回id
     * @param $params
     * @param array $paramJson
     * @return int
     */
    protected function jobReady($params, $paramJson = [])
    {
        $downParams = [
            'enterprise_id' => $params['enterprise_id'] ?? 0,
            'operator_id'   => $params['operate_id'] ?? 0,
            'operator_name' => $params['operate_name'] ?? '',
            'type'          => self::EXPORT_TYPE_EXCEL,
            'name'          => $this->name . '.' . self::EXTENSION,
            'route_url'     => $params['route_url'] ?? '',
            'tag'           => $params['tag'] ?? '',
            'param_json'    => $paramJson
        ];

        try {
            $res = (new DownloadCenter())->add($downParams);
        } catch (\Exception $exception) {
            Log::info('创建下载任务失败' . $exception->getMessage(), $downParams);
            throw new Error(100000, '创建下载任务失败');
        }

        return $res['data']['down_id'] ?? 0;
    }


    /**
     * 更新下载中心状态为成功或失败
     * @param bool $bool true-成功 false-失败
     * @param string $reason
     * @return $this
     */
    protected function jobStatus($bool = true, $reason = '')
    {
        $downParams = [
            'enterprise_id' => $this->params['enterprise_id'],
            'down_id'       => $this->downId,
            'status'        => $bool ? self::JOB_DOWN : self::JOB_FAIL,
            'file_name'     => $this->filePath,
            'reason'        => $reason ?? ''
        ];

        try {
            (new DownloadCenter())->status($downParams);
        } catch (\Exception $exception) {
            Log::info('更新下载任务状态失败:' . $exception->getMessage(), $downParams);
            throw new Error(100000, '更新下载任务状态失败');
        }

        return $this;
    }

    /**
     * 设置临时下载的时间个内存
     *
     * @param int $seconds
     * @param string $memory
     *
     * @return $this
     */
    public function initSet($seconds = 0, $memory = '1000M')
    {
        set_time_limit($seconds);
        ini_set('memory_limit', $memory);

        return $this;
    }

    /**
     * 生成文件名称
     * @return $this
     */
    protected function setFileName()
    {
        $saltStr   = '';
        $randomStr = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';

        while (strlen($saltStr) < 8) {
            $saltStr .= $randomStr[rand(0, strlen($randomStr) - 1)];
        }

        $this->fileName = md5(time() . rand(1000, 9999) . $saltStr);
        return $this;
    }

    /**
     * 导出 - 写入表头
     * @param array $header
     * @param string $sheet
     * @return $this
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Writer\Exception\InvalidSheetNameException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    protected function excelExportInit($header = [], $sheet = '')
    {
        $this->initSet()->setFileName();

        $tmpPath = date('Y/m/d/', time());

        $exportPath = getDownloadPath();

        if (!file_exists($exportPath) && !mkdir($exportPath, 0777, true) && !is_dir($exportPath)) {
            throw new \RuntimeException(sprintf('Directory "%s" was not created', $exportPath));
        }

        $fileName = $this->fileName . '.' . self::EXTENSION;

        if (!file_exists($exportPath . $tmpPath) && !mkdir($concurrentDirectory = $exportPath . $tmpPath, 0777, true) && !is_dir($concurrentDirectory)) {
            throw new \RuntimeException(sprintf('Directory "%s" was not created', $concurrentDirectory));
        }

        $filePath       = $exportPath . $tmpPath . $fileName;
        $this->filePath = $tmpPath . $fileName;
        // $this->filePath = $filePath;
        $this->writer->openToFile($filePath);

        // 设置sheet名称
        $sheet = empty($sheet) ? $this->name : $sheet;
        $this->writer->getCurrentSheet()->setName($sheet);

        // 写入表头
        $row = WriterEntityFactory::createRowFromArray($header);
        $this->writer->addRow($row);
        return $this;
    }

    /**
     * 写入数据-调用
     * @param array $data 二维数组
     * @return $this
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    protected function exportExcel($data = [])
    {
        foreach ($data as $row) {
            $row = WriterEntityFactory::createRowFromArray($row);
            $this->writer->addRow($row);
        }

        return $this;
    }

    /**
     * 关闭打开文件-通知下载中心导出成功
     * @return $this
     */
    protected function exportClose()
    {
        $this->writer->close();
        return $this;
    }

}