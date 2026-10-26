<?php
/*
 * @Author: your name
 * @Date: 2020-02-10 09:47:35
 * @LastEditTime: 2020-03-11 17:37:04
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \shop\src\app\Services\PasswordCouponService.php
 */


namespace App\Services;

use App\Common\Constant;
use App\Library\ImageTools;
use App\Library\Uc\Member;
use Illuminate\Support\Facades\DB;

class PosterService extends BaseService
{
    /**
     * 页面对应path
     * @var array
     */
    protected $marketingPath
        = [
            Constant::MARKETING_GOODS           => 'pages/goods-detail/goods-detail',
            Constant::MARKETING_COUPON          => 'pages/receive-coupons/receive-coupons',
            Constant::MARKETING_HELPER_POST     => 'pages/goods-detail/goods-detail',
            Constant::MARKETING_HELPER_BUY_GIVE => 'pages/goods-detail/goods-detail',
            Constant::MARKETING_SECKILL         => 'pages/goods-detail/goods-detail',
            Constant::MARKETING_LOTTERY         => 'pages/goods-detail/goods-detail',
            Constant::MARKETING_GROUP           => 'pages/goods-detail/goods-detail',
            Constant::MARKETING_BARGAIN         => 'pages/goods-detail/goods-detail',
        ];

    protected $releaseEnterprise = [13, 120, 102, 73, 86, 100244242505984, 102756213248256];

    private $imageTools;

    /**
     * PosterService constructor.
     */
    public function __construct()
    {
        $this->imageTools = new ImageTools();
    }

