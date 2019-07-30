
import en from 'i18n/en.json';

import es from 'i18n/es.json';

import {id as pluginId} from './manifest';

import {
    postDropdownMenuAction,
} from './actions';

function getTranslations(locale) {
    switch (locale) {
    case 'en':
        return en;
    case 'es':
        return es;
    }
    return {};
}

function getSubMenu() {
    const primary = 'Remind me about this';
    const secondary = [
        'In 20 minutes',
        'In 1 hour',
        'In 3 hours',
        'Tomorrow',
        'Next week',
    ];
    return {primary, secondary};
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
