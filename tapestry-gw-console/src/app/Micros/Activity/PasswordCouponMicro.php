<?php
/**
 * Created by PhpStorm.
 * User: wangpengfei
 * Date: 2020-03-19
 * Time: 16:56
 */

namespace App\Micros\Activity;

use CrsSdk\MicroActivity\PasswordCoupon;

class PasswordCouponMicro
{
    public function generatePosters($params){
        return (new PasswordCoupon())->generatePosters($params);
    }

    public function add($params){
        return (new PasswordCoupon())->add($params);
    }

    public function list($params){
        return (new PasswordCoupon())->list($params);
    }


    public function passwordListExport($params)
    {
        return (new PasswordCoupon())->passwordListExport($params);
    }

    public function DownloadPosterMoreStatus($params)
    {
        return (new PasswordCoupon())->DownloadPosterMoreStatus($params);
    }

    public function detail($params){
        return (new PasswordCoupon())->detail($params);
    }

    public function update($params){
        return (new PasswordCoupon())->update($params);
    }
}