    /**
     * 生成普通商品海报
     *
     * @param $params
     * @param $member
     * @param $operateType
     * @param string $domain
     * @return mixed
     * @throws \Exception
     * @internal param $appInfo
     */
    public function generateGoodsPoster($params, $member, $operateType, $domain = '')
    {
        $multiple     = 2;
        $jpeg_quality = 100;  // 图片质量90，满分为100

        $params['picture_url'] = $this->imageTools->getImageSource($params['picture_url'])[0];
        $params['sun_code']    = $this->imageTools->getImageSource($params['sun_code'])[0];

        $member['avatar'] = $this->imageTools->getImageSource($member['memFace'])[0];
        $member['avatar'] = $this->imageTools->toCircularSource($member['avatar']);

        //创建画布
        $im = imagecreatetruecolor(375 * $multiple, 632 * $multiple);
        //填充画布背景色
        $color = imagecolorallocate($im, 255, 255, 255);
        imagefill($im, 0, 0, $color);

        //字体文件
        $font_file = base_path('public/fonts/msyh.ttc');

        // 商品图
        $g_w = imagesx($params['picture_url']);
        $g_h = imagesy($params['picture_url']);
        imagecopyresampled($im, $params['picture_url'], 15 * $multiple, 155 * $multiple, 0, 0, 345 * $multiple, 345 * $multiple, $g_w, $g_h);
        // 员工头像
        $m_w = imagesx($member['avatar']);
        $m_h = imagesy($member['avatar']);
        imagecopyresampled($im, $member['avatar'], 15 * $multiple, 15 * $multiple, 0, 0, 60 * $multiple, 60 * $multiple, $m_w, $m_h);

        // 双引号
        $marks = base_path('public/sourceImages/marks@3x.png');
        $marks = imagecreatefrompng($marks);
        $m_w   = imagesx($marks);
        $m_h   = imagesy($marks);
        imagecopyresampled($im, $marks, 300 * $multiple, 25 * $multiple, 0, 0, 26.37 * $multiple, 25.05 * $multiple, $m_w, $m_h);

        // 认证员工
        $authMember = base_path('public/sourceImages/renzheng@3x.png');
        $authMember = imagecreatefrompng($authMember);
        $a_w        = imagesx($authMember);
        $a_h        = imagesy($authMember);
        imagecopyresampled($im, $authMember, 90 * $multiple, 49 * $multiple, 0, 0, 65 * $multiple, 19 * $multiple, $a_w, $a_h);

        // 太阳码
        $s_w = imagesx($params['sun_code']);
        $s_h = imagesy($params['sun_code']);
        imagecopyresampled($im, $params['sun_code'], 279 * $multiple, 508 * $multiple, 0, 0, 75 * $multiple, 75 * $multiple, $s_w, $s_h);

        // 员工名称
        $font_style = ImageColorAllocate($im, 127, 127, 127);
        $fontSize   = 11 * $multiple;
        imagettftext($im, $fontSize, 0, 90 * $multiple, 40 * $multiple, $font_style, $font_file, $member['memUsername']);

        $fontBox    = imagettfbbox($fontSize, 0, $font_file, $member['memUsername']);
        $font_style = ImageColorAllocate($im, 50, 50, 50);
        imagettftext($im, $fontSize, 0, $fontBox[2] + 90 * $multiple, 40 * $multiple, $font_style, $font_file, '为你推荐');

        // 商品名称
        $fontSize = 13.5 * $multiple;
        $text     = $this->Wrap([
            'font_size' => $fontSize,
            'string'    => $params['sku_name'],
            'font_file' => $font_file,
            'width'     => (375 - 15 * 2) * $multiple
        ]);
        imagettftext($im, $fontSize, 0, 15 * $multiple, 110 * $multiple, $font_style, $font_file, $text);

        // 商品售卖价
        $price      = explode('.', $params['price'] / 100);
        $font_style = ImageColorAllocate($im, 250, 64, 70);
        $fontSize   = 14 * $multiple;
        imagettftext($im, $fontSize, 0, 15 * $multiple, 540 * $multiple, $font_style, $font_file, '¥');

        $priceFirst = 0;
        if (isset($price[0])) {
            $fontSize = 25 * $multiple;
            imagettftext($im, $fontSize, 0, 27 * $multiple, 541 * $multiple, $font_style, $font_file, $price[0]);
            $fontBox    = imagettfbbox($fontSize, 0, $font_file, $price[0]);
            $priceFirst = $fontBox[2];
        }

        $priceSecond = 0;
        if (isset($price[1])) {
            $fontSize = 14 * $multiple;
            imagettftext($im, $fontSize, 0, $priceFirst + 27 * $multiple, 541 * $multiple, $font_style, $font_file, '.' . $price[1]);
            $fontBox     = imagettfbbox($fontSize, 0, $font_file, $price[1]);
            $priceSecond = $fontBox[2];
        }

        // 吊牌价
        if (!empty($params['sign_price'])) {
            $fontSize   = 14 * $multiple;
            $font_style = ImageColorAllocate($im, 166, 166, 166);
            imagettftext($im, $fontSize, 0, $priceFirst + $priceSecond + 45 * $multiple, 541 * $multiple, $font_style, $font_file, '¥' . $params['sign_price'] / 100);

            $fontBox = imagettfbbox($fontSize, 0, $font_file, '¥' . $params['sign_price'] / 100);

            // 删除线
            imageline($im, $priceFirst + $priceSecond + 45 * $multiple, 534 * $multiple, $priceFirst + $priceSecond + 45 * $multiple + $fontBox[2], 534 * $multiple, $font_style);
            imageline($im, $priceFirst + $priceSecond + 45 * $multiple, 534.5 * $multiple, $priceFirst + $priceSecond + 45 * $multiple + $fontBox[2], 534.5 * $multiple, $font_style);
        }

        // 时间不等人 长按识别抢购
        $fontSize   = 14 * $multiple;
        $font_style = ImageColorAllocate($im, 127, 127, 127);
        imagettftext($im, $fontSize, 0, 140 * $multiple, 615 * $multiple, $font_style, $font_file, '时间不等人 长按识别抢购');

        return $this->getMediaId($im, $domain, $operateType);
    }

