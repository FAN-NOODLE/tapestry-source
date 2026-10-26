package shorturl

var ShortHtml = `<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no, minimum-scale=1.0, maximum-scale=1.0">
  <title>%v</title>
  <meta name="description" content="%v" />
</head>
<style>
    a{-webkit-tap-highlight-color:transparent;}
    body,a,div,img {
      margin: 0;
      padding: 0;
    }
    #middle-page{
      width: 100%%;
      /* padding: 20px 15px; */
      box-sizing: border-box;
    }
    
    #middle-page .middle-page-box{
    }
    #middle-page .btn-box{
      margin-top: -15px;
      padding: 30px 0 15px 0;
      position: fixed;
      width: 100%%;
      bottom: 0;
      z-index: 2;
      background-color: #181818;
    }
    #middle-page .open-btn{
      width: 92%%;
      display: flex;
      height: 44px;
      background: #FFFFFE;
      border-radius: 4px;
      text-align: center;
      align-items:center;
      justify-content:center;
      font-size: 15px;
      color: #708DF5;
      margin: 0 auto 10px auto;
    }
    #middle-page .decoration{
      color: white;
      text-align: center;
      font-size: 12px;
    }
    #middle-page .product-introduce{
      width: 100%%;
      position: fixed;
      overflow: hidden;
      bottom: 100px;
      z-index: 4;
      /* height: 130px; */
    }
    #middle-page .product-content {
      padding-bottom: 100px;
      white-space: nowrap;
      height: 130px;
      overflow: hidden;
      overflow-x: scroll; /* 1 */
      -webkit-backface-visibility: hidden;
      -webkit-perspective: 1000;
      -webkit-overflow-scrolling: touch; /* 2 */
      text-align: justify; /* 3 */
    }
    -webkit-scrollbar {
      display: none;
    }
    #middle-page .introduce-box{
      overflow: hidden;
      padding: 15px;
      box-sizing: border-box;
      background-color: #fff;
      border-radius: 10px;
      display: inline-block;
      margin-right: 10px;
      width: 300px;
    }
    #middle-page .introduce-box:first-child {
      margin-left: 15px;
    }
    #middle-page .introduce-box:last-child {
      margin-right: 15px;
    }
    .introduce-box-one {
      width: 92%% !important;
      margin: 0 auto;
    }
    #middle-page .introduce-box .image{
      float: left;
      width: 100px;
      height: 100px;
    }
    #middle-page .introduce-box .image img{
      display: block;
      width: 100px;
      height: 100px;
    }
    #middle-page .introduce-box .goods-info{
      float: left;
      width: calc( 100%% - 110px );
      margin-left: 10px;
    }
    #middle-page .introduce-box .goods-info .spu-name{
      font-size: 14px;
      font-weight: bold;
      color: #333;
      text-overflow: -o-ellipsis-lastline;
      white-space: break-spaces;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
      line-height: 20px;
    }
    #middle-page .introduce-box .goods-info .product-sku{
      width: 100%%;
      font-size: 12px;
      color: #A6A6A6;
      box-sizing: border-box;
      line-height: 24px;
      margin:5px 0 12px 0;
      position: relative;
      padding-right: 20px;
      box-sizing: border-box;
      overflow: hidden;
      height: 24px;
      text-overflow:ellipsis;
      white-space: nowrap;
    }
	.price-wrap {
      line-height: 16px;
      color: #fa4046;
    }
    #middle-page .introduce-box .goods-info .price{
      font-size: 22px;
    }
    .price-icon {
      font-size: 14px;
    }
    .price-point {
      font-size: 14px;
    }
    #middle-page .banner-image{
      height: auto;
      border-radius: 10px;
      width: 92%%;
      left: 4%%;
      position: absolute;
      bottom: 20px;
      z-index: 4;
    }
    #middle-page .clickBtnA {
      position: fixed;
      top: 0;
      left: 0;
      width: 100%%;
      height: 100%%;
      z-index: 3;
    }
  </style>
<body>
<div id="middle-page">
 <img src='%v' alt="" style="display: none;">
</div>
</body>
<script>
var configData = {
  brand_img: '%v', // 品牌图片
  background_img: '%v', // 背景图片
  button_color: '%v', // 按钮的背景颜色
  button_name: '%v', // 按钮的文案
  button_font_color: '%v', // 按钮的字体颜色
  location_desc: '%v', // 跳转文案
  location_font_color: '%v', // 跳转文案的字体颜色
  is_hide: '%v', // 是否显示跳转按钮
  active_img: '%v', // 如果为空则表示没有banner /src/assets/images/login/login-bac.jpg
  active_page_url: '%v', // banner的详情页链接
  wx_active_page_url: '%v',
  spu_data: %v,
  wx_url: '%v', // 跳转链接
}
</script>
<script>
    // 根据参数来渲染页面
    if (configData.wx_url !== '') {
      window.location.href = configData.wx_url;
    }

    const middlePage = document.querySelector('#middle-page');
    const middlePageBox = document.querySelector('.middle-page-box');
    const nowElementA = document.createElement('a');
    nowElementA.classList.add('clickBtnA');
    nowElementA.href = configData.wx_url;
    middlePage.appendChild(nowElementA);
    const btnBox = document.querySelector('#middle-page .btn-box');
    if (configData.background_img) {
      const imgContent = document.createElement('div');
      const imgElement = document.createElement('img');
      imgElement.src = configData.background_img;
      imgContent.classList.add("img-temp");
      imgContent.setAttribute("style", "width: 100%%;position: fixed;overflow: auto;height: calc(100%% - 100px);border-radius: 0 0 14px 14px;z-index: 3;")
      imgElement.setAttribute("style","width: 100%%;");
      imgContent.appendChild(imgElement);
      middlePage.appendChild(imgContent);
    }
    // middlePage.style.background = 'url('+ configData.background_img +') no-repeat center';
    // middlePage.style.backgroundSize = '100%% auto';
    
    if (configData.active_img !== '') {
      const imgElement = document.createElement('img');
      imgElement.classList.add('banner-image');
      imgElement.src = configData.active_img;
      middlePage.appendChild(imgElement);
    }
    // if ((parseInt(configData.is_hide) === 1) && (configData.spu_data.length === 0)) {
    //   middlePageBox.style.display = 'none';
    // }
    let insertgoodInfo = '';
    if (configData.spu_data.length > 0) {
      const spuDataBox = configData.spu_data;
      spuDataBox.forEach(val => {
        insertgoodInfo += '<div class="' + (configData.spu_data.length == 1 ? 'introduce-box introduce-box-one' : 'introduce-box') + '"' + 'link-url='+ val.wx_spu_page_url + '>'+
        '<div class="image">'+
          '<img src="'+ val.spu_img +'" alt="">'+
        '</div>'+
        '<div class="goods-info">'+
          '<div class="spu-name">'+ val.spu_name +'</div>'+
          '<div class="product-sku">'+ val.sku_attribute +'</div>'+
          '<div class="price-wrap">' +
            '<span class="price-icon">¥</span>' +
            '<span class="price">'+ (val.sku_price + '').split('.')[0] + '</span>'+
            '<span class="price-point">' + ((val.sku_price + '').indexOf('.') != -1 ? '.' + (val.sku_price + '').split('.')[1] : '') + '</span>' +
          '</div>' +
        '</div>'+
      '</div>'
      })
    }
    const nowElement = document.createElement('div');
    const contentEle = document.createElement('div');
    const TemEle = document.createElement('div');
    contentEle.classList.add('product-content');
    nowElement.classList.add('product-introduce');
    TemEle.setAttribute('style', 'height: 130px; overflow: hidden;padding-bottom: 15px;border-radius: 0 0 14px 14px;z-index: 3;position: relative;');
    contentEle.innerHTML = insertgoodInfo;
    nowElement.appendChild(TemEle);
    TemEle.appendChild(contentEle)
    middlePage.appendChild(nowElement)
    
    if (parseInt(configData.is_hide) !== 1) {
      let btnBoxInfo = '<div class="open-btn">'
        + configData.button_name +
      '</div>'+
      '<div class="decoration">'+configData.location_desc+'</div>';
      const nowBtnElement = document.createElement('div');
      nowBtnElement.classList.add('btn-box');
      nowBtnElement.innerHTML = btnBoxInfo;
      middlePage.appendChild(nowBtnElement);
      const middleOpenBtn = document.querySelector('#middle-page .open-btn');
      const location_font_color = document.querySelector('#middle-page .decoration');
      middleOpenBtn.style.backgroundColor = configData.button_color;
      middleOpenBtn.style.color = configData.button_font_color;
      location_font_color.style.color = configData.location_font_color;
      const banner = document.querySelector('.banner-image');
      if (banner) banner.style.bottom = '115px';
    } else {
      const imgTemp = document.querySelector('.img-temp');
      imgTemp.style.height = '100%%';
      imgTemp.style.borderRadius = 0;
      nowElement.style.bottom = 0;
    }
    // 点击跳转
    middlePage.onclick = (event) => {
      const productIntroduce = document.querySelector('#middle-page .product-introduce');
      const bannerImage = document.querySelector('#middle-page .banner-image');
      if (productIntroduce && productIntroduce.contains(event.target)) {
        if (findEle(event.target) !== "") {
          window.location.href = findEle(event.target);
        }
      } else if (bannerImage && bannerImage.contains(event.target)) {
        if (configData.wx_active_page_url !== "") {
          window.location.href = configData.wx_active_page_url;
        }
      } else {
        if (configData.wx_url !== "") {
          window.location.href = configData.wx_url;
        }
      }
    }
    // 递归查找元素的属性。target为当前点击的元素
    function findEle (target) {
      if (target.getAttribute('link-url') !== null) return target.getAttribute('link-url');
      return findEle(target.parentNode);
    }
    </script>
</html>`
