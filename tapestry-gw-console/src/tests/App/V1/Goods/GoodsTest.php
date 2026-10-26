<?php
/**
 * Created by PhpStorm.
 * User: zq
 * Date: 19-3-29
 * Time: 上午9:30
 */

namespace Test\App\V1\Goods;

class GoodsTest extends BaseTest
{
    public function testSave ()
    {
        $this->login();
        $re = $this->randNameAndSpu();
        $data = [
            "spu_name" => $re['name'],
            "spu_picture" => [
                "https://oss-crs-dev.vchangyi.com/2019/03/29/51a8167eab92ccfdf7ed95ca11f345ef.jpg"
            ],
            "category_id" => 65,
            "content" => ["https://oss-crs-dev.vchangyi.com/2019/03/29/53bb8f54e8ec329b39e4f6e1c17117f7.jpg"],
            "attribute_ids" => [35, 36],
            "sku_data" => [
                [
                    "attribute_item_ids" => [302, 305],
                    "sku_picture" => "https://oss-crs-dev.vchangyi.com/2019/03/29/7da1166866aa2ebaa3fb9170fa830d6b.jpg",
                    "sell_price" => 10000,
                    "is_enabled" => 1,
                    "stock" => "100",
                    "market_price" => 10000
                ],
                [
                    "attribute_item_ids" => [304, 307],
                    "sku_picture" => "https://oss-crs-dev.vchangyi.com/2019/03/29/d7f71e038bc80ab12ac213182ff0be9c.jpg",
                    "sell_price" => 10000,
                    "is_enabled" => 1,
                    "stock" => "100",
                    "market_price" => 10000
                ]
            ],
            "spu_sn" => $re['spu_sn']
        ];
        $response = $this->call('POST', BaseTest::BASE_URL.'/goods/save', $data, $this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    public function testUpdate ()
    {
        $this->login();
        $re = $this->randNameAndSpu();
        // 只进行数据编辑，不添加新的sku信息
        $data = [
            "spu_id" => 207,
            "spu_sn" => $re['spu_sn'],
            "spu_name" => $re['name'],
            "spu_picture" => [
                "https://oss-crs-dev.vchangyi.com/2019/03/29/51a8167eab92ccfdf7ed95ca11f345ef.jpg"
            ],
            "category_id" => 65,
            "attribute_ids" => [36, 35],
            "content" => [
                "https://oss-crs-dev.vchangyi.com/2019/03/29/53bb8f54e8ec329b39e4f6e1c17117f7.jpg"
            ],
            "sku_data" => [
                [
                    "attribute_item_ids" => ["35:302", "36:305"],
                    "sku_picture" => "https://oss-crs-dev.vchangyi.com/2019/03/29/7da1166866aa2ebaa3fb9170fa830d6b.jpg",
                    "sell_price" => 10000,
                    "is_enabled" => 1,
                    "stock" => 100,
                    "market_price" => 10000,
                    "sku_id" => 361
                ],
                [
                    "attribute_item_ids" => ["35:304", "36:307"],
                    "sku_picture" => "https://oss-crs-dev.vchangyi.com/2019/03/29/d7f71e038bc80ab12ac213182ff0be9c.jpg",
                    "sell_price" => 10000,
                    "is_enabled" => 1,
                    "stock" => 100,
                    "market_price" => 10000,
                    "sku_id" => 362
                ]
            ]
        ];

        // 只进行数据编辑，添加新的sku信息
        $dataTwo = [
            "spu_id" => 207,
            "spu_sn" => $re['spu_sn'],
            "spu_name" => $re['name'],
            "spu_picture" => ["https://oss-crs-dev.vchangyi.com/2019/03/29/51a8167eab92ccfdf7ed95ca11f345ef.jpg"],
            "category_id" => 65,
            "attribute_ids" => [36, 35],
            "content" => ["https://oss-crs-dev.vchangyi.com/2019/03/29/53bb8f54e8ec329b39e4f6e1c17117f7.jpg"],
            "sku_data" => [
                [
                    "attribute_item_ids" => ["35:302", "36:305"],
                    "sku_picture" => "https://oss-crs-dev.vchangyi.com/2019/03/29/7da1166866aa2ebaa3fb9170fa830d6b.jpg",
                    "sell_price" => 10000,
                    "is_enabled" => 1,
                    "stock" => 100,
                    "market_price" => 10000,
                    "sku_id" => 361
                ],
                [
                    "attribute_item_ids" => ["35:304", "36:307"],
                    "sku_picture" => "https://oss-crs-dev.vchangyi.com/2019/03/29/d7f71e038bc80ab12ac213182ff0be9c.jpg",
                    "sell_price" => 10000,
                    "is_enabled" => 1,
                    "stock" => 100,
                    "market_price" => 10000,
                    "sku_id" => 362
                ],
                [
                    "attribute_item_ids" => ["35:303", "36:306"],
                    "sku_picture" => "https://oss-crs-dev.vchangyi.com/2019/03/29/2aa37bee27052f76538154077c61b08d.jpg",
                    "sell_price" => 30000,
                    "is_enabled" => 1,
                    "stock" => 300,
                    "market_price" => 30000,
                    "sku_id" => 0
                ]
            ]
        ];

        $response = $this->call('POST', BaseTest::BASE_URL.'/goods/update', $data, $this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    public function testLists ()
    {
        $this->login();
        $data = [
            'page' => 1,
            'size' => 10,
            'search' => [
                'spu_name' => '',
                'spu_sn' => '',
                'category_id' => '',
                'category_parent_id' => '',
                'is_sale' => [],
                'exc_at' => ''
            ]
        ];
        $response = $this->call('POST', BaseTest::BASE_URL.'/goods/update', $data, $this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    private function randNameAndSpu ()
    {
        $name = [
            '中式服装',
            '西式服装',
            '刺绣服装',
            '呢绒服装',
            '丝绸服装',
            '棉布服装',
            '毛皮服装',
            '针织服装',
            '羽绒服装',
            '整件装',
            '套装',
            '外套',
            '背心',
            '裙',
            '裤',
            '社交服',
            '日常服',
            '职业服',
            '运动服',
            '室内服',
            '舞台服'
        ];
        return ['name' => array_rand($name, 1) . rand(1, 9999), 'spu_sn' => md5(time())];
    }
}