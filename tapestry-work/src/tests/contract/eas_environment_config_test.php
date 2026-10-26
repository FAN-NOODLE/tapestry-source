<?php

$testEnvironment = [];

function env($key, $default = null)
{
    global $testEnvironment;
    return array_key_exists($key, $testEnvironment) ? $testEnvironment[$key] : $default;
}

require_once __DIR__ . '/../../app/Sdk/CDP/Coach/Config.php';

use App\Sdk\CDP\Coach\Config;

function assertConfigValue($expected, $actual, $message)
{
    if ($expected !== $actual) {
        fwrite(STDERR, $message . PHP_EOL);
        fwrite(STDERR, 'Expected: ' . var_export($expected, true) . PHP_EOL);
        fwrite(STDERR, 'Actual: ' . var_export($actual, true) . PHP_EOL);
        exit(1);
    }
}

function readEnvironmentFile(string $path): array
{
    $result = [];
    foreach (file($path, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES) as $line) {
        $line = trim($line);
        if ($line === '' || strpos($line, '#') === 0 || strpos($line, '=') === false) {
            continue;
        }
        [$key, $value] = explode('=', $line, 2);
        $result[trim($key)] = trim($value, " \t\n\r\0\x0B\"'");
    }
    return $result;
}

$uatEnvironment = readEnvironmentFile(__DIR__ . '/../../.env.uat');
assertConfigValue(
    'https://wecom-rec-uat.coach.com.cn/api/predict/sa_rec_service_prepub/api/rec/feed',
    $uatEnvironment['EAS_RECOMMEND_URL'] ?? null,
    'UAT 应配置完整推荐接口地址'
);
assertConfigValue(
    '5cf22c25bfbb2bed8a1aecb216b305f9354ea4863d6070fec34a3984a39100ad',
    hash('sha256', $uatEnvironment['EAS_RECOMMEND_AUTHORIZATION'] ?? ''),
    'UAT 应配置最新 Authorization'
);

$testEnvironment = $uatEnvironment;
$uatConfig = Config::getEasConfig();
assertConfigValue($uatEnvironment['EAS_RECOMMEND_URL'], $uatConfig['url'], '运行时应从环境变量读取完整推荐接口地址');
assertConfigValue($uatEnvironment['EAS_RECOMMEND_AUTHORIZATION'], $uatConfig['authorization'], '运行时应从环境变量读取 Authorization');

$devEnvironment = readEnvironmentFile(__DIR__ . '/../../.env.dev');
assertConfigValue($uatEnvironment['EAS_RECOMMEND_URL'], $devEnvironment['EAS_RECOMMEND_URL'] ?? null, '开发环境应复用 UAT 推荐接口地址');
assertConfigValue(
    hash('sha256', $uatEnvironment['EAS_RECOMMEND_AUTHORIZATION']),
    hash('sha256', $devEnvironment['EAS_RECOMMEND_AUTHORIZATION'] ?? ''),
    '开发环境应复用 UAT Authorization'
);

$prodEnvironment = readEnvironmentFile(__DIR__ . '/../../.env.prod');
assertConfigValue(
    'https://wecom-rec.coach.com.cn/api/predict/sa_rec_service/api/rec/feed',
    $prodEnvironment['EAS_RECOMMEND_URL'] ?? null,
    '生产应配置完整推荐接口地址'
);
assertConfigValue(
    '595070545c7443a54ab3c5104413dbdd9db15bedd768a749c8c7a2ac00adde86',
    hash('sha256', $prodEnvironment['EAS_RECOMMEND_AUTHORIZATION'] ?? ''),
    '生产应配置最新 Authorization'
);

$testEnvironment = $prodEnvironment;
$prodConfig = Config::getEasConfig();
assertConfigValue($prodEnvironment['EAS_RECOMMEND_URL'], $prodConfig['url'], '生产运行时应读取完整推荐接口地址');
assertConfigValue($prodEnvironment['EAS_RECOMMEND_AUTHORIZATION'], $prodConfig['authorization'], '生产运行时应读取 Authorization');

echo "eas_environment_config_test: PASS" . PHP_EOL;
