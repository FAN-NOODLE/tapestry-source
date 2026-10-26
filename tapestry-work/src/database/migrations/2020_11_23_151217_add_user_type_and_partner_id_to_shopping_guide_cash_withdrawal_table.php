<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddUserTypeAndPartnerIdToShoppingGuideCashWithdrawalTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('shopping_guide_cash_withdrawal')) {
            if (!Schema::hasColumn('shopping_guide_cash_withdrawal', 'user_type')) {
                Schema::table('shopping_guide_cash_withdrawal',function (Blueprint $table) {
                    $table->tinyInteger('user_type')->default('1')->comment('身份 1导购 2合伙人');
                    $table->unsignedInteger('partner_id')->default('0')->comment('合伙人ID');
                });
            }
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('shopping_guide_cash_withdrawal', function (Blueprint $table) {
            $table->dropColumn('user_type');
            $table->dropColumn('partner_id');
        });
    }
}
