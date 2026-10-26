<?php


namespace App\Jobs\Export;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Jobs\Job;
use Box\Spout\Writer\Common\Creator\WriterEntityFactory;
use CrsSdk\Shop\DownloadCenter;
use Illuminate\Support\Facades\Log;

class NotSaleListExportJob extends Job
{
    /**
     * 导出类型
     */
    const EXPORT_TYPE_EXCEL = 1;

    /**
     * 导出路径
     */
    const EXPORT_PATH = 'download/';

    /**
     * 下载中心任务状态 1-任务准备 2-等待下载 3-任务失败 4-已删除
     */
    const JOB_DOWN = 2;
    const JOB_FAIL = 3;

    /**
     * 导出参数
     */
    protected $params;

    /**
     * 导出类型后缀
     */
    const  EXTENSION = 'xlsx';

    /**
     * 生成文件名称
     */
    protected $fileName;

    /**
     * 生成的文件全路径--下载路径
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
     * 文件表头
     *
     * @var array
     */
    protected $header;

    /**
     * 要导出的待上架商品列表
     *
     * @var array
     */
    protected $datas;

    /**
     * NotSaleListExportJob constructor.相关数据初始化
     * @param array $datas
     * @param array $header
     * @param $params
     * @param string $fileName
     */
    public function __construct(array $datas, array $header, $params, string $fileName = '待上架商品信息')
    {
        // 初始化数据：队列名
        $this->queue = Constant::NOT_SALE_LIST_DOWNLOAD;

        // 初始化数据：待上架商品数据
        $this->datas = $datas;

        // 初始化数据：要下载的文件名
        $this->fileName = $fileName;

        // 表头
        $this->header = $header;

        // 微服务相关参数
        $this->params = $params;

        // 拿到下载中心给的id
        $this->downId = $this->downId = $this->jobReady();

        // 创建写excel的工具类
        $this->writer = WriterEntityFactory::createXLSXWriter();
    }

    /**
     * 通知下载中心下载正在准备-返回id
     *
     * @return int
     */
    protected function jobReady()
    {
        try {
            $res = (new DownloadCenter())->add($this->params);
        } catch (\Exception $exception) {
            Log::info('创建下载任务失败' . $exception->getMessage(), $this->params);
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
            'down_id' => $this->downId,
            'status' => $bool ? self::JOB_DOWN : self::JOB_FAIL,
            'file_name' => $this->filePath,
            'reason' => $reason ?? ''
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
     * 导出Excel列表
     *
     * @return bool
     */
    public function handle()
    {
        // 生成excel
        try {
            $this->excelExportInit($this->header)->listExport($this->datas);
        } catch (\Exception $exception) {
            Log::info('导出失败' . $exception->getMessage(), $this->datas);
            // 修改状态为失败
            $this->jobStatus(false, $exception->getMessage());
            throw new Error(100000, $exception->getMessage());
        }

        // 生成excel 修改状态为成功
        $this->jobStatus();
        return true;
    }

    /**
     * 导出 - 写入表头
     *
     * @param array $header
     * @param string $sheet
     * @return $this
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Writer\Exception\InvalidSheetNameException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    protected function excelExportInit($header = [], $sheet = 'sheet1')
    {
        // 设置文件名
        $this->initSet()->setFileName();
        $tmpPath = date('Y/m/d', time());

        // 生成文件
        $exportPath = storage_path(self::EXPORT_PATH);
        if (!file_exists($exportPath)) {
            mkdir($exportPath, 0777, true);
        }
        $tmpPath .= '/';
        $fileName = $this->fileName . '.' . self::EXTENSION;
        if (!file_exists($exportPath . $tmpPath)) {
            mkdir($exportPath . $tmpPath, 0777, true);
        }
        $filePath = $exportPath . $tmpPath . $fileName;
        $this->filePath = $tmpPath . $fileName;
        $this->writer->openToFile($filePath);


        // 设置sheet名称
        $this->writer->getCurrentSheet()->setName($sheet);

        // 写入表头
        $row = WriterEntityFactory::createRowFromArray($header);

        // 这里的row就是表头数据
        $this->writer->addRow($row);

        return $this;
    }

    /**
     * 插入每行数据
     *
     * @param $datas
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    private function listExport($datas)
    {
        // 插入excel
        $this->exportExcel($datas);
        $this->exportClose();
    }

    /**
     * 写入数据-调用
     *
     * @param array $data
     * @return $this
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    private function exportExcel($data = [])
    {
        foreach ($data as $row) {
            $row = WriterEntityFactory::createRowFromArray($row);
            $this->writer->addRow($row);
        }

        return $this;
    }

    /**
     * 关闭打开文件-稍后通知下载中心导出成功
     * @return $this
     */
    private function exportClose()
    {
        $this->writer->close();
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
     *
     * @return string
     */
    protected function setFileName()
    {
        return $this->fileName;
    }
}
