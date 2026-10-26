<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Document</title>
</head>
<body>
<br/>
<br/>
商城小程序授权
<form action="http://boss-api-dev.crs.dev-test.vchangyi.com/service/v1/auth/shop-auth-detail" method="post">
    企业id:<input name="enterprise_id" value="">
    <br/>
    回调地址:<input name="back_url"  value="">
    <br/>

    <input type="submit" value="提交">
</form>


<br/><br/>
企业微信小程序授权
<form action="http://boss-api-dev.crs.dev-test.vchangyi.com/service/v1/auth/work-auth-detail" method="post">
    企业id:<input name="enterprise_id" value="">
    <br/>
    回调地址:<input name="back_url"  value="">
    <br/>

    <input type="submit" value="提交">
</form>
</body>
</html>