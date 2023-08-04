var clearSelect, DeviceInfo;
if ("undefined" == typeof jQuery) throw new Error("Bootstrap's JavaScript requires jQuery");
DeviceInfo = function () {
    var r = typeof self != "undefined" ? self : this, i = r || {}, t = {
        navigator: typeof r.navigator != "undefined" ? r.navigator : {}, infoMap: {
            engine: ["WebKit", "Trident", "Gecko", "Presto"],
            browser: ["Safari", "Chrome", "Edge", "IE", "Firefox", "Firefox Focus", "Chromium", "Opera", "Vivaldi", "Yandex", "Arora", "Lunascape", "QupZilla", "Coc Coc", "Kindle", "Iceweasel", "Konqueror", "Iceape", "SeaMonkey", "Epiphany", "360", "360SE", "360EE", "UC", "QQBrowser", "QQ", "Baidu", "Maxthon", "Sogou", "LBBROWSER", "2345Explorer", "TheWorld", "XiaoMi", "Quark", "Qiyu", "Wechat", "Taobao", "Alipay", "Weibo", "Douban", "Suning", "iQiYi"],
            os: ["Windows", "Linux", "Mac OS", "Android", "Ubuntu", "FreeBSD", "Debian", "iOS", "Windows Phone", "BlackBerry", "MeeGo", "Symbian", "Chrome OS", "WebOS"],
            device: ["Mobile", "Tablet", "iPad"]
        }
    }, n = function () {
        return {
            getMatchMap: function (n) {
                return {
                    Trident: n.indexOf("Trident") > -1 || n.indexOf("NET CLR") > -1,
                    Presto: n.indexOf("Presto") > -1,
                    WebKit: n.indexOf("AppleWebKit") > -1,
                    Gecko: n.indexOf("Gecko/") > -1,
                    Safari: n.indexOf("Safari") > -1,
                    Chrome: n.indexOf("Chrome") > -1 || n.indexOf("CriOS") > -1,
                    IE: n.indexOf("MSIE") > -1 || n.indexOf("Trident") > -1,
                    Edge: n.indexOf("Edge") > -1,
                    Firefox: n.indexOf("Firefox") > -1 || n.indexOf("FxiOS") > -1,
                    "Firefox Focus": n.indexOf("Focus") > -1,
                    Chromium: n.indexOf("Chromium") > -1,
                    Opera: n.indexOf("Opera") > -1 || n.indexOf("OPR") > -1,
                    Vivaldi: n.indexOf("Vivaldi") > -1,
                    Yandex: n.indexOf("YaBrowser") > -1,
                    Arora: n.indexOf("Arora") > -1,
                    Lunascape: n.indexOf("Lunascape") > -1,
                    QupZilla: n.indexOf("QupZilla") > -1,
                    "Coc Coc": n.indexOf("coc_coc_browser") > -1,
                    Kindle: n.indexOf("Kindle") > -1 || n.indexOf("Silk/") > -1,
                    Iceweasel: n.indexOf("Iceweasel") > -1,
                    Konqueror: n.indexOf("Konqueror") > -1,
                    Iceape: n.indexOf("Iceape") > -1,
                    SeaMonkey: n.indexOf("SeaMonkey") > -1,
                    Epiphany: n.indexOf("Epiphany") > -1,
                    "360": n.indexOf("QihooBrowser") > -1 || n.indexOf("QHBrowser") > -1,
                    "360EE": n.indexOf("360EE") > -1,
                    "360SE": n.indexOf("360SE") > -1,
                    UC: n.indexOf("UC") > -1 || n.indexOf(" UBrowser") > -1,
                    QQBrowser: n.indexOf("QQBrowser") > -1,
                    QQ: n.indexOf("QQ/") > -1,
                    Baidu: n.indexOf("Baidu") > -1 || n.indexOf("BIDUBrowser") > -1,
                    Maxthon: n.indexOf("Maxthon") > -1,
                    Sogou: n.indexOf("MetaSr") > -1 || n.indexOf("Sogou") > -1,
                    LBBROWSER: n.indexOf("LBBROWSER") > -1,
                    "2345Explorer": n.indexOf("2345Explorer") > -1,
                    TheWorld: n.indexOf("TheWorld") > -1,
                    XiaoMi: n.indexOf("MiuiBrowser") > -1,
                    Quark: n.indexOf("Quark") > -1,
                    Qiyu: n.indexOf("Qiyu") > -1,
                    Wechat: n.indexOf("MicroMessenger") > -1,
                    Taobao: n.indexOf("AliApp(TB") > -1,
                    Alipay: n.indexOf("AliApp(AP") > -1,
                    Weibo: n.indexOf("Weibo") > -1,
                    Douban: n.indexOf("com.douban.frodo") > -1,
                    Suning: n.indexOf("SNEBUY-APP") > -1,
                    iQiYi: n.indexOf("IqiyiApp") > -1,
                    Windows: n.indexOf("Windows") > -1,
                    Linux: n.indexOf("Linux") > -1 || n.indexOf("X11") > -1,
                    "Mac OS": n.indexOf("Macintosh") > -1,
                    Android: n.indexOf("Android") > -1 || n.indexOf("Adr") > -1,
                    Ubuntu: n.indexOf("Ubuntu") > -1,
                    FreeBSD: n.indexOf("FreeBSD") > -1,
                    Debian: n.indexOf("Debian") > -1,
                    "Windows Phone": n.indexOf("IEMobile") > -1 || n.indexOf("Windows Phone") > -1,
                    BlackBerry: n.indexOf("BlackBerry") > -1 || n.indexOf("RIM") > -1,
                    MeeGo: n.indexOf("MeeGo") > -1,
                    Symbian: n.indexOf("Symbian") > -1,
                    iOS: n.indexOf("like Mac OS X") > -1,
                    "Chrome OS": n.indexOf("CrOS") > -1,
                    WebOS: n.indexOf("hpwOS") > -1,
                    Mobile: n.indexOf("Mobi") > -1 || n.indexOf("iPh") > -1 || n.indexOf("480") > -1,
                    Tablet: n.indexOf("Tablet") > -1 || n.indexOf("Nexus 7") > -1,
                    iPad: n.indexOf("iPad") > -1
                }
            }, matchInfoMap: function (i) {
                var e = t.navigator.userAgent || {}, o = n.getMatchMap(e), r, u, f;
                for (r in t.infoMap) for (u = 0; u < t.infoMap[r].length; u++) f = t.infoMap[r][u], o[f] && (i[r] = f)
            }, getOS: function () {
                var t = this;
                return n.matchInfoMap(t), t.os
            }, getOSVersion: function () {
                var i = this, n = t.navigator.userAgent || {}, r;
                return i.osVersion = "", r = {
                    Windows: function () {
                        var t = n.replace(/^.*Windows NT ([\d.]+);.*$/, "$1");
                        return {
                            "6.4": "10",
                            "6.3": "8.1",
                            "6.2": "8",
                            "6.1": "7",
                            "6.0": "Vista",
                            "5.2": "XP",
                            "5.1": "XP",
                            "5.0": "2000"
                        }[t] || t
                    }, Android: function () {
                        return n.replace(/^.*Android ([\d.]+);.*$/, "$1")
                    }, iOS: function () {
                        return n.replace(/^.*OS ([\d_]+) like.*$/, "$1").replace(/_/g, ".")
                    }, Debian: function () {
                        return n.replace(/^.*Debian\/([\d.]+).*$/, "$1")
                    }, "Windows Phone": function () {
                        return n.replace(/^.*Windows Phone( OS)? ([\d.]+);.*$/, "$2")
                    }, "Mac OS": function () {
                        return n.replace(/^.*Mac OS X ([\d_]+).*$/, "$1").replace(/_/g, ".")
                    }, WebOS: function () {
                        return n.replace(/^.*hpwOS\/([\d.]+);.*$/, "$1")
                    }
                }, r[i.os] && (i.osVersion = r[i.os](), i.osVersion == n && (i.osVersion = "")), i.osVersion
            }, GetOSBit: function () {
                return navigator.userAgent.indexOf("x64") > 0 ? "64位" : ""
            }, getOrientationStatu: function () {
                var n = window.matchMedia("(orientation: portrait)");
                return n.matches ? "竖屏" : "横屏"
            }, getDeviceType: function () {
                var t = this;
                return t.device = "PC", n.matchInfoMap(t), t.device
            }, getNetwork: function () {
                return navigator && navigator.connection && navigator.connection.effectiveType
            }, getLanguage: function () {
                var n = this;
                return n.language = function () {
                    var i = t.navigator.browserLanguage || t.navigator.language, n = i.split("-");
                    return n[1] && (n[1] = n[1].toUpperCase()), n.join("_")
                }(), n.language
            }, getBrowserInfo: function () {
                var u = this, o, c, s;
                n.matchInfoMap(u);
                var r = t.navigator.userAgent || {}, h = function (n, i) {
                    var r = t.navigator.mimeTypes;
                    for (var u in r) if (r[u][n] == i) return !0;
                    return !1
                }, f = n.getMatchMap(r), e = !1;
                if (i.chrome && (o = r.replace(/^.*Chrome\/([\d]+).*$/, "$1"), o > 36 && i.showModalDialog ? e = !0 : o > 45 && (e = h("type", "application/vnd.chromium.remoting-viewer"))), f.Baidu && f.Opera && (f.Baidu = !1), f.Mobile && (f.Mobile = !(r.indexOf("iPad") > -1)), e && (h("type", "application/gameplugin") ? f["360SE"] = !0 : t.navigator && typeof t.navigator.connection.saveData == "undefined" ? f["360SE"] = !0 : f["360EE"] = !0), f.IE || f.Edge) {
                    c = window.screenTop - window.screenY;
                    switch (c) {
                        case 102:
                            f["360EE"] = !0;
                            break;
                        case 104:
                            f["360SE"] = !0
                    }
                }
                return s = {
                    Safari: function () {
                        return r.replace(/^.*Version\/([\d.]+).*$/, "$1")
                    }, Chrome: function () {
                        return r.replace(/^.*Chrome\/([\d.]+).*$/, "$1").replace(/^.*CriOS\/([\d.]+).*$/, "$1")
                    }, IE: function () {
                        return r.replace(/^.*MSIE ([\d.]+).*$/, "$1").replace(/^.*rv:([\d.]+).*$/, "$1")
                    }, Edge: function () {
                        return r.replace(/^.*Edge\/([\d.]+).*$/, "$1")
                    }, Firefox: function () {
                        return r.replace(/^.*Firefox\/([\d.]+).*$/, "$1").replace(/^.*FxiOS\/([\d.]+).*$/, "$1")
                    }, "Firefox Focus": function () {
                        return r.replace(/^.*Focus\/([\d.]+).*$/, "$1")
                    }, Chromium: function () {
                        return r.replace(/^.*Chromium\/([\d.]+).*$/, "$1")
                    }, Opera: function () {
                        return r.replace(/^.*Opera\/([\d.]+).*$/, "$1").replace(/^.*OPR\/([\d.]+).*$/, "$1")
                    }, Vivaldi: function () {
                        return r.replace(/^.*Vivaldi\/([\d.]+).*$/, "$1")
                    }, Yandex: function () {
                        return r.replace(/^.*YaBrowser\/([\d.]+).*$/, "$1")
                    }, Arora: function () {
                        return r.replace(/^.*Arora\/([\d.]+).*$/, "$1")
                    }, Lunascape: function () {
                        return r.replace(/^.*Lunascape[\/\s]([\d.]+).*$/, "$1")
                    }, QupZilla: function () {
                        return r.replace(/^.*QupZilla[\/\s]([\d.]+).*$/, "$1")
                    }, "Coc Coc": function () {
                        return r.replace(/^.*coc_coc_browser\/([\d.]+).*$/, "$1")
                    }, Kindle: function () {
                        return r.replace(/^.*Version\/([\d.]+).*$/, "$1")
                    }, Iceweasel: function () {
                        return r.replace(/^.*Iceweasel\/([\d.]+).*$/, "$1")
                    }, Konqueror: function () {
                        return r.replace(/^.*Konqueror\/([\d.]+).*$/, "$1")
                    }, Iceape: function () {
                        return r.replace(/^.*Iceape\/([\d.]+).*$/, "$1")
                    }, SeaMonkey: function () {
                        return r.replace(/^.*SeaMonkey\/([\d.]+).*$/, "$1")
                    }, Epiphany: function () {
                        return r.replace(/^.*Epiphany\/([\d.]+).*$/, "$1")
                    }, "360": function () {
                        return r.replace(/^.*QihooBrowser\/([\d.]+).*$/, "$1")
                    }, "360SE": function () {
                        var n = r.replace(/^.*Chrome\/([\d]+).*$/, "$1");
                        return {
                            "63": "10.0", "55": "9.1", "45": "8.1", "42": "8.0", "31": "7.0", "21": "6.3"
                        }[n] || ""
                    }, "360EE": function () {
                        var n = r.replace(/^.*Chrome\/([\d]+).*$/, "$1");
                        return {
                            "69": "11.0", "63": "9.5", "55": "9.0", "50": "8.7", "30": "7.5"
                        }[n] || ""
                    }, Maxthon: function () {
                        return r.replace(/^.*Maxthon\/([\d.]+).*$/, "$1")
                    }, QQBrowser: function () {
                        return r.replace(/^.*QQBrowser\/([\d.]+).*$/, "$1")
                    }, QQ: function () {
                        return r.replace(/^.*QQ\/([\d.]+).*$/, "$1")
                    }, Baidu: function () {
                        return r.replace(/^.*BIDUBrowser[\s\/]([\d.]+).*$/, "$1")
                    }, UC: function () {
                        return r.replace(/^.*UC?Browser\/([\d.]+).*$/, "$1")
                    }, Sogou: function () {
                        return r.replace(/^.*SE ([\d.X]+).*$/, "$1").replace(/^.*SogouMobileBrowser\/([\d.]+).*$/, "$1")
                    }, LBBROWSER: function () {
                        var n = navigator.userAgent.replace(/^.*Chrome\/([\d]+).*$/, "$1");
                        return {
                            "57": "6.5",
                            "49": "6.0",
                            "46": "5.9",
                            "42": "5.3",
                            "39": "5.2",
                            "34": "5.0",
                            "29": "4.5",
                            "21": "4.0"
                        }[n] || ""
                    }, "2345Explorer": function () {
                        return r.replace(/^.*2345Explorer\/([\d.]+).*$/, "$1")
                    }, TheWorld: function () {
                        return r.replace(/^.*TheWorld ([\d.]+).*$/, "$1")
                    }, XiaoMi: function () {
                        return r.replace(/^.*MiuiBrowser\/([\d.]+).*$/, "$1")
                    }, Quark: function () {
                        return r.replace(/^.*Quark\/([\d.]+).*$/, "$1")
                    }, Qiyu: function () {
                        return r.replace(/^.*Qiyu\/([\d.]+).*$/, "$1")
                    }, Wechat: function () {
                        return r.replace(/^.*MicroMessenger\/([\d.]+).*$/, "$1")
                    }, Taobao: function () {
                        return r.replace(/^.*AliApp\(TB\/([\d.]+).*$/, "$1")
                    }, Alipay: function () {
                        return r.replace(/^.*AliApp\(AP\/([\d.]+).*$/, "$1")
                    }, Weibo: function () {
                        return r.replace(/^.*weibo__([\d.]+).*$/, "$1")
                    }, Douban: function () {
                        return r.replace(/^.*com.douban.frodo\/([\d.]+).*$/, "$1")
                    }, Suning: function () {
                        return r.replace(/^.*SNEBUY-APP([\d.]+).*$/, "$1")
                    }, iQiYi: function () {
                        return r.replace(/^.*IqiyiVersion\/([\d.]+).*$/, "$1")
                    }
                }, u.browserVersion = "", s[u.browser] && (u.browserVersion = s[u.browser](), u.browserVersion == r && (u.browserVersion = "")), u.browser == "Edge" && (u.engine = "EdgeHTML"), u.browser == "Chrome" && parseInt(u.browserVersion) > 27 && (u.engine = "Blink"), u.browser == "Opera" && parseInt(u.browserVersion) > 12 && (u.engine = "Blink"), u.browser == "Yandex" && (u.engine = "Blink"), {
                    Name: u.browser, Version: u.browserVersion, CoreType: u.engine
                }
            }
        }
    }(), u = function () {
        return {
            DeviceInfoObj: function () {
                return {
                    deviceType: n.getDeviceType(),
                    OS: {
                        Name: n.getOS(), Version: n.getOSVersion(), Bit: n.GetOSBit(), toString: function () {
                            return n.getOS() + " " + n.getOSVersion() + " " + n.GetOSBit()
                        }
                    },
                    screenHeight: i.screen.height,
                    screenWidth: i.screen.width,
                    language: n.getLanguage(),
                    netWork: n.getNetwork(),
                    orientation: n.getOrientationStatu(),
                    browserInfo: n.getBrowserInfo(),
                    userAgent: t.navigator.userAgent
                }
            }
        }
    }();
    return u.DeviceInfoObj()
}();
typeof Object.create != "function" && (Object.create = function (n) {
    function t() {
    }

    return t.prototype = n, new t
}), function (n, t, i, r) {
    var u = {
        init: function (t, i) {
            var r = this;
            if (r.elem = i, r.$elem = n(i), r.newsTagName = r.$elem.find(":first-child").prop("tagName"), r.newsClassName = r.$elem.find(":first-child").attr("class"), r.timer = null, r.resizeTimer = null, r.animationStarted = !1, r.isHovered = !1, typeof t == "string") {
                console && console.error("String property override is not supported");
                throw "String property override is not supported";
            } else r.options = n.extend({}, n.fn.bootstrapNews.options, t), r.prepareLayout(), r.options.autoplay && r.animate(), r.options.navigation && r.buildNavigation(), typeof r.options.onToDo == "function" && r.options.onToDo.apply(r, arguments)
        }, prepareLayout: function () {
            var i = this, r;
            n(i.elem).find("." + i.newsClassName).on("mouseenter", function () {
                i.onReset(!0)
            });
            n(i.elem).find("." + i.newsClassName).on("mouseout", function () {
                i.onReset(!1)
            });
            n.map(i.$elem.find(i.newsTagName), function (t, r) {
                r > i.options.newsPerPage - 1 ? n(t).hide() : n(t).show()
            });
            i.$elem.find(i.newsTagName).length < i.options.newsPerPage && (i.options.newsPerPage = i.$elem.find(i.newsTagName).length);
            r = 0;
            n.map(i.$elem.find(i.newsTagName), function (t, u) {
                u < i.options.newsPerPage && (r = parseInt(r) + parseInt(n(t).height()) + 10)
            });
            n(i.elem).css({
                "overflow-y": "hidden", height: r
            });
            n(t).resize(function () {
                i.resizeTimer !== null && clearTimeout(i.resizeTimer);
                i.resizeTimer = setTimeout(function () {
                    i.prepareLayout()
                }, 200)
            })
        }, findPanelObject: function () {
            for (var n = this.$elem; n.parent() !== r;) if (n = n.parent(), n.parent().hasClass("panel")) return n.parent();
            return r
        }, buildNavigation: function () {
            var t = this.findPanelObject(), i, r, u;
            if (t) {
                i = '<ul class="pagination pull-right" style="margin: 0px;"><li><a href="#" class="prev"><span class="glyphicon glyphicon-chevron-down"><\/span><\/a><\/li><li><a href="#" class="next"><span class="glyphicon glyphicon-chevron-up"><\/span><\/a><\/li><\/ul><div class="clearfix"><\/div>';
                r = n(t).find(".panel-footer")[0];
                r ? n(r).append(i) : n(t).append('<div class="panel-footer">' + i + "<\/div>");
                u = this;
                n(t).find(".prev").on("click", function (n) {
                    n.preventDefault();
                    u.onPrev()
                });
                n(t).find(".next").on("click", function (n) {
                    n.preventDefault();
                    u.onNext()
                })
            }
        }, onStop: function () {
        }, onPause: function () {
            var n = this;
            n.isHovered = !0;
            this.options.autoplay && n.timer && clearTimeout(n.timer)
        }, onReset: function (n) {
            var t = this;
            t.timer && clearTimeout(t.timer);
            t.options.autoplay && (t.isHovered = n, t.animate())
        }, animate: function () {
            var n = this;
            n.timer = setTimeout(function () {
                n.options.pauseOnHover || (n.isHovered = !1);
                n.isHovered || (n.options.direction === "up" ? n.onNext() : n.onPrev())
            }, n.options.newsTickerInterval)
        }, onPrev: function () {
            var t = this, i;
            if (t.animationStarted) return !1;
            t.animationStarted = !0;
            i = "<" + t.newsTagName + ' style="display:none;" class="' + t.newsClassName + '">' + n(t.$elem).find(t.newsTagName).last().html() + "<\/" + t.newsTagName + ">";
            n(t.$elem).prepend(i);
            n(t.$elem).find(t.newsTagName).first().slideDown(t.options.animationSpeed, function () {
                n(t.$elem).find(t.newsTagName).last().remove()
            });
            n(t.$elem).find(t.newsTagName + ":nth-child(" + parseInt(t.options.newsPerPage + 1) + ")").slideUp(t.options.animationSpeed, function () {
                t.animationStarted = !1;
                t.onReset(t.isHovered)
            });
            n(t.elem).find("." + t.newsClassName).on("mouseenter", function () {
                t.onReset(!0)
            });
            n(t.elem).find("." + t.newsClassName).on("mouseout", function () {
                t.onReset(!1)
            })
        }, onNext: function () {
            var t = this, i;
            if (t.animationStarted) return !1;
            t.animationStarted = !0;
            i = "<" + t.newsTagName + ' style="display:none;" class=' + t.newsClassName + ">" + n(t.$elem).find(t.newsTagName).first().html() + "<\/" + t.newsTagName + ">";
            n(t.$elem).append(i);
            n(t.$elem).find(t.newsTagName).first().slideUp(t.options.animationSpeed, function () {
                n(this).remove()
            });
            n(t.$elem).find(t.newsTagName + ":nth-child(" + parseInt(t.options.newsPerPage + 1) + ")").slideDown(t.options.animationSpeed, function () {
                t.animationStarted = !1;
                t.onReset(t.isHovered)
            });
            n(t.elem).find("." + t.newsClassName).on("mouseenter", function () {
                t.onReset(!0)
            });
            n(t.elem).find("." + t.newsClassName).on("mouseout", function () {
                t.onReset(!1)
            })
        }
    };
    n.fn.bootstrapNews = function (n) {
        return this.each(function () {
            var t = Object.create(u);
            t.init(n, this)
        })
    };
    n.fn.bootstrapNews.options = {
        newsPerPage: 4,
        navigation: !0,
        autoplay: !0,
        direction: "up",
        animationSpeed: "normal",
        newsTickerInterval: 4e3,
        pauseOnHover: !0,
        onStop: null,
        onPause: null,
        onReset: null,
        onPrev: null,
        onNext: null,
        onToDo: null
    }
}(jQuery, window, document);
var radius = 90, d = 200, dtr = Math.PI / 180, mcList = [], lasta = 1, lastb = 1, distr = !0, tspeed = 11, size = 200,
    mouseX = 0, mouseY = 10, howElliptical = 1, aA = null, oDiv = null;
window.onload = function () {
    var n = 0, t = null;
    if (oDiv = document.getElementById("tagscloud"), oDiv) {
        for (aA = oDiv.getElementsByTagName("a"), n = 0; n < aA.length; n++) t = {}, aA[n].onmouseover = function (n) {
            return function () {
                n.on = !0;
                this.style.zIndex = 9999;
                this.style.color = "#fff";
                this.style.padding = "5px 5px";
                this.style.filter = "alpha(opacity=100)";
                this.style.opacity = 1
            }
        }(t), aA[n].onmouseout = function (n) {
            return function () {
                n.on = !1;
                this.style.zIndex = n.zIndex;
                this.style.color = "#fff";
                this.style.padding = "5px";
                this.style.filter = "alpha(opacity=" + 100 * n.alpha + ")";
                this.style.opacity = n.alpha;
                this.style.zIndex = n.zIndex
            }
        }(t), t.offsetWidth = aA[n].offsetWidth, t.offsetHeight = aA[n].offsetHeight, mcList.push(t);
        sineCosine(0, 0, 0);
        positionAll(), function () {
            update();
            setTimeout(arguments.callee, 40)
        }()
    }
};
$(function () {
    var n = $("#rocket-to-top"), r = $(document).scrollTop(), i, t = !0;
    $(window).scroll(function () {
        var i = $(document).scrollTop();
        i == 0 ? n.css("background-position") == "0px 0px" ? n.fadeOut("slow") : t && (t = !1, $(".level-2").css("opacity", 1), n.delay(100).animate({
            marginTop: "-1000px"
        }, "normal", function () {
            n.css({
                "margin-top": "-125px", display: "none"
            });
            t = !0
        })) : n.fadeIn("slow")
    });
    n.hover(function () {
        $(".level-2").stop(!0).animate({
            opacity: 1
        })
    }, function () {
        $(".level-2").stop(!0).animate({
            opacity: 0
        })
    });
    $(".level-3").click(function () {
        function r() {
            var r = n.css("background-position");
            if (n.css("display") == "none" || t == 0) {
                clearInterval(i);
                n.css("background-position", "0px 0px");
                return
            }
            switch (r) {
                case "0px 0px":
                    n.css("background-position", "-298px 0px");
                    break;
                case "-298px 0px":
                    n.css("background-position", "-447px 0px");
                    break;
                case "-447px 0px":
                    n.css("background-position", "-596px 0px");
                    break;
                case "-596px 0px":
                    n.css("background-position", "-745px 0px");
                    break;
                case "-745px 0px":
                    n.css("background-position", "-298px 0px")
            }
        }

        t && (i = setInterval(r, 50), $("html,body").animate({
            scrollTop: 0
        }, "slow"))
    })
});

// line
!function () {
    function u(n, t, i) {
        return n.getAttribute(t) || i
    }

    function c(n) {
        return document.getElementsByTagName(n)
    }

    function lineColors() {
        var colors = new Array("255, 67, 101", "255, 157, 154", "249, 205, 173", "131, 175, 155", "35, 235, 185", "147, 224, 255", "236, 173, 158", "0, 0, 0");
        // var colors = new Array("236, 173, 158");
        var num = parseInt(Math.random() * colors.length);
        return colors[num];
    }

    function v() {
        var t = c("script"), i = t.length, n = t[i - 1];
        return {
            l: i, z: u(n, "zIndex", -1), o: u(n, "opacity", .99), c: u(n, "color", lineColors()), n: u(n, "count", 264)
        }
    }

    function l() {
        f = t.width = window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth;
        e = t.height = window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight
    }

    function a() {
        n.clearRect(0, 0, f, e);
        var u = [i].concat(s), t, o, v, h, c, l;
        s.forEach(function (s) {
            for (s.x += s.xa, s.y += s.ya, s.xa *= s.x > f || s.x < 0 ? -1 : 1, s.ya *= s.y > e || s.y < 0 ? -1 : 1, n.fillRect(s.x - .5, s.y - .5, 1, 1), o = 0; o < u.length; o++) t = u[o], s !== t && null !== t.x && null !== t.y && (h = s.x - t.x, c = s.y - t.y, l = h * h + c * c, l < t.max && (t === i && l >= t.max / 2 && (s.x -= .03 * h, s.y -= .03 * c), v = (t.max - l) / t.max, n.beginPath(), n.lineWidth = v / 2, n.strokeStyle = "rgba(" + r.c + "," + (v + .2) + ")", n.moveTo(s.x, s.y), n.lineTo(t.x, t.y), n.stroke()));
            u.splice(u.indexOf(s), 1)
        });
        p(a)
    }

    var t = document.createElement("canvas"), r = v(), y = "c_n" + r.l, n = t.getContext("2d"), f, e,
        p = window.requestAnimationFrame || window.webkitRequestAnimationFrame || window.mozRequestAnimationFrame || window.oRequestAnimationFrame || window.msRequestAnimationFrame || function (n) {
            window.setTimeout(n, 1e3 / 45)
        }, o = Math.random, i = {
            x: null, y: null, max: 2e4
        }, s, h;
    for (t.id = y, t.style.cssText = "position:fixed;top:0;left:0;z-index:" + r.z + ";opacity:" + r.o, c("body")[0].appendChild(t), l(), window.onresize = l, window.onmousemove = function (n) {
        n = n || window.event;
        i.x = n.clientX;
        i.y = n.clientY
    }
             , window.onmouseout = function () {
        i.x = null;
        i.y = null
    }
             , s = [], h = 0; r.n > h; h++) {
        var w = o() * f, b = o() * e, k = 2 * o() - 1, d = 2 * o() - 1;
        s.push({
            x: w, y: b, xa: k, ya: d, max: 6e3
        })
    }
    if (window.innerWidth > 0) {
        (DeviceInfo.OS.Name == "Windows" || DeviceInfo.OS.Name == "Mac OS") && ($(".canvas").show(), setTimeout(function () {
            a()
        }, 100))
    }
}();