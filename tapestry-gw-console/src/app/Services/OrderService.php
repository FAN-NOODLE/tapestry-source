<?php
/**
 * Created by PhpStorm.
 * User: cjp
 * Date: 2019/6/4
 * Time: 11:40
 */
namespace App\Services;

class OrderService extends BaseService
{
    /**
     * 允许上传的文件类型
     *
     * @var array
     */
    const ALLOW_FILE = ['xls', 'xlsx'];

    /**
     * 允许上传的文件类型
     *
     * @var array
     */
    const FILE_SIZE = 1024 * 1024 * 5;

    /**
     * 文件中文表头
     *
     * @var array
     */
    public $valueHand = ['小程序订单号', '物流公司标识', '物流单号'];

    /**
     * 文件英文表头
     *
     * @var array
     */
    public $keyHand = ['order_no', 'logistic_identity', 'logistic_no'];

    /**
     * excel文件上件
     * @param array $file 客户端对象的集合
     * @return array
     */
    public function uploadFile($file)
    {

        if ($file->getClientSize() > self::FILE_SIZE) notice(320253);

        // 生成随机的文件名
        $filename = date('Ymd', time()) . mt_rand(10000, 99999) . '.xls';

        // 校验文件扩展名
        if (!in_array($file->getClientOriginalExtension(), self::ALLOW_FILE)) notice(320252);

        // 移动文件
        if (!file_exists(resource_path('order'))) mkdir(resource_path('order'));
        $path = $file->move(resource_path('order'), $filename);

        if (empty($path)) notice(320254);

        return $path;

    }

    /**
     * excel导入
     *
     * @param $filename
     * @param int $num
     * @param array $valueHand
     * @return array|void
     * @throws \PHPExcel_Exception
     * @throws \PHPExcel_Reader_Exception
     */
    public function excelImport($filename, $num = 5000,$valueHand=[])
    {
        if (!file_exists($filename)) notice(320258);

        $objPHPExcel   = \PHPExcel_IOFactory::load($filename);
        $sheet         = $objPHPExcel->getSheet(0);
        $highestRow    = $sheet->getHighestRow(); // 总行数
        $highestColumn = $sheet->getHighestColumn(); // 最后一列
//        $highestColumn_num = \PHPExcel_Cell::columnIndexFromString($highestColumn);	//列数

        $char = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'];

        // 删除上传的文件
        unlink($filename);

        if ($highestRow >= $num) notice(320262);

        //超过26列 最大列为两个以上字母
        if (strlen($highestColumn) > 1) notice(320256);

        if (empty($valueHand)) {
            $valueHand = $this->valueHand;
        }

        //验证表头
        foreach ($valueHand as $key => $val) {
            $excelName = trim($sheet->getCell($char[$key] . '1')->getValue());
            if ($excelName != $val) notice(320257);
        }

        $import_data = [];
        for ($row = 1; $row <= $highestRow; $row++) {
            $data = [];
            for ($column = 0; $column < count($valueHand); $column++) {
                $val = $sheet->getCellByColumnAndRow($column, $row)->getValue();
                if($column==2){
                    $val = $this->numToStr($val);
                }
                $val = trim($val);
                $data[] = $val;
            }

            $import_data[] = $data;
        }

        return $import_data;
    }

    public function numToStr($num)
    {
        $result = "";
        if (stripos($num, 'e') === false) {
            return $num;
        }
        while ($num > 0) {
            $v = $num - floor($num / 10) * 10;
            $num = floor($num / 10);
            $result = $v . $result;
        }
        //处理物流单号默认带E的单号
        if(empty($result)){
            $result = $num;
        }
        return $result;
    }

	/**
	 * 数据罗盘数据概览格式转换
	 *
	 * @param $data
	 * @param string $key1
	 * @param string $key2
	 * @param string $key3
	 * @param string $key4
     * @param string $key5
     * @param string $key6
	 *
	 * @return array
	 * @author  baizhongyang <baizhongyang@vchangyi.com>
	 * @date    2020/8/21 20:53
	 */
	public function transformDataFormat($data, $key1 = 'extension_num', $key2 = 'extension_amount', $key3 = 'service_num', $key4 = 'service_amount', $key5 = 'staff_order_num', $key6 = 'partner_order_num')
	{
		$x = array_keys($data['detail_data']);
		$y = [
			'data1' => array_column($data['detail_data'], $key1) ?? [],
			'data2' => array_column($data['detail_data'], $key2) ?? [],
			'data3' => array_column($data['detail_data'], $key3) ?? [],
			'data4' => array_column($data['detail_data'], $key4) ?? [],
		];

		$result = [
			'data1' => $data[$key1],
			'data2' => $data[$key2],
			'data3' => $data[$key3],
			'data4' => $data[$key4],
			'data5' => $data[$key5],
			'data6' => $data[$key6],
			'detail_data' => [
				'x' => $x,
				'y' => $y,
			],
		];

		return $result;
	}
}