    /**
     * 生成优惠券海报
     *
     * @param $appInfo
     * @param $params
     * @param $operateType
     * @param string $domain
     * @return mixed
     * @throws \Exception
     * @internal param $appInfo
     */
    public function generateCouponPoster($appInfo, $params, $operateType, $domain = '')
    {
        $appInfo['logo'] = $this->imageTools->getImageSource($appInfo['logo'])[0];
        $appInfo['logo'] = $this->imageTools->toCircularSource($appInfo['logo']);

        $params['sun_code'] = $this->imageTools->getImageSource($params['sun_code'])[0];
        $params['sun_code'] = $this->imageTools->toCircularSource($params['sun_code']);

        $time      = explode('至', $params['validity']);
        $startTime = '';
        $endTime   = '';
        if ($time[0]) {
            $startTime = date('Y.m.d', strtotime($time[0]));
        }
        if ($time[1]) {
            $endTime = date('Y.m.d', strtotime($time[1]));
        }

        $multiple     = 2;
        $jpeg_quality = 100;  // 图片质量90，满分为100

        // 海报背景图片
        $bg_file = base_path('public/sourceImages/coupon@2x.png');
        //字体文件
        $font_file = base_path('public/fonts/msyh.ttf');

        $im    = imagecreatefrompng($bg_file);
        $width = imagesx($im);

        //Logo
        $l_w = imagesx($appInfo['logo']);
        $l_h = imagesy($appInfo['logo']);
        imagecopyresampled($im, $appInfo['logo'], 33 * $multiple, 469.5 * $multiple, 0, 0, 55 * $multiple, 55 * $multiple, $l_w, $l_h);

        //太阳码
        $l_w = imagesx($params['sun_code']);
        $l_h = imagesy($params['sun_code']);
        imagecopyresampled($im, $params['sun_code'], 265 * $multiple, 455 * $multiple, 0, 0, 84.12 * $multiple, 84.12 * $multiple, $l_w, $l_h);

        $font_style = ImageColorAllocate($im, 253, 246, 234);
        $fontSize   = $multiple * 48.68;
        $fontBox    = imagettfbbox($fontSize, 0, $font_file, '¥' . $params['coupon_value'] / 100);
        imagettftext($im, $fontSize, 0, ceil(($width - $fontBox[2]) / 2), $multiple * 270, $font_style, $font_file, '¥' . $params['coupon_value'] / 100);

        // 使用条件
        $font_style = ImageColorAllocate($im, 240, 213, 160);
        $fontSize   = $multiple * 22;
        $fontBox    = imagettfbbox($fontSize, 0, $font_file, $params['condition']);
        imagettftext($im, $fontSize, 0, ceil(($width - $fontBox[2]) / 2), $multiple * 355, $font_style, $font_file, $params['condition']);

        // 使用期限
        $string   = $startTime . '-' . $endTime;
        $fontSize = $multiple * 16;
        $fontBox  = imagettfbbox($fontSize, 0, $font_file, $string);
        imagettftext($im, $fontSize, 0, ceil(($width - $fontBox[2]) / 2), $multiple * 400, $font_style, $font_file, $string);

        return $this->getMediaId($im, $domain, $operateType);
    }

