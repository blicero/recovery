// Time-stamp: <2022-04-02 15:19:43 krylon>
// -*- mode: javascript; coding: utf-8; -*-
// Copyright 2015-2020 Benjamin Walkenhorst <krylon@gmx.net>
//
// This file has grown quite a bit larger than I had anticipated.
// It is not a /big/ problem right now, but in the long run, I will have to
// break this thing up into several smaller files.

"use strict";

function defined(x) {
    return undefined != x && null != x;
}

function fmtDateNumber(n) {
    if (n < 10) {
        return "0" + n.toString();
    } else {
        return n.toString();
    }
} // function fmtDateNumber(n)

function timeStampString(t) {
    if (typeof(t) == "string") {
        return t;
    }

    // (1900 + d.getYear()) + "-" + d.getMonth() + "-" + d.getDate() + " " + d.getHours() + ":" + d.getMinutes() + ":" + d.getSeconds()
    var year = t.getYear() + 1900;
    var month = fmtDateNumber(t.getMonth() + 1);
    var day = fmtDateNumber(t.getDate());
    var hour = fmtDateNumber(t.getHours());
    var minute = fmtDateNumber(t.getMinutes());
    var second = fmtDateNumber(t.getSeconds());

    var s =
        year + "-" + month + "-" + day +
        " " + hour + ":" + minute + ":" + second;
    return s;
} // function timeStampString(t)

function fmtDuration(seconds) {
    var minutes = 0, hours = 0;

    while (seconds > 3599) {
        hours++;
        seconds -= 3600;
    }

    while (seconds > 59) {
        minutes++;
        seconds -= 60;
    }

    if (hours > 0) {
        return `${hours}h${minutes}m${seconds}s`;
    } else if (minutes > 0) {
        return `${minutes}m${seconds}s`;
    } else {
        return `${seconds}s`;
    }
} // function fmtDuration(seconds)

function beaconLoop() {
    try {
        if (settings.beacon.active) {
            var req = $.get("/ajax/beacon",
                            {},
                            function(response) {
                                var status = "";

                                if (response.Status) {
                                    status =
                                        response.Message +
                                        " running on " +
                                        response.Hostname + 
                                        " is alive at " +
                                        response.Timestamp;
                                } else {
                                    status = "Server is not responding";
                                }

                                var beaconDiv = $("#beacon")[0];

                                if (defined(beaconDiv)) {
                                    beaconDiv.innerHTML = status;
                                    beaconDiv.classList.remove("error");
                                } else {
                                    console.log("Beacon field was not found");
                                }
                            },
                            "json"
                           ).fail(function() {
                               var beaconDiv = $("#beacon")[0];
                               beaconDiv.innerHTML = "Server is not responding";
                               beaconDiv.classList.add("error");
                               //logMsg("ERROR", "Server is not responding");
                           });
        }
    }
    finally {
        window.setTimeout(beaconLoop, settings.beacon.interval);
    }
} // function beaconLoop()

function beaconToggle() {
    settings.beacon.active = !settings.beacon.active;
    saveSetting("beacon", "active", settings.beacon.active);

    if (!settings.beacon.active) {
        var beaconDiv = $("#beacon")[0];
        beaconDiv.innerHTML = "Beacon is suspended";
        beaconDiv.classList.remove("error");
    }

} // function beaconToggle()


/*
The ???content??? attribute of Window objects is deprecated.  Please use ???window.top??? instead. interact.js:125:8
Ignoring get or set of property that has [LenientThis] because the ???this??? object is incorrect. interact.js:125:8

*/

function db_maintenance() {
    const maintURL = "/ajax/db_maint";

    var req = $.get(
        maintURL,
        {},
        function(res) {
            if (!res.Status) {
                console.log(res.Message);
                postMessage(new Date(), "ERROR", res.Message);
            } else {
                const msg = "Database Maintenance performed without errors";
                console.log(msg);
                postMessage(new Date(), "INFO", msg);
            }
        },
        "json"
    ).fail(function() {
        var msg = "Error performing DB maintenance";
        console.log(msg);
        postMessage(new Date(), "ERROR", msg);
    });
} // function db_maintenance()

function msgCheckSum(timestamp, level, msg) {
    var line = [ timeStampString(timestamp), level, msg ].join("##");

    var cksum = sha512(line);
    return cksum;
}

var curMessageCnt = 0;

function post_test_msg() {
    var user = $("#msgTestText")[0];
    var msg = user.value;
    var now = new Date();

    postMessage(now, "DEBUG", msg);
} // function post_tst_msg()

function postMessage(timestamp, level, msg) {
    var row = '<tr id="msg_' +
        msgCheckSum(timestamp, level, msg) +
        '"><td>' +
        timeStampString(timestamp) +
        '</td><td>' +
        level +
        '</td><td>' +
        msg +
        '</td></tr>\n';

    msgRowAdd(row);
} // function postMessage(timestamp, level, msg)

function adjustMsgMaxCnt() {
    var cntField = $("#max_msg_cnt")[0];
    var newMax = cntField.valueAsNumber;

    if (newMax < curMessageCnt) {
        var rows = $("#msg_body")[0].children;

        while (rows.length > newMax) {
            rows[rows.length - 1].remove();
            curMessageCnt--;
        }

    }
    
    saveSetting("messages", "maxShow", newMax);
} // function adjustMaxMsgCnt()

