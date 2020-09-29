const TOPICS = {
    RANKING_LOADED: 'RANKING_LOADED',
    ROLE_UPDATED: 'ROLE_UPDATED'
};

let ACTIONS = {
    articleScroll: 'article_scroll',
    articleClick: 'article_click',
    articleLoad: 'article_load'
};

let EVENTS_DISPATCH = {
    MENU: (topic, value) => {
        if (topic) {
            window.events.next({
                topic: topic,
                value: value
            });
        }
    },
    ARTICLE: (path, permalink, action, solution) => {
        if (path) {
            window.events.next({
                path: path,
                permalink: permalink,
                action: action,
                solution: solution
            });
        }
    }
};

export {EVENTS_DISPATCH, ACTIONS, TOPICS};