    /**
     * 生成活动商品海报
     *
     * @param        $appInfo
     * @param        $params
     * @param        $member
     * @param        $operateType
     * @param string $domain
     *
     * @return mixed
     * @throws \Exception
     * @internal param $appInfo
     */
    public function generateActivityPoster($appInfo, $params, $member, $operateType, $domain = '')
    {
        $multiple     = 2;
        $jpeg_quality = 100;  // 图片质量90，满分为100

        $params['sku_picture'] = $this->imageTools->getImageSource($params['sku_picture'])[0];
        $params['sun_code']    = $this->imageTools->getImageSource($params['sun_code'])[0];
        $member['avatar']      = $this->imageTools->getImageSource($member['memFace'])[0];
        $member['avatar']      = $this->imageTools->toCircularSource($member['avatar']);

        $appInfo['logo'] = $this->imageTools->getImageSource($appInfo['logo'])[0];
        $appInfo['logo'] = $this->imageTools->toCircularSource($appInfo['logo']);

        //创建画布
        $im = imagecreatetruecolor(375 * $multiple, 657 * $multiple);
        //填充画布背景色
        $color = imagecolorallocate($im, 255, 255, 255);
        imagefill($im, 0, 0, $color);

        $color = imagecolorallocate($im, $appInfo['style_rgb'][0], $appInfo['style_rgb'][1], $appInfo['style_rgb'][2]);
//        $color = imagecolorallocate($im, 191, 0, 23);
        imagefilledrectangle($im, 0, 0, 375 * $multiple, 70 * $multiple, $color);

        $width = imagesx($im);

        //字体文件
        $font_file = base_path('public/fonts/msyh.ttc');

        // logo
        $m_w = imagesx($appInfo['logo']);
        $m_h = imagesy($appInfo['logo']);
        imagecopyresampled($im, $appInfo['logo'], 15 * $multiple, 20 * $multiple, 0, 0, 30 * $multiple, 30 * $multiple, $m_w, $m_h);

        // 小程序名称
        $font_style = ImageColorAllocate($im, 255, 255, 255);
        $fontSize   = $multiple * 16;
        imagettftext($im, $fontSize, 0, 55 * $multiple, $multiple * 43, $font_style, $font_file, $appInfo['name']);

        // 商品图
        $g_w = imagesx($params['sku_picture']);
        $g_h = imagesy($params['sku_picture']);
        imagecopyresampled($im, $params['sku_picture'], 15 * $multiple, 213 * $multiple, 0, 0, 345 * $multiple, 345 * $multiple, $g_w, $g_h);
        // 员工头像
        $m_w = imagesx($member['avatar']);
        $m_h = imagesy($member['avatar']);
        imagecopyresampled($im, $member['avatar'], 15 * $multiple, 88 * $multiple, 0, 0, 45 * $multiple, 45 * $multiple, $m_w, $m_h);

        // 双引号
        $marks = base_path('public/sourceImages/marks@3x.png');
        $marks = imagecreatefrompng($marks);
        $m_w   = imagesx($marks);
        $m_h   = imagesy($marks);
        imagecopyresampled($im, $marks, 300 * $multiple, 91 * $multiple, 0, 0, 26.37 * $multiple, 25.05 * $multiple, $m_w, $m_h);

        // 认证员工
        $authMember = base_path('public/sourceImages/renzheng@3x.png');
        $authMember = imagecreatefrompng($authMember);
        $a_w        = imagesx($authMember);
        $a_h        = imagesy($authMember);
        imagecopyresampled($im, $authMember, 72 * $multiple, 112 * $multiple, 0, 0, 65 * $multiple, 19 * $multiple, $a_w, $a_h);

        // 员工名称
        $font_style = ImageColorAllocate($im, 127, 127, 127);
        $fontSize   = 11 * $multiple;
        imagettftext($im, $fontSize, 0, 72 * $multiple, 103 * $multiple, $font_style, $font_file, $member['memUsername']);

        $fontBox    = imagettfbbox($fontSize, 0, $font_file, $member['memUsername']);
        $font_style = ImageColorAllocate($im, 50, 50, 50);
        imagettftext($im, $fontSize, 0, $fontBox[2] + 72 * $multiple, 103 * $multiple, $font_style, $font_file, '为你推荐');

        // 商品名称
        $fontSize = 13.5 * $multiple;
        $text     = $this->Wrap([
            'font_size' => $fontSize,
            'string'    => $params['sku_name'],
            'font_file' => $font_file,
            'width'     => (375 - 15 * 2) * $multiple
        ]);
        imagettftext($im, $fontSize, 0, 15 * $multiple, 167 * $multiple, $font_style, $font_file, $text);

        // 商品售卖价
        $price      = explode('.', $params['activity_price'] / 100);
        $font_style = ImageColorAllocate($im, 250, 64, 70);
        $fontSize   = 14 * $multiple;
        imagettftext($im, $fontSize, 0, 15 * $multiple, 605 * $multiple, $font_style, $font_file, '¥');
        $priceFirst = 0;
        if (isset($price[0])) {
            $fontSize = 25 * $multiple;
            imagettftext($im, $fontSize, 0, 27 * $multiple, 605 * $multiple, $font_style, $font_file, $price[0]);
            $fontBox    = imagettfbbox($fontSize, 0, $font_file, $price[0]);
            $priceFirst = $fontBox[2];
        }

        $priceSecond = 0;
        if (isset($price[1])) {
            $fontSize = 14 * $multiple;
            imagettftext($im, $fontSize, 0, $priceFirst + 27 * $multiple, 605 * $multiple, $font_style, $font_file, '.' . $price[1]);
            $fontBox     = imagettfbbox($fontSize, 0, $font_file, $price[1]);
            $priceSecond = $fontBox[2];
        }

        $activityIcon = $this->getActivityIcon($params['activity_type']);
        // 上方icon
        $topIcon = $activityIcon[0];
        $topIcon = imagecreatefrompng($topIcon);
        $t_w     = imagesy($topIcon) * $activityIcon[2];
        $t_h     = imagesy($topIcon);

        // 下方icon
        $bottomIcon = $activityIcon[1];
        $bottomIcon = imagecreatefrompng($bottomIcon);
        $d_w        = imagesy($bottomIcon) * $activityIcon[3];
        $d_h        = imagesy($bottomIcon);


        imagecopyresampled($im, $topIcon, 273 * $multiple, 24 * $multiple, 0, 0, 88 * $multiple, 22 * $multiple, $t_w, $t_h);

        imagecopyresampled($im, $bottomIcon, $priceFirst + $priceSecond + (27 + 15) * $multiple, 580 * $multiple, 0, 0, 64.5 * $multiple, 28.5 * 1.8, $d_w, $d_h);

        // 吊牌价
        if (!empty($params['sell_price'])) {
            $fontSize   = 12 * $multiple;
            $font_style = ImageColorAllocate($im, 166, 166, 166);
            imagettftext($im, $fontSize, 0, 15 * $multiple, 630 * $multiple, $font_style, $font_file, '¥' . $params['sell_price'] / 100);

            $fontBox = imagettfbbox($fontSize, 0, $font_file, '¥' . $params['sell_price'] / 100);

            // 删除线
            imageline($im, 15 * $multiple, 624 * $multiple, 16 * $multiple + $fontBox[2], 624 * $multiple, $font_style);
            imageline($im, 15 * $multiple, 624.5 * $multiple, 16 * $multiple + $fontBox[2], 624.5 * $multiple, $font_style);
        }

        // 太阳码
        $s_w = imagesx($params['sun_code']);
        $s_h = imagesy($params['sun_code']);
        imagecopyresampled($im, $params['sun_code'], 280 * $multiple, 568 * $multiple, 0, 0, 75 * $multiple, 75 * $multiple, $s_w, $s_h);

        return $this->getMediaId($im, $domain, $operateType);
    }

