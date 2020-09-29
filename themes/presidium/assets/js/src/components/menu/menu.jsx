import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { createStore, applyMiddleware } from 'redux';
import ReduxPromise from 'redux-promise';
import rootReducer from '../../reducers/index';
import MenuItem from './menu-item';
import Versions from '../versions/versions';
import { ACTIONS, EVENTS_DISPATCH, TOPICS } from '../../util/events';
import { isInViewport, markArticleAsViewed } from '../../util/articles';

/**
 * Locale storage key
 */
const SELECTED_ROLE = 'role.selected';

const store = createStore(
    rootReducer,
    {},
    applyMiddleware(ReduxPromise)
);

const ELEMENTS = {
    BODY: 'body',
    SCROLLABLE_CONTAINER: 'presidium-scrollable-container',
    CONTENT_CONTAINER: 'presidium-content'
};

const getDomElement = (elem) => {
    if (elem === ELEMENTS.BODY) return document.getElementsByTagName('BODY')[0];
    return document.getElementById(elem)
}

/**
 * Root navigation menu.
 */
class Menu extends Component {
    constructor(props) {
        super(props);
        this.state = {
            children: this.props.menu.children,
            roles: this.roleFilter(),
            expanded: false,
            containerHeight: 0
        };
        this.filterByRole(this.state.roles.selected);
        this.mountContainerListeners = this.mountContainerListeners.bind(this)
        this.unMountContainerListeners = this.unMountContainerListeners.bind(this)
    }

    mountContainerListeners() {
        /**
         * Manually handle click events in Presidium Container
         * that might desync the offset of scroll spy element
         * (see https://developer.mozilla.org/en-US/docs/Web/HTML/Element/details)
         */
        const _contentContainer = getDomElement(ELEMENTS.CONTENT_CONTAINER);
        _contentContainer
            .addEventListener('click', (e) => {
                //If container did resize
                if (this.state.containerHeight !== _contentContainer.clientHeight) {
                    this.setState({ containerHeight: _contentContainer.clientHeight })
                }
            });

        window.addEventListener('scroll', (e) => {
            let articles = [...document.querySelectorAll('.article')];
            articles.map((article) => {
                if (isInViewport(article)) {
                    const articleId = article.querySelector('span[data-id]').getAttribute('data-id');
                    let permalink = article.querySelector('.permalink a')
                    // Section titles do not have a permalink
                    if (permalink) {
                        markArticleAsViewed(articleId, permalink.getAttribute('href'), ACTIONS.articleScroll);
                    }
                }
            })
        })
    }

    unMountContainerListeners() {
        getDomElement(ELEMENTS.CONTENT_CONTAINER)
            .removeEventListener('click');
    }

    componentDidMount() {
        this.mountContainerListeners();
    }

    componentWillUnmount() {
        this.unMountContainerListeners();
    }

    roleFilter() {
        let selected;
        let roles = this.props.menu.roles;
        if (roles.options.length > 0) {
            selected = sessionStorage.getItem(SELECTED_ROLE);
            if (!selected) {
                selected = roles.all;
                sessionStorage.setItem(SELECTED_ROLE, selected);
            }
        } else {
            selected = roles.all;
        }
        return {
            label: roles.label,
            all: roles.all,
            selected: selected,
            options: [roles.all, ...roles.options]
        }
    }

    brandUrl() {
        if (this.props.menu.brandUrl) return this.props.menu.brandUrl;
        else if (this.props.menu.baseUrl) return this.props.menu.baseUrl;
        else return '#';
    }

    render() {
        const menu = this.props.menu;
        return (
            <div
                id='presidium-scrollable-container'
                className='scrollable-container'>
                <nav>
                    <div className='navbar-header'>
                        <a href={this.brandUrl()} className='brand'>
                            <img src={menu.logo} alt='' />
                        </a>
                        {this.props.menu.brandName &&
                            <div>
                                <p className='brand-name'>{this.props.menu.brandName}</p>
                                <Versions store={store} />
                            </div>
                        }
                        <button className='toggle' onClick={() => this.toggleMenu()}>
                            <span className='sr-only'>Toggle navigation</span>
                            <span className='icon-bar' />
                            <span className='icon-bar' />
                            <span className='icon-bar' />
                        </button>
                    </div>

                    <div className={'navbar-items' + (this.state.expanded == true ? ' expanded' : '')}>
                        {this.renderFilter()}
                        <ul>
                            {this.state.children.map(item => {
                                return <MenuItem
                                    containerHeight={this.state.containerHeight}
                                    key={item.id}
                                    baseUrl={this.props.menu.baseUrl}
                                    item={item}
                                    roles={this.state.roles} onNavigate={() => this.collapseMenu()} />
                            })}
                        </ul>
                    </div>
                </nav>
            </div>
        )
    }

    toggleMenu() {
        this.setState({ expanded: !this.state.expanded })
    }

    collapseMenu() {
        this.setState({ expanded: false })
    }

    renderFilter() {
        return this.state.roles.selected && (
            <div className='filter form-group'>
                {this.state.roles.label &&
                    <label className='control-label' htmlFor='roles-select'>{this.state.roles.label}:</label>}
                <select ref='roleselector'
                    id='roles-select'
                    className='form-control'
                    value={this.state.roles.selected}
                    onChange={(e) => this.onFilterRole(e)}>
                    {this.state.roles.options.map(role => {
                        return <option key={role} value={role}>{role}</option>
                    })}
                </select>
            </div>)
    }

    onFilterRole(e) {
        let selected = e.target.value;
        const roles = Object.assign({}, this.state.roles, { selected: selected });

        this.filterByRole(selected);
        this.setState({ roles: roles });

        sessionStorage.setItem(SELECTED_ROLE, selected);
        EVENTS_DISPATCH.MENU(TOPICS.ROLE_UPDATED, selected)
    }

    filterByRole(selected) {
        const articles = document.querySelectorAll('#presidium-content .article');
        let articlesFound = false;
        articles.forEach(article => {
            if (selected == this.state.roles.all) {
                article.style.display = 'block';
                articlesFound = true;
                return;
            }
            const roles = article.getAttribute('data-roles').split(',');
            if (roles.includes(selected) || roles.includes(this.state.roles.all)) {
                article.style.display = 'block';
                articlesFound = true;
            } else {
                article.style.display = 'none';
            }
        });

        if (articlesFound || articles.length === 0) {
            document.getElementById('no-content-warning').style.display = 'none';
        } else {
            document.getElementById('no-content-warning').style.display = 'block';
        }

    }
}

Menu.propTypes = {
    menu: React.PropTypes.shape({
        brandName: React.PropTypes.string,
        roles: React.PropTypes.object
    }).isRequired,
};

function loadMenu(menu = {}, element = 'presidium-navigation') {
    ReactDOM.render(<Menu menu={menu} />, document.getElementById(element));
}

export { Menu, loadMenu };