function adjustMsgCheckInterval() {
    var intervalField = $("#msg_check_interval")[0];
    if (intervalField.checkValidity()) {
        var interval = intervalField.valueAsNumber;
        // intervalField.setInterval(interval); // ???
        saveSetting("messages", "interval", interval);
    }
} // function adjustMsgCheckInterval()

function toggleCheckMessages() {
    var box = $("#msg_check_switch")[0];
    var newVal = box.checked;

    saveSetting("messages", "queryEnabled", newVal);
} // function toggleCheckMessages()

function getNewMessages() {
    const msgURL = "/ajax/get_messages";

    try {
        if (!settings.messages.queryEnabled) {
            return;
        }
        
        var req = $.get(
            msgURL,
            {},
            function(res) {
                if (!res.Status) {
                    var msg = msgURL +
                        " failed: " +
                        res.Message;

                    console.log(msg)
                    alert(msg);
                } else {
                    for (var i = 0; i < res.Messages.length; i++) {
                        var item = res.Messages[i];
                        var rowid =
                            "msg_" +
                            msgCheckSum(item.Time, item.Level, item.Message);
                        var row = '<tr id="' +
                            rowid +
                            '"><td>' +
                            item.Time +
                            '</td><td>' +
                            item.Level +
                            '</td><td>' +
                            item.Message +
                            '</td><td>' +
                            '<input type="button" value="Delete" onclick="msgRowDelete(\'' +
                            rowid +
                            '\');" />' +
                            '</td></tr>\n';

                        msgRowAdd(row);
                    }
                }
            },
            "json"
        );
    }
    finally {
        window.setTimeout(getNewMessages, settings.messages.interval);
    }

} // function getNewMessages()

function logMsg(level, msg) {
    var timestamp = timeStampString(new Date());
    var rowID = "msg_" + sha512(msgCheckSum(timestamp, level, msg));
    var row = '<tr id="' +
        rowID +
        '"><td>' +
        timestamp +
        '</td><td>' +
        level +
        '</td><td>' +
        msg +
        '</td><td>' +
        '<input type="button" value="Delete" onclick="msgRowDelete(\'' +
        rowID +
        '\');" />' +
        '</td></tr>\n';

    $("#msg_display_tbl")[0].innerHTML += row;
} // function logMsg(level, msg)

function msgRowAdd(row) {
    var msgBody = $("#msg_body")[0];

    msgBody.innerHTML = row + msgBody.innerHTML;

    if (++curMessageCnt > settings.messages.maxShow) {
        msgBody.children[msgBody.children.length - 1].remove();
    }

    var tbl = $("#msg_tbl")[0];
    if (tbl.hidden) {
        tbl.hidden = false;
    }
} // function msgRowAdd(row)

function msgRowDelete(rowID) {
    var row = $("#" + rowID)[0];

    if (row != undefined) {
        row.remove();
        if (--curMessageCnt == 0) {
            var tbl = $("#msg_tbl")[0];
            tbl.hidden = true;
        }
    }
} // function msgRowDelete(rowID)

function msgRowDeleteAll() {
    var msgBody = $("#msg_body")[0];
    msgBody.innerHTML = '';
    curMessageCnt = 0;

    var tbl = $("#msg_tbl")[0];
    tbl.hidden = true;
} // function msgRowDeleteAll()

function requestTestMessages() {
    const urlRoot = "/ajax/rnd_message/";

    var cnt = $("#msg_cnt")[0].valueAsNumber;
    var rounds = $("#msg_round_cnt")[0].valueAsNumber;
    var delay = $("#msg_round_delay")[0].valueAsNumber;

    if (cnt == 0) {
        console.log("Generate *0* messages? Alrighty then...");
        return;
    }

    var reqURL = urlRoot + cnt;

    $.get(
        reqURL,
        {
            "Rounds": rounds,
            "Delay": delay,
        },
        function(res) {
            if (!res.Status) {
                console.log(res.Message);
                alert(res.Message);
            }
        },
        "json"
    ).fail(function() {
            const msg = "Requesting test messages failed.";
            console.log(msg);
            //alert(msg);
            logMsg("ERROR", msg);
        });
} // function requestTestMessages()

function toggleMsgTestDisplayVisible() {
    var tbl = $("#test_msg_cfg")[0];

    if (tbl.hidden) {
        tbl.hidden = false;

        var checkbox = $("#msg_check_switch")[0];
        settings.messages.queryEnabled = checkbox.checked;
    } else {
        settings.messages.queryEnabled = false;
        tbl.hidden = true;
    }
} // function toggleMsgTmpDisplayVisible()

function toggleMsgDisplayVisible() {
    var display = $("#msg_display_div")[0];

    display.hidden = !display.hidden;
} // function toggleMsgDisplayVisible()

function update_score_display(name) {
    const input_id = `#${name}_score`;
    const display_id = `#${name}_score_display`;

    const val = $(input_id)[0].value;

    $(display_id)[0].innerHTML = val;
} // function update_score_display()
