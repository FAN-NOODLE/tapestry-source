<?php

require_once __DIR__ . '/../../app/Sdk/CDP/Coach/RecommendationFeed.php';

use App\Sdk\CDP\Coach\RecommendationFeed;

function assertSameValue($expected, $actual, $message)
{
    if ($expected !== $actual) {
        fwrite(STDERR, $message . PHP_EOL);
        fwrite(STDERR, 'Expected: ' . var_export($expected, true) . PHP_EOL);
        fwrite(STDERR, 'Actual: ' . var_export($actual, true) . PHP_EOL);
        exit(1);
    }
}

$payload = RecommendationFeed::buildPayload([
    'uid' => 'union-001',
    'size' => 9,
    'features' => [
        'sa_id' => 'guide-001',
        'store_id' => 'store-001',
    ],
]);

assertSameValue('union-001', $payload['uid'], 'uid 应使用客户 Union ID');
assertSameValue(9, $payload['size'], '推荐数量应为 9');
assertSameValue('sa_rec', $payload['scene_id'], '推荐场景应默认为 sa_rec');
assertSameValue('guide-001', $payload['features']['sa_id'], 'features.sa_id 应使用 members 表的 User ID');
assertSameValue('store-001', $payload['features']['store_id'], 'features.store_id 应传递门店编码');

$options = RecommendationFeed::buildOptions($payload, 'fixed-authorization');
assertSameValue('fixed-authorization', $options['headers']['Authorization'], '应使用固定 Authorization 请求头');
assertSameValue('application/json', $options['headers']['Content-Type'], '应以 JSON 格式请求');
assertSameValue($payload, json_decode($options['body'], true), '请求正文应为推荐参数');

$items = RecommendationFeed::normalizeResponse([
    'code' => 299,
    'msg' => 'items size not enough',
    'items' => [
        ['item_id' => 'P001', 'score' => 0.99, 'retrieve_id' => 'u2i_recall'],
        ['item_id' => '', 'score' => 0.50, 'retrieve_id' => 'u2i_recall'],
        ['item_id' => 'P002', 'score' => 0.80, 'retrieve_id' => 'context_recall'],
    ],
]);

assertSameValue([
    ['productCode' => 'P001', 'score' => 0.99, 'retrieve_id' => 'u2i_recall'],
    ['productCode' => 'P002', 'score' => 0.80, 'retrieve_id' => 'context_recall'],
], $items, '数量不足时仍应按顺序返回有效商品编码');

assertSameValue([], RecommendationFeed::normalizeResponse([
    'code' => 500,
    'items' => [['item_id' => 'P003']],
]), '异常业务码不应向页面返回商品');

echo "recommendation_feed_test: PASS" . PHP_EOL;
