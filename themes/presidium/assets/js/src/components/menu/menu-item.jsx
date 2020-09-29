import React, {Component} from 'react';
import gumshoe from './scroll-spy';
import classnames from 'classnames';

import {MENU_TYPE} from './menu-structure';
import {ACTIONS, TOPICS} from '../../util/events';
import {markArticleAsViewed, slugify} from '../../util/articles'

let timeout;

/**
 * Menu item that may have one or more articles or nested groups of articles.
 */
export default class MenuItem extends Component {

    constructor(props) {
        super(props);
        let itemURL = this.props.item.url;
        if (!itemURL.endsWith('/')) itemURL = `${itemURL}/`;

        const onPage = itemURL === window.location.pathname;
        const inSection = this.inSection();
        const hasChildren = props.item.children.length > 0;
        this.state = {
            onPage: onPage,
            inSection: onPage || inSection,
            alwaysExpanded: this.props.item.alwaysExpanded,
            isCollapsed: this.props.item.collapsed,
            hasChildren: hasChildren,
            activeArticle: this.props.activeArticle,
            isExpanded: (inSection || this.props.item.alwaysExpanded) && hasChildren,
            selectedRole: this.props.roles.selected,
            newTab: this.props.item.newTab,
        };
    }

    inSection() {
        const base = this.props.item.url;
        const reference = window.location.pathname;
        if (base === this.props.baseUrl) {
            return base === reference;
        }
        return reference.startsWith(base);
    }

    resetScrollSpyHeights() {
        if (this.state.isExpanded) gumshoe.setDistances();
    }

    componentDidMount() {
        window.events.subscribe({
            next: (event) => {
                if (event.path === TOPICS.RANKING_LOADED) this.resetScrollSpyHeights();
            }
        });
        if (this.state.onPage) {
            clearTimeout(timeout);
            if (this.props.item.children.length === 1) {
                let action = ACTIONS.articleLoad;
                if (sessionStorage.getItem('article.clicked') === this.props.item.children[0].id) {
                    action = ACTIONS.articleClick;
                    sessionStorage.removeItem('article.clicked');
                }
                timeout = setTimeout(function () {
                    const id = this.props.item.children[0].id;
                    const permalink = document
                        .querySelector(`span[data-id='${id}']`)
                        .parentElement.querySelector('.permalink a')
                        .href;
                    markArticleAsViewed(id, permalink, action);
                }.bind(this), 2000);
            }
            this.initializeScrollSpy()
        }
    }

    componentWillReceiveProps(props) {
        //If there's a new click event, trigger the recalc of scroll spy offset
        if (this.props.containerHeight !== props.containerHeight) this.resetScrollSpyHeights();

        //Propagate active article and roles down the menu chain
        const activeArticle = this.state.onPage ? this.state.activeArticle : props.activeArticle;
        this.setState({
            activeArticle: activeArticle,
            selectedRole: props.roles.selected
        });
    }

    componentDidUpdate(prevProps, prevState) {
        if (this.state.onPage && prevState.selectedRole !== this.state.selectedRole) {
            this.initializeScrollSpy()
        }
    }

    determineScrollSpyOffset() {
        let defaultOffset = 100; //TODO: Pass 'loaded' event from Presidium JS Enterprise and remove this

        let topBar = document.getElementById('presidium-enterprise-toolbar');
        let solutionBar = document.getElementById('presidium-solution-search');
        let blueBar = document.getElementById('presidium-currently-viewing-marquee');

        if (topBar) defaultOffset += topBar.clientHeight;
        if (solutionBar) defaultOffset += solutionBar.clientHeight;
        if (blueBar) defaultOffset += blueBar.clientHeight;

        return defaultOffset;
    }

    initializeScrollSpy() {
        gumshoe.init({
            selector: '[data-spy] a',
            selectorTarget: '#presidium-content .article > .anchor',
            container: window,
            offset: this.determineScrollSpyOffset(),
            activeClass: 'on-article',
            callback: (active) => {
                //Update active article on scroll. Ignore hidden articles (with distance = 0)
                const activeArticle = active && active.distance > 0 ? active.nav.getAttribute('data-id') : undefined;
                if (activeArticle && this.state.activeArticle !== activeArticle) {
                    clearTimeout(timeout);
                    timeout = setTimeout(function () {
                        if (this.state.activeArticle === activeArticle) this.resetScrollSpyHeights();
                    }.bind(this), 2000);
                    sessionStorage.removeItem('article.clicked');
                    this.setState({activeArticle: activeArticle});
                }
            }
        });
    }

