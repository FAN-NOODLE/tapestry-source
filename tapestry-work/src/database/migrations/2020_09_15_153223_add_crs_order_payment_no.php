<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddCrsOrderPaymentNo extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('crs_order')) {
            Schema::table('crs_order', function (Blueprint $table) {
                $table->string('third_payment_no', 255)->default('')->after('read_status')->comment('第三方支付流水号');
                $table->unsignedInteger('logistics_price')->default(0)->after('read_status')->comment('物流邮费');
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
