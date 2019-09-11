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
    return {
        id: 'submenu.message',
        text: (
            <FormattedMessage
                id='submenu.message'
                defaultMessage='Remind me about this'
            />
        ),
        subMenu: [
            {
                id: 'submenu.20min',
                text: (
                    <FormattedMessage
                        id='submenu.20min'
                        key='submenu.20min'
                        defaultMessage='In 20 minutes'
                    />
                ),
            },
            {
                id: 'submenu.1hr',
                text: (
                    <FormattedMessage
                        id='submenu.1hr'
                        key='submenu.1hr'
                        defaultMessage='In 1 hour'
                    />
                ),
            },
            {
                id: 'submenu.3hr',
                text: (
                    <FormattedMessage
                        id='submenu.3hr'
                        key='submenu.3hr'
                        defaultMessage='In 3 hours'
                    />
                ),
            },
            {
                id: 'submenu.tomorrow',
                text: (
                    <FormattedMessage
                        id='submenu.tomorrow'
                        key='submenu.tomorrow'
                        defaultMessage='Tomorrow'
                    />
                ),
            },
            {
                id: 'submenu.nextweek',
                text: (
                    <FormattedMessage
                        id='submenu.nextweek'
                        key='submenu.nextweek'
                        defaultMessage='Next week'
                    />
                ),
            },
        ],
    };
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
