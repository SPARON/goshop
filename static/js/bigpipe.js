var bigpipe = function () {
    var _self = this;
    var js_list = {};
    this.css_handler = function (css_url_array) {
        for (i = 0; i < css_url_array.length; i++) {
            jQuery("head").append("<link href='" + css_url_array[i] + "' rel='stylesheet' type='text/css' />");
        }
    }
    this.js_handler = function (js_url_array) {
        for (i = 0; i < js_url_array.length; i++) {
            jQuery("head").append("<script language='javascript' src='" + js_url_array[i] + "'></script>");
        }
    }
    this.html_handler = function (id, html) {
        html = html.replace(/\\"/g, '"').replace(/<\\\/script>/g, '<\/script>'); //数组转义会原本数据
        jQuery(id).html(html);
    }
    this.onPageletArrive = function (data) {
        if (!('id' in data)) {
            return;
        }
        var id = data.id;
        if ('html' in data) {
            _self.html_handler(id, data.html);
        }
        if ('css_url' in data) {
            _self.css_handler(data.css_url);
        }
        if ('js_url' in data) {
            _self.js_handler(data.js_url);
        }
    }
    return this;
}();