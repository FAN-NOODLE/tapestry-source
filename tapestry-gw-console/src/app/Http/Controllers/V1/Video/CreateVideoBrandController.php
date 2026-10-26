<?php


namespace App\Http\Controllers\V1\Video;


use App\Http\Controllers\Controller;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\Video;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class CreateVideoBrandController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $params = $this->getParams();
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new Video())->handlerVideoBrand($params));
    }

    /**
     * 获取参数
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate($this->request, [
            'id' => 'integer',
            'license' => 'string',
            'brand_audit_type' => 'integer',
            'trademark_type' => 'string',
            'brand_management_type' => 'integer',
            'commodity_origin_type' => 'integer',
            'brand_wording' => 'string',
            'sale_authorization' => 'array',
            'trademark_registration_certificate' => 'array',
            'trademark_change_certificate' => 'array',
            'trademark_registrant' => 'string',
            'trademark_registrant_nu' => 'string',
            'trademark_authorization_period' => 'string',
            'trademark_registration_application' => 'array',
            'trademark_applicant' => 'string',
            'trademark_application_time' => 'string',
            'imported_goods_form' => 'array',
        ], [
            'id.*' => '参数[id]不合法',
            'license.*' => '参数[license]不合法',
            'brand_audit_type.*' => '参数[brand_audit_type]不合法',
            'trademark_type.*' => '参数[trademark_type]不合法',
            'brand_management_type.*' => '参数[brand_management_type]不合法',
            'commodity_origin_type.*' => '参数[commodity_origin_type]不合法',
            'brand_wording.*' => '参数[brand_wording]不合法',
            'sale_authorization.*' => '参数[sale_authorization]不合法',
            'trademark_registration_certificate.*' => '参数[trademark_registration_certificate]不合法',
            'trademark_change_certificate.*' => '参数[trademark_change_certificate]不合法',
            'trademark_registrant.*' => '参数[trademark_registrant]不合法',
            'trademark_registrant_nu.*' => '参数[trademark_registrant_nu]不合法',
            'trademark_authorization_period.*' => '参数[trademark_authorization_period]不合法',
            'trademark_registration_application.*' => '参数[trademark_registration_application]不合法',
            'trademark_applicant.*' => '参数[trademark_applicant]不合法',
            'trademark_application_time.*' => '参数[trademark_application_time]不合法',
            'imported_goods_form.*' => '参数[imported_goods_form]不合法',
        ]);
    }
}