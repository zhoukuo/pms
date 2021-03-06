var hostname = "127.0.0.1:8088";

// 对Date的扩展，将 Date 转化为指定格式的String
// 月(M)、日(d)、小时(h)、分(m)、秒(s)、季度(q) 可以用 1-2 个占位符， 
// 年(y)可以用 1-4 个占位符，毫秒(S)只能用 1 个占位符(是 1-3 位的数字) 
// 例子： 
// (new Date()).Format("yyyy-MM-dd hh:mm:ss.S") ==> 2006-07-02 08:09:04.423 
// (new Date()).Format("yyyy-M-d h:m:s.S")      ==> 2006-7-2 8:9:4.18 
Date.prototype.Format = function(fmt) { //author: meizz 
	var o = {
		"M+": this.getMonth() + 1, //月份 
		"d+": this.getDate(), //日 
		"h+": this.getHours(), //小时 
		"m+": this.getMinutes(), //分 
		"s+": this.getSeconds(), //秒 
		"q+": Math.floor((this.getMonth() + 3) / 3), //季度 
		"S": this.getMilliseconds() //毫秒 
	};
	if(/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
	for(var k in o)
		if(new RegExp("(" + k + ")").test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
	return fmt;
}

var GetList = function(urlPare, callback, ErrorText) {
	$.ajax({
		url: "http://" + hostname + urlPare,
		type: "GET",
		dataType: "json",
		data: "",
		success: function(response) {
			if(response.Status != 200) {
				$("#ErrorText").text(response.StatusText + "(" + response.Status + ")");
				$("#ErrorMessage").slideDown("500").delay("2000").slideUp("500");
				return;
			}
			callback(response);
		},
		error: function(XmlHttpRequest, textStatus, errorThrown) {
			//调用失败
			$("#ErrorText").text("获取" + ErrorText + "信息失败，请联系管理员！")
			$("#ErrorMessage").slideDown("500").delay("2000").slideUp("500");
		}
	});
}

var PostCase = function(urlPare, jsonDate, callback, ErrorText) {
	$.ajax({
		url: "http://" + hostname + urlPare,
		type: "POST",
		dataType: "json",
		data: jsonDate,
		success: function(response) {
			if(response.Status != 200) {
				$("#ErrorText").text(response.StatusText + "(" + response.Status + ")");
				$("#ErrorMessage").slideDown("500").delay("2000").slideUp("500");
				return;
			}
			$("#SuccessMessage").slideDown("500").delay("2000").slideUp("500");
			callback(response);
		},
		error: function(XmlHttpRequest, textStatus, errorThrown) {
			//调用失败
			$("#ErrorText").text("提交" + ErrorText + "信息失败，请联系管理员！")
			$("#ErrorMessage").slideDown("500").delay("2000").slideUp("500");
		}
	});
}

var a = location.pathname;
var b = a.split("/");
var sleepFun = function() {
	if($.cookie("status") != "ok") {
		if(b[b.length - 1] != "login.html") {
			document.location = "../../login.html"
		}
		return;
	}
	$("span[name='uname']").text($.cookie("username"))
}

//setTimeout("sleepFun()", 1000);
//if($("span[name='uname']").length === 0) {
//	setTimeout("sleepFun()", 1000);
//	if($("span[name='uname']").length === 0) {
//		setTimeout("sleepFun()", 2000);
//	}
//}

//替换左侧、头部、底部html
if(b[b.length - 1] != "index.html") {
	$("#nav").load("../../pages/main/index.html #nav",function(){setTimeout("sleepFun()", 100)});
	$("#aside").load("../../pages/main/index.html #aside",function(){setTimeout("sleepFun()", 100)});
	$("#footer").load("../../pages/main/index.html #footer",function(){setTimeout("sleepFun()", 100)});
}else{
	setTimeout("sleepFun()", 100)
}
