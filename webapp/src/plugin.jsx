import {FormattedMessage} from 'react-intl';

import en from 'i18n/en.json';

import {id as pluginId} from './manifest';

import {
    postDropdownMenuAction,
} from './actions';

function getTranslations(locale) {
    switch (locale) {
    case 'en':
        return en;
    }
    return {};
}

function getSubMenu() {
    const menu = (
        <FormattedMessage
            id='submenu.message'
            defaultMessage='Remind me about this'
        />
    );
    const subMenu = [
        (
            <FormattedMessage
                id='submenu.20min'
                key='submenu.20min'
                defaultMessage='In 20 minutes'
            />
        ),
        (
            <FormattedMessage
                id='submenu.1hr'
                key='submenu.1hr'
                defaultMessage='In 1 hour'
            />
        ),
        (
            <FormattedMessage
                id='submenu.3hr'
                key='submenu.3hr'
                defaultMessage='In 3 hours'
            />
        ),
        (
            <FormattedMessage
                id='submenu.tomorrow'
                key='submenu.tomorrow'
                defaultMessage='Tomorrow'
            />
        ),
        (
            <FormattedMessage
                id='submenu.nextweek'
                key='submenu.nextweek'
                defaultMessage='Next week'
            />
        ),
    ];
    return {menu, subMenu};
}

export default class RemindPlugin {
    initialize(registry, store) {
        registry.registerPostDropdownMenuAction(
            getSubMenu(),
            (postId, item) => store.dispatch(postDropdownMenuAction(postId, item)),
        );

        registry.registerTranslations(getTranslations);
    }

    uninitialize() {
        //eslint-disable-next-line no-console
        console.log(pluginId + '::uninitialize()');
    }
}
