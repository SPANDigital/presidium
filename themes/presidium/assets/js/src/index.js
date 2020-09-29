import React from 'react';
import {loadMenu} from './components/menu/menu';
import {init as initModal} from './components/image/modal';
// import {handleQueryString, checkSessionStorageConfig} from './util/config';
import {Subject} from 'rxjs/Subject'

initModal();
// TODO: Find a solution that is easier on the local storage if needed for edit mode.
//       Safari returns and error with 303 when local storage exceeds the limit for the domain
//       https://macreports.com/safari-kcferrordomaincfnetwork-error-blank-page-fix/
//       https://www.quora.com/Why-does-Safari-give-me-a-kCFErrorDomainCFNetwork-error-303-when-browsing-some-sites
// handleQueryString(); //Check query string arguments and persist to sessionStorage
// checkSessionStorageConfig(); //Check sessionStorage for known configurations and apply them

var presidium = {
    menu: {
        load: loadMenu,
    },
    modal: {
        init: initModal,
    }
};

window.presidium = presidium;
window.events = new Subject();