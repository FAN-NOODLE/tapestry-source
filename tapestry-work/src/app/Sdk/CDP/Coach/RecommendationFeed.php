<?php

namespace App\Sdk\CDP\Coach;

class RecommendationFeed
{
    public static function buildPayload(array $params): array
    {
        return [
            'uid' => $params['uid'] ?? '',
            'size' => (int)($params['size'] ?? 9),
            'scene_id' => $params['scene_id'] ?? 'sa_rec',
            'features' => [
                'sa_id' => $params['features']['sa_id'] ?? '',
                'store_id' => $params['features']['store_id'] ?? '',
            ],
        ];
    }

    public static function buildOptions(array $payload, string $authorization): array
    {
        return [
            'headers' => [
                'Authorization' => $authorization,
                'Content-Type' => 'application/json',
                'Cache-Control' => 'no-cache',
            ],
            'body' => json_encode($payload),
        ];
    }

    public static function normalizeResponse(array $response): array
    {
        if (!in_array($response['code'] ?? 0, [200, 299], true)) {
            return [];
        }

        $result = [];
        foreach ($response['items'] ?? [] as $item) {
            $productCode = trim((string)($item['item_id'] ?? ''));
            if ($productCode === '') {
                continue;
            }
            $result[] = [
                'productCode' => $productCode,
                'score' => $item['score'] ?? 0,
                'retrieve_id' => $item['retrieve_id'] ?? '',
            ];
        }
        return $result;
    }
}
