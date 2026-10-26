<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFailReasonToSubscribeMessageTouchTemplateTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasColumn('subscribe_message_touch_template', 'fail_reason')) {
            Schema::table('subscribe_message_touch_template', function (Blueprint $table) {
                $table->string('fail_reason')->default('')->after('pri_tmpl_id')->comment('失败原因');
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
        Schema::table('subscribe_message_touch_template', function (Blueprint $table) {
            $table->dropColumn(['fail_reason']);
        });
    }
}
