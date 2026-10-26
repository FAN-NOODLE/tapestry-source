export default {
  init: function () {
    const mapUrl = "https://map.qq.com/api/gljs?v=1.exp&key=AGFBZ-CM7L4-7F3UY-DY7NO-5P3KV-FGFIY&callback=onMapCallback";
    return new Promise((resolve) => {
      (window as any).onMapCallback = function (e) {
        console.log(e, 111);
        
        // resolve(qq.maps);
      };
      let scriptLabel = document.createElement("script");
      scriptLabel.setAttribute("type", "text/javascript");
      scriptLabel.setAttribute("src", mapUrl);
      document.body.appendChild(scriptLabel);
    });
  }
};