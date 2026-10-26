<?php

namespace App\Traits;

use Box\Spout\Common\Entity\Style\CellAlignment;
use Box\Spout\Reader\Common\Creator\ReaderEntityFactory;
use Box\Spout\Writer\Common\Creator\Style\StyleBuilder;
use Box\Spout\Writer\Common\Creator\WriterEntityFactory;
use DateTime;

trait ExcelTrait
{
    /**
     * 读Excel返回信息
     *
     * @param string $savePath 文件完整路径
     * @return array
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\UnsupportedTypeException
     * @throws \Box\Spout\Reader\Exception\ReaderNotOpenedException
     */
    public function excelRead($savePath)
    {
        $reader = ReaderEntityFactory::createReaderFromFile($savePath);
        $reader->open($savePath);

        // 仅读第一个sheet
        $iterator = $reader->getSheetIterator();
        $iterator->rewind();
        $sheet1 = $iterator->current();

        $data = [];
        foreach ($sheet1->getRowIterator() as $rowIndex => $row) {
            $rowData = [];
            foreach ($row->getCells() as $cell) {
                if (is_object($cell->getValue())) {
                    if ($cell->getValue() instanceof DateTime) {
                        $rowData[] = $cell->getValue()->format('Y-m-d H:i:s');
                    }
                } else {
                    $rowData[] = (string)$cell->getValue();
                }
            }
            $data[] = $rowData;
        }
        $reader->close();

        // 处理wps部分列无值读出来的每行字段数与标题不符
        if (!empty($data)) {
            $titleCount = count($data[0]);
            for ($i = 0; $i < $titleCount; $i++) {
                foreach ($data as $key => $value) {
                    if (!isset($value[$i])) {
                        $data[$key][$i] = '';
                    }
                }
            }
        }

        return $data;
    }

    /**
     * 导出Excel
     *
     * @param array $title 表头
     * @param array $data 数据
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    public function excelExport(array $title = [], array $data = [], $fileName = '')
    {

        list($writer,$style) = $this->setExcelTitle($title, $fileName);

        // 内容
        $this->setExceldata($writer,$style,$data);
        $this->closeExcel($writer);
    }


    /**
     * 导出Excel
     *
     * @param array $title 表头
     * @param string $fileName
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    public function setExcelTitle(array $title = [], string $fileName = ''): array
    {
        $style = (new StyleBuilder())->setFontSize(12)->setFontName('微软雅黑')->setCellAlignment(CellAlignment::CENTER)->build();

        $writer = WriterEntityFactory::createXLSXWriter();
        // 保存到本地
        if ($fileName) {
            // 浏览器输出
            $writer->openToBrowser($fileName . '-' . date('YmdHis').'.xlsx');
        } else {
            // 浏览器输出
            $writer->openToBrowser('export-'.date('YmdHis').'.xlsx');
        }

        // 表头
        $row = WriterEntityFactory::createRowFromArray($title, $style);
        $writer->addRow($row);
        return [$writer,$style];
    }

    /**
     * 导出Excel 多次写入
     *
     * @param $writer
     * @param $style
     * @param array $data 数据
     */
    public function setExcelData($writer, $style, array $data = []){
        // 内容
        $rows = [];
        foreach ($data as $value) {
            $rows[] = WriterEntityFactory::createRowFromArray($value, $style);
        }
        $writer->addRows($rows);
    }

    /**
     * @param $writer
     */
    public function closeExcel($writer)
    {
        $writer->close();
    }
}
