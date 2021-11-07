// ==UserScript==
// @name         New World Map Live Position
// @namespace    http://tampermonkey.net/
// @version      0.3
// @description  try to take over the world!
// @author       You
// @match        https://www.newworld-map.com/*
// @icon         https://www.google.com/s2/favicons?domain=newworld-map.com
// @grant        none
// @run-on       document-end
// @downloadURL  https://raw.githubusercontent.com/seler/new-world-map-coordinates-sse/master/userscripts/NewWorldMapLivePosition.user
// @updateURL    https://raw.githubusercontent.com/seler/new-world-map-coordinates-sse/master/userscripts/NewWorldMapLivePosition.user
// ==/UserScript==

const LAT_OFFSET = -14336; // taken from script source at newworld-map.com
const LNG_OFFSET = 170; // to cover for left pane when closed
const HISTORY_TIME = 10 * 60 * 1000 // 10 minutes
const SSE_URL = "http://localhost:5000/events"; // stream of coordinates in {lat: float, lng: float} json format

function htmlToElement(html) {
    var template = document.createElement('template');
    html = html.trim(); // Never return a text node of whitespace as the result
    template.innerHTML = html;
    return template.content.firstChild;
}

(function() {
    'use strict';

    var map, marker, es, path;
    var history = [];
    var pathPoints = [];
    var tracking = false;
    var following = false;

    function updateHistoryPath(){
        const now = Math.floor(Date.now());
        history = history.filter(point => point.timestamp > (now - HISTORY_TIME));
        path.setLatLngs(history.map(({lat, lng}) => [lat, lng]));
        path.redraw();
    }


    function enableTracking() {
        tracking = true;
        map = document.getElementById('map').__vue__.mapObject;
        const icon = window.L.divIcon({
            html: `<svg xmlns="http://www.w3.org/2000/svg" viewbox="0 0 100 100" style="filter: drop-shadow(0px 0px 10px rgb(0, 0, 0, 1));">
                <polygon fill="currentColor" points="50,0 90,100 50,75 10,100"/>
            </svg>`,
            iconSize: [24, 24],
            iconAnchor: [12, 12],
            className: 'playerPositionMarker',
        });
        marker = window.L.marker(
            {lat: 0, lng: 0}, 
            {title: "Player position", zIndexOffset: 9999999, riseOnHover: true, icon: icon}
        );
        marker.addTo(map);

        var iconElement = document.getElementsByClassName('playerPositionMarker')[0]
        iconElement.style.background = null;
        iconElement.style.border = null;
        iconElement.style.color = "yellow";
        iconElement.style.transformOrigin = "50% 50%";
        var iconSVGElement = document.querySelector('.playerPositionMarker svg');

        path = window.L.polyline(pathPoints, {color: 'rgba(255, 0, 0, 1)', weight: 1.5});
        path.bringToFront();
        path.addTo(map);
        map.setZoom(7);

        es = new EventSource(SSE_URL);

        es.addEventListener('message', event => {
            const latLng = JSON.parse(event.data);
            const lat = parseFloat(latLng.lat) + LAT_OFFSET;
            const lng = parseFloat(latLng.lng)
            marker.setLatLng({lat, lng});
            history.push({lat, lng, timestamp: Math.floor(Date.now())});
            updateHistoryPath();
            
            if (!!history.length) {
                var angleRadians = -1 * Math.atan2(lat - history[history.length-2].lat, lng - history[history.length-2].lng) + Math.PI/2;
                iconSVGElement.style.transform = `rotate(${angleRadians}rad)`;
            }
            
            if (following) {
                map.panTo({lat, lng});
            }
        });

        es.addEventListener('open', console.log);
        es.addEventListener('error', console.log);
        console.log('tracking enabled');

        setInterval(updateHistoryPath, 1000);
    }

    function disableTracking() {
        tracking = false;
        marker.remove();
        es.close();
        console.log('tracking disabled');
    }

    function enableFollowing() {
        following = true;
        console.log('following enabled');
    }

    function disableFollowing() {
        following = false;
        console.log('following disabled');
    }

    function init() {
        document.getElementById('main_bg').insertAdjacentHTML('afterend', `
        <div data-v-28c05aa0="" data-v-2559b0c1="" class="v-card v-card--flat v-sheet theme--light" id="main_bg" style="min-height: 107px; min-width: 107px;">
            <div data-v-28c05aa0="" style="overflow: hidden;" class="container ma-0 pa-0 text-center primary--text container--fluid fill-height">
                <div data-v-28c05aa0="" style="" class="v-list-item__title text-center mt-4 mb-4" dark="">
                    <span data-v-28c05aa0="" class="white--text">Position tracking</span>
                </div>

                <div data-v-28c05aa0="" class="row no-gutters nwmlp-control-buttons"></div>
            </div>
        </div>
        `);

        function insertButton(query, text) {
            var buttons = document.querySelector(query);
            var buttonContainer = htmlToElement('<div data-v-28c05aa0="" class="col col-12" />');
            buttons.appendChild(buttonContainer);
            var buttonContent = htmlToElement('<div data-v-28c05aa0="" tabindex="0" class="primary--text v-card v-card--flat v-card--link v-sheet theme--light rounded-0" id="cat_btn" />');
            buttonContainer.appendChild(buttonContent);
            buttonContent.textContent = text;
            return buttonContent;
        }

        var trackingButton = insertButton('.nwmlp-control-buttons', 'Enable tracking');
        trackingButton.textContent = tracking ? "Disable tracking" : "Enable tracking";
        trackingButton.addEventListener("click", event => {
            tracking ? disableTracking() : enableTracking();
            trackingButton.textContent = tracking ? "Disable tracking" : "Enable tracking";
        });

        var followingButton = insertButton('.nwmlp-control-buttons', 'Enable following');
        followingButton.addEventListener("click", event => {
            following ? disableFollowing() : enableFollowing();
            followingButton.textContent = tracking ? "Disable following" : "Enable following";
        });
    }

    // dunno how to wait for map to initialize
    setTimeout(init, 1000);

    document.body.parentElement.style.overflow = 'hidden'; // fix for empty scrollbar
})();
