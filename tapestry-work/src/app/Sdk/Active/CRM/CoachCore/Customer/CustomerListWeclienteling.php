<?php

namespace App\Sdk\Active\CRM\CoachCore\Customer;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * 会员列表聚合搜索
 */
class CustomerListWeclienteling extends Base
{
    public $errStr = 'query member info error';
    public $api    = '/weclientws/ws/rest/customer/showCustomerListWeclienteling';


    private $storeCode;
    private $regionCode;
    private $curPage;
    private $pageSize;

    private $saAccount;
    private $memberIds;
    private $name;
    private $phoneNumber;
    private $email;
    private $city;
    private $region;
    private $firstPurchase;
    private $lastPurchase;
    private $collection;
    private $style;
    private $color;
    private $priceRange;
    private $department;
    private $vipLevel;
    private $gender;
    private $birthMonth;
    private $type;

    public function __construct(array $params)
    {
        parent::__construct();

        $this->storeCode  = $params['storeCode'] ?? '';
        $this->curPage    = $params['curPage'] ?? 1;
        $this->pageSize   = $params['pageSize'] ?? 10;
        $this->regionCode = !empty($params['regionCode']) ? $params['regionCode'] : 'CN';

        $this->saAccount     = $params['saAccount'] ?? '';
        $this->memberIds     = $params['memberIds'] ?? '';
        $this->name          = $params['name'] ?? '';
        $this->phoneNumber   = $params['phoneNumber'] ?? '';
        $this->email         = $params['email'] ?? '';
        $this->city          = $params['city'] ?? '';
        $this->region        = $params['region'] ?? '';
        $this->firstPurchase = $params['firstPurchase'] ?? '';
        $this->lastPurchase  = $params['lastPurchase'] ?? '';
        $this->collection    = $params['collection'] ?? '';
        $this->style         = $params['style'] ?? '';
        $this->color         = $params['color'] ?? '';
        $this->priceRange    = $params['priceRange'] ?? '';
        $this->department    = $params['department'] ?? '';
        $this->vipLevel      = $params['vipLevel'] ?? '';
        $this->gender        = $params['gender'] ?? '';
        $this->birthMonth    = $params['birthMonth'] ?? '';
        $this->type          = $params['type'] ?? '';

        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        // 门店 storeCode=""  curPage="1" pageSize="100"
        //导购code saAccount="" 会员id（多个） memberIds="" 会员姓名 name="" 顾客手机号 phoneNumber="" email="" 城市city="" region="" 购买开始时间firstPurchase="" 购买结束时间lastPurchase="" collection="" 款式style="" 颜色color="" 消费总金额区间范围 priceRange="" 产品类型department="" 会员等级vipLevel="" 会员性别gender="" brithMonth=""  会员类型type=""

        $strMemberIds = "";
        // memberIds 为空 不带次参数
        if ($this->memberIds) {
            $strMemberIds = ' memberIds="' . $this->memberIds.'"';
        }
        $this->queryParam = '<?xml version="1.0" encoding="UTF-8" standalone="no"?>' .
            '<request storeCode="' . $this->storeCode . '" regionCode="' . $this->regionCode . '" curPage="' . $this->curPage . '" pageSize="' . $this->pageSize . '">' .
            '<weCustomerListReq saAccount="' . $this->saAccount . '" ' . $strMemberIds . ' name="' . $this->name . '" phoneNumber="' . $this->phoneNumber . '" email="' . $this->email . '" city="' . $this->city . '" region="' . $this->region . '" firstPurchase="' . $this->firstPurchase . '" lastPurchase="' . $this->lastPurchase . '" collection="' . $this->collection . '" style="' . $this->style . '" color="' . $this->color . '" priceRange="' . $this->priceRange . '" department="' . $this->department . '" vipLevel="' . $this->vipLevel . '" gender="' . $this->gender . '" birthMonth="' . $this->birthMonth . '" type="' . $this->type . '"/>' .
            '</request>';


        // № Name of data Character type Size
        //(Length) Description Definition of code
        //Element Of request
        //1 storeCode alphanumeric 10 SA's Store Code
        //Exmaple: J101 or OC01
        //2 regionCode alphabet 2
        //SA's Region Code
        //Example: JP or KR
        //JP：Japan、KR：Korea、
        //TW：Taiwan、 HK：Hongkong、CN：
        //China、MY：Malaysia、
        //SG：Singapore
        //3 curPage numeric 10 Specify page number
        //Example: 2
        //4 pageSize numeric 10
        //Set record number which is included
        //in 1 page
        //Example: 10
        //Element Of wecustomerListReq(request)
        //1
        //saAccount(Unique
        //Code) alphanumeric 8
        //If this value is empty, response
        //should be all customers. If this value
        //is not empty, response should be
        //customers who has this SA employee
        //ID.
        //When "UASALL" is set, response
        //should be unassigned customers.
        //UASALL：Unassigned
        //Customer
        //2 name string 100 Search by customer Name
        //3 memberIds text
        //Search by Client's member IDs
        //Example: COH651487,COH651488
        //4 phoneNumber numeric 20 Search by Phone number. Needs to
        //search both home phone and cell
        //20
        //phone
        //Example: 08012341234
        //5 email string 200
        //Search by Email address. Needs to
        //search both PC and Mobil
        //6 city string 50 Search by City
        //7 region string 50 Search by Region (Prefecture)
        //8 firstPurchase numeric 6
        //Search bv first purchase year &
        //month
        //Example: 201401
        //9 lastPurchase numeric 6
        //Search bv last purchase year &
        //month
        //Example: 201408
        //10 collection string 100
        //Search bv style_desc (product long
        //description)
        //(*Wildcard search should be
        //possbile)
        //Example: Madison -> In your
        //program, it needs to be treated as
        //"*Madison*"
        //11 style alphanumeric 6
        //Search bv Style
        //Example: Q6556
        //12 color alphanumeric 5
        //Search bv Color Code
        //Example: BK/BK
        //13 priceRange numeric 1
        //Search by Price Range (LAST YEAR
        //TOTAL AMOUNT)
        //Example: 1 (Under RMB 3000)
        //LAST YEAR TOTAL AMOUNT
        //CHINA (CC/MO/HK)
        //0: All
        //1: 3,000 以下
        //2: 3,000-4,999(包含)
        //3: 5,000-7,999(包含)
        //4: 8,000-9,999(包含)
        //5: 10,000-19,999(包含)
        //6: 20,000-29,999(包含)
        //7: 30,000 及以上
        //14 department alphanumeric 3
        //Search by Department Code
        //(product category code)
        //Example: D01 (WOMENS BAGS)
        //0: All
        //D01: WOMENS BAGS
        //D02: WOMENS WALLETS
        //D03: WOMENS SLGs
        //D04: MENS BAGS
        //D05: MENS WALLETS
        //D06: WOMENS TRAVEL
        //D07: WOMENS READY-
        //21
        //TO-WEAR
        //D08: MENS SPECIALTY
        //ACCESSORIES
        //D09: WOMENS
        //WATCHES
        //D10: MENS TRAVEL
        //D11: WOMENS SHOES
        //D12: WOMENS
        //EYEWEAR
        //D13: JEWELRY
        //D14: OTHER
        //D15: KIDS/BABY
        //D16: MENS READY-TOWEAR
        //D17: WOMENS
        //FRAGRANCE
        //D18: MENS WATCHES
        //D19: MENS SHOES
        //D20: MENS EYEWEAR
        //D21: MENS FRAGRANCE
        //D22: WOMENS SOFT
        //ACCESSORIES
        //D23: WOMENS
        //SPECIALTY ACCESSORIES
        //15 vipLevel numeric 1
        //Search by VIP level
        //Example: 1 (Silver)
        //0：All、1：Silver、2：Gold、 3：Dimand
        //16 gender numeric 1
        //Search by Gender
        //Example: 1 (Female)
        //0：All、1：Female、2：
        //Male、3：Unknown
        //17 birthMonth numeric 2
        //Search by Birthday Month
        //Example: 02 (February)
        //00：All、01：January、02：
        //Febuary、
        //03：March、04：April、
        //05：May、06：June、
        //07：July、08：August、
        //09：September、
        //10：October、11：
        //November、12：
        //December
        //18 type numeric 1
        //Customer's Type
        //Example:
        //1 (Top Client)
        //1: Top Client
        //2: New High Potential
        //3: Top Client or New
        //High Potential
        //0: All, include TC, NHP or
        //others

//        $this->queryParam = [
//            'storeCode'  => $this->storeCode,
//            'regionCode' => $this->regionCode,
//            'curPage'    => $this->curPage,
//            'pageSize'   => $this->pageSize,
//        ];
//        $this->queryParam = arrayToXmlRecur([
//            'request' => [
//                'storeCode'         => $this->storeCode,
//                'regionCode'        => $this->regionCode,
//                'curPage'           => $this->curPage,
//                'pageSize'          => $this->pageSize,
//                'weCustomerListReq' => ['saAccount' => "", 'name' => "", 'phoneNumber' => "", 'email' => "", 'city' => "", 'region' => "", 'firstPurchase' => "", 'lastPurchase' => "", 'collection' => "", 'style' => "", 'color' => "", 'priceRange' => "", 'department' => "", 'vipLevel' => "", 'gender' => "", 'brithMonth' => "", 'memberIds' => "", 'type' => "",]
//            ]
//        ]);

    }
}