    /**
     * 文本换行
     * @param $params
     * @return string
     */
    public function Wrap($params)
    {
        $length = $len = mb_strlen($params['string'], "utf-8");

        $text = '';

        while (true) {
            $text = mb_substr($params['string'], 0, $len);

            $fontBox = imagettfbbox($params['font_size'], 0, $params['font_file'], $text);

            if ($fontBox[2] - $fontBox[0] < $params['width']) {
                break;
            }

            $len--;
        }

        return $text . PHP_EOL . mb_substr($params['string'], $len, $length);
    }

    /**
     * 获取素材id
     *
     * @param $im
     * @param $domain
     * @param $optionType
     *
     * @return mixed
     * @throws \Exception
     */
    public function getMediaId($im, $domain, $optionType)
    {
        if ($optionType == Constant::OPERATE_TYPE_DOWNLOAD) {
            ob_start();
            imagepng($im);
            $imageData = ob_get_contents();
            ob_end_clean();
            $base64String = base64_encode($imageData);

            //释放空间
            imagedestroy($im);

            return $base64String;

//            header('Content-type: image/jpeg');
//            imagejpeg($im,null, 100); // 输出图象到浏览器或文件
        }

        //输出图片
        $path = base_path('temp');
        if (!is_dir($path)) {
            mkdir($path, 0777, true);
        }
        $fileName = $path . '/' . rand(1, 99999) . '.png';
        imagepng($im, $fileName);
        $file = fopen($fileName, 'r');

        // 获取素材id
        $multipart = [
            [
                'name'     => 'file',
                'contents' => $file,
            ],
            [
                'name'     => 'fileType',
                'contents' => 1, // 图片
            ],
            [
                'name'     => 'filename',
                'contents' => rand(10000, 99999),
            ],
        ];

        $result = (new Member())->imageToMediaId($domain, 'contact', $multipart);

        unlink($fileName);

        //释放空间
        imagedestroy($im);

        return $result['mediaId'];
    }

