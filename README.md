# rhuidean-mirror
Golang/webview reboot of stoneoftear-mirror

Requires [manetheren-server](https://github.com/thehildebreezy/manetheren-server) and [manetheren-serial](https://github.com/thehildebreezy/manetheren-serial) to function properly.

Rhuidean-mirror is the display application for a smart mirror based on a Raspberry Pi zero.

Its goal was to be able to run the html/js based display on the old armv6l architecture, something that I couldn't accomplish with node.js/electron with any sort of ease. Fortunately, the capabilities of [webview](https://github.com/webview/webview) were more than sufficient, and even much better suited to the task.