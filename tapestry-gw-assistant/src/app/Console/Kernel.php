<?php

namespace App\Console;

use App\Console\Commands\ActivityCenter\Haggle\HaggleCompatData;
use App\Console\Commands\ActivityCenter\discountCard\DiscountCardRuleEnd;
use App\Console\Commands\ActivityCenter\discountCard\TicketExpired;
use App\Console\Commands\ActivityCenter\Haggle\HaggleInitFail;
use App\Console\Commands\ActivityCenter\Haggle\HaggleRuleEnd;
use App\Console\Commands\ActivityCenter\Haggle\HaggleRuleEndSkuRollBack;
use App\Console\Commands\ActivityCenter\Luck\LotteryDrawExt;
use App\Console\Commands\ActivityCenter\Group\GroupCompatData;
use App\Console\Commands\ActivityCenter\Group\GroupFail;
use App\Console\Commands\ActivityCenter\Group\GroupRobot;
use App\Console\Commands\ActivityCenter\Group\GroupRuleEnd;
use App\Console\Commands\ActivityCenter\Group\GroupTicket;
use App\Console\Commands\ActivityEnd;
use App\Console\Commands\ActivityHelpDataUpdate;
use App\Console\Commands\AdjustUserScore;
use App\Console\Commands\AutoCancelGiftCertificate;
use App\Console\Commands\AutoCancelOrder;
use App\Console\Commands\AutoChangeLabel;
use App\Console\Commands\AutoChangeRegionalTotal;
use App\Console\Commands\AutoChangeUserLevel;
use App\Console\Commands\AutoCheckPushRefundApply;
use App\Console\Commands\AutoCloseDoubleCouponActivity;
use App\Console\Commands\AutoCourierSign;
use App\Console\Commands\AutoChangeMiniTpl;
use App\Console\Commands\AutoMoveOverdueUserCoupon;
use App\Console\Commands\AutoLiveCheck;
use App\Console\Commands\AutoRefundApply;
use App\Console\Commands\AutoReleaseVersion;
use App\Console\Commands\AutoSuccessOrder;
use App\Console\Commands\AutoUpdateOrderUser;
use App\Console\Commands\AutoSystemLabel;
use App\Console\Commands\BatchSale;
use App\Console\Commands\ChangeUserAll;
use App\Console\Commands\CouponInvalid;
use App\Console\Commands\Document;
use App\Console\Commands\ExportAchievementDpList;
use App\Console\Commands\ExportAchievementMemberList;
use App\Console\Commands\ExportAchievementOrderList;
use App\Console\Commands\ExportOrder;
use App\Console\Commands\ExportStaffOrderList;
use App\Console\Commands\InitOrderUser;
use App\Console\Commands\InitUserAll;
use App\Console\Commands\InvalidScoreDeal;
use App\Console\Commands\ActivityCenter\Luck\LotteryDraw;
use App\Console\Commands\QueryLogistic;
use App\Console\Commands\ResetUserData;
use App\Console\Commands\SendWxCoupons;
use App\Console\Commands\SpuPostageTemplate;
use App\Console\Commands\TemplatePageVersionRelease;
use App\Console\Commands\Test;
use App\Console\Commands\TransferOrderLogistics;
use App\Console\Commands\UnbindStaffShare;
use App\Console\Commands\UpdateFlashSaleStatus;
use App\Console\Commands\UpdateGiftAndSortData;
use App\Console\Commands\UpdateHelpData;
use App\Console\Commands\UpdateHelpStatus;
use App\Console\Commands\UpdateOrderData;
use App\Console\Commands\UpdateRelationData;
use App\Console\Commands\SysHelpOrder;
use App\Console\Commands\UpdateSkuSales;
use App\Console\Commands\UpdateUserPrizeStatus;
use App\Console\Commands\WxMall\GoodsSaleOutToPending;
use App\Console\Commands\WxMall\UpdateGoodsToWxMall;
use App\Console\Commands\WxMall\UpdateOrderWxMallStatusToCompleted;
use App\Console\Commands\WxMall\UpdateOrderWxMallStatusToDeleted;
use App\Console\Commands\WxMall\UpdateOrderWxMallStatusToInit;
use App\Console\Commands\WxMall\UpdateOrderWxMallStatusToPayed;
use App\Console\Commands\WxMall\UpdateOrderWxMallStatusToRefund;
use App\Console\Commands\WxMall\UpdateOrderWxMallStatusToSendOut;
use App\Console\Commands\UpdateRelationGroupData;
use Illuminate\Console\Scheduling\Schedule;
use Illuminate\Support\Facades\Log;
use Laravel\Lumen\Console\Kernel as ConsoleKernel;
use App\Console\Commands\ActivityCenter\Luck\LuckActivityEnd;
use App\Console\Commands\Goods\SynCdpGoodsSale;

class Kernel extends ConsoleKernel
{
    /**
     * The Artisan commands provided by your application.
     *
     * @var array
     */
    protected $commands = [
        'doc'                                    => Document::class,
        'test'                                    => Test::class,
    ];

    /**
     * Define the application's command schedule.
     *
     * @param \Illuminate\Console\Scheduling\Schedule $schedule
     *
     * @return void
     */
    protected function schedule(Schedule $schedule)
    {
        Log::debug('任务调度开始:' . date('Y-m-d H:i:s'));
    }
}
