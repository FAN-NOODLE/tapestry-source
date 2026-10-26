<?php
/**
 * Created by sublime.
 * User: zhao qi
 * Date: 2019/01/17
 * Time: 10:00
 */

namespace App\Http\Controllers\V1\Label;

use App\Http\Controllers\V1\BaseController;

abstract class LabelBaseController extends BaseController
{
    /**
     * 标签状态 1: 启用
     */
    const LABEL_STATUS_OPEN = 1;

    /**
     * 标签状态 2: 禁用
     */
    const LABEL_STATUS_CLOSE = 2;
}

