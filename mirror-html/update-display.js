/**
 * Setup holder to hold our logic
 * Should be electron independent so we won't have to worry too
 * much in our webview port about this changing on us
 */

const Update = {
    display : ( arg ) => {
        console.log('updating display')

        console.log(arg)

        jsonObj = JSON.parse( arg )
        if( jsonObj == null ){
            console.log('Bad JSON object')
            return
        }

        if( document.readyState !== "complete" ){
            console.log("DOM not ready")
            return
        }

        // update display
        // first get type
        var type = jsonObj['type'];
        if( type == null ){
            console.log('Bad message input')
        } else if( type == "other" ){
            Update.other( jsonObj['message'] )
        } else if( type == "time" ){
            Update.time( jsonObj['message'] )
        } else if( type == "weather" ){
            Update.weather( jsonObj['message'] )
        } else if( type == "forecast" ){
            Update.forecast( jsonObj['message'] )
        }

        else console.log('No valid type to process')
    },

    time: ( message ) => {
            
        var utcSec = parseInt(message)
        var utcUSec = utcSec * 1000

        Store.date = new Date( utcUSec )

        clearInterval(Store.interval)
        Store.interval = setInterval(Update.ticker,1000)
    },

    weather: ( message ) => {

        // make sure I can get the object
        var el = document.getElementById('weather')
        if( el == null ){
            console.log( "Could not find element #weather" )
        }

        // write this html to the element
        var HTML = "";

        // get weather information

        HTML += "<div id='temp'>"
        HTML += '<span class="big"><i class="wi wi-owm-'+ message['weather'][0]['id'] +'"></i>'
        HTML += Utility.ktof( message['main']['temp'] ) + "&deg;F</span>";
        
        //HTML += ''
        HTML += "</div>"
        HTML += "<div id='loc'>"
        HTML += message['name'] + ' / ' + message['weather'][0]['description']
        HTML += "</div>"

        // finaly act is to update the display
        el.innerHTML = HTML;
    },

    forecast : ( message ) => {
        // make sure I can get the object
        var el = document.getElementById('forecast')
        if( el == null ){
            console.log( "Could not find element" )
        }

        var HTML = ""

        // get forecast information
        HTML += "<table id='forecast'>"

        
        
        for( var i=0; i< message['list'].length; i++ ){
            HTML += "<tr><td>"
            
            var d = new Date( message['list'][i]['dt_txt'] )
            HTML += Utility.days[d.getDay()]+":"

            HTML += "</td><td>"
            HTML += '<i class="wi wi-owm-'+ message['list'][i]['weather'][0]['id'] +'"></i>'
            HTML += "</td><td>"
            HTML += Utility.ktof( message['list'][i]['main']['temp'] ) + "&deg;F";
            HTML += "</td></tr>"
        }

        HTML += "</table>"

        // finally update element
        el.innerHTML = HTML;
    },

    ticker : () => {

        // make sure I can get the object
        var el = document.getElementById('time')
        if( el == null ){
            console.log( "Could not find element #time" )
        }

        if( Store.date == null ) return;

        Store.date.setSeconds( Store.date.getSeconds() + 1)
        
        var h = Store.date.getHours()
        var m = Store.date.getMinutes()
        var seconds = Store.date.getSeconds()

        var year = Store.date.getFullYear()
        var month = Store.date.getMonth()
        var daynum = Store.date.getDate()

        var day = Store.date.getDay()

        var ampm = "pm"

        if ( m < 10 ){
            m = "0" + m
        }

        if( h == 0 ) {
            h = 12
            ampm = "am"
        } else if ( h < 12 ) {
            ampm = "am"
        } else if( h > 12 ) {
            h = h-12
            ampm = "pm"
        } else {
            ampm = "pm"
        }

        var color = (seconds%2 == 0) ? 'black' : 'white';

        var HTML = "";
        HTML += "<div class='big'>"
        HTML += h
        HTML += "<span id='sep' style='color:"+color+"'>:</span>"
        HTML += m + ampm
        HTML += "</div><div>"
        HTML += Utility.months[month] + " " + daynum + ", " + year
        HTML += "</div><div>"
        HTML += Utility.fullDays[day]
        HTML += "</div>"

        el.innerHTML = HTML
    }
}

var Store = {
    date : null,
    interval : null
}

const Utility = {
    ktof : ( Kelvin ) => {
        console.log(Kelvin)
        var K = parseInt(Kelvin)
        if( K == null ) {
            console.log("bad temperature input")
            return
        }
        return parseInt(Math.floor(32 + (((K-273.15) * 9)/5)))
    },

    
    days : ['S','M','T','W','T','F','S'],
    months : ["January","February","March","April","May","June","July","August","September","October","November","December"],
    fullDays : ["Sunday","Monday","Tuesday","Wednesday","Thursday","Friday","Saturday"]
}