    render() {
        const item = this.props.item;
        return (
            <li key={item.id} className={this.parentStyle(item)}>
                <div onClick={(e) => this.clickParent(e)} className={'menu-row ' + this.levelClass(item.level)}>
                    <div className="menu-expander">
                        {this.expander()}
                    </div>
                    <div className="menu-title">
                        <a target={item.newTab ? '_blank' : '_self'} data-id={item.id} href={item.url}>{item.title}</a>
                    </div>
                </div>

                {/* Normal sub-menu items */}
                {!this.state.isCollapsed &&
                <ul
                    {...this.spyOnMe()}
                    className={this.state.isExpanded ? 'dropdown expanded' : 'dropdown'}>
                    {this.children()}
                </ul>
                }

                {/* Hidden (for scroll-spy) sub-menu items */}
                {this.state.isCollapsed &&
                <ul {...this.spyOnMe()} className="hidden">
                    {this.children()}
                </ul>
                }
            </li>
        )
    }

    children() {
        return this.props.item.children.map(item => {
            switch (item.type) {
                case MENU_TYPE.CATEGORY:
                    return <MenuItem key={item.title}
                                     item={item}
                                     activeArticle={this.state.activeArticle}
                                     roles={this.props.roles}
                                     baseUrl={this.props.baseUrl}
                                     onNavigate={this.props.onNavigate}/>;
                case MENU_TYPE.ARTICLE:
                    return <li key={item.id} className={this.childStyle(item)}>
                        <div onClick={() => this.clickChild(item.url, item.id)}
                             className={'menu-row ' + this.articleStyle(item)}>
                            <div className="menu-expander"></div>
                            <div className="menu-title">
                                <a data-id={item.id} href={`#${item.slug}`}>{item.title}</a>
                            </div>
                        </div>
                    </li>;
            }
        });
    }

    expander() {
        if (!this.state.isCollapsed && this.state.hasChildren) {
            return <span onClick={(e) => this.toggleExpand(e)}
                         className={this.state.isExpanded ? 'glyphicon glyphicon-chevron-down' : 'glyphicon glyphicon-chevron-right'}/>
        } else {
            return ''
        }
    }

    spyOnMe() {
        return this.state.onPage ? {'data-spy': ''} : {};
    }

    parentStyle(item) {
        return classnames(`menu-parent_${slugify(item.title)}`, {
            'in-section': this.state.inSection || this.containsArticle(),
            'expanded': this.state.isExpanded,
            'on-page': this.state.onPage,
            'on-article': this.state.activeArticle === item.id,
            'hidden': !this.hasRole(item)
        })
    }

    containsArticle() {
        if (!this.state.activeArticle || !this.state.hasChildren) {
            return false;
        }
        return this.containsNested(this.state.activeArticle, this.props.item.children);
    }

    containsNested(activeArticle, children) {
        return children.find(child => {
            if (child.type == MENU_TYPE.ARTICLE && child.id == activeArticle) {
                return true;
            } else if (child.children) {
                return this.containsNested(activeArticle, child.children);
            }
        })
    }

    childStyle(item) {
        return classnames({
            'on-article': this.state.activeArticle === item.id,
            'hidden': !this.hasRole(item)
        })
    }

    articleStyle(item) {
        return this.levelClass(item.level);
    }

    levelClass(level) {
        switch (level) {
            case 1:
                return ' level-one';
            case 2:
                return ' level-two';
            case 3:
                return ' level-three';
            case 4:
                return ' level-four';
            case 5:
                return ' level-five';
        }
        return '';
    }

    hasRole(item) {
        return this.props.roles.selected == this.props.roles.all ||
            item.roles.indexOf(this.props.roles.all) >= 0 ||
            item.roles.indexOf(this.props.roles.selected) >= 0;
    }

    toggleExpand(e) {
        e.stopPropagation();
        if (this.state.hasChildren && !this.state.alwaysExpanded) {
            this.setState({isExpanded: !this.state.isExpanded})
        }
    }

    clickParent(e) {
        if (!this.props.item.newTab) {
            if (this.state.onPage) {
                e.preventDefault();
                e.stopPropagation();
            } else {
                sessionStorage.setItem('article.clicked', this.props.item.children.length === 1 ? this.props.item.children[0].id : this.props.item.id);
                this.props.onNavigate();
                window.location = this.props.item.url;
            }
        }
    }

    clickChild(path, id) {
        sessionStorage.setItem('article.clicked', id);
        this.props.onNavigate();
        window.location = path;
    }
}

MenuItem.propTypes = {
    item: React.PropTypes.object.isRequired,
    baseUrl: React.PropTypes.string.isRequired,
    activeArticle: React.PropTypes.string,
    onNavigate: React.PropTypes.func,
    roles: React.PropTypes.object
};