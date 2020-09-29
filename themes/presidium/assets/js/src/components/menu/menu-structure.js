import {path} from '../../util/articles';

const LEVEL_1 = 1;
const LEVEL_2 = 2;
const LEVEL_3 = 3;
const LEVEL_4 = 4;

export const MENU_TYPE = {
    SECTION: 'section',
    CATEGORY: 'category',
    ARTICLE: 'article'
};

function menuSection(section, defaultRole) {
    return {
        type: MENU_TYPE.SECTION,
        id: section.path,
        level: LEVEL_1,
        expandable: section.expandable,
        title: section.title,
        slug: section.slug,
        path: section.path,
        categories: {},
        children: [],
        roles: new Set([defaultRole])
    }
}

function menuCategory(key, id, path, level) {
    return {
        type: MENU_TYPE.CATEGORY,
        id: id,
        level: level,
        expandable: true,
        title: key,
        slug: path,
        path: path,
        categories: {},
        children: [],
        roles: new Set()
    };
}

function menuArticle(article, level, defaultRole) {
    return {
        type: MENU_TYPE.ARTICLE,
        id: article.id,
        path: article.path,
        slug: article.slug,
        title: article.title,
        level: level,
        expandable: false,
        roles: article.roles.length > 0 ? new Set(article.roles) : new Set([defaultRole])
    }
}

/**
 * Returns a new Set of the merged current and additional filters.
 * Merges the default filters if no additional filters are provided.
 */
function mergeSets(current, additional, defaultFilter) {
    if (additional.length > 0) {
        return new Set([...current, ...additional])
    } else {
        return defaultFilter ? new Set([...current, defaultFilter]) : current;
    }
}

/**
 * Creates or gets a category.
 */
function getOrCreateCategory(section, key, path, level) {
    let category;
    if (section.categories[key]) {
        category = section.categories[key];
    } else {
        const id = path.concat(section.id, key);
        category = menuCategory(key, id, path, level);
        section.categories[key] = category;
        section.children.push(category);
    }
    return category;
}

function hasSub(categories) {
    return categories.length > 1;
}

/**
 *  Build the menu structure maintaining the provided order.
 *  Group articles in a section by a distinct category and optional sub category
 *  Filters for each subsection are merged to a parents for filtering.
 */
export function groupByCategory(root, defaultRole) {

    const section = menuSection(root, defaultRole);

    root.articles.forEach(article => {

        if (article.id.endsWith('index')) {
            return;
        }

        section.roles = mergeSets(section.roles, article.roles, defaultRole);

        if (!article.category) {
            section.children.push(menuArticle(article, LEVEL_2, defaultRole));
        } else {
            const categories = article.category.split('/');
            const category = getOrCreateCategory(section, categories[0], article.path, LEVEL_2);
            category.roles = mergeSets(category.roles, article.roles, defaultRole);

            if (!hasSub(categories)) {
                category.children.push(menuArticle(article, LEVEL_3, defaultRole));
            } else {
                const subCategory = getOrCreateCategory(category, categories[1], article.path, LEVEL_3);
                subCategory.roles = mergeSets(subCategory.roles, article.roles, defaultRole);
                subCategory.children.push(menuArticle(article, LEVEL_4, defaultRole));
            }
        }
    });
    return section;
}

