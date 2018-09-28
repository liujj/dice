Date.prototype.Format = function (fmt) {
    var o = {
        "M+": this.getMonth() + 1, //月份
        "d+": this.getDate(), //日
        "h+": this.getHours(), //小时
        "m+": this.getMinutes(), //分
        "s+": this.getSeconds(), //秒
        "q+": Math.floor((this.getMonth() + 3) / 3), //季度
        "S": this.getMilliseconds() //毫秒
    };
    if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (const k in o)
        if (new RegExp("(" + k + ")").test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length === 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
    return fmt;
};

function getIPs(callback) {
    var ip_dups = {};
    var RTCPeerConnection = window.RTCPeerConnection || window.mozRTCPeerConnection || window.webkitRTCPeerConnection;
    var mediaConstraints = {
        optional: [{RtpDataChannels: true}]
    };
    var servers = undefined;
    var pc = new RTCPeerConnection(servers, mediaConstraints);
    pc.onicecandidate = function (ice) {
        if (ice.candidate) {
            var ip_regex = /([0-9]{1,3}(\.[0-9]{1,3}){3})/;
            var ip_addr = ip_regex.exec(ice.candidate.candidate)[1];
            if (ip_dups[ip_addr] === undefined)
                callback(ip_addr);
            ip_dups[ip_addr] = true;
        }
    };
    pc.createDataChannel("");
    pc.createOffer(function (result) {
        pc.setLocalDescription(result, function () {
        }, function () {
        });
    }, function () {
    });
}

Promise.prototype.finally = function (callback) {
    var Promise = this.constructor;
    return this.then(
        function (value) {
            Promise.resolve(callback()).then(
                function () {
                    return value;
                }
            );
        },
        function (reason) {
            Promise.resolve(callback()).then(
                function () {
                    throw reason;
                }
            );
        }
    );
};