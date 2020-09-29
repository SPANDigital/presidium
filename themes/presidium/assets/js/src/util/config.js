import {cache, CACHE_KEYS} from './cache';

/**
 * Get key/value pairs from the query string
 * @returns {string[]}
 */
const getUrlVars = () => {
    return window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
};

/**
 * Save query string key/value pairs to sessionStorage
 */
const handleQueryString = () => {
    const hashes = getUrlVars();
    hashes.map((hash) => {
        const vals = hash.split('=');
        cache.set(`presidium.urlArgs.${vals[0]}`, vals[1]);
    })
};

const checkSessionStorageConfig = () => {
    const articleTitles = [...document.querySelectorAll('[data-url-variable]')];

    //Check localStorage cache for presidium entries
    const localStorageKeys = Object.keys(localStorage).filter((item) => {
        return item.startsWith('presidium.urlArgs.')
    });

    //Check DOM elements that have the corresponding key as a data-attribute
    localStorageKeys.map((item) => {
        if (cache.get(item) === true) {
            articleTitles.filter((article) => {
                const urlProp = item.split('presidium.urlArgs.')[1];
                if (article.dataset.urlVariable === urlProp) article.style.display = 'flex';
            })
        }
    });
};

export {
    handleQueryString,
    checkSessionStorageConfig
};