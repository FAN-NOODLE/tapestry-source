<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

/**
 * 添加user_id字段到partner_enterprise_map表
 * Created by PhpStorm.
 * User: Hui Yang
 * Email: 2361313833@qq.com
 * DateTime: 2020年12月15日 14:26:14
 * Class AddMemTypeInWithdrawalAutomaticReviewTable
 */
class AddUserIdToPartnerEnterpriseMapTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        //
        if (Schema::hasTable('partner_enterprise_map')) {
            Schema::table('partner_enterprise_map', function (Blueprint $table) {
                if (!Schema::hasColumn('partner_enterprise_map', 'user_id')) {
                    $table->bigInteger('user_id')
                        ->default(0)
                        ->after('worker_uid')
                        ->comment('合伙人的商城用户ID');
                }
            });
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
    }
}