    /**
     * 获取活动icon
     * @param $type
     * @return array
     */
    public function getActivityIcon($type)
    {
        switch ($type) {
            case Constant::MARKETING_HELPER_POST :
            case Constant::MARKETING_HELPER_BUY_GIVE :
            case Constant::MARKETING_LOTTERY :
                $return = [
                    base_path('public/sourceImages/zhuliTop@3x.png'),
                    base_path('public/sourceImages/zhuliBottom@3x.png'),
                    4.6,
                    2.8

                ];
                break;
            case Constant::MARKETING_SECKILL :
                $return = [
                    base_path('public/sourceImages/xianshiTop@3x.png'),
                    base_path('public/sourceImages/xianshiBottom@3x.png'),
                    3.7,
                    2.25

                ];
                break;
            case Constant::MARKETING_GROUP :
                $return = [
                    base_path('public/sourceImages/pintuanTop@3x.png'),
                    base_path('public/sourceImages/pintuanBottom@3x.png'),
                    4,
                    2.25

                ];
                break;
            case Constant::MARKETING_BARGAIN :
                $return = [
                    base_path('public/sourceImages/kanjiaTop@3x.png'),
                    base_path('public/sourceImages/kanjiaBottom@3x.png'),
                    3.8,
                    2.25

                ];
                break;
            default :
                $return = [
                    base_path('public/sourceImages/kanjiaTop@3x.png'),
                    base_path('public/sourceImages/kanjiaBottom@3x.png'),
                    4,
                    2.25

                ];
                break;
        }

        return $return;
    }

    /**
     * 获取二维码配置
     * @param $type
     * @param $params
     * @return array
     */
    public function getQrCodeParams($type, $params)
    {
        // 查询新版路径
        $posterPathData = DB::table("poster_release_enterprise")
            ->where("enterprise_id", $params["enterprise_id"])
            ->first();
        if (!empty($posterPathData)) {
            $marketingPath = json_decode($posterPathData->marketing_path,true);
        } else {
            $marketingPath = $this->marketingPath;
        }

        $path   = '';
        $return = [];
        switch ($type) {
            case Constant::MARKETING_GOODS :
                $return = [
                    'spu_id'   => $params['spu_id'],
                    'sku_id'   => $params['id'],
                    'mem_uid'  => $params['mem_uid'],
                    'store_id' => $params['store_id'] ?? 0,
                    'source'   => $params['source'] ?? '',
                    'share_uid' => $params['share_uid'] ?? '',
                ];
                $path   = $marketingPath[$type];
                break;
            case Constant::MARKETING_COUPON :
                $return = [
                    'id'       => $params['id'],
                    'store_id' => $params['store_id'] ?? 0,
                    'source'   => $params['source'] ?? '',
                    'share_uid' => $params['share_uid'] ?? '',
                ];
                if (isset($params['mem_uid']) && !empty($params['mem_uid'])) {
                    $return['mem_uid'] = $params['mem_uid'];
                }
                $path   = $marketingPath[$type];
                break;
            case Constant::MARKETING_HELPER_POST :
            case Constant::MARKETING_HELPER_BUY_GIVE :
            case Constant::MARKETING_SECKILL :
            case Constant::MARKETING_LOTTERY :
            case Constant::MARKETING_BARGAIN :
            case Constant::MARKETING_GROUP :
                $return = [
                    'spu_id'   => $params['spu_id'],
                    'sku_id'   => $params['sku_id'],
                    'type'     => $type,
                    'rule_id'  => $params['activity_rule_id'],
                    'mem_uid'  => $params['mem_uid'],
                    'store_id' => $params['store_id'] ?? 0,
                    'source'   => $params['source'] ?? '',
                    'share_uid' => $params['share_uid'] ?? '',
                ];
                $path   = $marketingPath[$type];
                break;
            default :
                break;
        }

        return [$path, $return];
    }
}
